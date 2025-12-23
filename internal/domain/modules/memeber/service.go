package memeber

import (
	"github.com/google/uuid"
)

type FilterParams struct {
	AgglomerationID *uuid.UUID
	AccountID       *uuid.UUID
	Username        *string
	RoleID          *uuid.UUID
	PermissionCode  *string
}
