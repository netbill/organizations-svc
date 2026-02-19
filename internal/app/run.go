package app

import (
	"context"
	"sync"

	eventpg "github.com/netbill/eventbox/pg"
	"github.com/netbill/organizations-svc/internal/core/modules/invite"
	"github.com/netbill/organizations-svc/internal/core/modules/member"
	"github.com/netbill/organizations-svc/internal/core/modules/organization"
	"github.com/netbill/organizations-svc/internal/core/modules/profile"
	"github.com/netbill/organizations-svc/internal/core/modules/role"
	"github.com/netbill/organizations-svc/internal/messenger/inbound/handlers"
	"github.com/netbill/organizations-svc/internal/messenger/outbound/sender"
	"github.com/netbill/organizations-svc/internal/repository"
	"github.com/netbill/organizations-svc/internal/repository/pg"
	"github.com/netbill/organizations-svc/internal/rest/controller"
)

func (a *App) Run(ctx context.Context) {
	var wg sync.WaitGroup
	run := func(f func()) {
		wg.Add(1)
		go func() {
			defer wg.Done()
			f()
		}()
	}

	repo := &repository.Repository{
		Transactioner:             pg.NewTransaction(a.db),
		OrganizationsSql:          pg.NewOrganizationsQ(a.db),
		OrgMembersSql:             pg.NewOrgMembersQ(a.db),
		OrgMemberRolesSql:         pg.NewOrgMemberRolesQ(a.db),
		OrgRolesSql:               pg.NewOrgRolesQ(a.db),
		OrgRolePermissionLinksSql: pg.NewOrgRolePermissionLinksQ(a.db),
		OrgRolePermissionsSql:     pg.NewOrgRolePermissionsQ(a.db),
		OrgInvitesSql:             pg.NewOrgInvitesQ(a.db),
		ProfilesSql:               pg.NewProfilesQ(a.db),
	}

	producer := a.BuildProducer()
	defer producer.Close()

	kafkaProducer := sender.New("", eventpg.NewOutbox(a.db), producer)

	bucket := a.BuildBucket()

	orgSvc := organization.New(repo, kafkaProducer, bucket)
	memberSvc := member.New(repo, kafkaProducer)
	roleSvc := role.New(repo, kafkaProducer)
	inviteSvc := invite.New(repo, kafkaProducer)
	profileSvc := profile.New(repo)

	inbound := handlers.New(handlers.Modules{Profile: profileSvc})

	run(func() {
		a.RunRest(ctx, &controller.Modules{
			Organization: orgSvc,
			Role:         roleSvc,
			Invite:       inviteSvc,
			Member:       memberSvc,
		})
	})

	run(func() {
		a.RunOutbox(ctx, producer)
	})

	run(func() {
		a.RunInbox(ctx, inbound)
	})

	run(func() {
		a.RunConsumer(ctx)
	})
}
