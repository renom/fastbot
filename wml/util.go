// This file is part of Fastbot.
//
// Fastbot is free software: you can redistribute it and/or modify
// it under the terms of the GNU General Public License as published by
// the Free Software Foundation, either version 3 of the License, or
// (at your option) any later version.
//
// Fastbot is distributed in the hope that it will be useful,
// but WITHOUT ANY WARRANTY; without even the implied warranty of
// MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
// GNU General Public License for more details.
//
// You should have received a copy of the GNU General Public License
// along with Fastbot.  If not, see <https://www.gnu.org/licenses/>.

package wml

import (
	"regexp"
	"strings"
)

func EscapeString(s string) string {
	return strings.ReplaceAll(s, "\"", "\"\"")
}

func IndentString(text string, indent uint) string {
	r, _ := regexp.Compile(`(?m)^([^$])`)
	return r.ReplaceAllString(text, strings.Repeat("\t", int(indent))+"$1")
}