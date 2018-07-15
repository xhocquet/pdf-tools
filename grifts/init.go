package grifts

import (
	"github.com/gobuffalo/buffalo"
	"github.com/xhocquet/pdf_tool/actions"
)

func init() {
	buffalo.Grifts(actions.App())
}
