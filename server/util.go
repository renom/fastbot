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

package server

import (
	"math/rand"
	"regexp"
	"strings"
	"time"

	"github.com/renom/fastbot/config"
	serverTypes "github.com/renom/fastbot/server/types"
	"github.com/renom/fastbot/wesnoth"
	"github.com/renom/wml"
)

func SplitMessage(text string) []string {
	upperLimit := 256

	result := []string{}

	for pos := 0; pos < len(text); {
		from := pos
		var to int
		if pos+upperLimit < len(text) {
			splitIndex := strings.LastIndex(text[pos:pos+upperLimit+1], "\n")
			if splitIndex == -1 {
				if text[pos+upperLimit-1] == '"' && strings.Count(text[pos:pos+upperLimit], "\"\"")%2 != 0 {
					splitIndex = pos + upperLimit
				}
			}
			if splitIndex != -1 {
				to = pos + splitIndex
			} else {
				to = pos + upperLimit
			}
		} else {
			to = len(text)
		}
		result = append(result, text[from:to])
		pos = to
	}
	return result
}

func insertFaction(side *serverTypes.Side, faction wml.Data, textdomain string) wml.Data {
	var leaders = []string{}
	if faction.Contains("random_leader") {
		leaders = strings.Split(faction["random_leader"].(string), ",")
	} else {
		leaders = strings.Split(faction["leader"].(string), ",")
	}
	rand.Seed(time.Now().UTC().UnixNano())
	leader := leaders[rand.Int31n(int32(len(leaders)))]
	r, _ := regexp.Compile(`(?U)\[unit_type\]\n` +
		`(?:(?:[\t ]*[0-9a-z_]+=.*\n)|[^\[\]]*\n)*` +
		`[\t ]*gender="([a-z,]+)"\n` +
		`(?:[\t ]*[0-9a-z_]+=.*\n)*` +
		`[\t ]*id="` + leader + `"\n`)

	var gender string
	if subString := r.FindSubmatch(wesnoth.Preprocess(config.Units, nil)); len(subString) == 2 {
		gender = string(subString[1])
		if genders := strings.Split(gender, ","); len(genders) == 2 {
			rand.Seed(time.Now().UTC().UnixNano())
			gender = genders[rand.Int31n(2)]
		}
	} else {
		gender = "male"
	}
	sideData := wml.Data{"insert": wml.Data{
		"chose_random":   true,
		"color":          side.Color,
		"current_player": side.Player,
		"faction":        faction["id"],
		"gender":         gender,
		"is_host":        false,
		"is_local":       false,
		"leader":         faction["leader"],
		"name":           side.Player,
		"player_id":      side.Player,
		"recruit":        faction["recruit"],
		"terrain_liked":  faction["terrain_liked"],
		"type":           leader,
		// Not necessary since already defined in the scenario:
		//"user_team_name": "whatever",
	},
		"delete":       wml.Data{"random_faction": "x"},
		"insert_child": wml.Data{"index": 1, "ai": faction["ai"]},
	}
	if textdomain != "" {
		sideData["insert"].(wml.Data)["faction_name"] = wml.Domain{faction["name"], textdomain}
	} else {
		sideData["insert"].(wml.Data)["faction_name"] = faction["name"]
	}
	if faction.Contains("random_leader") {
		sideData["insert"].(wml.Data)["random_leader"] = faction["random_leader"]
	}
	return sideData
}
