package requests

import (
	"strings"

	"github.com/NdoleStudio/httpsms/pkg/entities"

	"github.com/NdoleStudio/httpsms/pkg/repositories"

	"github.com/NdoleStudio/httpsms/pkg/services"
)

// MessageThreadIndex is the payload fetching entities.MessageThread sent between 2 numbers
type MessageThreadIndex struct {
	request
	IsArchived string `json:"is_archived" query:"is_archived" example:"false"`
	Skip       string `json:"skip" query:"skip"`
	Query      string `json:"query" query:"query"`
	Limit      string `json:"limit" query:"limit"`
	Owner      string `json:"owner" query:"owner"`
}

// Sanitize sets defaults to MessageOutstanding
func (input *MessageThreadIndex) Sanitize() MessageThreadIndex {
	if strings.TrimSpace(input.Limit) == "" {
		input.Limit = "20"
	}

	if strings.TrimSpace(input.IsArchived) == "" {
		input.IsArchived = "false"
	}

	input.IsArchived = input.sanitizeBool(input.IsArchived)
	input.Query = strings.TrimSpace(input.Query)
	input.Owner = input.sanitizeAddress(input.Owner)

	input.Skip = strings.TrimSpace(input.Skip)
	if input.Skip == "" {
		input.Skip = "0"
	}

	return *input
}

// ToGetParams converts MessageThreadIndex into services.MessageThreadGetParams
func (input *MessageThreadIndex) ToGetParams(userID entities.UserID) services.MessageThreadGetParams {
	return services.MessageThreadGetParams{
		IndexParams: repositories.IndexParams{
			Skip:  input.getInt(input.Skip),
			Query: input.Query,
			Limit: input.getInt(input.Limit),
		},
		UserID:     userID,
		IsArchived: input.getBool(input.IsArchived),
		Owner:      input.Owner,
	}
}
