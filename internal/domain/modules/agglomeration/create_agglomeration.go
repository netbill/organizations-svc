package agglomeration

import (
	"context"
	"fmt"

	"github.com/google/uuid"
	"github.com/umisto/cities-svc/internal/domain/errx"
	"github.com/umisto/cities-svc/internal/domain/models"
)

type CreateParams struct {
	Name string
	Icon *string
}

func (s Service) CreateAgglomeration(
	ctx context.Context,
	accountID uuid.UUID,
	params CreateParams,
) (agglo models.Agglomeration, err error) {
	if err = s.repo.Transaction(ctx, func(ctx context.Context) error {
		agglo, err = s.repo.CreateAgglomeration(ctx, params)
		if err != nil {
			return errx.ErrorInternal.Raise(
				fmt.Errorf("failed to create agglomeration: %w", err),
			)
		}

		err = s.messenger.WriteAgglomerationCreated(ctx, agglo)
		if err != nil {
			return errx.ErrorInternal.Raise(
				fmt.Errorf("failed to publish agglomeration create event: %w", err),
			)
		}

		role, err := s.createRoleHead(ctx, agglo.ID)
		if err != nil {
			return err
		}

		_, err = s.createMemberHead(ctx, accountID, agglo.ID, role.ID)
		if err != nil {
			return err
		}

		return nil
	}); err != nil {
		return models.Agglomeration{}, err
	}

	return agglo, nil
}

func (s Service) createRoleHead(ctx context.Context, agglomerationID uuid.UUID) (role models.Role, err error) {
	role, err = s.repo.CreateHeadRole(ctx, agglomerationID)
	if err != nil {
		return models.Role{}, errx.ErrorInternal.Raise(
			fmt.Errorf("failed to create role: %w", err),
		)
	}

	err = s.messenger.WriteRoleCreated(ctx, role)
	if err != nil {
		return models.Role{}, errx.ErrorInternal.Raise(
			fmt.Errorf("failed to publish role create event: %w", err),
		)
	}

	return role, nil
}

func (s Service) createMemberHead(
	ctx context.Context,
	accountID uuid.UUID,
	agglomerationID uuid.UUID,
	roleID uuid.UUID,
) (member models.Member, err error) {
	member, err = s.repo.CreateMember(ctx, accountID, agglomerationID)
	if err != nil {
		return models.Member{}, errx.ErrorInternal.Raise(
			fmt.Errorf("failed to create member: %w", err),
		)
	}

	err = s.repo.AddMemberRole(ctx, member.ID, roleID)
	if err != nil {
		return models.Member{}, errx.ErrorInternal.Raise(
			fmt.Errorf("failed to assign head role to member: %w", err),
		)
	}

	return member, nil
}
