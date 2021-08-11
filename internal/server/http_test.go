package server_test

import (
	"bytes"
	"context"
	"io/ioutil"
	"net"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"

	"github.com/hugocorbucci/onde-2a-dose-backend/internal/server"

	_ "github.com/go-sql-driver/mysql"
	"github.com/stretchr/testify/require"
)

type HTTPClient interface {
	Do(r *http.Request) (*http.Response, error)
}
type InMemoryHTTPClient struct {
	server *server.Server
}

func (c *InMemoryHTTPClient) Do(r *http.Request) (*http.Response, error) {
	w := httptest.NewRecorder()
	c.server.ServeHTTP(w, r)

	return &http.Response{
		Header:     w.Header(),
		StatusCode: w.Code,
		Body:       ioutil.NopCloser(bytes.NewReader(w.Body.Bytes())),
	}, nil
}

func TestHomeReturns404(t *testing.T) {
	withDependencies(t, func(t *testing.T, deps *TestDependencies) {
		httpReq, err := http.NewRequest(http.MethodGet, deps.BaseURL+"/", nil)
		require.NoError(t, err, "could not create GET / request")

		resp, err := deps.HTTPClient.Do(httpReq)
		require.NoError(t, err, "error making request %+v", httpReq)

		require.Equal(t, http.StatusNotFound, resp.StatusCode, "expected status code to match for req %+v", httpReq)
		_, err = readBodyFrom(resp)
		require.NoError(t, err, "unexpected error reading response body")
	})
}

// TestDependencies encapsulates the dependencies needed to run a test
type TestDependencies struct {
	BaseURL    string
	HTTPClient HTTPClient
}

func withDependencies(baseT *testing.T, test func(*testing.T, *TestDependencies)) {
	if len(os.Getenv("TARGET_URL")) == 0 {
		testStates := map[string]func(*testing.T) (*TestDependencies, func()){
			"unitServerTest":        unitDependencies,
			"integrationServerTest": integrationDependencies,
		}
		for name, dep := range testStates {
			baseT.Run(name, func(t *testing.T) {
				deps, stop := dep(t)
				defer stop()
				test(t, deps)
			})
		}
	} else {
		test(baseT, smokeDependencies(baseT))
	}
}

type testStructure struct {
	test     func()
	tearDown func()
}

func do(test func()) *testStructure {
	return &testStructure{
		test: test,
	}
}

func withTearDown(tearDown func()) *testStructure {
	return &testStructure{
		tearDown: tearDown,
	}
}

func (s *testStructure) do(test func()) *testStructure {
	copy := &testStructure{}
	*copy = *s
	copy.test = test
	return copy
}

func (s *testStructure) withTearDown(tearDown func()) *testStructure {
	copy := &testStructure{}
	*copy = *s
	copy.tearDown = tearDown
	return copy
}

func (s *testStructure) Now() {
	if s.tearDown != nil {
		defer s.tearDown()
	}

	if s.test != nil {
		s.test()
	}
}

func unitDependencies(*testing.T) (*TestDependencies, func()) {
	s := server.NewHTTPServer()
	httpClient := &InMemoryHTTPClient{server: s}
	return &TestDependencies{
		BaseURL:    "",
		HTTPClient: httpClient,
	}, func() {}
}

func integrationDependencies(t *testing.T) (*TestDependencies, func()) {
	baseURL, stop := startTestingHTTPServer(t)
	http.DefaultClient.CheckRedirect = func(req *http.Request, via []*http.Request) error {
		return http.ErrUseLastResponse
	}

	return &TestDependencies{
		BaseURL:    baseURL,
		HTTPClient: http.DefaultClient,
	}, stop
}

func smokeDependencies(_ *testing.T) *TestDependencies {
	http.DefaultClient.CheckRedirect = func(req *http.Request, via []*http.Request) error {
		return http.ErrUseLastResponse
	}
	return &TestDependencies{
		BaseURL:    os.Getenv("TARGET_URL"),
		HTTPClient: http.DefaultClient,
	}
}

func startTestingHTTPServer(t *testing.T) (string, func()) {
	ctx := context.Background()
	s := server.NewHTTPServer()

	listener, err := net.Listen("tcp", "localhost:0")
	if err != nil {
		t.Fatalf("could not listen for HTTP requests: %+v", err)
	}
	baseURL := "http://" + listener.Addr().String()
	srvr := http.Server{Addr: baseURL, Handler: s}

	go srvr.Serve(listener)
	return baseURL, func() {
		if err := srvr.Shutdown(ctx); err != nil {
			t.Logf("could not shutdown http server: %+v", err)
		}
	}
}

func readBodyFrom(resp *http.Response) (string, error) {
	bodyBytes, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}
	return string(bodyBytes), nil
}
