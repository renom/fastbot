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

package types

import "regexp"

type StringList []string

func (p *StringList) ContainsValue(value string) bool {
	for _, v := range *p {
		if v == value {
			return true
		}
	}
	return false
}

func (p *StringList) IndexOf(value string) int {
	for i, v := range *p {
		if v == value {
			return i
		}
	}
	return -1
}

func (p *StringList) Delete(i int) {
	*p = append((*p)[:i], (*p)[i+1:]...)
}

func (p *StringList) DeleteValue(value string) {
	p.Delete(p.IndexOf(value))
}

func (p *StringList) Unique() StringList {
	result := StringList{}
	for _, v := range *p {
		if !result.ContainsValue(v) {
			result = append(result, v)
		}
	}
	return result
}

func (p *StringList) Match(expr string) bool {
	r, err := regexp.Compile(expr)
	if err != nil {
		return false
	}
	for _, v := range *p {
		if r.MatchString(v) == false {
			return false
		}
	}
	return true
}
