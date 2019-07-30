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
	"strconv"
	"strings"
)

func ParseTag(text string) Tag {
	r, _ := regexp.Compile(`^\[([0-9a-z_]+)\]\n` +
		`((?:` + `(?:[\t ]*#textdomain [0-9a-z_-]+\n)?` + `[\t ]*[0-9a-z_]+[\t ]*=.+\n` + `)*)` +
		`((?:` + `(?:[\t ]*#textdomain [0-9a-z_-]+\n)?` + `[\t ]*\[[0-9a-z_]+\]\n` + `(?:.+\n)*` + `[\t ]*\[/[0-9a-z_]+\]\n` + `)*)` +
		`[\t ]*\[/[0-9a-z_]+\]$`)

	submatches := r.FindStringSubmatch(text)
	submatches = append(submatches, make([]string, 4-len(submatches))...)

	return Tag{submatches[1], MergeData(parseAttributes(submatches[2]), parseTags(submatches[3]))}
}

func ParseData(text string) Data {
	r, _ := regexp.Compile(`^` + `((?:` + `(?:[\t ]*#textdomain [0-9a-z_-]+\n)?` +
		`[\t ]*[0-9a-z_]+[\t ]*=.+(?:\n|$)` + `)*)` +
		`((?:` + `(?:[\t ]*#textdomain [0-9a-z_-]+\n)?` +
		`[\t ]*\[[0-9a-z_]+\]\n` + `(?:.+\n)*` + `[\t ]*\[/[0-9a-z_]+\](?:\n|$)` + `)*)`)

	submatches := r.FindStringSubmatch(text)
	submatches = append(submatches, make([]string, 3-len(submatches))...)

	return MergeData(parseAttributes(submatches[1]), parseTags(submatches[2]))
}

func parseAttributes(text string) Data {
	r, _ := regexp.Compile(`(?:[\t ]*#textdomain ([0-9a-z_-]+)\n)?` + `[\t ]*([0-9a-z_]+)[\t ]*=[\t ]*(_?)[\t ]*("?)(.*?)"?[\t ]*(?:\n|$)`)

	data := make(Data)
	float, _ := regexp.Compile(`^\d+\.\d+$`)
	integer, _ := regexp.Compile(`^\d+$`)
	for _, v := range r.FindAllStringSubmatch(text, -1) {
		var value interface{}
		switch {
		case v[3] == "" && v[4] == "\"":
			value = v[5]
		case v[3] == "_" && v[4] == "\"":
			value = Tr(v[5])
		case v[4] == "" && v[5] == "yes":
			value = true
		case v[4] == "" && v[5] == "no":
			value = false
		case v[4] == "" && integer.MatchString(v[5]):
			value, _ = strconv.Atoi(v[5])
		case v[4] == "" && float.MatchString(v[5]):
			value, _ = strconv.ParseFloat(v[5], 64)
		default:
			value = Raw(v[5])
		}
		if v[1] == "" {
			data[v[2]] = value
		} else {
			data[v[2]] = Domain{value, v[1]}
		}
	}

	return data
}

func parseTags(text string) Data {
	i, _ := regexp.Compile(`^(?:[\t ]*#textdomain [0-9a-z_-]+\n)?` + `([\t ]*)\[[0-9a-z_]+\]`)
	submatches := i.FindStringSubmatch(text)
	indent := append(submatches, make([]string, 2-len(submatches))...)[1]
	text = strings.ReplaceAll(strings.Replace(text, indent, "", 1), "\n"+indent, "\n")
	r, _ := regexp.Compile(`(?U)` + `(?:[\t ]*#textdomain ([0-9a-z_-]+)\n)?` + `\[([0-9a-z_]+)\]\n` + `((?:.+\n)*)` + `\[/[0-9a-z_]+\]`)

	data := make(Data)
	for _, v := range r.FindAllStringSubmatch(text, -1) {
		textdomain := v[1]
		key := v[2]
		var value interface{}
		if textdomain == "" {
			value = ParseData(v[3])
		} else {
			value = Domain{ParseData(v[3]), textdomain}
		}
		if data.Contains(key) {
			switch data[key].(type) {
			case Multiple:
				data[key] = append(data[key].(Multiple), value)
			default:
				data[key] = Multiple{data[key], value}
			}
		} else {
			data[key] = value
		}
	}

	return data
}
