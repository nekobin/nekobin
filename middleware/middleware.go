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

package middleware

import (
	"github.com/labstack/echo/v4"
	"io/ioutil"
	"log"
	"nekobin/config"
	"nekobin/database"
	"nekobin/limiter"
	"nekobin/response"
	"net/http"
)

// Middleware to add the configuration in handlers
func Config(cfg *config.Config) echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(ctx echo.Context) error {
			ctx.Set("cfg", cfg)
			return next(ctx)
		}
	}
}

// Middleware to add Database context in handlers
func Database(db *database.Database) echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(ctx echo.Context) error {
			ctx.Set("db", db)
			return next(ctx)
		}
	}
}

// Middleware to make the About document available in handlers
func About() echo.MiddlewareFunc {
	file, err := ioutil.ReadFile("./README.md")
	if err != nil {
		log.Fatal(err)
	}

	about := &database.Document{
		Key:     "about",
		Content: string(file),
	}

	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(ctx echo.Context) error {
			ctx.Set("about", about)
			return next(ctx)
		}
	}
}

// Middleware to limit requests
func Limiter(limits []limiter.Limit) echo.MiddlewareFunc {
	lim := limiter.NewLimiter(limits...)

	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(ctx echo.Context) error {
			if !lim.IsAllowed(ctx.RealIP()) {
				return ctx.JSON(
					http.StatusTooManyRequests,
					response.ErrorTooFast,
				)
			}

			return next(ctx)
		}
	}
}
