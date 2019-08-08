# Fastbot

Fastbot is a Golang app which implements a bot for Battle for Wesnoth that intended to automatically host multiplayer games, particularly to host Fast tournament games.

## Requirements

Installed Battle for Wesnoth 1.14+.

## Installation

```bash
git clone https://github.com/renom/fastbot
cd fastbot
go install
```

## Usage

```
fastbot -host example.com -admins admin1,admin2,admin3 -accounts nickname1[:password1],nickname2[:password2],...,nicknameN[:passwordN] -title "{Player1} vs {Player2}" -timer 270,270,270,0 -baseDir /usr/share/wesnoth/data/multiplayer/scenarios player1,player2:2p_Clearing_Gushes.cfg,2p_The_Freelands.cfg,2p_Elensefar_Courtyard.cfg,2p_Sablestone_Delta.cfg,2p_Den_of_Onis.cfg player3,player4:2p_Caves_of_the_Basilisk.cfg
```

## License

[GPLv3](https://www.gnu.org/licenses/gpl-3.0.txt)
