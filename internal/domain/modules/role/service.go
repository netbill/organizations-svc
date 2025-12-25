package role

import (
	"github.com/google/uuid"
)

type CreateParams struct {
	ID              uuid.UUID
	AgglomerationID uuid.UUID
	Head            bool
	Editable        bool
	Rank            int32
	Name            string
}

type FilterParams struct {
	AgglomerationID uuid.UUID
	MemberID        *uuid.UUID
	PermissionCodes []string
}

type UpdateParams struct {
	Name *string
}

type FilterPermissionsParams struct {
	Description *string
	Code        *string
}
