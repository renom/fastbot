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
	"strconv"
	"strings"
	"sync"

	"github.com/renom/fastbot/config"
	"github.com/renom/fastbot/scenario"
	"github.com/renom/fastbot/server"
)

func main() {
	config.LoadFromArgs()

	var wg sync.WaitGroup
	for i, v := range config.Games {
		t := config.Title
		for j, w := range v.Players {
			t = strings.ReplaceAll(t, "{Player"+strconv.Itoa(j+1)+"}", w)
		}
		var scenarios []scenario.Scenario
		for _, x := range v.Scenarios {
			scenarios = append(scenarios, scenario.FromPath(x.Path, x.Defines))
		}
		srv := server.NewServer(
			config.Hostname,
			config.Port,
			config.Version,
			config.Accounts[i].Username,
			config.Accounts[i].Password,
			config.Era,
			t,
			scenarios,
			config.Admins,
			v.Players,
			v.PickingPlayer,
			config.Timer.Enabled,
			config.Timer.InitTime,
			config.Timer.TurnBonus,
			config.Timer.ReservoirTime,
			config.Timer.ActionBonus,
			config.Timeout)
		wg.Add(1)
		go func() {
			srv.Connect()
			srv.HostGame()
			srv.Listen()
			wg.Done()
		}()
	}
	wg.Wait()
}
