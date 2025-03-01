package models

import (
	"github.com/google/uuid"
	"time"
)

type File struct {
	Uuid         uuid.UUID `db:"uuid" json:"uuid"`
	BucketUuid   uuid.UUID `db:"bucket_uuid" json:"bucketUuid"`
	FullFileName string    `db:"full_file_name" json:"fullFileName"`
	Size         int       `db:"size" json:"size"` // in KB
	MimeType     string    `db:"mime_type" json:"mimeType"`
	CreatedBy    uuid.UUID `db:"created_by" json:"createdBy"`
	UpdatedBy    uuid.UUID `db:"updated_by" json:"updatedBy"`
	CreatedAt    time.Time `db:"created_at" json:"createdAt"`
	UpdatedAt    time.Time `db:"updated_at" json:"updatedAt"`
}

func (u File) GetTableName() string {
	return "storage.files"
}
