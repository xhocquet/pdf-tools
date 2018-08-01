package models

import (
	"encoding/json"
	"path/filepath"
	"time"

	"github.com/gobuffalo/pop"
	"github.com/gobuffalo/uuid"
	"github.com/gobuffalo/validate"
)

type Document struct {
	ID          uuid.UUID `json:"id" db:"id"`
	Name        string    `json:"name" db:"name"`
	CreatedAt   time.Time `json:"created_at" db:"created_at"`
	UpdatedAt   time.Time `json:"updated_at" db:"updated_at"`
	PageCount   int       `json:"page_count" db:"page_count"`
	Status      int       `json:"status" db:"status"`
	IsEncrypted bool      `json:"is_encrypted" db:"is_encrypted"`
	Metadata    string    `json:"metadata" db:"metadata"`
}

// String is not required by pop and may be deleted
func (d Document) String() string {
	jd, _ := json.Marshal(d)
	return string(jd)
}

// Documents is not required by pop and may be deleted
type Documents []Document

// String is not required by pop and may be deleted
func (d Documents) String() string {
	jd, _ := json.Marshal(d)
	return string(jd)
}

// Validate gets run every time you call a "pop.Validate*" (pop.ValidateAndSave, pop.ValidateAndCreate, pop.ValidateAndUpdate) method.
// This method is not required and may be deleted.
func (d *Document) Validate(tx *pop.Connection) (*validate.Errors, error) {
	return validate.NewErrors(), nil
}

// ValidateCreate gets run every time you call "pop.ValidateAndCreate" method.
// This method is not required and may be deleted.
func (d *Document) ValidateCreate(tx *pop.Connection) (*validate.Errors, error) {
	return validate.NewErrors(), nil
}

// ValidateUpdate gets run every time you call "pop.ValidateAndUpdate" method.
// This method is not required and may be deleted.
func (d *Document) ValidateUpdate(tx *pop.Connection) (*validate.Errors, error) {
	return validate.NewErrors(), nil
}

func (d *Document) FilePath() string {
	uuid := d.ID.String()
	return filepath.Join(UploadsPath(), uuid)
}

func (d *Document) PagesPath() string {
	uuid := d.ID.String()
	pages_folder := uuid + "-pages"
	return filepath.Join(UploadsPath(), pages_folder)
}

// private

func UploadsPath() string {
	return filepath.Join(".", "public/uploads")
}
