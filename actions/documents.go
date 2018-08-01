package actions

import (
	"io"
	"net/http"
	"os"
	"path/filepath"

	"github.com/gobuffalo/buffalo"
	"github.com/gobuffalo/buffalo/worker"
	"github.com/gobuffalo/pop"
	"github.com/pkg/errors"
	"github.com/thedevsaddam/renderer"
	"github.com/xhocquet/pdf_tool/models"
)

func DocumentsShow(c buffalo.Context) error {
	tx := c.Value("tx").(*pop.Connection)
	uuid := c.Param("uuid")

	// Query docs
	document := models.Document{}
	query := tx.Where("ID = ?", uuid)
	err := query.First(&document)
	if err != nil {
		return errors.WithStack(err)
	}

	c.Set("document", document)

	return c.Render(200, r.HTML("documents/show.html"))
}

func DocumentPreview(c buffalo.Context) error {
	rnd := renderer.New()
	tx := c.Value("tx").(*pop.Connection)
	uuid := c.Param("uuid")

	document := models.Document{}
	query := tx.Where("ID = ?", uuid)
	err := query.First(&document)
	if err != nil {
		return errors.WithStack(err)
	}

	return rnd.FileView(c.Response(), http.StatusOK, document.FilePath(), uuid)
}

func DocumentDownload(c buffalo.Context) error {
	rnd := renderer.New()
	tx := c.Value("tx").(*pop.Connection)
	uuid := c.Param("uuid")

	document := models.Document{}
	query := tx.Where("ID = ?", uuid)
	err := query.First(&document)
	if err != nil {
		return errors.WithStack(err)
	}

	return rnd.FileDownload(c.Response(), http.StatusOK, document.FilePath(), uuid)
}

func DocumentsIndex(c buffalo.Context) error {
	tx := c.Value("tx").(*pop.Connection)

	documents := []models.Document{}
	err := tx.All(&documents)
	if err != nil {
		return errors.WithStack(err)
	}

	c.Set("documents", documents)
	return c.Render(200, r.HTML("documents/index.html"))
}

func DocumentsCreate(c buffalo.Context) error {
	file, err := c.File("uploadedFile")
	if err != nil {
		return errors.WithStack(err)
	}

	tx := c.Value("tx").(*pop.Connection)
	filename := file.FileHeader.Filename

	// Create new record
	document := models.Document{}
	document.Name = filename
	err = tx.Create(&document)

	if err != nil {
		return errors.WithStack(err)
	}

	dir := filepath.Join(".", "public/uploads")
	if err := os.MkdirAll(dir, 0755); err != nil {
		return errors.WithStack(err)
	}

	// Create system file
	f, err := os.Create(document.FilePath())
	if err != nil {
		return errors.WithStack(err)
	}
	defer f.Close()
	// Write data to file
	_, err = io.Copy(f, file)
	if err != nil {
		return errors.WithStack(err)
	}

	// Create job to process PDF
	w.Perform(worker.Job{
		Queue:   "default",
		Handler: "process_pdf",
		Args: worker.Args{
			"document_id": document.ID,
		},
	})

	return c.Redirect(302, "/documents/%s", document.ID)
}
