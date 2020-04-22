/*
 * MIT License
 *
 * Copyright (c) 2020 Dan <https://github.com/delivrance>
 *
 * Permission is hereby granted, free of charge, to any person obtaining a copy
 * of this software and associated documentation files (the "Software"), to deal
 * in the Software without restriction, including without limitation the rights
 * to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
 * copies of the Software, and to permit persons to whom the Software is
 * furnished to do so, subject to the following conditions:
 *
 * The above copyright notice and this permission notice shall be included in all
 * copies or substantial portions of the Software.
 *
 * THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
 * IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
 * FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
 * AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
 * LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
 * OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN THE
 * SOFTWARE.
 */

package main

import (
	"fmt"
	"html/template"
	"io"
	"os"

	"github.com/labstack/echo/v4"
	mw "github.com/labstack/echo/v4/middleware"
	_ "github.com/lib/pq"

	"github.com/nekobin/nekobin/config"
	"github.com/nekobin/nekobin/database"
	"github.com/nekobin/nekobin/handlers"
	"github.com/nekobin/nekobin/middleware"
)

type Template struct {
	templates *template.Template
}

func (t *Template) Render(w io.Writer, name string, data interface{}, _ echo.Context) error {
	return t.templates.ExecuteTemplate(w, name, data)
}

func main() {
	e := echo.New()

	e.HideBanner = true
	e.Renderer = &Template{
		templates: template.Must(
			template.ParseGlob("./assets/templates/*"),
		),
	}

	cfg := config.Load("config.yaml")
	db := database.NewDatabase(&cfg.Database)

	e.Use(
		mw.LoggerWithConfig(
			mw.LoggerConfig{
				Format: "[${time_rfc3339}] ${status} ${method} ${path} (${remote_ip}) ${latency_human}\n",
				Output: os.Stdout,
			},
		),
		mw.Recover(),
		middleware.Config(cfg),
		middleware.Database(db),
		middleware.About(),
	)

	e.Static("/static", "./assets/static")

	root := e.Group("")
	{
		root.GET("/", handlers.GetRoot)
		root.GET("/:key", handlers.GetRoot)

		getLimiter := middleware.Limiter(cfg.Limits.Documents.Get)
		postLimiter := middleware.Limiter(cfg.Limits.Documents.Post)

		api := root.Group("/api")
		{
			documents := api.Group("/documents")
			{
				documents.GET("/about.md", handlers.GetAbout)
				documents.GET("/:key", handlers.GetDocument, getLimiter)
				documents.POST("", handlers.PostDocument, postLimiter)
			}

			api.GET("/ping", handlers.Pong)
		}

		raw := root.Group("/raw")
		{
			raw.GET("/:key", handlers.GetRawDocument, getLimiter)
		}
	}

	e.Logger.Fatal(e.Start(fmt.Sprintf("%v:%v", cfg.Nekobin.Host, cfg.Nekobin.Port)))
}
