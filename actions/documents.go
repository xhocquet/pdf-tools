package actions

import (
  "github.com/gobuffalo/buffalo"
  "github.com/pkg/errors"
  "github.com/thanhpk/randstr"
  "github.com/thedevsaddam/renderer"
  "io"
  "net/http"
  "os"
  "path/filepath"
)

// DocumentsShow default implementation.
func DocumentsShow(c buffalo.Context) error {
  uuid := c.Param("uuid")
  c.Set("uuid", uuid)
  return c.Render(200, r.HTML("documents/show.html"))
}

func DocumentPreview(c buffalo.Context) error {
  rnd := renderer.New()
  uploads_dir := filepath.Join(".", "public/uploads")
  uuid := c.Param("uuid")
  file_path := filepath.Join(uploads_dir, uuid)

  return rnd.FileView(c.Response(), http.StatusOK, file_path, uuid)
}

func DocumentDownload(c buffalo.Context) error {
  rnd := renderer.New()
  uploads_dir := filepath.Join(".", "public/uploads")
  uuid := c.Param("uuid")
  file_path := filepath.Join(uploads_dir, uuid)

  return rnd.FileDownload(c.Response(), http.StatusOK, file_path, uuid)
}

// DocumentsIndex default implementation.
func DocumentsIndex(c buffalo.Context) error {
  return c.Render(200, r.HTML("documents/index.html"))
}

// DocumentsCreate default implementation.
func DocumentsCreate(c buffalo.Context) error {
  random_token := randstr.Hex(16)
  file, err := c.File("uploadedFile")
  if err != nil {
    return errors.WithStack(err)
  }
  dir := filepath.Join(".", "public/uploads")
  if err := os.MkdirAll(dir, 0755); err != nil {
    return errors.WithStack(err)
  }
  f, err := os.Create(filepath.Join(dir, random_token))
  if err != nil {
    return errors.WithStack(err)
  }
  defer f.Close()
  _, err = io.Copy(f, file)
  if err != nil {
    return errors.WithStack(err)
  }
  return c.Redirect(302, "/documents/%s", random_token)
}
