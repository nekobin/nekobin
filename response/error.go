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

package response

type Error struct {
	Ok    bool   `json:"ok"`
	Error string `json:"error"`
}

func NewError(error string) *Error {
	return &Error{
		Ok:    false,
		Error: error,
	}
}

var (
	ErrorDocumentNotFound = NewError("DOCUMENT_NOT_FOUND")
	ErrorInvalidData      = NewError("INVALID_DATA")
	ErrorTitleTooLong     = NewError("TITLE_TOO_LONG")
	ErrorAuthorTooLong    = NewError("AUTHOR_TOO_LONG")
	ErrorContentEmpty     = NewError("CONTENT_EMPTY")
	ErrorContentTooLong   = NewError("CONTENT_TOO_LONG")
	ErrorTooFast          = NewError("TOO_FAST")
)
