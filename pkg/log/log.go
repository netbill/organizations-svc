package log

import (
	"context"
	"errors"
	"log/slog"
	"net/http"
	"os"
	"strings"

	"github.com/google/uuid"
	"github.com/netbill/ape"
	"github.com/netbill/eventbox/headers"
	"github.com/netbill/logium"
	"github.com/netbill/organizations-svc/internal/core/models"
	"github.com/segmentio/kafka-go"
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
)

type Logger struct {
	base *slog.Logger
}

func New(level string, format string, serviceName string) *Logger {
	lvl := parseLevel(level)

	var handler slog.Handler

	switch strings.ToLower(strings.TrimSpace(format)) {
	case "json":
		handler = slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{
			Level: lvl,
		})

	default:
		handler = logium.NewAlignedTextHandler(os.Stdout, logium.AlignedTextOptions{
			Level:      lvl,
			TimeFormat: "2006-01-02 15:04:05",
			MsgWidth:   55,
			Colors:     true,
		})
	}

	base := slog.New(handler).
		With(slog.String(ServiceField, serviceName))

	return &Logger{base: base}
}

func (l *Logger) With(args ...any) logium.Logger {
	return &Logger{base: l.base.With(args...)}
}

func (l *Logger) WithFields(fields map[string]any) logium.Logger {
	if len(fields) == 0 {
		return l
	}
	args := make([]any, 0, len(fields))
	for k, v := range fields {
		args = append(args, slog.Any(k, v))
	}
	return &Logger{base: l.base.With(args...)}
}

func (l *Logger) WithField(key string, value any) logium.Logger {
	return &Logger{base: l.base.With(slog.Any(key, value))}
}

func (l *Logger) WithError(err error) logium.Logger {
	if err == nil {
		return l
	}

	var ae *ape.Error
	if errors.As(err, &ae) {
		return &Logger{base: l.base.With(
			slog.String("error", ae.Error()),
			slog.Any("ape", ae),
		)}
	}

	return &Logger{base: l.base.With(slog.String("error", err.Error()))}
}

func (l *Logger) WithRequest(r *http.Request) logium.Logger {
	if r == nil {
		return l
	}
	return &Logger{base: l.base.With(
		slog.String(HTTPMethodField, r.Method),
		slog.String(HTTPPathField, r.URL.Path),
	)}
}

func (l *Logger) WithOperation(operation string) logium.Logger {
	return &Logger{base: l.base.With(slog.String(OperationField, operation))}
}

func (l *Logger) WithComponent(component string) logium.Logger {
	return &Logger{base: l.base.With(slog.String(ComponentField, component))}
}

func (l *Logger) Debug(msg string, args ...any) {
	l.base.Debug(msg, args...)
}

func (l *Logger) Info(msg string, args ...any) {
	l.base.Info(msg, args...)
}

func (l *Logger) Warn(msg string, args ...any) {
	l.base.Warn(msg, args...)
}

func (l *Logger) Error(msg string, args ...any) {
	l.base.Error(msg, args...)
}

func (l *Logger) DebugContext(ctx context.Context, msg string, args ...any) {
	l.base.DebugContext(ctx, msg, args...)
}
func (l *Logger) InfoContext(ctx context.Context, msg string, args ...any) {
	l.base.InfoContext(ctx, msg, args...)
}
func (l *Logger) WarnContext(ctx context.Context, msg string, args ...any) {
	l.base.WarnContext(ctx, msg, args...)
}
func (l *Logger) ErrorContext(ctx context.Context, msg string, args ...any) {
	l.base.ErrorContext(ctx, msg, args...)
}

// --- your models sugar stays returning *Logger if you want ---

func (l *Logger) WithAccountAuthClaims(auth interface {
	GetAccountID() uuid.UUID
	GetSessionID() uuid.UUID
}) *Logger {
	if auth == nil {
		return l
	}
	return &Logger{base: l.base.With(
		slog.String(AccountIDField, auth.GetAccountID().String()),
		slog.String(AccountSessionIDField, auth.GetSessionID().String()),
	)}
}

func (l *Logger) WithOrganization(organization models.Organization) *Logger {
	return &Logger{base: l.base.With(slog.String(OrganizationIDField, organization.ID.String()))}
}

func (l *Logger) WithMember(member models.Member) *Logger {
	return &Logger{base: l.base.With(
		slog.String(OrganizationIDField, member.OrganizationID.String()),
		slog.String(MemberIDField, member.ID.String()),
	)}
}

func (l *Logger) WithInvite(invite models.Invite) *Logger {
	return &Logger{base: l.base.With(
		slog.String(OrganizationIDField, invite.OrganizationID.String()),
		slog.String(InviteIDField, invite.ID.String()),
	)}
}

func (l *Logger) WithTopic(topic string) *Logger {
	return &Logger{base: l.base.With(slog.String(EventTopicField, topic))}
}

func (l *Logger) WithMessage(msg *kafka.Message) *Logger {
	if msg == nil {
		return l
	}

	args := []any{
		slog.String(EventTopicField, msg.Topic),
		slog.String(EventIDField, "unknown"),
		slog.String(EventTypeField, "unknown"),
		slog.String(EventProducerField, "unknown"),
	}

	hs, err := headers.ParseMessageRequiredHeaders(msg.Headers)
	if err == nil {
		args = append(args,
			slog.String(EventIDField, hs.EventID.String()),
			slog.Int(EventVersionField, int(hs.EventVersion)),
			slog.String(EventProducerField, hs.Producer),
		)
	}

	return &Logger{base: l.base.With(args...)}
}

func parseLevel(level string) slog.Level {
	switch strings.ToLower(strings.TrimSpace(level)) {
	case "debug":
		return slog.LevelDebug
	case "warn", "warning":
		return slog.LevelWarn
	case "error":
		return slog.LevelError
	default:
		return slog.LevelInfo
	}
}
