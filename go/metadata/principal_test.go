//revive:disable:package-comments
package metadata

import (
	"context"
	"errors"
	"testing"

	"connectrpc.com/connect/v2"
	"connectrpc.com/connect/v2/connectinprocess"
	"google.golang.org/protobuf/types/known/wrapperspb"
)

// whoSpec names the one method these tests serve: it answers with the
// principal ExtractPrincipal read off the call, or the error it answered with.
var whoSpec = connect.Spec{
	StreamType: connect.StreamTypeUnary,
	Procedure:  "/example.WhoService/Who",
}

// who serves whoSpec in-process and answers with a client on it, so a request
// reaches ExtractPrincipal the way a real one does: through the call connect
// attaches to the handler's context, headers included.
func who() *connect.Client {
	rpc := connect.NewServer()
	rpc.Register(connect.Method{
		Spec: whoSpec,
		Handler: func(ctx context.Context, _ connect.Spec, stream connect.ServerStream) error {
			var request wrapperspb.StringValue
			if err := stream.Receive(&request); err != nil {
				return err
			}

			principalID, _, err := ExtractPrincipal(ctx)
			if err != nil {
				return connect.NewError(connect.CodeUnauthenticated, err.Error())
			}

			return stream.Send(wrapperspb.String(principalID))
		},
	})

	return connect.NewClient(connectinprocess.New(rpc))
}

// ask calls who over ctx and answers with the principal it read.
func ask(ctx context.Context, t *testing.T, client *connect.Client) (string, error) {
	t.Helper()

	var response wrapperspb.StringValue

	err := client.CallUnary(ctx, whoSpec, wrapperspb.String("who am I"), &response)

	return response.GetValue(), err
}

// assertUnauthenticated reports that err is the Unauthenticated answer
// carrying want's text.
func assertUnauthenticated(t *testing.T, err, want error) {
	t.Helper()

	var connectErr *connect.Error
	if !errors.As(err, &connectErr) || connectErr.Code() != connect.CodeUnauthenticated {
		t.Fatalf("error = %v, want unauthenticated", err)
	}
	if connectErr.Message() != want.Error() {
		t.Errorf("message = %q, want %q", connectErr.Message(), want)
	}
}

func TestExtractPrincipal(t *testing.T) {
	t.Run("reads the principal the request carries", func(t *testing.T) {
		got, err := ask(InjectPrincipalID(t.Context(), "principal-123"), t, who())
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if got != "principal-123" {
			t.Errorf("principal = %q, want principal-123", got)
		}
	})

	t.Run("refuses a request naming no principal", func(t *testing.T) {
		_, err := ask(t.Context(), t, who())

		assertUnauthenticated(t, err, ErrNoPrincipal)
	})

	t.Run("refuses a request naming more than one", func(t *testing.T) {
		ctx, info := connect.NewClientContext(t.Context())
		info.RequestHeader().SetValues(PrincipalIDKey, []string{"principal-1", "principal-2"})

		_, err := ask(ctx, t, who())

		assertUnauthenticated(t, err, ErrMultiplePrincipal)
	})

	t.Run("refuses a context that carries no call", func(t *testing.T) {
		_, updated, err := ExtractPrincipal(t.Context())
		if !errors.Is(err, ErrNoCall) {
			t.Fatalf("error = %v, want %v", err, ErrNoCall)
		}
		if updated == nil {
			t.Fatal("expected a context back even on failure")
		}
	})
}

func TestInjectPrincipalID(t *testing.T) {
	t.Run("makes a client context when given none", func(t *testing.T) {
		ctx := InjectPrincipalID(t.Context(), "principal-789")

		info, ok := connect.CallInfoForClientContext(ctx)
		if !ok {
			t.Fatal("expected a client call on the context")
		}
		if got := info.RequestHeader().Values(PrincipalIDKey); len(got) != 1 || got[0] != "principal-789" {
			t.Errorf("header = %v, want the one principal", got)
		}
	})

	t.Run("sets the principal on the client context it is given", func(t *testing.T) {
		ctx, info := connect.NewClientContext(t.Context())
		info.RequestHeader().Set("X-Other", "kept")

		if got := InjectPrincipalID(ctx, "principal-789"); got != ctx {
			t.Fatal("expected the same context back")
		}
		if info.RequestHeader().Get(PrincipalIDKey) != "principal-789" || info.RequestHeader().Get("X-Other") != "kept" {
			t.Errorf("headers = %v, want the principal beside what was there", info.RequestHeader())
		}
	})
}
