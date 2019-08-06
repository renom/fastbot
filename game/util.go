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

package game

import (
	"fmt"
	"math/rand"
	"regexp"
	"time"

	"github.com/renom/fastbot/wml"
)

func replaceSide(scenario string, side wml.Tag, indent uint) string {
	r, _ := regexp.Compile(`(?U)[\t ]*\[side\]\n([^\[\]]*\n)*[\t ]*side=` + string(side.Data["side"].(int)) + `\n(.*\n)*[\t ]*\[/side\]\n`)
	return r.ReplaceAllString(scenario, side.Indent(indent))
}

func randomSeed() string {
	rand.Seed(time.Now().UTC().UnixNano())
	seed := fmt.Sprintf("%x", rand.Int63n(4294967295+1))
	return seed
}
