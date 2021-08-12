package server_test

import (
	"bytes"
	"context"
	"encoding/json"
	"io/ioutil"
	"net"
	"net/http"
	"net/http/httptest"
	"os"
	"strings"
	"testing"

	prefeituraclient "github.com/hugocorbucci/onde-2a-dose-backend/internal/clients/prefeitura"
	deps "github.com/hugocorbucci/onde-2a-dose-backend/internal/dependencies"
	"github.com/hugocorbucci/onde-2a-dose-backend/internal/dependencies/dependenciesfakes"
	"github.com/hugocorbucci/onde-2a-dose-backend/internal/dependencies/prefeitura"
	"github.com/hugocorbucci/onde-2a-dose-backend/internal/server"
	"github.com/stretchr/testify/assert"
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
	withDependencies(t, func(t *testing.T, ctx context.Context, deps *TestDependencies) {
		httpReq, err := http.NewRequestWithContext(ctx, http.MethodGet, deps.BaseURL+"/", nil)
		require.NoError(t, err, "could not create GET / request")

		resp, err := deps.HTTPClient.Do(httpReq)
		require.NoError(t, err, "error making request %+v", httpReq)

		require.Equal(t, http.StatusNotFound, resp.StatusCode, "expected status code to match for req %+v", httpReq)
		_, err = readBodyFrom(resp)
		require.NoError(t, err, "unexpected error reading response body")
	})
}

func TestPostDataRawWithoutBodyReturnsError(t *testing.T) {
	withDependencies(t, func(t *testing.T, ctx context.Context, deps *TestDependencies) {
		httpReq, err := http.NewRequestWithContext(ctx, http.MethodPost, deps.BaseURL+"/data.raw", nil)
		httpReq.Header.Add(prefeituraclient.ContentTypeHeader, prefeituraclient.FormContentType)
		require.NoError(t, err, "could not create POST / request")

		resp, err := deps.HTTPClient.Do(httpReq)
		require.NoError(t, err, "error making request %+v", httpReq)

		require.Equal(t, http.StatusBadRequest, resp.StatusCode, "expected status code to match for req %+v", httpReq)
		body, err := readBodyFrom(resp)
		require.NoError(t, err, "unexpected error reading response body")
		assert.Equal(t, "{\"error\":\"missing body\"}", body, "expected body to match")
	})
}

func TestPostDataRawWithBodyButNoFormEncodingHeaderReturnsError(t *testing.T) {
	withDependencies(t, func(t *testing.T, ctx context.Context, deps *TestDependencies) {
		httpReq, err := http.NewRequestWithContext(ctx, http.MethodPost, deps.BaseURL+"/data.raw", strings.NewReader("dada=a"))
		require.NoError(t, err, "could not create POST / request")

		resp, err := deps.HTTPClient.Do(httpReq)
		require.NoError(t, err, "error making request %+v", httpReq)

		require.Equal(t, http.StatusBadRequest, resp.StatusCode, "expected status code to match for req %+v", httpReq)
		body, err := readBodyFrom(resp)
		require.NoError(t, err, "unexpected error reading response body")
		assert.Equal(t, "{\"error\":\"missing body\"}", body, "expected body to match")
	})
}

func TestPostDataRawWithIncorrectBodyKeyReturnsError(t *testing.T) {
	withDependencies(t, func(t *testing.T, ctx context.Context, deps *TestDependencies) {
		httpReq, err := http.NewRequestWithContext(ctx, http.MethodPost, deps.BaseURL+"/data.raw", strings.NewReader("dada=a"))
		httpReq.Header.Add(prefeituraclient.ContentTypeHeader, prefeituraclient.FormContentType)
		require.NoError(t, err, "could not create POST / request")

		resp, err := deps.HTTPClient.Do(httpReq)
		require.NoError(t, err, "error making request %+v", httpReq)

		require.Equal(t, http.StatusBadRequest, resp.StatusCode, "expected status code to match for req %+v", httpReq)
		body, err := readBodyFrom(resp)
		require.NoError(t, err, "unexpected error reading response body")
		assert.Equal(t, "{\"error\":\"invalid body\"}", body, "expected body to match")
	})
}

