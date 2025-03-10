package email

import (
	"fmt"
	"strings"
)

func (s *EmailService) extractFormFields(content string) map[string]string {
	fields := make(map[string]string)
	lines := strings.Split(content, "\n")

	for _, line := range lines {
		line = strings.TrimSpace(line)
		if line == "" {
			continue
		}

		// Check for form field patterns
		patterns := []string{" : ", ": ", ":", " :", " = ", "="}
		for _, pattern := range patterns {
			if parts := strings.SplitN(line, pattern, 2); len(parts) == 2 {
				key := strings.TrimSpace(parts[0])
				value := strings.TrimSpace(parts[1])
				if key != "" && value != "" {
					fields[key] = value
					break
				}
			}
		}
	}

	return fields
}

func (s *EmailService) formatFormFields(fields map[string]string) string {
	var formattedLines []string

	// Define field order for better readability
	fieldOrder := []string{
		"Naam",
		"E-mail",
		"Telefoonnummer",
		"Geboortedatum",
		"Geslacht",
		"Adres",
		"Postcode",
		"Woonplaats",
		"Afstand",
		"Vereniging",
		"Inschrijving_voor",
		"Betaalmethode",
		"lange_tekst",
		"Heb_je_een_vraag_of_opmerking_neem_dan_contact_op_met_ons",
	}

	// First add fields in the preferred order
	for _, key := range fieldOrder {
		if value, exists := fields[key]; exists {
			formattedLines = append(formattedLines, fmt.Sprintf("%s : %s", key, value))
			delete(fields, key)
		}
	}

	// Then add any remaining fields
	for key, value := range fields {
		formattedLines = append(formattedLines, fmt.Sprintf("%s : %s", key, value))
	}

	return strings.Join(formattedLines, "\n")
}
