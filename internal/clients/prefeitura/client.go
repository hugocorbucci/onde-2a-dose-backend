package prefeitura

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"strings"

	deps "github.com/hugocorbucci/onde-2a-dose-backend/internal/dependencies"
	prefeituradeps "github.com/hugocorbucci/onde-2a-dose-backend/internal/dependencies/prefeitura"
)

const (
	prefeituraURL = "https://deolhonafila.prefeitura.sp.gov.br/processadores/dados.php"
	bodyKey = "dados"
	bodyValue = "dados"
)

type Client struct {
	HTTPClient deps.HTTPClient
}

func (c *Client) Fetch(ctx context.Context) ([]*prefeituradeps.DeOlhoNaFilaUnit, error) {
	req, err := http.NewRequestWithContext(ctx, http.MethodPost, prefeituraURL, strings.NewReader(fmt.Sprintf("%s=%s", bodyKey, bodyValue)))
	if err != nil {
		return nil, err
	}

	resp, err := c.HTTPClient.Do(req)
	if err != nil {
		return nil, err
	}
	if resp.StatusCode != http.StatusOK {
		return nil, errors.New("invalid response status code")
	}
	if resp.Body == nil {
		return nil, errors.New("empty body")
	}

	results := []*prefeituradeps.DeOlhoNaFilaUnit{}
	err = json.NewDecoder(resp.Body).Decode(&results)
	if err != nil {
		return nil, err
	}

	return results, nil
}