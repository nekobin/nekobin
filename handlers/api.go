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

package handlers

import (
	"net/http"
	"strings"

	"github.com/labstack/echo/v4"

	"github.com/nekobin/nekobin/config"
	"github.com/nekobin/nekobin/database"
	"github.com/nekobin/nekobin/response"
)

func GetAbout(ctx echo.Context) error {
	about := ctx.Get("about").(*database.Document)

	return ctx.JSON(
		http.StatusOK,
		response.NewResult(about),
	)
}

func GetDocument(ctx echo.Context) error {
	db := ctx.Get("db").(*database.Database)
	key := strings.Split(ctx.Param("key"), ".")[0]
	doc, err := db.Documents.Select(key)

	if err != nil {
		return ctx.JSON(
			http.StatusBadRequest,
			response.ErrorDocumentNotFound,
		)
	}

	go db.Documents.IncrementViews(key, ctx.RealIP())

	return ctx.JSON(
		http.StatusOK,
		response.NewResult(doc),
	)
}

func PostDocument(ctx echo.Context) error {
	doc := &database.Document{}

	if err := ctx.Bind(doc); err != nil {
		return ctx.JSON(
			http.StatusBadRequest,
			response.ErrorInvalidData,
		)
	}

	title, author, content := doc.Title, doc.Author, doc.Content

	cfg := ctx.Get("cfg").(*config.Config)

	if title != nil {
		switch length := len(*title); {
		case length == 0:
			title = nil
		case length > cfg.Nekobin.MaxTitleLength:
			return ctx.JSON(
				http.StatusBadRequest,
				response.ErrorTitleTooLong,
			)
		}
	}

	if author != nil {
		switch length := len(*author); {
		case length == 0:
			author = nil
		case length > cfg.Nekobin.MaxAuthorLength:
			return ctx.JSON(
				http.StatusBadRequest,
				response.ErrorAuthorTooLong,
			)
		}
	}

	if len(content) == 0 {
		return ctx.JSON(
			http.StatusBadRequest,
			response.ErrorContentEmpty,
		)
	}

	if len(content) > cfg.Nekobin.MaxContentLength {
		return ctx.JSON(
			http.StatusBadRequest,
			response.ErrorContentTooLong,
		)
	}

	db := ctx.Get("db").(*database.Database)
	doc, err := db.Documents.Insert(title, author, content)

	if err != nil {
		return err
	}

	go db.Documents.IncrementViews(doc.Key, ctx.RealIP())

	return ctx.JSON(
		http.StatusCreated,
		response.NewResult(doc),
	)
}

func Pong(ctx echo.Context) error {
	return ctx.JSON(
		http.StatusOK,
		response.NewResult("pong"),
	)
}
