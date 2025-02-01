package models

import (
	"encoding/json"
	"fluxton/types"
	"fmt"
	"strings"
	"time"
)

type JSONColumns []types.TableColumn // important for reading from db

type Table struct {
	ID        uint        `db:"id"`
	ProjectID uint        `db:"project_id"`
	Name      string      `db:"name"`
	Columns   JSONColumns `db:"columns" json:"columns"`
	CreatedAt time.Time   `db:"created_at"`
	UpdatedAt time.Time   `db:"updated_at"`
}

func (t Table) GetTableName() string {
	return "tables"
}

func (t Table) GetColumns() string {
	return "id, project_id, name, columns, created_at, updated_at"
}

func (t Table) GetColumnsWithAlias(alias string) string {
	columns := strings.Split(t.GetColumns(), ", ")
	for i, column := range columns {
		columns[i] = alias + "." + column
	}

	return strings.Join(columns, ", ")
}

func (t Table) MarshalJSONColumns() ([]byte, error) {
	return json.Marshal(t.Columns)
}

func (t Table) UnmarshalJSONColumns(data []byte) error {
	return json.Unmarshal(data, &t.Columns)
}

func (j *JSONColumns) Scan(value interface{}) error {
	byteData, ok := value.([]byte)
	if !ok {
		return fmt.Errorf("cannot convert database value to []byte")
	}
	return json.Unmarshal(byteData, j)
}
