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
	"github.com/emersion/go-message/mail"
)

func (s *EmailService) processMessage(msg *imap.Message, accountName string) (*models.Email, error) {
	log.Printf("=== Start processing message ===")
	log.Printf("Message ID: %s", msg.Envelope.MessageId)
	log.Printf("Subject: %s", msg.Envelope.Subject)
	log.Printf("From: %v", msg.Envelope.From)

	var body, html string
	metadata := make(map[string]string)
	headers := make(map[string]string)
	var attachments []models.EmailAttachment

	// Get the whole message body
	section := &imap.BodySectionName{}
	r := msg.GetBody(section)
	if r == nil {
		log.Printf("Error: Server didn't return message body")
		return nil, fmt.Errorf("server didn't return message body")
	}

	// Read the entire message body
	var buf bytes.Buffer
	if _, err := io.Copy(&buf, r); err != nil {
		log.Printf("Error: Failed to read message body: %v", err)
		return nil, fmt.Errorf("error reading message body: %v", err)
	}

	log.Printf("Raw message size: %d bytes", buf.Len())

	// Create a new reader from the buffer
	mr, err := mail.CreateReader(&buf)
	if err != nil {
		log.Printf("Warning: Error creating mail reader: %v", err)
		log.Printf("Attempting content recovery...")

		// Try to recover the content as plain text
		content := buf.String()

		// Try UTF-8 first
		if utf8.ValidString(content) {
			log.Printf("Content is valid UTF-8")
			body = content
		} else {
			log.Printf("Content is not valid UTF-8, attempting ISO-8859-1 decode...")
			decoded, err := s.decodeCharset(content, "iso-8859-1")
			if err != nil {
				log.Printf("Warning: Failed to decode as ISO-8859-1: %v", err)
				body = content
			} else {
				body = decoded
			}
		}
	} else {
		log.Printf("Successfully created mail reader")

		// Read all headers
		header := mr.Header
		commonHeaders := []string{
			"From", "To", "Cc", "Bcc", "Subject", "Date",
			"Message-ID", "In-Reply-To", "References",
			"Content-Type", "Content-Transfer-Encoding",
			"MIME-Version", "Received", "Return-Path",
			"Delivered-To", "Reply-To", "Sender",
			"Authentication-Results", "DKIM-Signature",
		}
		for _, name := range commonHeaders {
			if value := header.Get(name); value != "" {
				headers[name] = value
			}
		}

		// Extract common headers
		if date, err := header.Date(); err == nil {
			metadata["date"] = date.Format(time.RFC1123Z)
			log.Printf("Email Date: %s", metadata["date"])
		}
		if from, err := header.AddressList("From"); err == nil && len(from) > 0 {
			metadata["from"] = from[0].String()
			log.Printf("From Header: %s", metadata["from"])
		}
		if to, err := header.AddressList("To"); err == nil && len(to) > 0 {
			metadata["to"] = to[0].String()
			metadata["delivered-to"] = to[0].Address
			log.Printf("To Header: %s", metadata["to"])
		}
		if subject, err := header.Subject(); err == nil {
			metadata["subject"] = subject
			log.Printf("Subject Header: %s", metadata["subject"])
		}

		// Read message content
		var textPart, htmlPart string
		partCount := 0

		// Process each part of the message
		for {
			p, err := mr.NextPart()
			if err == io.EOF {
				break
			}
			if err != nil {
				log.Printf("Warning: Error reading message part: %v", err)
				continue
			}
			partCount++

			contentType := p.Header.Get("Content-Type")
			contentDisposition := p.Header.Get("Content-Disposition")
			contentID := p.Header.Get("Content-ID")
			contentCharset := s.extractCharset(contentType)

			log.Printf("Processing part %d - Type: %s, Disposition: %s, ID: %s, Charset: %s",
				partCount, contentType, contentDisposition, contentID, contentCharset)

			// Lees de inhoud van dit deel
			partContent, err := io.ReadAll(p.Body)
			if err != nil {
				log.Printf("Warning: Error reading part body: %v", err)
				continue
			}
			log.Printf("Part %d size: %d bytes", partCount, len(partContent))

			// Als dit een attachment is
			if contentDisposition != "" && strings.Contains(contentDisposition, "attachment") {
				filename := s.extractFilename(contentDisposition, contentType)
				attachment := models.EmailAttachment{
					Filename:    filename,
					ContentType: contentType,
					Size:        int64(len(partContent)),
					Content:     partContent,
					ContentID:   strings.Trim(contentID, "<>"),
				}
				attachments = append(attachments, attachment)
				log.Printf("Found attachment: %s (%d bytes)", filename, len(partContent))
				continue
			}

			// Als dit een inline image is
			if contentID != "" && strings.HasPrefix(contentType, "image/") {
				attachment := models.EmailAttachment{
					Filename:    fmt.Sprintf("inline-%s", strings.Trim(contentID, "<>")),
					ContentType: contentType,
					Size:        int64(len(partContent)),
					Content:     partContent,
					ContentID:   strings.Trim(contentID, "<>"),
				}
				attachments = append(attachments, attachment)
				log.Printf("Found inline image: %s (%d bytes)", contentID, len(partContent))
				continue
			}

			// Verwerk tekst content
			content := string(partContent)
			if contentCharset != "" {
				decoded, err := s.decodeCharset(content, contentCharset)
				if err != nil {
					log.Printf("Warning: Failed to decode charset %s: %v", contentCharset, err)
				} else {
					content = decoded
				}
			}

			// Store content based on type
			if strings.Contains(contentType, "text/html") {
				log.Printf("Found HTML content in part %d", partCount)
				htmlPart = content
			} else if strings.Contains(contentType, "text/plain") {
				log.Printf("Found plain text content in part %d", partCount)
				textPart = content
			}
		}

		log.Printf("Processed %d parts total", partCount)

		// Als we HTML hebben, gebruik dat als primaire inhoud
		if htmlPart != "" {
			log.Printf("Using HTML content as primary")
			html = htmlPart
			if textPart != "" {
				body = s.processPlainContent(textPart)
			} else {
				body = s.ProcessHTML(htmlPart)
			}
		} else if textPart != "" {
			log.Printf("Using plain text content")
			body = s.processPlainContent(textPart)
		}
	}

	// Clean up the body text
	if body != "" {
		originalLength := len(body)
		body = strings.TrimSpace(body)
		body = strings.ReplaceAll(body, "\r\n", "\n")
		body = strings.ReplaceAll(body, "\n\n\n", "\n\n")
		log.Printf("Body text cleanup: %d bytes -> %d bytes", originalLength, len(body))
	}

	// Extract form fields if this is a form submission
	formFields := s.extractFormFields(body)
	if len(formFields) > 0 {
		log.Printf("Found %d form fields", len(formFields))
		body = s.formatFormFields(formFields)
	}

	// Verzamel alle geadresseerden
	var toAddresses []string
	var ccAddresses []string
	var bccAddresses []string
	var replyToAddresses []string

	if msg.Envelope.To != nil {
		for _, addr := range msg.Envelope.To {
			toAddresses = append(toAddresses, addr.Address())
		}
	}
	if msg.Envelope.Cc != nil {
		for _, addr := range msg.Envelope.Cc {
			ccAddresses = append(ccAddresses, addr.Address())
		}
	}
	if msg.Envelope.Bcc != nil {
		for _, addr := range msg.Envelope.Bcc {
			bccAddresses = append(bccAddresses, addr.Address())
		}
	}
	if msg.Envelope.ReplyTo != nil {
		for _, addr := range msg.Envelope.ReplyTo {
			replyToAddresses = append(replyToAddresses, addr.Address())
		}
	}

	// Get References from headers if available
	var references []string
	if refsHeader := headers["References"]; refsHeader != "" {
		references = strings.Fields(refsHeader)
	}

	// Create the email object with all data
	email := &models.Email{
		ID:          fmt.Sprintf("%d", msg.Uid),
		Sender:      msg.Envelope.From[0].Address(),
		Subject:     msg.Envelope.Subject,
		Body:        body,
		HTML:        html,
		Account:     accountName,
		MessageID:   msg.Envelope.MessageId,
		CreatedAt:   msg.Envelope.Date.Format(time.RFC3339),
		Read:        hasFlag(msg.Flags, imap.SeenFlag),
		Metadata:    metadata,
		To:          toAddresses,
		Cc:          ccAddresses,
		Bcc:         bccAddresses,
		ReplyTo:     replyToAddresses,
		InReplyTo:   msg.Envelope.InReplyTo,
		References:  references,
		Attachments: attachments,
		Headers:     headers,
	}

	log.Printf("=== Finished processing message ===")
	log.Printf("Final body length: %d bytes", len(email.Body))
	log.Printf("Final HTML length: %d bytes", len(email.HTML))
	log.Printf("Number of attachments: %d", len(email.Attachments))

	return email, nil
}

