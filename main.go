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

// fastbot project main.go
package main

import (
	"flag"
	"path/filepath"
	"strconv"
	"strings"
	"sync"
	"time"

	"github.com/renom/fastbot/game"
	"github.com/renom/fastbot/server"
	"github.com/renom/fastbot/types"
)

// Default Bot's parameters
var (
	hostname        = "127.0.0.1"
	port     uint16 = 15000
	version         = "1.14.7"
	accounts        = AccountList{Guest("wl_bot_1"), Guest("wl_bot_2")}
	timeout         = time.Second * 30
	scenario        = "/usr/share/wesnoth/data/multiplayer/scenarios/2p_Den_of_Onis.cfg"
	defines         = []string{}
	era             = "default"
	title           = "Game 1"
	admins          = types.StringList{}
	players         = types.StringList{}
	baseDir         = ""
)

func main() {
	// String parameters
	flag.StringVar(&hostname, "host", hostname, "The server domain or IP address")
	flag.StringVar(&version, "version", version, "The game version")
	flag.StringVar(&era, "era", era, "The era name")
	flag.StringVar(&title, "title", title, "The game title")
	flag.StringVar(&baseDir, "baseDir", baseDir, "A base dir for scenarios")
	flag.StringVar(&game.Wesnoth, "wesnoth", game.Wesnoth, "The path to the wesnoth binary")
	flag.StringVar(&server.Units, "units", server.Units, "The path to units.cfg")
	flag.StringVar(&game.Eras, "eras", game.Eras, "The path to eras.cfg")
	flag.StringVar(&game.Path, "tmpDir", game.Path, "The path to the tmp folder")
	// Parameters that require extra check
	portUint := flag.Uint("port", uint(port), "The port")
	scenarioString := flag.String("scenario", "", "The default scenario")
	accountsString := flag.String("accounts", "", "The accounts")
	adminsString := flag.String("admins", "", "The admins")
	playersString := flag.String("players", "", "The players")
	flag.Parse()
	// Extra check
	if 0 < int(*portUint) && int(*portUint) <= 65535 {
		port = uint16(*portUint)
	}
	if *scenarioString != "" {
		s := strings.Split(*scenarioString, ":")
		scenario = s[0]
		if len(s) > 1 {
			defines = strings.Split(s[1], ",")
		}
	}
	if *accountsString != "" {
		a := ParseAccounts(*accountsString)
		if len(a) > 0 {
			accounts = a
		}
	}
	if *adminsString != "" {
		admins = strings.Split(*adminsString, ",")
	}
	if *playersString != "" {
		players = strings.Split(*playersString, ",")
	}

	var wg sync.WaitGroup
	i := 0
	for _, v := range flag.Args() {
		fields := strings.Split(v, ":")
		if len(fields) > 0 {
			p := strings.Split(fields[0], ",")
			var s string
			var d []string // defines
			if len(fields) > 1 {
				s = fields[1]
			} else {
				s = scenario
			}
			if len(fields) > 2 {
				d = strings.Split(fields[2], ",")
			} else {
				d = defines
			}
			if baseDir != "" {
				s = filepath.Clean(baseDir + "/" + s)
			}

			game := game.NewGame(title, s, d, era, version)
			t := title
			for j, w := range p {
				t = strings.ReplaceAll(t, "{Player"+strconv.Itoa(j+1)+"}", w)
			}
			srv := server.NewServer(hostname, port, version, accounts[i].Username, accounts[i].Password, t, game.Bytes(), admins, p, timeout)
			wg.Add(1)
			go func() {
				srv.Connect()
				srv.HostGame()
				srv.Listen()
				wg.Done()
			}()
			i++
		}
	}
	wg.Wait()
}
