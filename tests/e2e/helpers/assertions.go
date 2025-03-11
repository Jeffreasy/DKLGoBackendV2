package helpers

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"testing"

	"github.com/stretchr/testify/assert"
)

// AssertJSONResponse controleert of de response een geldige JSON response is
// en decodeert deze naar de gegeven struct
func AssertJSONResponse(t *testing.T, resp *http.Response, expectedStatus int, target interface{}) {
	assert.Equal(t, expectedStatus, resp.StatusCode, "Unexpected status code")
	assert.Equal(t, "application/json", resp.Header.Get("Content-Type"), "Content-Type is not application/json")

	body, err := io.ReadAll(resp.Body)
	assert.NoError(t, err, "Failed to read response body")

	err = json.Unmarshal(body, target)
	assert.NoError(t, err, fmt.Sprintf("Failed to decode JSON response: %s", string(body)))
}

// AssertEmailReceived controleert of een email is ontvangen door de gegeven ontvanger
func AssertEmailReceived(t *testing.T, mailClient *MailhogClient, recipient string, subjectContains string) {
	emails, err := mailClient.GetEmailsTo(recipient)
	assert.NoError(t, err, "Failed to get emails")

	found := false
	for _, email := range emails {
		if contains(email.Subject, subjectContains) {
			found = true
			break
		}
	}

	assert.True(t, found, fmt.Sprintf("No email with subject containing '%s' found for recipient '%s'", subjectContains, recipient))
}

// contains controleert of een string een andere string bevat
func contains(s, substr string) bool {
	return s != "" && substr != "" && s != substr && len(s) >= len(substr) && s[0:len(substr)] == substr
}
