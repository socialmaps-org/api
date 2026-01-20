package method

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"codeberg.org/socialmaps/api/internal/model"
	"codeberg.org/socialmaps/api/internal/web"
	"github.com/stretchr/testify/require"
)

type TestAuthenticator struct {
	web.Authenticator
	t                *testing.T
	nIntrospectCalls int
}

func NewTestAuthenticator(t *testing.T, osmSubs ...string) *TestAuthenticator {
	clientID := "my-client-id"
	clientSecret := "my-client-secret"
	token := "my-auth-token"

	authSrv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		username, password, ok := r.BasicAuth()
		require.True(t, ok)

		require.Equal(t, clientID, username)
		require.Equal(t, clientSecret, password)
		require.Equal(t, token, r.FormValue("token"))

		w.Header().Set("Content-Type", "application/json")

		if len(osmSubs) != 0 {
			osmSub := osmSubs[0]
			osmSubs = osmSubs[1:]
			w.Write([]byte(fmt.Sprintf(`{"active": true, "openstreetmap_sub": "%s"}`, osmSub)))
		} else {
			w.Write([]byte(`{"active": false}`))
		}
	}))
	t.Cleanup(authSrv.Close)

	return &TestAuthenticator{
		Authenticator: web.NewAuthenticator(authSrv.URL, clientID, clientSecret),
		t:             t,
	}
}

func (a *TestAuthenticator) Introspect(token string) (*web.Authentication, error) {
	a.nIntrospectCalls++
	return a.Authenticator.Introspect(token)
}

func NewTestServer(t *testing.T, authr web.Authenticator, qs *model.Queries, overpassEndpoint string) *httptest.Server {
	handler := Mux(authr, qs, overpassEndpoint)
	srv := httptest.NewServer(handler)
	t.Cleanup(srv.Close)

	return srv
}
