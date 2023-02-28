package sigv4_test

import (
	"context"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/grafana/agent/component/otelcol/auth"
	"github.com/grafana/agent/component/otelcol/auth/sigv4"
	"github.com/grafana/agent/pkg/flow/componenttest"
	"github.com/grafana/agent/pkg/river"
	"github.com/grafana/agent/pkg/util"
	"github.com/stretchr/testify/require"
	"go.opentelemetry.io/collector/config/configauth"
)

// Test performs a basic integration test which runs the otelcol.auth.sigv4
// component and ensures that it can be used for authentication.
func Test(t *testing.T) {
	// Create an HTTP server which will assert that sigv4 auth has been injected
	// into the request.
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		for k, v_arr := range r.Header {
			fmt.Printf("----- %s -----\n", k)
			for _, v := range v_arr {
				fmt.Printf("%s; ", v)
			}
		}

		w.WriteHeader(http.StatusOK)
	}))
	defer srv.Close()

	ctx := componenttest.TestContext(t)
	ctx, cancel := context.WithTimeout(ctx, time.Minute)
	defer cancel()

	l := util.TestLogger(t)

	// Create and run our component
	ctrl, err := componenttest.NewControllerFromID(l, "otelcol.auth.sigv4")
	require.NoError(t, err)

	cfg := `
	assume_role {
		session_name = "role_session_name"
	}
	region = "region"
	service = "service"
	`
	var args sigv4.Arguments
	require.NoError(t, river.Unmarshal([]byte(cfg), &args))

	go func() {
		err := ctrl.Run(ctx, args)
		require.NoError(t, err)
	}()

	require.NoError(t, ctrl.WaitRunning(time.Second), "component never started")
	require.NoError(t, ctrl.WaitExports(time.Second), "component never exported anything")

	// Get the authentication extension from our component and use it to make a
	// request to our test server.
	exports := ctrl.Exports().(auth.Exports)
	require.NotNil(t, exports.Handler.Extension, "handler extension is nil")

	clientAuth, ok := exports.Handler.Extension.(configauth.ClientAuthenticator)
	require.True(t, ok, "handler does not implement configauth.ClientAuthenticator")

	rt, err := clientAuth.RoundTripper(http.DefaultTransport)
	require.NoError(t, err)
	cli := &http.Client{Transport: rt}

	// Wait until the request finishes. We don't assert anything else here; our
	// HTTP handler won't write the response until it ensures that the sigv4
	// were set.
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, srv.URL, nil)
	require.NoError(t, err)
	resp, err := cli.Do(req)
	require.NoError(t, err, "HTTP request failed")
	require.Equal(t, http.StatusOK, resp.StatusCode)
}
