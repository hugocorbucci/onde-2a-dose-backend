package prefeitura

import (
	"fmt"
	"strconv"
	"time"
)

const (
	// DateLayout is the datetime layout used for the data_hora json payload
	DateLayout = "2006-01-02 15:04:05.999Z07:00"
)
// {"equipamento":"GRCS ESCOLA DE SAMBA VAI-VAI","endereco":"Rua S\u00e3o Vicente, n\u00ba 276 - Bela Vista","tipo_posto":"POSTO VOLANTE","id_tipo_posto":"4","id_distrito":"1","distrito":"Bela Vista","id_crs":"1","crs":"CENTRO","data_hora":"2021-08-11 07:50:49.173","indice_fila":"5","status_fila":"N\u00c3O FUNCIONANDO","coronavac":"0","astrazeneca":"0","pfizer":"0","id_tb_unidades":"1571"}

// DeOlhoNaFilaUnit represents the payload for a single entity provided by São Paulo's city hall for
// COVID-19 vacination via https://deolhonafila.prefeitura.sp.gov.br/
type DeOlhoNaFilaUnit struct {
	IDStr string `json:"id_tb_unidades"`
	Name string `json:"equipamento"`
	Address string `json:"endereco"`

	TypeName string `json:"tipo_posto"`
	TypeIDStr string `json:"id_tipo_posto"`
	NeighborhoodName string `json:"distrito"`
	NeighborhoodIDStr string `json:"id_distrito"`
	RegionName string `json:"crs"`
	RegionIDStr string `json:"id_crs"`
	LastUpdatedAtStr string `json:"data_hora"`
	LineIndexStr string `json:"indice_fila"`
	LineStatus string `json:"status_fila"`

	CoronaVacStr string `json:"coronavac"`
	AstraZenecaStr string `json:"astrazeneca"`
	PfizerStr string `json:"pfizer"`
}

// ID returns the ID of the unit as an int or 0 if not parseable
func (u *DeOlhoNaFilaUnit) ID() int {
	return parseInt(u.IDStr)
}

// TypeID returns the TypeID of the unit as an int or 0 if not parseable
func (u *DeOlhoNaFilaUnit) TypeID() int {
	return parseInt(u.TypeIDStr)
}

// NeighborhoodID returns the NeighborhoodID of the unit as an int or 0 if not parseable
func (u *DeOlhoNaFilaUnit) NeighborhoodID() int {
	return parseInt(u.NeighborhoodIDStr)
}

// RegionID returns the RegionID of the unit as an int or 0 if not parseable
func (u *DeOlhoNaFilaUnit) RegionID() int {
	return parseInt(u.RegionIDStr)
}

// LineIndex returns the LineIndex of the unit as an int or 0 if not parseable
func (u *DeOlhoNaFilaUnit) LineIndex() int {
	return parseInt(u.LineIndexStr)
}

// HasCoronaVac returns whether CoronaVac is available at the unit
func (u *DeOlhoNaFilaUnit) HasCoronaVac() bool {
	return parseBool(u.CoronaVacStr)
}

// HasAstraZeneca returns whether AstraZeneca is available at the unit
func (u *DeOlhoNaFilaUnit) HasAstraZeneca() bool {
	return parseBool(u.AstraZenecaStr)
}

// HasPfizer returns whether Pfizer is available at the unit
func (u *DeOlhoNaFilaUnit) HasPfizer() bool {
	return parseBool(u.PfizerStr)
}

// LastUpdatedAt returns the last time information on this unit has been updated at
func (u *DeOlhoNaFilaUnit) LastUpdatedAt() time.Time {
	t, err := time.Parse(DateLayout, fmt.Sprintf("%s-03:00", u.LastUpdatedAtStr))
	if err != nil {
		return time.Unix(0, 0)
	}
	return t
}

func parseBool(v string) bool {
	id, err := strconv.Atoi(v)
	if err != nil {
		return false
	}
	return id == 1
}

func parseInt(v string) int {
	id, err := strconv.Atoi(v)
	if err != nil {
		return 0
	}
	return id
}