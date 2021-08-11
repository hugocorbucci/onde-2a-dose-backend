package dependencies

import (
	"context"

	"github.com/hugocorbucci/onde-2a-dose-backend/internal/dependencies/prefeitura"
)

//counterfeiter:generate . DeOlhoNaFila
type DeOlhoNaFila interface {
	Fetch(ctx context.Context) ([]*prefeitura.DeOlhoNaFilaUnit, error)
}
