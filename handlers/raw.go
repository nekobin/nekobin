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
	"github.com/labstack/echo/v4"
	"nekobin/database"
	"nekobin/response"
	"net/http"
	"strconv"
	"strings"
)

func GetRawDocument(ctx echo.Context) error {
	db := ctx.Get("db").(*database.Database)
	key := strings.Split(ctx.Param("key"), ".")[0]
	doc, err := db.Documents.Select(key)

	if err != nil {
		return ctx.String(
			http.StatusBadRequest,
			response.ErrorDocumentNotFound.Error,
		)
	}

	go db.Documents.IncrementViews(key, ctx.RealIP())

	if doc.Title != nil {
		ctx.Response().Header().Set("Document-Title", *doc.Title)
	}

	if doc.Author != nil {
		ctx.Response().Header().Set("Document-Author", *doc.Author)
	}

	ctx.Response().Header().Set("Document-Date", strconv.Itoa(doc.Date))
	ctx.Response().Header().Set("Document-Views", strconv.Itoa(doc.Views))
	ctx.Response().Header().Set("Document-length", strconv.Itoa(doc.Length))

	return ctx.String(
		http.StatusOK,
		doc.Content,
	)
}
