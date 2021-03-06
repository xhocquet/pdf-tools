package workers

import (
  "encoding/json"
  "os"

  "github.com/gobuffalo/buffalo/worker"
  "github.com/gobuffalo/pop"
  "github.com/xhocquet/pdf_tool/models"

  pdf "github.com/unidoc/unidoc/pdf/model"
)

func ProcessPDF(args worker.Args) {
  uuid := args["document_id"].(string)
  document := models.Document{}

  err := models.DB.Transaction(func(tx *pop.Connection) error {
    query := tx.Where("ID = ?", uuid)
    err := query.First(&document)

    if err != nil {
      return err
    }

    f, err := os.Open(document.FilePath())
    if err != nil {
      return err
    }

    defer f.Close()

    pdfReader, err := pdf.NewPdfReader(f)
    if err != nil {
      return err
    }

    isEncrypted, err := pdfReader.IsEncrypted()
    if err != nil {
      return err
    }

    numPages, err := pdfReader.GetNumPages()
    if err != nil {
      return err
    }

    objTypes, err := pdfReader.Inspect()
    if err != nil {
      return err
    }

    marshalled_metadata, err := json.Marshal(objTypes)

    document.IsEncrypted = isEncrypted
    document.PageCount = numPages
    document.Status = 1
    document.Metadata = string(marshalled_metadata)
    err = tx.Update(&document)

    return err
  })

  if err != nil {
    return
  }

  return
}
