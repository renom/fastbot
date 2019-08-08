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

import (
	"math/rand"
	"time"
)

type Side struct {
	Side   int
	Player string
	Color  string
	Ready  bool
}

type SideList []*Side

func (s *SideList) HasSide(side int) bool {
	for _, v := range *s {
		if v.Side == side {
			return true
		}
	}
	return false
}

func (s *SideList) HasPlayer(player string) bool {
	for _, v := range *s {
		if v.Player == player {
			return true
		}
	}
	return false
}

func (s *SideList) HasColor(color string) bool {
	for _, v := range *s {
		if v.Color == color {
			return true
		}
	}
	return false
}

func (s *SideList) Side(side int) *Side {
	for _, v := range *s {
		if v.Side == side {
			return v
		}
	}
	return &Side{}
}

func (s *SideList) Find(player string) *Side {
	for _, v := range *s {
		if v.Player == player {
			return v
		}
	}
	return &Side{}
}

func (s *SideList) FreeSlots() int {
	result := 0
	for _, v := range *s {
		if v.Player == "" {
			result++
		}
	}
	return result
}

func (s *SideList) MustStart() bool {
	for _, v := range *s {
		if v.Ready == false {
			return false
		}
	}
	return true
}

func (s *SideList) Shuffle() {
	rand.Seed(time.Now().UTC().UnixNano())
	rand.Shuffle(len(*s), func(i, j int) {
		(*s)[i], (*s)[j] = (*s)[j], (*s)[i]
		(*s)[i].Side, (*s)[j].Side = (*s)[j].Side, (*s)[i].Side
	})
}
