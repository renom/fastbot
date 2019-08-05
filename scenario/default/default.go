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

package default

const Default = Scenario{
	name: "WL Bot",
	body: `[multiplayer]
	id="WL_Bot"
	map_data="border_size=1
usage=map

Xv, Xv, Xv, Xv, Xv, Xv, Xv, Xv, Xv, Xv, Wot, Wot, Wot, Wot, Wot, Wot, Wot, Wot, Wot, Wot, Wot, Wot, Wot, Wot, Wot, Wot, Wot, Wot, Wot, Xv, Xv, Xv, Xv, Xv, Xv, Xv, Xv, Xv, Xv, Xv
Xv, Xv, Xv, Xv, Xv, Xv, Xv, Xv, Wot, Wot, Wog, Wot, Wog, Wot, Wog, Wot, Wog, Wot, Wog, Wot, Wog, Wot, Wog, Wot, Wog, Wot, Wog, Wot, Wog, Wot, Wot, Xv, Xv, Xv, Xv, Xv, Xv, Xv, Xv, Xv
Xv, Xv, Xv, Xv, Xv, Xv, Wot, Wot, Wog, Wog, Wog, Wog, Wog, Wog, Wog, Wog, Wog, Wog, Wog, Wog, Wog, Wog, Wog, Wog, Wog, Wog, Wog, Wog, Wog, Wog, Wog, Wot, Wot, Xv, Xv, Xv, Xv, Xv, Xv, Xv
Xv, Xv, Xv, Xv, Wot, Wot, Wog, Wog, Wog, Wog, Wog, Wog, Wog, Wog, Wog, Wog, Wog, Wog, Wog, Wog, Wog, Wog, Wog, Wog, Wog, Wog, Wog, Wog, Wog, Wog, Wog, Wog, Wog, Wot, Wot, Xv, Xv, Xv, Xv, Xv
Xv, Xv, Wot, Wot, Wog, Wog, Wog, Wog, Wog, Wog, Wog, Wog, Wog, Wog, Wog, Wog, Wog, Wog, Wog, Wog, Wog, Wog, Wog, Wog, Wog, Wog, Wog, Wog, Wog, Wog, Wog, Wog, Wog, Wog, Wog, Wot, Wot, Xv, Xv, Xv
Xv, Wot, Wog, Wog, Wog, Wog, Wog, Wog, Wog, Wog, Wog, Wog, Wog, Wog, Wog, Wog, Wog, Wog, Wog, Wog, Wog, Wog, Wog, Wog, Wog, Wog, Wog, Wog, Wog, Wog, Wog, Wog, Wog, Wog, Wog, Wog, Wog, Wot, Xv, Xv
Xv, Wot, Wog, Wog, Wog, Wog, Wog, Wog, Wog, Wog, Wog, Wog, Wog, Wog, Wog, Wog, Wog, Wog, Wog, Wog, Wog, Wog, Wog, Wog, Wog, Wog, Wog, Wog, Wog, Wog, Wog, Wog, Wog, Wog, Wog, Wog, Wog, Wot, Xv, Xv
Xv, Wot, Wog, Wog, Wog, Wog, Wwt, Wog, Wog, Wog, Wog, Wog, Wog, Wog, Wog, Wog, Wog, Wog, Wwt, Wog, Wog, Wog, Wog, Wog, Wog, Wog, Wwt, Wog, Wog, Wog, Wog, Wog, Wog, Wog, Wog, Wog, Wog, Wot, Xv, Xv
Xv, Wot, Wog, Wog, Wwt, Wwt, Wwt, Wwt, Wog, Wog, Wog, Wog, Wog, Wog, Wog, Wog, Wog, Wwt, Wwt, Wwt, Wwt, Wog, Wog, Wog, Wwt, Wwt, Wwt, Wwt, Wog, Wog, Wog, Wog, Wog, Wog, Wog, Wog, Wog, Wot, Xv, Xv
Xv, Wot, Wog, Wog, Wwt, Wwt, Wwt, Wwt, Wog, Wog, Wog, Wog, Wog, Wog, Wog, Wog, Wog, Wwt, Wwt, Wwt, Wwt, Wog, Wog, Wog, Wwt, Wwt, Wwt, Wwt, Wog, Wog, Wog, Wog, Wog, Wog, Wog, Wog, Wog, Wot, Xv, Xv
Xv, Wot, Wog, Wog, Wwt, Wwt, Wwt, Wwt, Wog, Wog, Wog, Wog, Wog, Wog, Wog, Wog, Wog, Wwt, Wwt, Wwt, Wwt, Wog, Wog, Wog, Wwt, Wwt, Wwt, Wwt, Wog, Wog, Wog, Wog, Wog, Wog, Wog, Wog, Wog, Wot, Xv, Xv
Xv, Wot, Wog, Wog, Wwt, Wwt, Wwt, Wwt, Wog, Wog, Wog, Wog, Wwt, Wog, Wog, Wog, Wog, Wwt, Wwt, Wwt, Wwt, Wog, Wog, Wog, Wwt, Wwt, Wwt, Wwt, Wog, Wog, Wog, Wog, Wog, Wog, Wog, Wog, Wog, Wot, Xv, Xv
Xv, Wot, Wog, Wog, Wwt, Wwt, Wwt, Wwt, Wog, Wog, Wog, Wwt, Wwt, Wwt, Wog, Wog, Wog, Wwt, Wwt, Wwt, Wwt, Wog, Wog, Wog, Wwt, Wwt, Wwt, Wwt, Wog, Wog, Wog, Wog, Wog, Wog, Wog, Wog, Wog, Wot, Xv, Xv
Xv, Wot, Wog, Wog, Wwt, Wwt, Wwt, Wwt, Wog, Wog, Wog, Wwt, Wwt, Wwt, Wog, Wog, Wog, Wwt, Wwt, Wwt, Wwt, Wog, Wog, Wog, Wwt, Wwt, Wwt, Wwt, Wog, Wog, Wog, Wog, Wog, Wog, Wog, Wog, Wog, Wot, Xv, Xv
Xv, Wot, Wog, Wog, Wwt, Wwt, Wwt, Wwt, Wwt, Wog, Wwt, Wwt, Wwt, Wwt, Wwt, Wog, Wwt, Wwt, Wwt, Wwt, Wwt, Wog, Wog, Wog, Wwt, Wwt, Wwt, Wwt, Wog, Wog, Wog, Wog, Wog, Wog, Wog, Wog, Wog, Wot, Xv, Xv
Xv, Wot, Wog, Wog, Wwt, Wwt, Wwt, Wwt, Wwt, Wwt, Wwt, Wwt, Wwt, Wwt, Wwt, Wwt, Wwt, Wwt, Wwt, Wwt, Wwt, Wog, Wog, Wog, Wwt, Wwt, Wwt, Wwt, Wwt, Wwt, Wwt, Wwt, Wwt, Wwt, Wwt, Wwt, Wog, Wot, Xv, Xv
Xv, Wot, Wog, Wog, Wog, Wwt, Wwt, Wwt, Wwt, Wwt, Wwt, Wwt, Wwt, Wwt, Wwt, Wwt, Wwt, Wwt, Wwt, Wwt, Wog, Wog, Wog, Wog, Wwt, Wwt, Wwt, Wwt, Wwt, Wwt, Wwt, Wwt, Wwt, Wwt, Wwt, Wwt, Wog, Wot, Xv, Xv
Xv, Wot, Wog, Wog, Wog, Wog, Wog, Wwt, Wwt, Wwt, Wwt, Wwt, Wog, Wwt, Wwt, Wwt, Wwt, Wwt, Wog, Wog, Wog, Wog, Wog, Wog, Wwt, Wwt, Wwt, Wwt, Wwt, Wwt, Wwt, Wwt, Wwt, Wwt, Wwt, Wwt, Wog, Wot, Xv, Xv
Xv, Wot, Wog, Wog, Wog, Wog, Wog, Wog, Wog, Wwt, Wog, Wog, Wog, Wog, Wog, Wwt, Wog, Wog, Wog, Wog, Wog, Wog, Wog, Wog, Wog, Wwt, Wog, Wwt, Wog, Wwt, Wog, Wwt, Wog, Wwt, Wog, Wog, Wog, Wot, Xv, Xv
Xv, Wot, Wog, Wog, Wog, Wog, Wog, Wog, Wog, Wog, Wog, Wog, Wog, Wog, Wog, Wog, Wog, Wog, Wog, Wog, Wog, Wog, Wog, Wog, Wog, Wog, Wog, Wog, Wog, Wog, Wog, Wog, Wog, Wog, Wog, Wog, Wog, Wot, Xv, Xv
Xv, Wot, Wog, Wog, Wog, Wog, Wog, Wog, Wog, Wog, Wog, Wog, Wog, Wog, Wog, Wog, Wog, Wog, Wog, Wog, Wog, Wog, Wog, Wog, Wog, Wog, Wog, Wog, Wog, Wog, Wog, Wog, Wog, Wog, Wog, Wog, Wog, Wot, Xv, Xv
Xv, Wot, Wot, Wog, Wog, Wog, Wog, Wog, Wog, Wog, Wog, Wog, Wog, Wog, Wog, Wog, Wog, Wog, Wog, Wog, Wog, Wog, Wog, Wog, Wog, Wog, Wog, Wog, Wog, Wog, Wog, Wog, Wog, Wog, Wog, Wog, Wot, Wot, Xv, Xv
Xv, Xv, Xv, Wot, Wot, Wog, Wog, Wog, Wog, Wog, Wog, Wog, Wog, Wog, Wog, Wog, Wog, Wog, Wog, Wog, Wog, Wog, Wog, Wog, Wog, Wog, Wog, Wog, Wog, Wog, Wog, Wog, Wog, Wog, Wot, Wot, Xv, Xv, Xv, Xv
Xv, Xv, Xv, Xv, Xv, Wot, Wot, Wog, Wog, Wog, Wog, Wog, Wog, Wog, Wog, Wog, Wog, Wog, Wog, Wog, Wog, Wog, Wog, Wog, Wog, Wog, Wog, Wog, Wog, Wog, Wog, Wog, Wot, Wot, Xv, Xv, Xv, Xv, Xv, Xv
Xv, Xv, Xv, Xv, Xv, Xv, Xv, Wot, Wot, Wot, Wot, Wot, Wot, Wot, Wot, Wot, Wot, Wot, Wot, Wot, Wot, Wot, Wot, Wot, Wot, Wot, Wot, Wot, Wot, Wot, Wot, Wot, Xv, Xv, Xv, Xv, Xv, Xv, Xv, Xv
Xv, Xv, Xv, Xv, Xv, Xv, Xv, Xv, Xv, Xv, Xv, Xv, Xv, Xv, Xv, Xv, Xv, Xv, Xv, Xv, Xv, Xv, Xv, Xv, Xv, Xv, Xv, Xv, Xv, Xv, Xv, Xv, Xv, Xv, Xv, Xv, Xv, Xv, Xv, Xv
"
	name=_"WL Bot"
	random_start_time=no
	[side]
		canrecruit=yes
		controller="human"
		fog=yes
		side=1
		team_name="south"
		user_team_name=_"teamname^South"
		[ai]
			villages_per_scout=8
		[/ai]
	[/side]
	[side]
		canrecruit=yes
		controller="human"
		fog=yes
		side=2
		team_name="north"
		user_team_name=_"teamname^North"
		[ai]
			villages_per_scout=8
		[/ai]
	[/side]
	[event]
		name="prestart"
		[endlevel]
			carryover_add=false
			carryover_percentage=0
			carryover_report=false
			linger_mode=false
			next_scenario="scenario"
			replay_save=no
			result="victory"
			save=no
		[/endlevel]
	[/event]
[/multiplayer]
`,
}
