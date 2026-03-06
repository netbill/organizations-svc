package place

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/netbill/organizations-svc/internal/errx"
)

type Service struct {
	repo      repo
	tombstone tombstone
	tx        transactor
}

type ServiceDeps struct {
	Repo      repo
	Tombstone tombstone
	Tx        transactor
}

func NewPlaceModule(deps ServiceDeps) *Service {
	return &Service{
		repo:      deps.Repo,
		tombstone: deps.Tombstone,
		tx:        deps.Tx,
	}
}

type CreateParams struct {
	ID             uuid.UUID `json:"id"`
	OrganizationID uuid.UUID `json:"organization_id"`
	CreatedAt      time.Time `json:"created_at"`
}

func (s *Service) Create(
	ctx context.Context,
	params CreateParams,
) error {
	exists, err := s.repo.Exists(ctx, params.ID)
	if err != nil {
		return err
	}
	if exists {
		return errx.ErrorPlaceAlreadyExists.Raise(
			errors.New("place with given id already exists"),
		)
	}

	buried, err := s.tombstone.PlaceIsBuried(ctx, params.ID)
	if err != nil {
		return err
	}
	if buried {
		return errx.ErrorPlaceDeleted.Raise(
			errors.New("place with given id is deleted"),
		)
	}

	return s.repo.Create(ctx, params)
}

func (s *Service) Delete(
	ctx context.Context,
	id uuid.UUID,
) error {
	buried, err := s.tombstone.PlaceIsBuried(ctx, id)
	if err != nil {
		return err
	}
	if buried {
		return errx.ErrorPlaceDeleted.Raise(
			fmt.Errorf("place with id %s is already deleted", id),
		)
	}

	return s.tx.Transaction(ctx, func(ctx context.Context) error {
		if err := s.tombstone.BuryPlace(ctx, id); err != nil {
			return err
		}

		return s.repo.Delete(ctx, id)
	})
}
