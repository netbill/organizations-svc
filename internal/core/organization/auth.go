package organization

import (
	"context"
	"errors"
	"fmt"

	"github.com/google/uuid"
	"github.com/netbill/organizations-svc/internal/errx"
	"github.com/netbill/organizations-svc/internal/models"
)

func (s *Service) AuthorizeOrgHead(
	ctx context.Context,
	actor models.AccountActor,
	organizationID uuid.UUID,
) (models.Member, error) {
	member, err := s.member.GetForAccountAndOrg(ctx, actor, organizationID)
	if errors.Is(err, errx.ErrorMemberNotExists) {
		return models.Member{}, errx.ErrorInitiatorNotMemberOfOrganization.Raise(
			fmt.Errorf("initiator with account id %s is not a member of organization %s", actor, organizationID),
		)
	}
	if err != nil {
		return models.Member{}, err
	}

	if !member.Head {
		return models.Member{}, errx.ErrorNotOrganizationHead.Raise(
			fmt.Errorf(
				"only organization head member can manage organization, but member %s is not head", member.ID,
			),
		)
	}

	return member, nil
}

func (s *Service) AuthorizeOrgMember(
	ctx context.Context,
	actor models.AccountActor,
	organizationID uuid.UUID,
) (models.Member, error) {
	initiator, err := s.member.GetForAccountAndOrg(ctx, actor, organizationID)
	if errors.Is(err, errx.ErrorMemberNotExists) {
		return models.Member{}, errx.ErrorInitiatorNotMemberOfOrganization.Raise(
			fmt.Errorf("initiator with account id %s is not a member of organization %s", actor, organizationID),
		)
	}
	if err != nil {
		return models.Member{}, err
	}

	return initiator, nil
}

func (s *Service) ValidateOrg(ctx context.Context, organizationID uuid.UUID) (models.Organization, error) {
	org, err := s.GetByID(ctx, organizationID)
	if err != nil {
		return models.Organization{}, err
	}

	if org.Status == models.OrganizationStatusSuspended {
		return models.Organization{}, errx.ErrorOrganizationIsSuspended.Raise(
			fmt.Errorf("organization with id %s is suspended", organizationID),
		)
	}

	return org, nil
}
