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

package config

import (
	"flag"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/renom/fastbot/types"
)

// Default Bot's parameters
var (
	// Accessed from the outside
	Hostname        = "127.0.0.1"
	Port     uint16 = 15000
	Version         = "1.14.7"
	Accounts        = AccountList{Guest("wl_bot_1"), Guest("wl_bot_2")}
	Era             = "default"
	Title           = "Game 1"
	Timer           = TimerConfig{false, 300, 300, 300, 0}
	Admins          = types.StringList{}
	BaseDir         = ""
	Games           = []GameConfig{}
	// Not need to be accessed from the outside
	scenarios = []ScenarioConfig{ScenarioConfig{Path: "/usr/share/wesnoth/data/multiplayer/scenarios/2p_Den_of_Onis.cfg"}}

	// Game distro related confs and timeouts (accessed from the outside)
	Wesnoth = "/usr/bin/wesnoth"
	Eras    = "/usr/share/wesnoth/data/multiplayer/eras.cfg"
	TmpDir  = os.TempDir() + "/fastbot"
	Units   = "/usr/share/wesnoth/data/core/units.cfg"
	Timeout = time.Second * 30
)

type GameConfig struct {
	Players       []string
	PickingPlayer string
	Scenarios     []ScenarioConfig
}

type ScenarioConfig struct {
	Path    string
	Defines []string
}

type TimerConfig struct {
	Enabled       bool
	InitTime      int
	TurnBonus     int
	ReservoirTime int
	ActionBonus   int
}

func LoadFromArgs() {
	// String parameters
	flag.StringVar(&Hostname, "host", Hostname, "The server domain or IP address")
	flag.StringVar(&Version, "version", Version, "The game version")
	flag.StringVar(&Era, "era", Era, "The era name")
	flag.StringVar(&Title, "title", Title, "The game title")
	flag.StringVar(&BaseDir, "baseDir", BaseDir, "A base dir for scenarios")
	flag.StringVar(&Wesnoth, "wesnoth", Wesnoth, "The path to the wesnoth binary")
	flag.StringVar(&Units, "units", Units, "The path to units.cfg")
	flag.StringVar(&Eras, "eras", Eras, "The path to eras.cfg")
	flag.StringVar(&TmpDir, "tmpDir", TmpDir, "The path to the tmp folder")
	// Parameters that require extra check
	portUint := flag.Uint("port", uint(Port), "The port")
	timerString := flag.String("timer", "", "The timer values, comma-separated sequence: init_time,turn_bonus,reservoir_time,action_bonus")
	scenarioString := flag.String("scenarios", "", "The default scenario")
	accountsString := flag.String("accounts", "", "The bot accounts")
	adminsString := flag.String("admins", "", "The admin usernames")
	flag.Parse()
	// Extra check
	if 0 < int(*portUint) && int(*portUint) <= 65535 {
		Port = uint16(*portUint)
	}
	if *adminsString != "" {
		timer := strings.Split(*timerString, ",")
		if len(timer) == 4 {
			var initTime, turnBonus, reservoirTime, actionBonus int
			initTime = types.ParseInt(timer[0], -1)
			turnBonus = types.ParseInt(timer[1], -1)
			reservoirTime = types.ParseInt(timer[2], -1)
			actionBonus = types.ParseInt(timer[3], -1)
			if initTime != -1 && turnBonus != -1 && reservoirTime != -1 && actionBonus != -1 {
				Timer.Enabled = true
				Timer.InitTime = initTime
				Timer.TurnBonus = turnBonus
				Timer.ReservoirTime = reservoirTime
				Timer.ActionBonus = actionBonus
			}
		}
	}
	if *scenarioString != "" {
		s := strings.Split(*scenarioString, ":")
		var d []string
		if len(s) > 1 {
			d = strings.Split(s[1], ",")
		}
		scenarios = []ScenarioConfig{}
		for _, value := range strings.Split(s[0], ",") {
			if BaseDir != "" {
				value = filepath.Clean(BaseDir + "/" + value)
			}
			scenarios = append(scenarios, ScenarioConfig{value, d})
		}
	}
	if *accountsString != "" {
		a := ParseAccounts(*accountsString)
		if len(a) > 0 {
			Accounts = a
		}
	}
	if *adminsString != "" {
		Admins = strings.Split(*adminsString, ",")
	}

	// Load the args
	for _, v := range flag.Args() {
		fields := strings.Split(v, ":")
		if len(fields) > 0 {
			p := strings.Split(fields[0], ",")
			pickingPlayer := ""
			for i, v := range p {
				if v[len(v)-1] == '^' {
					p[i] = v[:len(v)-1]
					pickingPlayer = v[:len(v)-1]
					break
				}
			}
			var d []string // defines
			if len(fields) > 2 {
				d = strings.Split(fields[2], ",")
			}
			var s []ScenarioConfig
			if len(fields) > 1 {
				for _, value := range strings.Split(fields[1], ",") {
					if BaseDir != "" {
						value = filepath.Clean(BaseDir + "/" + value)
					}
					s = append(s, ScenarioConfig{value, d})
				}
			} else {
				s = append(s[:0:0], scenarios...)
			}
			Games = append(Games, GameConfig{p, pickingPlayer, s})
		}
	}
}
