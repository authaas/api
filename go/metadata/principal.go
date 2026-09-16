//revive:disable:package-comments
package metadata

import (
	"context"
	"errors"
	"log/slog"

	"connectrpc.com/connect/v2"

	"git.sonicoriginal.software/logger"
)

const (
	// PrincipalIDKey is the request header carrying the authenticated principal
	// ID. The gateway's authenticator sets it from the verified token's "sub"
	// and removes any value a client supplied; services read it to identify
	// the authenticated principal.
	PrincipalIDKey = "x-principal-id"
)

// The ways a request fails to name its principal.
var (
	ErrNoCall            = errors.New("no call in context")
	ErrNoPrincipal       = errors.New("no principal ID in request")
	ErrMultiplePrincipal = errors.New("multiple principal IDs in request")
)

func unauthenticated(ctx context.Context, log *slog.Logger, err error) (string, context.Context, error) {
	return "", logger.ContextWithLogger(ctx, log.With("authenticated", false)), err
}

// ExtractPrincipal reads the authenticated principal ID off the request the
// handler is serving, from the call connect attached to ctx. It answers with
// the ID and a context whose logger carries "requester" and
// "authenticated=true"; on failure the logger carries "authenticated=false"
// and the error says why.
func ExtractPrincipal(
	ctx context.Context,
) (
	principalID string, updatedCtx context.Context, err error,
) {
	log := logger.FromContext(ctx)

	info, ok := connect.CallInfoForServerContext(ctx)
	if !ok {
		return unauthenticated(ctx, log, ErrNoCall)
	}

	values := info.RequestHeader().Values(PrincipalIDKey)

	switch len(values) {
	case 0:
		return unauthenticated(ctx, log, ErrNoPrincipal)
	case 1:
	default:
		return unauthenticated(ctx, log, ErrMultiplePrincipal)
	}

	log = log.With("requester", values[0], "authenticated", true)

	return values[0], logger.ContextWithLogger(ctx, log), nil
}

// InjectPrincipalID sets the principal ID on the outgoing call ctx carries,
// so a service calling another on the principal's behalf forwards who it is
// acting for. ctx is a client context from connect.NewClientContext; one that
// is not becomes one. The returned context is what the call is made with.
func InjectPrincipalID(ctx context.Context, principalID string) context.Context {
	info, ok := connect.CallInfoForClientContext(ctx)
	if !ok {
		ctx, info = connect.NewClientContext(ctx)
	}

	info.RequestHeader().Set(PrincipalIDKey, principalID)

	return ctx
}
