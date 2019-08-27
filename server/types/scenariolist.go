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

import sc "github.com/renom/go-wesnoth/scenario"

type Scenario struct {
	Skip     bool
	Scenario sc.Scenario
}

type ScenarioList []Scenario

func (s *ScenarioList) PickedScenario() *sc.Scenario {
	if s.MustStart() == true {
		for _, v := range *s {
			if v.Skip == false {
				return &v.Scenario
			}
		}
	}
	return &sc.Scenario{}
}

func (s *ScenarioList) PickedIndex() int {
	if s.MustStart() == true {
		for i, v := range *s {
			if v.Skip == false {
				return i
			}
		}
	}
	return -1
}

func (s *ScenarioList) MustStart() bool {
	var count int = 0
	for _, v := range *s {
		if v.Skip == false {
			count++
		}
	}
	if count == 1 {
		return true
	} else {
		return false
	}
}
