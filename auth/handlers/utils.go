package handlers

import (
	"github.com/google/uuid"
)

// parseUUID parst een string naar een UUID
func parseUUID(id string) (uuid.UUID, error) {
	return uuid.Parse(id)
}
