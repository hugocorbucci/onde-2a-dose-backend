package prefeitura_test

import (
	"context"
	"errors"
	"io/ioutil"
	"net/http"
	"strings"
	"testing"

	"github.com/hugocorbucci/onde-2a-dose-backend/internal/clients/prefeitura"
	deps "github.com/hugocorbucci/onde-2a-dose-backend/internal/dependencies"
	"github.com/hugocorbucci/onde-2a-dose-backend/internal/dependencies/dependenciesfakes"
	"github.com/stretchr/testify/assert"

	"github.com/stretchr/testify/require"
)

var _ deps.DeOlhoNaFila = &prefeitura.Client{}

func TestClient_FetchErrorsWhenDownstreamErrors(t *testing.T) {
	fakeClient := &dependenciesfakes.FakeHTTPClient{}
	client := &prefeitura.Client{
		HTTPClient: fakeClient,
	}

	fakeClient.DoReturns(nil, errors.New("error"))
	_, err := client.Fetch(context.Background())
	require.Error(t, err, "expected error to match")
}

func TestClient_FetchErrorsWhenDownstreamReturnsNonOKStatusCode(t *testing.T) {
	fakeClient := &dependenciesfakes.FakeHTTPClient{}
	client := &prefeitura.Client{
		HTTPClient: fakeClient,
	}

	fakeClient.DoReturns(&http.Response{
		Status:           "I'm a teapot",
		StatusCode:       http.StatusTeapot,
	}, nil)
	_, err := client.Fetch(context.Background())
	require.Error(t, err, "expected error to match")
}

func TestClient_FetchErrorsWhenDownstreamReturnsEmptyBody(t *testing.T) {
	fakeClient := &dependenciesfakes.FakeHTTPClient{}
	client := &prefeitura.Client{
		HTTPClient: fakeClient,
	}

	fakeClient.DoReturns(&http.Response{
		Status:           "OK",
		StatusCode:       http.StatusOK,
	}, nil)
	_, err := client.Fetch(context.Background())
	require.Error(t, err, "expected error to match")
}

func TestClient_FetchErrorsWhenDownstreamReturnsNonParseableResponse(t *testing.T) {
	fakeClient := &dependenciesfakes.FakeHTTPClient{}
	client := &prefeitura.Client{
		HTTPClient: fakeClient,
	}

	fakeClient.DoReturns(&http.Response{
		Status:           "OK",
		StatusCode:       http.StatusOK,
		Body: ioutil.NopCloser(strings.NewReader("{}")),
	}, nil)
	_, err := client.Fetch(context.Background())
	require.Error(t, err, "expected error to match")
}

func TestClient_FetchWorksWhenDownstreamReturnsValidEmptyResponse(t *testing.T) {
	fakeClient := &dependenciesfakes.FakeHTTPClient{}
	client := &prefeitura.Client{
		HTTPClient: fakeClient,
	}

	fakeClient.DoReturns(&http.Response{
		Status:           "OK",
		StatusCode:       http.StatusOK,
		Body: ioutil.NopCloser(strings.NewReader("[]")),
	}, nil)
	res, err := client.Fetch(context.Background())
	require.NoError(t, err, "expected error to match")
	assert.Len(t, res, 0, "expected length to match")
}

func TestClient_FetchWorksWhenDownstreamReturnsValidResponse(t *testing.T) {
	fakeClient := &dependenciesfakes.FakeHTTPClient{}
	client := &prefeitura.Client{
		HTTPClient: fakeClient,
	}

	fakeClient.DoReturns(&http.Response{
		Status:           "OK",
		StatusCode:       http.StatusOK,
		Body: ioutil.NopCloser(strings.NewReader(`[
{"equipamento":"GRCS ESCOLA DE SAMBA VAI-VAI","endereco":"Rua S\u00e3o Vicente, n\u00ba 276 - Bela Vista","tipo_posto":"POSTO VOLANTE","id_tipo_posto":"4","id_distrito":"1","distrito":"Bela Vista","id_crs":"1","crs":"CENTRO","data_hora":"2021-08-11 07:50:49.173","indice_fila":"5","status_fila":"N\u00c3O FUNCIONANDO","coronavac":"1","astrazeneca":"0","pfizer":"false","id_tb_unidades":"1571"}
]`)),
	}, nil)
	res, err := client.Fetch(context.Background())
	require.NoError(t, err, "expected error to match")
	if assert.Len(t, res, 1, "expected length to match") {
		unit := res[0]
		assert.Equal(t, "GRCS ESCOLA DE SAMBA VAI-VAI", unit.Name, "expected name to match")
		assert.Equal(t, "Rua S\u00e3o Vicente, n\u00ba 276 - Bela Vista", unit.Address, "expected address to match")
		assert.Equal(t, false, unit.HasAstraZeneca(), "expected astrazeneca to match")
		assert.Equal(t, false, unit.HasPfizer(), "expected pfizer to match")
		assert.Equal(t, true, unit.HasCoronaVac(), "expected coronavac to match")
		// TODO: Assert on time
	}
}