package controller

import (
	"context"

	"github.com/google/uuid"
	"github.com/umisto/agglomerations-svc/internal/domain/models"
	"github.com/umisto/agglomerations-svc/internal/domain/modules/role"
	"github.com/umisto/logium"
	"github.com/umisto/pagi"
)

type RoleController struct {
	domain Role
	log    logium.Logger
}
