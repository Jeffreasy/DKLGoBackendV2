package email

import (
	"context"
	"crypto/tls"
	"dklautomationgo/models"
	"fmt"
	"log"
	"sync"
	"time"

	"github.com/emersion/go-imap"
	"github.com/emersion/go-imap/client"
)

// AccountCache houdt de cache bij voor één account
type AccountCache struct {
	emails     []*models.Email
	lastFetch  time.Time
	cacheMutex sync.RWMutex
}

func (s *EmailService) FetchEmails(options *models.EmailFetchOptions) ([]*models.Email, error) {
	var allEmails []*models.Email
	var mu sync.Mutex
	var wg sync.WaitGroup
	errChan := make(chan error, len(s.config.Accounts))

	// Check cache voor elk account
	if s.config.Cache.Enabled {
		for _, cache := range s.accountCaches {
			cache.cacheMutex.RLock()
			if time.Since(cache.lastFetch) < s.config.Cache.Duration && len(cache.emails) > 0 {
				allEmails = append(allEmails, cache.emails...)
				cache.cacheMutex.RUnlock()
				continue
			}
			cache.cacheMutex.RUnlock()
		}

		if len(allEmails) == len(s.config.Accounts) {
			return s.filterEmails(allEmails, options), nil
		}
	}

	// Context met timeout voor alle operaties
	ctx, cancel := context.WithTimeout(context.Background(), s.config.FetchTimeout)
	defer cancel()

	// Parallel ophalen van emails voor elk account
	for accountName, config := range s.config.Accounts {
		wg.Add(1)
		go func(accName string, cfg *EmailConfig) {
			defer wg.Done()

			emails, err := s.fetchEmailsFromAccount(ctx, accName, cfg, options)
			if err != nil {
				log.Printf("[ERROR] Account %s: Failed to fetch emails: %v", accName, err)
				errChan <- fmt.Errorf("account %s: %w", accName, err)
				return
			}

			// Update cache en voeg emails toe aan resultaat
			if s.config.Cache.Enabled {
				cache := s.accountCaches[accName]
				cache.cacheMutex.Lock()
				cache.emails = emails
				cache.lastFetch = time.Now()
				cache.cacheMutex.Unlock()
			}

			mu.Lock()
			allEmails = append(allEmails, emails...)
			mu.Unlock()
		}(accountName, config)
	}

	wg.Wait()
	close(errChan)

	// Controleer op fouten
	var errors []error
	for err := range errChan {
		errors = append(errors, err)
	}

	// Als alle accounts faalden, geef een fout terug
	if len(errors) == len(s.config.Accounts) {
		return nil, fmt.Errorf("all accounts failed: %v", errors)
	}

	// Als er geen emails zijn gevonden, geef een lege lijst terug
	if len(allEmails) == 0 {
		return []*models.Email{}, nil
	}

	return s.filterEmails(allEmails, options), nil
}

func (s *EmailService) filterEmails(emails []*models.Email, options *models.EmailFetchOptions) []*models.Email {
	if options == nil {
		return emails
	}

	// Filter by read status if specified
	filtered := emails
	if options.Read != nil {
		temp := make([]*models.Email, 0)
		for _, email := range filtered {
			if email.Read == *options.Read {
				temp = append(temp, email)
			}
		}
		filtered = temp
	}

	// Apply offset and limit
	start := options.Offset
	if start >= len(filtered) {
		return []*models.Email{}
	}

	end := len(filtered)
	if options.Limit > 0 {
		end = start + options.Limit
		if end > len(filtered) {
			end = len(filtered)
		}
	}

	return filtered[start:end]
}

func (s *EmailService) fetchEmailsFromAccount(ctx context.Context, accountName string, config *EmailConfig, options *models.EmailFetchOptions) ([]*models.Email, error) {
	log.Printf("[EMAIL] %s: Starting fetch", accountName)

	// Maak IMAP client met context en verbeterde TLS configuratie
	tlsConfig := &tls.Config{
		ServerName:         config.IMAPHost,
		InsecureSkipVerify: true,
		MinVersion:         tls.VersionTLS12,
	}

	// Connect to IMAP server
	c, err := client.DialTLS(fmt.Sprintf("%s:%d", config.IMAPHost, config.IMAPPort), tlsConfig)
	if err != nil {
		return nil, fmt.Errorf("IMAP connection failed: %w", err)
	}
	defer func() {
		if err := c.Logout(); err != nil {
			log.Printf("[EMAIL] %s: Logout failed: %v", accountName, err)
		}
	}()

	// Check context before proceeding
	select {
	case <-ctx.Done():
		return nil, ctx.Err()
	default:
	}

	// Login
	if err := c.Login(config.Email, config.Password); err != nil {
		return nil, fmt.Errorf("IMAP login failed: %w", err)
	}

	// Select INBOX
	mbox, err := c.Select("INBOX", false)
	if err != nil {
		return nil, fmt.Errorf("IMAP select inbox failed: %w", err)
	}

	if mbox.Messages == 0 {
		log.Printf("[EMAIL] %s: Inbox empty", accountName)
		return []*models.Email{}, nil
	}

	// Calculate message range
	from := uint32(1)
	to := mbox.Messages

	if options != nil && options.Limit > 0 {
		if uint32(options.Offset) >= to {
			return []*models.Email{}, nil
		}
		from = to - uint32(options.Limit+options.Offset)
		if from < 1 {
			from = 1
		}
		to = to - uint32(options.Offset)
	}

	log.Printf("[EMAIL] %s: Fetching messages %d-%d of %d", accountName, from, to, mbox.Messages)

	seqSet := new(imap.SeqSet)
	seqSet.AddRange(from, to)

	// Define fetch items
	items := []imap.FetchItem{
		imap.FetchEnvelope,
		imap.FetchFlags,
		imap.FetchBody,
		imap.FetchBodyStructure,
		"BODY[]",
	}

	messages := make(chan *imap.Message, 100)
	done := make(chan error, 1)

	go func() {
		done <- c.Fetch(seqSet, items, messages)
	}()

	var emails []*models.Email
	processedCount := 0
	errorCount := 0

	for msg := range messages {
		select {
		case <-ctx.Done():
			return nil, ctx.Err()
		default:
			email, err := s.processMessage(msg, accountName)
			if err != nil {
				errorCount++
				continue
			}
			emails = append(emails, email)
			processedCount++
		}
	}

	if err := <-done; err != nil {
		return nil, fmt.Errorf("IMAP fetch failed: %w", err)
	}

	log.Printf("[EMAIL] %s: Completed. Processed %d messages (%d errors)",
		accountName, processedCount, errorCount)

	return emails, nil
}
