package resources

import (
	"fluxton/models"
)

type ColumnResponse struct {
	Name     string `json:"name"`
	Position int    `json:"position"`
	NotNull  bool   `json:"notNull"`
	Type     string `json:"type"`
	Default  string `json:"defaultValue"`
	Primary  bool   `json:"primary"`
	Unique   bool   `json:"unique"`
}

func ColumnResource(column *models.Column) ColumnResponse {
	return ColumnResponse{
		Name:     column.Name,
		Position: column.Position,
		NotNull:  column.NotNull,
		Type:     column.Type,
		Default:  column.Default,
		Primary:  column.Primary,
		Unique:   column.Unique,
	}
}

func ColumnResourceCollection(columns []models.Column) []ColumnResponse {
	resourceColumns := make([]ColumnResponse, len(columns))
	for i, column := range columns {
		resourceColumns[i] = ColumnResource(&column)
	}

	return resourceColumns
}
