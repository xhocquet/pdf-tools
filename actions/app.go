package actions

import (
	"github.com/gobuffalo/buffalo"
	"github.com/gobuffalo/buffalo/middleware"
	"github.com/gobuffalo/buffalo/middleware/ssl"
	"github.com/gobuffalo/buffalo/worker"
	"github.com/gobuffalo/envy"
	"github.com/unrolled/secure"

	"github.com/gobuffalo/buffalo/middleware/csrf"
	"github.com/gobuffalo/buffalo/middleware/i18n"
	"github.com/gobuffalo/packr"
	"github.com/xhocquet/pdf_tool/models"
	"github.com/xhocquet/pdf_tool/workers"

	"github.com/gobuffalo/gocraft-work-adapter"
	"github.com/gomodule/redigo/redis"
)

// ENV is used to help switch settings based on where the
// application is being run. Default is "development".
var ENV = envy.Get("GO_ENV", "development")
var app *buffalo.App
var T *i18n.Translator
var w worker.Worker

// App is where all routes and middleware for buffalo
// should be defined. This is the nerve center of your
// application.
func App() *buffalo.App {
	if app == nil {
		app = buffalo.New(buffalo.Options{
			Env:         ENV,
			SessionName: "_pdf_tool_session",
			Worker: gwa.New(gwa.Options{
				Pool: &redis.Pool{
					MaxActive: 5,
					MaxIdle:   5,
					Wait:      true,
					Dial: func() (redis.Conn, error) {
						return redis.Dial("tcp", ":6379")
					},
				},
				Name:           "pdftool",
				MaxConcurrency: 25,
			}),
		})

		app.Use(forceSSL())
		app.Use(csrf.New)

		if ENV == "development" {
			app.Use(middleware.ParameterLogger)
		}

		app.Use(middleware.PopTransaction(models.DB))

		app.Use(translations())

		w = app.Worker // Get a ref to the previously defined Worker
		w.Register("process_pdf", func(args worker.Args) error {
			workers.ProcessPDF(args)
			return nil
		})

		app.GET("/", HomeHandler)
		app.GET("/documents/{uuid}", DocumentsShow)
		app.GET("/documents/{uuid}/download", DocumentDownload)
		app.GET("/documents/{uuid}/preview", DocumentPreview)
		app.GET("/documents", DocumentsIndex)
		app.POST("/documents/create", DocumentsCreate)

		app.ServeFiles("/", assetsBox) // serve files from the public directory
	}

	return app
}

// translations will load locale files, set up the translator `actions.T`,
// and will return a middleware to use to load the correct locale for each
// request.
// for more information: https://gobuffalo.io/en/docs/localization
func translations() buffalo.MiddlewareFunc {
	var err error
	if T, err = i18n.New(packr.NewBox("../locales"), "en-US"); err != nil {
		app.Stop(err)
	}
	return T.Middleware()
}

// forceSSL will return a middleware that will redirect an incoming request
// if it is not HTTPS. "http://example.com" => "https://example.com".
// This middleware does **not** enable SSL. for your application. To do that
// we recommend using a proxy: https://gobuffalo.io/en/docs/proxy
// for more information: https://github.com/unrolled/secure/
func forceSSL() buffalo.MiddlewareFunc {
	return ssl.ForceSSL(secure.Options{
		SSLRedirect:     ENV == "production",
		SSLProxyHeaders: map[string]string{"X-Forwarded-Proto": "https"},
	})
}
