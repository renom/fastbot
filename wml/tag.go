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
	"strings"
)

type Tag struct {
	Name string
	Data Data
}

func EmptyTag(name string) Tag {
	return Tag{name, Data{}}
}

func MergeTags(first Tag, second Tag) Tag {
	return Tag{second.Name, MergeData(first.Data, second.Data)}
}

func (t Tag) Bytes() []byte {
	return []byte(t.String())
}

func (t Tag) String() string {
	return t.Indent(0)
}

func (t *Tag) Indent(indent uint) string {
	tabulation := strings.Repeat("\t", int(indent))
	return tabulation + "[" + t.Name + "]\n" +
		t.Data.Indent(indent+1) +
		tabulation + "[/" + t.Name + "]\n"
}
