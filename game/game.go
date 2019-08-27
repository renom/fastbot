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
	"regexp"
	"strings"

	"github.com/renom/fastbot/era"
	"github.com/renom/fastbot/scenario"
	"github.com/renom/fastbot/wesnoth"
	"github.com/renom/go-wml"
)

var (
	sideData = wml.Data{
		"allow_changes":   true,
		"allow_player":    true,
		"chose_random":    false,
		"faction":         "Random",
		"faction_name":    wml.Domain{wml.Tr("Random"), "wesnoth-multiplayer"},
		"fog":             true,
		"gender":          "null",
		"gold":            100,
		"income":          0,
		"is_host":         false,
		"is_local":        false,
		"random_faction":  true,
		"shroud":          false,
		"type":            "null",
		"village_gold":    2,
		"village_support": 1,
		"default_faction": wml.Data{},
		// Must be defined inside a real scenario:
		//"canrecruit":     true,
		//"controller":     "human",
		//"side":            1,
		//"team_name":       "whatever",
		//"user_team_name":  "whatever",
		//"ai":              wml.Data{"villages_per_scout": 8},
		// Must be manually defined:
		//"color":           "red",
	}
)

type Game struct {
	Title      string
	Scenario   scenario.Scenario
	Era        era.Era
	Version    string
	NotNewGame bool   // To set up manually
	Player1    string // To set up manually
	Player2    string // To set up manually
	Id         string // Obtained by Parse()
	Name       string // Obtained by Parse()
	scenario   string // Obtained by Parse()
	// Timer-related config
	TimerEnabled  bool
	InitTime      int
	TurnBonus     int
	ReservoirTime int
	ActionBonus   int
}

func NewGame(title string, scenario scenario.Scenario, era era.Era,
	timerEnabled bool, initTime int, turnBonus int, reservoirTime int,
	actionBonus int, version string) Game {
	game := Game{Title: title, Scenario: scenario, Era: era, Version: version,
		TimerEnabled: timerEnabled, InitTime: initTime, TurnBonus: turnBonus,
		ReservoirTime: reservoirTime, ActionBonus: actionBonus}
	game.Parse()
	return game
}

func (g *Game) Parse() {
	replacer := strings.NewReplacer("[multiplayer]", "[scenario]",
		"[/multiplayer]", "[/scenario]")
	g.scenario = replacer.Replace(string(wesnoth.Preprocess(g.Scenario.Path(), g.Scenario.Defines())))
	s, _ := regexp.Compile(`(?U)\[scenario\]\n(?:[^\[\]]*\n)*\tid="(.*)"\n(?:.*\n)*\tname=_?"(.*)"\n(?:.*\n)*\[/scenario\]`)
	g.Id = s.FindStringSubmatch(g.scenario)[1]
	g.Name = s.FindStringSubmatch(g.scenario)[2]
}

func (g Game) Bytes() []byte {
	return []byte(g.String())
}

func (g Game) String() string {
	return g.topLevel() +
		g.scenarioBlock() +
		g.carryoverBlock() +
		g.multiplayerBlock() +
		g.eraBlock()
}

func (g *Game) topLevel() string {
	topLevel := wml.Data{
		"abbrev":                 "",
		"campaign":               "",
		"campaign_define":        "",
		"campaign_extra_defines": "",
		"campaign_name":          "",
		"campaign_type":          "multiplayer",
		"difficulty":             "NORMAL",
		"end_credits":            true,
		"end_text":               "",
		"end_text_duration":      0,
		"era_define":             "",
		"label":                  g.Name,
		"mod_defines":            "",
		"oos_debug":              false,
		"random_mode":            "",
		"scenario_define":        "",
		"version":                g.Version,
		"replay":                 wml.Data{"upload_log": wml.Data{}},
	}
	return topLevel.String()
}

func (g *Game) scenarioBlock() string {
	r, _ := regexp.Compile(`(?U)\[side\]\n(?:[^\[\]]*\n)*[\t ]*controller="human"\n(?:.*\n)*([\t ]*)\[/side\]`)

	sides := r.FindAllString(g.scenario, -1)
	indent := uint(strings.Count(r.FindStringSubmatch(g.scenario)[1], "\t"))

	scenario := g.scenario
	scenario = replaceSide(scenario, sideTag(wml.MergeData(sideData, wml.ParseTag(sides[0]).Data), "red", g.Player1), indent)
	scenario = replaceSide(scenario, sideTag(wml.MergeData(sideData, wml.ParseTag(sides[1]).Data), "blue", g.Player2), indent)
	r_attributes, _ := regexp.Compile(`(?sU)^(\[scenario\]\n)(.*\n)([\t ]*\[/?[0-9a-z_]+\])`)
	matches := r_attributes.FindStringSubmatch(scenario)
	attributes := wml.ParseData(matches[2])
	if g.NotNewGame == true {
		attributes["allow_new_game"] = false
		attributes["disallow_recall"] = true
	}
	attributes["experience_modifier"] = 70
	attributes["has_mod_events"] = true
	attributes["objectives"] = "<big>Victory:</big>\n<span color='#00ff00'>• Defeat enemy leader(s)</span>"
	attributes["turns"] = -1
	events := wml.Multiple{}
	for _, v := range g.Era.Events {
		events = append(events, v)
	}
	attributes["event"] = events
	scenario = r_attributes.ReplaceAllString(scenario, matches[1]+attributes.Indent(1)+matches[3])

	return scenario
}

func (g *Game) carryoverBlock() string {
	carryover := wml.Tag{"carryover_sides_start", wml.Data{
		"next_scenario": g.Id,
		"random_calls":  0,
		"random_seed":   randomSeed(),
		"variables":     wml.Data{},
	}}
	return carryover.String()
}

func (g *Game) multiplayerBlock() string {
	multiplayer := wml.Tag{"multiplayer", wml.Data{
		"active_mods":                 "",
		"difficulty_define":           "NORMAL",
		"experience_modifier":         70,
		"hash":                        "",
		"mp_campaign":                 "",
		"mp_campaign_name":            "",
		"mp_countdown":                g.TimerEnabled,
		"mp_countdown_action_bonus":   g.ActionBonus,
		"mp_countdown_init_time":      g.InitTime,
		"mp_countdown_reservoir_time": g.ReservoirTime,
		"mp_countdown_turn_bonus":     g.TurnBonus,
		"mp_era":                      "era_" + g.Era.Id,
		"mp_era_name":                 g.Era.Name,
		"mp_fog":                      true,
		"mp_num_turns":                -1,
		"mp_random_start_time":        false,
		"mp_scenario":                 g.Id,
		"mp_scenario_name":            g.Name,
		"mp_shroud":                   false,
		"mp_use_map_settings":         true,
		"mp_village_gold":             2,
		"mp_village_support":          1,
		"observer":                    true,
		"random_faction_mode":         "No Mirror",
		"registered_users_only":       false,
		"savegame":                    false,
		"scenario":                    g.Title,
		"shuffle_sides":               true,
		"side_users":                  "",
		"options":                     wml.Data{},
	}}
	return multiplayer.String()
}

func (g *Game) eraBlock() string {
	return g.Era.Body
}
