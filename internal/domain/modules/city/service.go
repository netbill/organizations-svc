package city

import (
	"context"
	"fmt"

	"github.com/google/uuid"
	"github.com/paulmach/orb"
	"github.com/umisto/cities-svc/internal/domain/errx"
	"github.com/umisto/cities-svc/internal/domain/models"
	"github.com/umisto/pagi"
)

type Service struct {
	repo      repo
	messanger messanger
}

func New(repo repo, messanger messanger) Service {
	return Service{
		repo:      repo,
		messanger: messanger,
	}
}

type repo interface {
	GetAgglomerationByID(ctx context.Context, ID uuid.UUID) (models.Agglomeration, error)

	CreateCity(ctx context.Context, params CreateParams) (models.City, error)

	GetCityByID(ctx context.Context, ID uuid.UUID) (models.City, error)
	GetCityBySlug(ctx context.Context, slug string) (models.City, error)

	UpdateCity(ctx context.Context, ID uuid.UUID, params UpdateParams) (models.City, error)
	UpdateCityStatus(ctx context.Context, ID uuid.UUID, status string) (models.City, error)
	UpdateCitySlug(ctx context.Context, ID uuid.UUID, slug *string) (models.City, error)
	UpdateCityAgglomeration(
		ctx context.Context,
		cityID uuid.UUID,
		agglomerationID *uuid.UUID,
	) (models.City, error)

	FilterCities(
		ctx context.Context,
		params FilterParams,
		offset, limit uint,
	) (pagi.Page[[]models.City], error)
	FilterCitiesNearest(
		ctx context.Context,
		filter FilterParams,
		point orb.Point,
		offset, limit uint,
	) (pagi.Page[map[float64]models.City], error)

	CheckAccountHavePermissionByCode(
		ctx context.Context,
		accountID, agglomerationID uuid.UUID,
		permissionKey string,
	) (bool, error)

	Transaction(ctx context.Context, fn func(ctx context.Context) error) error
}

type messanger interface {
	WriteCityCreated(ctx context.Context, city models.City) error

	WriteCityUpdated(ctx context.Context, city models.City) error

	WriteCitySlugUpdated(ctx context.Context, city models.City, newSlug *string) error
	WriteCityAgglomerationUpdated(ctx context.Context, city models.City, newAggloId *uuid.UUID) error

	WriteCityActivated(ctx context.Context, city models.City) error
	WriteCityDeactivated(ctx context.Context, city models.City) error
	WriteCityArchived(ctx context.Context, city models.City) error

	WriteCityDeleted(ctx context.Context, city models.City) error
}

func (s Service) checkAgglomerationIsActiveAndExists(ctx context.Context, agglomerationID uuid.UUID) (models.Agglomeration, error) {
	agglo, err := s.repo.GetAgglomerationByID(ctx, agglomerationID)
	if err != nil {
		return models.Agglomeration{}, errx.ErrorInternal.Raise(
			fmt.Errorf("failed to get agglomeration by id: %w", err),
		)
	}
	if agglo.IsNil() {
		return models.Agglomeration{}, errx.ErrorAgglomerationNotFound.Raise(
			fmt.Errorf("agglomeration with id %s not found", agglomerationID),
		)
	}

	if agglo.Status != models.AgglomerationStatusActive {
		return models.Agglomeration{}, errx.ErrorAgglomerationIsNotActive.Raise(
			fmt.Errorf("agglomeration with id %s is not active", agglomerationID),
		)
	}

	return agglo, nil
}

func (s Service) checkPermissionForManageCity(
	ctx context.Context,
	accountID uuid.UUID,
	agglomerationID uuid.UUID,
) error {
	access, err := s.repo.CheckAccountHavePermissionByCode(
		ctx,
		accountID,
		agglomerationID,
		models.RolePermissionManageCities.String(),
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
