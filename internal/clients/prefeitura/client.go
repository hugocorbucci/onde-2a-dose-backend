package prefeitura

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"strings"

	deps "github.com/hugocorbucci/onde-2a-dose-backend/internal/dependencies"
	prefeituradeps "github.com/hugocorbucci/onde-2a-dose-backend/internal/dependencies/prefeitura"
)

const (
	prefeituraURL = "https://deolhonafila.prefeitura.sp.gov.br/processadores/dados.php"
	bodyKey = "dados"
	bodyValue = "dados"

	// ContentTypeHeader is the header name for content-type
	ContentTypeHeader = "Content-Type"
	// FormContentType is the value of the header for www form encoded content
	FormContentType = "application/x-www-form-urlencoded"
)

type Client struct {
	HTTPClient deps.HTTPClient
}

func (c *Client) Fetch(ctx context.Context) ([]*prefeituradeps.DeOlhoNaFilaUnit, error) {
	req, err := http.NewRequestWithContext(ctx, http.MethodPost, prefeituraURL, strings.NewReader(fmt.Sprintf("%s=%s", bodyKey, bodyValue)))
	req.Header.Add(ContentTypeHeader, FormContentType)
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

	defer resp.Body.Close()
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	results := []*prefeituradeps.DeOlhoNaFilaUnit{}
	err = json.NewDecoder(bytes.NewReader(body)).Decode(&results)
	if err != nil {
		fmt.Println("error decoding", string(body))
		return nil, err
	}

	return results, nil
}