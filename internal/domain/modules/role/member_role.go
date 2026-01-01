package role

import (
	"context"
	"fmt"

	"github.com/google/uuid"
	"github.com/umisto/agglomerations-svc/internal/domain/errx"
	"github.com/umisto/agglomerations-svc/internal/domain/models"
)

func (s Service) GetMemberRoles(ctx context.Context, memberID uuid.UUID) ([]models.Role, error) {
	roles, err := s.repo.GetMemberRoles(ctx, memberID)
	if err != nil {
		return nil, errx.ErrorInternal.Raise(
			fmt.Errorf("failed to get member roles: %w", err),
		)
	}

	return roles, nil
}

func (s Service) AddMemberRoleByUser(
	ctx context.Context,
	accountID, memberID, roleID uuid.UUID,
) error {
	member, err := s.repo.GetMember(ctx, memberID)
	if err != nil {
		return errx.ErrorInternal.Raise(
			fmt.Errorf("failed to get member: %w", err),
		)
	}
	if member.IsNil() {
		return errx.ErrorMemberNotFound.Raise(
			fmt.Errorf("member with id %s not found", memberID),
		)
	}

	initiator, err := s.repo.GetMemberByAccountAndAgglomeration(ctx, accountID, member.AgglomerationID)
	if err != nil {
		return errx.ErrorInternal.Raise(
			fmt.Errorf("failed to get initiator member: %w", err),
		)
	}
	if initiator.IsNil() {
		return errx.ErrorNotEnoughRights.Raise(
			fmt.Errorf("initiator member with account id %s and agglomeration id %s not found",
				accountID, member.AgglomerationID),
		)
	}

	access, err := s.repo.CanInteract(ctx, initiator.ID, roleID)
	if err != nil {
		return errx.ErrorInternal.Raise(
			fmt.Errorf("failed to check permissions: %w", err),
		)
	}
	if !access {
		return errx.ErrorNotEnoughRights.Raise(
			fmt.Errorf("member %s does not have enough rights to assign role %s", initiator.ID, roleID),
		)
	}

	if err = s.repo.AddMemberRole(ctx, memberID, roleID); err != nil {
		return errx.ErrorInternal.Raise(
			fmt.Errorf("role Service AddMemberRoleByUser: repo AddMemberRole: %w", err),
		)
	}

	return nil
}

func (s Service) DeleteMemberRoleByUser(
	ctx context.Context,
	accountID, memberID, roleID uuid.UUID,
) error {
	member, err := s.repo.GetMember(ctx, memberID)
	if err != nil {
		return errx.ErrorInternal.Raise(
			fmt.Errorf("failed to get member: %w", err),
		)
	}
	if member.IsNil() {
		return errx.ErrorMemberNotFound.Raise(
			fmt.Errorf("member with id %s not found", memberID),
		)
	}

	initiator, err := s.repo.GetMemberByAccountAndAgglomeration(ctx, accountID, member.AgglomerationID)
	if err != nil {
		return errx.ErrorInternal.Raise(
			fmt.Errorf("failed to get initiator member: %w", err),
		)
	}
	if initiator.IsNil() {
		return errx.ErrorNotEnoughRights.Raise(
			fmt.Errorf("initiator member with account id %s and agglomeration id %s not found",
				accountID, member.AgglomerationID),
		)
	}

	access, err := s.repo.CanInteract(ctx, initiator.ID, roleID)
	if err != nil {
		return errx.ErrorInternal.Raise(
			fmt.Errorf("failed to check permissions: %w", err),
		)
	}
	if !access {
		return errx.ErrorNotEnoughRights.Raise(
			fmt.Errorf("member %s does not have enough rights to remove role %s", initiator.ID, roleID),
		)
	}

	if err = s.repo.DeleteMemberRole(ctx, memberID, roleID); err != nil {
		return errx.ErrorInternal.Raise(
			fmt.Errorf("role Service DeleteMemberRoleByUser: repo DeleteMemberRole: %w", err),
		)
	}

	return nil
}
