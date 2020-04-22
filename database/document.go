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

package database

import (
	"fmt"
	"log"
	"time"

	"github.com/jmoiron/sqlx"

	"github.com/nekobin/nekobin/keygen"
)

type Document struct {
	Key     string  `json:"key"`
	Title   *string `json:"title"`
	Author  *string `json:"author"`
	Date    int     `json:"date"`
	Views   int     `json:"views"`
	Length  int     `json:"length"`
	Content string  `json:"content"`
}

type DocumentsQuery interface {
	Select(key string) (doc *Document, err error)
	Insert(title, author *string, content string) (doc *Document, err error)
	Exists(key string) (exists bool, err error)
	IncrementViews(key, ip string)
}

type ViewIPsKey struct {
	documentKey string
	ipAddress   string
}

type Documents struct {
	*sqlx.DB

	keygen  keygen.Keygen
	viewIPs map[ViewIPsKey]time.Time
}

func NewDocuments(db *sqlx.DB) *Documents {
	return &Documents{
		DB:      db,
		keygen:  keygen.NewPhoneticKeygen(),
		viewIPs: make(map[ViewIPsKey]time.Time),
	}
}

func (docs *Documents) Select(key string) (doc *Document, err error) {
	row := docs.QueryRowx(`
		SELECT
			key, title, author,
			extract(EPOCH FROM date AT TIME ZONE 'utc')::INT date,
			views, length, content
		FROM documents
		WHERE key = $1
		LIMIT 1`,
		key,
	)

	doc = &Document{}
	err = row.StructScan(doc)

	return
}

func (docs *Documents) Insert(title, author *string, content string) (doc *Document, err error) {
	if title != nil && *title == "" {
		title = nil
	}

	if author != nil && *author == "" {
		author = nil
	}

	var key string
	for {
		key = docs.keygen.GenerateKey()
		exists, err := docs.Exists(key)

		if err != nil {
			log.Println(err)
			return nil, err
		}

		if !exists {
			break
		}
	}

	rows, err := docs.Query(
		"INSERT INTO documents (key, title, author, length, content) VALUES ($1, $2, $3, $4, $5)",
		key, title, author, len(content), content,
	)

	if err == nil {
		defer func() {
			err := rows.Close()

			if err != nil {
				log.Println(err)
			}
		}()
	}

	doc, err = docs.Select(key)

	return
}

func (docs *Documents) Exists(key string) (exists bool, err error) {
	row := docs.QueryRowx("SELECT EXISTS(SELECT 1 FROM documents WHERE key = $1)", key)
	err = row.Scan(&exists)

	return
}

func (docs *Documents) IncrementViews(key, ip string) {
	viewIPsKey := ViewIPsKey{key, ip}
	value, exists := docs.viewIPs[viewIPsKey]

	if exists && time.Now().Sub(value).Minutes() < 30 {
		return
	}

	docs.viewIPs[viewIPsKey] = time.Now()

	rows, err := docs.Query("UPDATE documents SET views = views + 1 WHERE key = $1", key)

	if err != nil {
		fmt.Println(err)
		return
	}

	err = rows.Close()

	if err != nil {
		log.Println(err)
	}
}
