package city

import (
	"context"
	"fmt"

	"github.com/google/uuid"
	"github.com/paulmach/orb"
	"github.com/umisto/cities-svc/internal/domain/entity"
	"github.com/umisto/cities-svc/internal/domain/errx"
	"github.com/umisto/pagi"
)

type Service struct {
	repo      repo
	messanger messanger
}

type repo interface {
	GetAgglomerationByID(ctx context.Context, ID uuid.UUID) (entity.Agglomeration, error)

	CreateCity(ctx context.Context, params CreateParams) error

	GetCityByID(ctx context.Context, ID uuid.UUID) (entity.City, error)
	GetCityBySlug(ctx context.Context, slug string) (entity.City, error)

	UpdateCity(ctx context.Context, ID uuid.UUID, params UpdateParams) error

	UpdateCityStatus(ctx context.Context, ID uuid.UUID, status string) error

	FilterCities(
		ctx context.Context,
		params FilterParams,
		pagination pagi.Params,
	) (pagi.Page[[]entity.City], error)
	FilterCitiesNearest(
		ctx context.Context,
		point orb.Point,
		filter FilterParams,
		pagination pagi.Params,
	) (pagi.Page[map[int64]entity.City], error)

	CheckAccountHavePermissionByCode(ctx context.Context, accountID, agglomerationID uuid.UUID, permissionKey string) (bool, error)

	Transaction(ctx context.Context, fn func(ctx context.Context) error) error
}

type messanger interface {
	CreateCity(ctx context.Context, city entity.City) error

	UpdateCity(ctx context.Context, city entity.City) error

	ActivateCity(ctx context.Context, city entity.City) error
	DeactivateCity(ctx context.Context, city entity.City) error
	ArchivedCity(ctx context.Context, city entity.City) error

	DeleteCity(ctx context.Context, city entity.City) error
}

func (s Service) checkAgglomerationIsActiveAndExists(ctx context.Context, agglomerationID uuid.UUID) (entity.Agglomeration, error) {
	agglo, err := s.repo.GetAgglomerationByID(ctx, agglomerationID)
	if err != nil {
		return entity.Agglomeration{}, errx.ErrorInternal.Raise(
			fmt.Errorf("failed to get agglomeration by id: %w", err),
		)
	}
	if agglo.IsNil() {
		return entity.Agglomeration{}, errx.ErrorAgglomerationNotFound.Raise(
			fmt.Errorf("agglomeration with id %s not found", agglomerationID),
		)
	}

	if agglo.Status != entity.AgglomerationStatusActive {
		return entity.Agglomeration{}, errx.ErrorAgglomerationIsNotActive.Raise(
			fmt.Errorf("agglomeration with id %s is not active", agglomerationID),
		)
	}

	return agglo, nil
}

func (s Service) checkSlugIsAvailable(ctx context.Context, slug string) error {
	existingCity, err := s.repo.GetCityBySlug(ctx, slug)
	if err != nil {
		return errx.ErrorInternal.Raise(
			fmt.Errorf("failed to get city by slug: %w", err),
		)
	}
	if !existingCity.IsNil() {
		return errx.ErrorCityWithSlugAlreadyExists.Raise(
			fmt.Errorf("city with slug %s already exists", slug),
		)
	}
	return nil
}

func (s Service) checkPermissionByCode(
	ctx context.Context,
	accountID uuid.UUID,
	agglomerationID uuid.UUID,
	permissionKey entity.CodeRolePermission,
) error {
	access, err := s.repo.CheckAccountHavePermissionByCode(
		ctx,
		accountID,
		agglomerationID,
		permissionKey.String(),
	)
	if err != nil {
		return errx.ErrorInternal.Raise(
			fmt.Errorf("failed to check initiator permissions: %w", err))
	}
	if !access {
		return errx.ErrorNotEnoughRights.Raise(
			fmt.Errorf("initiator has no access to activate agglomeration"),
		)
	}

	return nil
}
