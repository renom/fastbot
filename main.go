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
	"strings"
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
	username        = "wl_bot"
	password        = ""
	timeout         = time.Second * 30
	scenario        = "/usr/share/wesnoth/data/multiplayer/scenarios/2p_Den_of_Onis.cfg"
	era             = "default"
	title           = "Game 1"
	admins          = types.StringList{}
	players         = types.StringList{}
)

func main() {
	// String parameters
	flag.StringVar(&hostname, "host", hostname, "The server domain or IP address")
	flag.StringVar(&version, "version", version, "The game version")
	flag.StringVar(&username, "username", username, "The username")
	flag.StringVar(&password, "password", password, "The password")
	flag.StringVar(&era, "era", era, "The era name")
	flag.StringVar(&title, "title", title, "The game title")
	flag.StringVar(&game.Wesnoth, "wesnoth", game.Wesnoth, "The path to the wesnoth binary")
	flag.StringVar(&server.Units, "units", server.Units, "The path to units.cfg")
	flag.StringVar(&game.Eras, "eras", game.Eras, "The path to eras.cfg")
	flag.StringVar(&game.Path, "tmpDir", game.Path, "The path to the tmp folder")
	// Parameters that require extra check
	portUint := flag.Uint("port", uint(port), "The port")
	adminsString := flag.String("admins", "", "The admins")
	playersString := flag.String("players", "", "The players")
	flag.Parse()
	// Extra check
	if 0 < int(*portUint) && int(*portUint) <= 65535 {
		port = uint16(*portUint)
	}
	if *adminsString != "" {
		admins = strings.Split(*adminsString, ",")
	}
	if *playersString != "" {
		players = strings.Split(*playersString, ",")
	}
	if len(flag.Args()) > 0 {
		scenario = flag.Arg(0)
	}

	game := game.NewGame(title, scenario, era, version)
	s := server.NewServer(hostname, port, version, username, password, title, game.Bytes(), admins, players, timeout)
	s.Connect()
	s.HostGame()
	s.Listen()
}
