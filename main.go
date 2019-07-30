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
	"time"

	"github.com/renom/fastbot/game"
	"github.com/renom/fastbot/server"
	"github.com/renom/fastbot/types"
)

// Bot's parameters
const (
	hostname = "127.0.0.1"
	port     = 15000
	version  = "1.14.7"
	username = "wl_bot"
	password = ""
	timeout  = time.Second * 30
	scenario = "/usr/share/wesnoth/data/multiplayer/scenarios/2p_Den_of_Onis.cfg"
	era      = "default"
	title    = "Game 1"
)

// The same but arrays
var admins = types.StringList{"Player1"}
var players = types.StringList{"Player2", "Player3"}

func main() {
	game := game.NewGame(title, scenario, era, version)
	s := server.NewServer(hostname, port, version, username, password, title, game.Bytes(), admins, players, timeout)
	s.Connect()
	s.HostGame()
	s.Listen()
}