func TestPostDataRawWithCorrectEmptyBodyKeyReturnsError(t *testing.T) {
	withDependencies(t, func(t *testing.T, ctx context.Context, deps *TestDependencies) {
		httpReq, err := http.NewRequestWithContext(ctx, http.MethodPost, deps.BaseURL+"/data.raw", strings.NewReader("dados="))
		httpReq.Header.Add(prefeituraclient.ContentTypeHeader, prefeituraclient.FormContentType)
		require.NoError(t, err, "could not create POST / request")

		resp, err := deps.HTTPClient.Do(httpReq)
		require.NoError(t, err, "error making request %+v", httpReq)

		require.Equal(t, http.StatusBadRequest, resp.StatusCode, "expected status code to match for req %+v", httpReq)
		body, err := readBodyFrom(resp)
		require.NoError(t, err, "unexpected error reading response body")
		assert.Equal(t, "{\"error\":\"invalid body\"}", body, "expected body to match")
	})
}

func TestPostDataRawWithCorrectBodyReturnsRawOriginJSON(t *testing.T) {
	withDependencies(t, func(t *testing.T, ctx context.Context, deps *TestDependencies) {
		httpReq, err := http.NewRequestWithContext(ctx, http.MethodPost, deps.BaseURL+"/data.raw", strings.NewReader("dados=dados"))
		httpReq.Header.Add(prefeituraclient.ContentTypeHeader, prefeituraclient.FormContentType)
		require.NoError(t, err, "could not create POST / request")

		if deps.PrefeituraFake != nil {
			deps.PrefeituraFake.FetchReturns([]*prefeitura.DeOlhoNaFilaUnit{
				{
					IDStr:             "1",
					Name:              "Teste",
					Address:           "Rua dos bobos, 0",
					TypeName:          "POSTO VOLANTE",
					TypeIDStr:         "4",
					NeighborhoodName:  "Imaginário",
					NeighborhoodIDStr: "0",
					RegionName:        "CENTRO",
					RegionIDStr:       "1",
					LastUpdatedAtStr:  "2021-08-11 20:00:00.000",
					LineIndexStr:      "1",
					LineStatus:        "SEM FILA",
					CoronaVacStr:      "0",
					AstraZenecaStr:    "1",
					PfizerStr:         "0",
				},
			}, nil)
		}

		resp, err := deps.HTTPClient.Do(httpReq)
		require.NoError(t, err, "error making request %+v", httpReq)

		require.Equal(t, http.StatusOK, resp.StatusCode, "expected status code to match for req %+v", httpReq)
		body := []map[string]interface{}{}
		err = json.NewDecoder(resp.Body).Decode(&body)
		require.NoError(t, err, "unexpected error reading response body")
		assert.Len(t, body, 1, "expected body size to match")
	})
}

// TestDependencies encapsulates the dependencies needed to run a test
type TestDependencies struct {
	BaseURL    string
	HTTPClient HTTPClient

	PrefeituraFake *dependenciesfakes.FakeDeOlhoNaFila
}

func withDependencies(baseT *testing.T, test func(*testing.T, context.Context, *TestDependencies)) {
	ctx := context.Background()
	if len(os.Getenv("TARGET_URL")) == 0 {
		testStates := map[string]func(*testing.T) (*TestDependencies, func()){
			"unitServerTest":        unitDependencies,
			"integrationServerTest": integrationDependencies,
		}
		for name, dep := range testStates {
			baseT.Run(name, func(t *testing.T) {
				deps, stop := dep(t)
				defer stop()
				test(t, ctx, deps)
			})
		}
	} else {
		test(baseT, ctx, smokeDependencies(baseT))
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
	prefeituraClient := &dependenciesfakes.FakeDeOlhoNaFila{}
	s := server.NewHTTPServer(prefeituraClient)
	httpClient := &InMemoryHTTPClient{server: s}
	return &TestDependencies{
		BaseURL:    "",
		HTTPClient: httpClient,
		PrefeituraFake: prefeituraClient,
	}, func() {}
}

func integrationDependencies(t *testing.T) (*TestDependencies, func()) {
	prefeituraClient := &dependenciesfakes.FakeDeOlhoNaFila{}
	baseURL, stop := startTestingHTTPServer(t, prefeituraClient)
	http.DefaultClient.CheckRedirect = func(req *http.Request, via []*http.Request) error {
		return http.ErrUseLastResponse
	}

	return &TestDependencies{
		BaseURL:    baseURL,
		HTTPClient: http.DefaultClient,
		PrefeituraFake: prefeituraClient,
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

func startTestingHTTPServer(t *testing.T, prefeitura deps.DeOlhoNaFila) (string, func()) {
	ctx := context.Background()
	s := server.NewHTTPServer(prefeitura)

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
