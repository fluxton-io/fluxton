package models

type Table struct {
	BaseModel
	Id            int    `db:"id"`
	Name          string `db:"name"`
	Schema        string `db:"schema"`
	EstimatedRows int    `db:"estimated_rows"`
	TotalSize     string `db:"total_size"`
}