func (s *EmailService) processPlainContent(content string) string {
	// Remove any potential HTML content
	if strings.Contains(content, "<") && strings.Contains(content, ">") {
		content = s.ProcessHTML(content)
	}

	// Clean up whitespace
	content = strings.TrimSpace(content)
	content = strings.ReplaceAll(content, "\r\n", "\n")
	content = strings.ReplaceAll(content, "\n\n\n", "\n\n")

	return content
}

func hasFlag(flags []string, flag string) bool {
	for _, f := range flags {
		if f == flag {
			return true
		}
	}
	return false
}

func (s *EmailService) ProcessContactForm(content string) (map[string]string, error) {
	// Process HTML content first
	content = s.ProcessHTML(content)

	// Extract form fields
	fields := make(map[string]string)
	lines := strings.Split(content, "\n")

	for _, line := range lines {
		line = strings.TrimSpace(line)
		if line == "" {
			continue
		}

		// Try different separators
		parts := strings.SplitN(line, ":", 2)
		if len(parts) != 2 {
			parts = strings.SplitN(line, "=", 2)
		}
		if len(parts) != 2 {
			continue
		}

		key := strings.TrimSpace(parts[0])
		value := strings.TrimSpace(parts[1])
		if key != "" && value != "" {
			fields[key] = value
		}
	}

	return fields, nil
}

func (s *EmailService) ExtractFormData(content string) (map[string]string, error) {
	// Clean HTML first
	content = s.ProcessHTML(content)

	// Extract form fields
	fields := make(map[string]string)
	lines := strings.Split(content, "\n")

	for _, line := range lines {
		line = strings.TrimSpace(line)
		if line == "" {
			continue
		}

		// Try different separators
		parts := strings.SplitN(line, ":", 2)
		if len(parts) != 2 {
			parts = strings.SplitN(line, "=", 2)
		}
		if len(parts) != 2 {
			continue
		}

		key := strings.TrimSpace(parts[0])
		value := strings.TrimSpace(parts[1])
		if key != "" && value != "" {
			fields[key] = value
		}
	}

	return fields, nil
}

func (s *EmailService) extractFilename(disposition, contentType string) string {
	// Probeer eerst uit Content-Disposition
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

	// Probeer uit Content-Type als fallback
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
