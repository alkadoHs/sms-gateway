package repositories

import (
	"context"

	"github.com/NdoleStudio/httpsms/pkg/entities"
)

// Integration3CxRepository loads and persists an entities.Integration3CX
type Integration3CxRepository interface {
	// Save an entities.Integration3CX
	Save(ctx context.Context, heartbeat *entities.Integration3CX) error

	// Load an entities.Integration3CX based on the entities.UserID
	Load(ctx context.Context, userID entities.UserID) (*entities.Integration3CX, error)

	// DeleteAllForUser deletes all entities.Integration3CX for a user
	DeleteAllForUser(ctx context.Context, userID entities.UserID) error
}
