package log

import (
	"errors"
	"net/http"

	"github.com/google/uuid"
	"github.com/netbill/ape"
	"github.com/netbill/eventbox"
	"github.com/netbill/eventbox/headers"
	"github.com/netbill/organizations-svc/internal/core/models"
	"github.com/segmentio/kafka-go"
	"github.com/sirupsen/logrus"
)

const (
	ServiceField   = "service"
	OperationField = "operation"
	ComponentField = "component"

	AccountIDField        = "account_id"
	AccountSessionIDField = "account_session_id"

	OrganizationIDField = "organization_id"
	MemberIDField       = "org_member_id"
	RoleIDField         = "org_role_id"
	InviteIDField       = "org_invite_id"

	HTTPMethodField = "http_method"
	HTTPPathField   = "http_path"

	EventIDField       = "event_id"
	EventTypeField     = "event_type"
	EventTopicField    = "event_topic"
	EventVersionField  = "event_version"
	EventProducerField = "event_producer"
	EventAttemptField  = "event_attempt"
)

type Entry struct {
	*logrus.Entry
}

func New(
	level string,
	format string,
	serviceName string,
) *Entry {
	log := logrus.New()

	lvl, err := logrus.ParseLevel(level)
	if err != nil {
		lvl = logrus.InfoLevel
		log.WithField("bad_level", level).Warn("unknown log level, fallback to info")
	}

	log.SetLevel(lvl)

	switch {
	case format == "json":
		log.SetFormatter(&logrus.JSONFormatter{})
	default:
		log.SetFormatter(&logrus.TextFormatter{
			FullTimestamp:   true,
			TimestampFormat: "2006-01-02 15:04:05",
		})
	}

	return NewEntry(log.WithField(ServiceField, serviceName))
}

func NewEntry(entry *logrus.Entry) *Entry {
	return &Entry{entry}
}

func (e *Entry) WithFields(fields map[string]any) *Entry {
	return &Entry{Entry: e.Entry.WithFields(fields)}
}

func (e *Entry) WithField(key string, value any) *Entry {
	return &Entry{Entry: e.Entry.WithField(key, value)}
}

func (e *Entry) WithError(err error) *Entry {
	var ae *ape.Error
	if errors.As(err, &ae) {
		return &Entry{Entry: e.Entry.WithError(ae)}
	}
	return &Entry{Entry: e.Entry.WithError(err)}
}

func (e *Entry) WithRequest(r *http.Request) *Entry {
	return e.WithFields(map[string]any{
		HTTPMethodField: r.Method,
		HTTPPathField:   r.URL.Path,
	})
}

func (e *Entry) WithOperation(operation string) *Entry {
	return e.WithField(OperationField, operation)
}

func (e *Entry) WithComponent(component string) *Entry {
	return e.WithField(ComponentField, component)
}

type accountAuthClaims interface {
	GetAccountID() uuid.UUID
	GetSessionID() uuid.UUID
}

func (e *Entry) WithAccountAuthClaims(auth accountAuthClaims) *Entry {
	return e.WithFields(map[string]any{
		AccountIDField:        auth.GetAccountID(),
		AccountSessionIDField: auth.GetSessionID(),
	})
}

func (e *Entry) WithOrganization(organization models.Organization) *Entry {
	return e.WithFields(map[string]any{
		OrganizationIDField: organization.ID,
	})
}

func (e *Entry) WithMember(member models.Member) *Entry {
	return e.WithFields(map[string]any{
		OrganizationIDField: member.OrganizationID,
		MemberIDField:       member.ID,
	})
}

func (e *Entry) WithRole(role models.Role) *Entry {
	return e.WithFields(map[string]any{
		OrganizationIDField: role.OrganizationID,
		RoleIDField:         role.ID,
	})
}

func (e *Entry) WithInvite(invite models.Invite) *Entry {
	return e.WithFields(map[string]any{
		OrganizationIDField: invite.OrganizationID,
		InviteIDField:       invite.ID,
	})
}

func (e *Entry) WithTopic(topic string) *Entry {
	return e.WithField(EventTopicField, topic)
}

func (e *Entry) WithMessage(msg *kafka.Message) *Entry {
	res := map[string]any{
		EventTopicField:    msg.Topic,
		EventIDField:       "unknown",
		EventTypeField:     "unknown",
		EventProducerField: "unknown",
	}

	hs, err := headers.ParseMessageRequiredHeaders(msg.Headers)
	if err == nil {
		res[EventIDField] = hs.EventID
		res[EventVersionField] = hs.EventVersion
		res[EventProducerField] = hs.Producer
	}

	return e.WithFields(res)
}

func (e *Entry) WithOutboxEvent(ev eventbox.OutboxEvent) *Entry {
	return e.WithFields(map[string]any{
		EventIDField:       ev.EventID,
		EventTopicField:    ev.Topic,
		EventTypeField:     ev.Type,
		EventVersionField:  ev.Version,
		EventProducerField: ev.Producer,
		EventAttemptField:  ev.Attempts,
	})
}

func (e *Entry) WithInboxEvent(ev eventbox.InboxEvent) *Entry {
	return e.WithFields(map[string]any{
		EventIDField:       ev.EventID,
		EventTopicField:    ev.Topic,
		EventTypeField:     ev.Type,
		EventVersionField:  ev.Version,
		EventProducerField: ev.Producer,
		EventAttemptField:  ev.Attempts,
	})
}
