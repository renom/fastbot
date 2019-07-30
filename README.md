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
fastbot -host example.com -username nickname -password somepass -title "New game" -players=Player1,Player2 /usr/share/wesnoth/data/multiplayer/scenarios/2p_Clearing_Gushes.cfg
```

## License

[GPLv3](https://www.gnu.org/licenses/gpl-3.0.txt)
