package email

import (
	"bytes"
	"dklautomationgo/models"
	"fmt"
	"io"
	"log"
	"strings"
	"time"
	"unicode/utf8"

	"github.com/emersion/go-imap"
	"github.com/emersion/go-message"
	"golang.org/x/text/encoding/charmap"
)

func (s *EmailService) processMessage(msg *imap.Message, accountName string) (*models.Email, error) {
	// Only log message ID and subject
	if msg.Envelope != nil && msg.Envelope.Subject != "" {
		log.Printf("[EMAIL] Processing: %s", msg.Envelope.Subject)
	}

	// Get the message body
	section := &imap.BodySectionName{}
	r := msg.GetBody(section)
	if r == nil {
		return nil, fmt.Errorf("server didn't return message body")
	}

	// Read the message body
	var buf bytes.Buffer
	if _, err := io.Copy(&buf, r); err != nil {
		return nil, fmt.Errorf("failed to read message body: %w", err)
	}

	// Create mail reader
	mr, err := message.Read(&buf)
	if err != nil {
		// Try content recovery if mail reader fails
		body := buf.Bytes()
		if !utf8.Valid(body) {
			if decoded, err := charmap.ISO8859_1.NewDecoder().Bytes(body); err == nil {
				body = decoded
			}
		}
		return &models.Email{
			ID:        fmt.Sprintf("%s:%d", accountName, msg.Uid),
			Subject:   msg.Envelope.Subject,
			Body:      string(body),
			Account:   accountName,
			CreatedAt: msg.Envelope.Date.Format(time.RFC3339),
			Read:      hasFlag(msg.Flags, imap.SeenFlag),
		}, nil
	}

	// Extract headers
	headers := make(map[string]string)
	commonHeaders := []string{
		"From", "To", "Cc", "Subject", "Date",
		"Message-ID", "In-Reply-To", "References",
		"Content-Type", "Content-Transfer-Encoding",
	}
	for _, name := range commonHeaders {
		if value := mr.Header.Get(name); value != "" {
			headers[name] = value
		}
	}

	// Get message parts
	var textBody, htmlBody string
	var attachments []models.EmailAttachment

	// Process each part
	mpr := mr.MultipartReader()
	if mpr != nil {
		for {
			part, err := mpr.NextPart()
			if err == io.EOF {
				break
			}
			if err != nil {
				continue
			}

			contentType := part.Header.Get("Content-Type")
			contentDisposition := part.Header.Get("Content-Disposition")
			contentID := part.Header.Get("Content-ID")

			// Read part content
			var partBuf bytes.Buffer
			if _, err := io.Copy(&partBuf, part.Body); err != nil {
				continue
			}
			content := partBuf.String()

			// Handle different content types
			switch {
			case strings.HasPrefix(contentType, "text/plain"):
				textBody = content
			case strings.HasPrefix(contentType, "text/html"):
				htmlBody = content
			case contentDisposition != "" && strings.Contains(contentDisposition, "attachment"):
				filename := s.extractFilename(contentDisposition, contentType)
				attachments = append(attachments, models.EmailAttachment{
					Filename:    filename,
					ContentType: contentType,
					Size:        int64(len(content)),
					Content:     partBuf.Bytes(),
					ContentID:   strings.Trim(contentID, "<>"),
				})
			}
		}
	} else {
		// Single part message
		var buf bytes.Buffer
		io.Copy(&buf, mr.Body)
		textBody = buf.String()
	}

	// Create email object
	email := &models.Email{
		ID:          fmt.Sprintf("%s:%d", accountName, msg.Uid),
		Sender:      msg.Envelope.From[0].Address(),
		Subject:     msg.Envelope.Subject,
		Body:        textBody,
		HTML:        htmlBody,
		Account:     accountName,
		MessageID:   msg.Envelope.MessageId,
		CreatedAt:   msg.Envelope.Date.Format(time.RFC3339),
		Read:        hasFlag(msg.Flags, imap.SeenFlag),
		Headers:     headers,
		Attachments: attachments,
	}

	// Only log total attachments if there are any
	if len(email.Attachments) > 0 {
		log.Printf("[EMAIL] Found %d attachments", len(email.Attachments))
	}

	return email, nil
}

func hasFlag(flags []string, flag string) bool {
	for _, f := range flags {
		if f == flag {
			return true
		}
	}
	return false
}

func (s *EmailService) extractFilename(disposition, contentType string) string {
	// Try Content-Disposition first
	if disposition != "" {
		if start := strings.Index(disposition, "filename="); start != -1 {
			filename := disposition[start+9:]
			if strings.HasPrefix(filename, `"`) {
				if end := strings.Index(filename[1:], `"`); end != -1 {
					return filename[1 : end+1]
				}
			}
			if end := strings.Index(filename, `;`); end != -1 {
				return filename[:end]
			}
			return filename
		}
	}

	// Try Content-Type as fallback
	if contentType != "" {
		if start := strings.Index(contentType, "name="); start != -1 {
			filename := contentType[start+5:]
			if strings.HasPrefix(filename, `"`) {
				if end := strings.Index(filename[1:], `"`); end != -1 {
					return filename[1 : end+1]
				}
			}
			if end := strings.Index(filename, `;`); end != -1 {
				return filename[:end]
			}
			return filename
		}
	}

	return "unnamed-attachment"
}
