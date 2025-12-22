package city

import (
	"github.com/google/uuid"
)

type CreateParams struct {
	AgglomerationID *uuid.UUID
	Slug            *string
	Name            string
	Icon            *string
	Banner          *string
}

type UpdateParams struct {
	AgglomerationID *uuid.UUID
	Name            *string
	Slug            *string
	Icon            *string
	Banner          *string
}

type FilterParams struct {
	AgglomerationID *uuid.UUID
	Status          *string
	NameLike        *string
}
