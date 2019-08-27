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
	"bytes"
	"compress/gzip"
	"encoding/binary"
	"errors"
	"fmt"
	"io/ioutil"
	"math/rand"
	"net"
	"regexp"
	"strconv"
	"strings"
	"time"

	"github.com/renom/fastbot/picker"
	serverTypes "github.com/renom/fastbot/server/types"
	"github.com/renom/fastbot/types"
	e "github.com/renom/go-wesnoth/era"
	"github.com/renom/go-wesnoth/game"
	"github.com/renom/go-wesnoth/scenario"
	"github.com/renom/go-wml"
)

type Server struct {
	hostname      string
	port          uint16
	version       string
	username      string
	password      string
	era           e.Era
	title         string
	game          []byte
	scenarios     serverTypes.ScenarioList
	lastSkip      string // player name
	admins        types.StringList
	players       types.StringList
	observers     types.StringList
	timeout       time.Duration
	err           error
	conn          net.Conn
	disconnecting bool
	sides         serverTypes.SideList
	picking       bool
	pickingPlayer string
	// Timer-related config
	TimerEnabled  bool
	InitTime      int
	TurnBonus     int
	ReservoirTime int
	ActionBonus   int
}

var colors = types.StringList{"red", "green", "purple", "orange", "white", "teal"}

func NewServer(hostname string, port uint16, version string, username string,
	password string, era string, title string, scenarios []scenario.Scenario,
	admins types.StringList, players types.StringList, pickingPlayer string, timerEnabled bool,
	initTime int, turnBonus int, reservoirTime int, actionBonus int,
	timeout time.Duration) Server {
	s := Server{
		hostname:      hostname,
		port:          port,
		version:       version,
		username:      username,
		password:      password,
		era:           e.Parse(era),
		title:         title,
		admins:        admins,
		players:       players,
		pickingPlayer: pickingPlayer,
		timeout:       timeout,
		TimerEnabled:  timerEnabled,
		InitTime:      initTime,
		TurnBonus:     turnBonus,
		ReservoirTime: reservoirTime,
		ActionBonus:   actionBonus,
	}
	var scenarioList serverTypes.ScenarioList
	for _, v := range scenarios {
		scenarioList = append(scenarioList, serverTypes.Scenario{false, v})
	}
	s.scenarios = scenarioList
	s.sides = serverTypes.SideList{&serverTypes.Side{Side: 1, Color: "red"}, &serverTypes.Side{Side: 2, Color: "blue"}}
	var path string
	var defines []string
	if len(scenarios) > 1 {
		s.picking = true
		path = picker.Scenario().Path()
		defines = append(defines[:0:0], picker.Scenario().Defines()...)
	} else {
		s.picking = false
		path = s.scenarios[0].Scenario.Path()
		defines = append(defines[:0:0], s.scenarios[0].Scenario.Defines()...)
	}
	g := game.NewGame(s.title, scenario.FromPath(path, defines), s.era,
		s.TimerEnabled, s.InitTime, s.TurnBonus, s.ReservoirTime, s.ActionBonus,
		s.version)
	s.game = g.Bytes()
	return s
}

func (s *Server) Connect() error {
	// Set up a TCP connection
	if s.conn, s.err = net.Dial("tcp", s.hostname+":"+strconv.Itoa(int(s.port))); s.err != nil {
		return s.err
	}
	//s.conn.SetDeadline(time.Now().Add(s.timeout))
	// Init the connection to the server
	s.conn.Write([]byte{0, 0, 0, 0})
	var buffer []byte
	if buffer = s.read(4); s.err != nil {
		return s.err
	}
	fmt.Println(binary.BigEndian.Uint32(buffer))
	// Expects the server to ask for a version, otherwise return an error
	if data := s.receiveData(); bytes.Equal(data, wml.EmptyTag("version").Bytes()) {
		s.sendData((&wml.Tag{"version", wml.Data{"version": s.version}}).Bytes())
	} else {
		return errors.New("Expects the server to request a version, but it doesn't.")
	}
	// Expects the server to require the log in step, otherwise return an error
	{
		rawData := s.receiveData()
		data := wml.ParseData(string(rawData))
		switch {
		case bytes.Equal(rawData, wml.EmptyTag("mustlogin").Bytes()):
			s.sendData((&wml.Tag{"login", wml.Data{"selective_ping": "1", "username": s.username}}).Bytes())
		case data.Contains("redirect"):
			if redirect, ok := data["redirect"].(wml.Data); ok {
				host, okHost := redirect["host"].(string)
				port, okPort := redirect["port"].(string)
				if okHost && okPort {
					portInt, err := strconv.Atoi(port)
					if err == nil {
						s.hostname = host
						s.port = uint16(portInt)
						return s.Connect()
					}
				}
			}
			fallthrough
		default:
			return errors.New("Expects the server to require a log in step, but it doesn't.")
		}
	}
	rawData := s.receiveData()
	data := wml.ParseData(string(rawData))
	switch {
	case data.Contains("error"):
		if errorTag, ok := data["error"].(wml.Data); ok {
			if code, ok := errorTag["error_code"].(string); ok {
				switch code {
				case "200":
					if errorTag["password_request"].(string) == "yes" && errorTag["phpbb_encryption"].(string) == "yes" {
						salt := errorTag["salt"].(string)
						s.sendData((&wml.Tag{"login", wml.Data{"username": s.username, "password": Sum(s.password, salt)}}).Bytes())
						goto nextCase
					}
				case "105":
					if message, ok := errorTag["message"].(string); ok {
						return errors.New(message)
					} else {
						return errors.New("The nickname is not registered. This server disallows unregistered nicknames.")
					}
				}
			}
		}
		break
	nextCase:
		fallthrough
	case bytes.Equal(rawData, wml.EmptyTag("join_lobby").Bytes()):
		return nil
	default:
		return errors.New("An unknown error occurred")
	}
	return nil
}

func (s *Server) HostGame() {
	s.sendData((&wml.Tag{"create_game", wml.Data{"name": s.title, "password": ""}}).Bytes())
	s.sendData(s.game)
	s.sides.Side(1).Color = "orange"
	s.ChangeSide(1, "insert", wml.Data{"color": "orange"})
	s.sides.Side(2).Color = "purple"
	s.ChangeSide(2, "insert", wml.Data{"color": "purple"})
}

func (s *Server) StartGame() {
	s.sendData(wml.EmptyTag("stop_updates").Bytes())
	r, _ := regexp.Compile(`(?Us)\[era\](.*)\[multiplayer_side\]`)
	rDomain, _ := regexp.Compile(`[\t ]*#textdomain ([0-9a-z_-]+)\n`)
	textdomains := rDomain.FindAllSubmatch(r.FindSubmatch(s.game)[1], -1)
	textdomain := string(textdomains[len(textdomains)-1][1])
	s.sides.Shuffle()
	rand.Seed(time.Now().UTC().UnixNano())
	index := rand.Int31n(6)
	faction1 := s.era.Factions[index]
	faction2 := append(append([]wml.Data{}, s.era.Factions[:index]...), s.era.Factions[index+1:]...)[rand.Int31n(5)]
	data := wml.Data{"scenario_diff": wml.Data{"change_child": wml.Data{
		"index": 0,
		"scenario": wml.Data{"change_child": wml.Multiple{
			wml.Data{"index": 0, "side": insertFaction(s.sides.Side(1), faction1, textdomain)},
			wml.Data{"index": 1, "side": insertFaction(s.sides.Side(2), faction2, textdomain)},
		}},
	},
	}}
	s.sendData(data.Bytes())
	s.sendData(wml.EmptyTag("start_game").Bytes())
}

func (s *Server) StoreNext() {
	s.sendData(wml.EmptyTag("update_game").Bytes())
	pickedScenario := *s.scenarios.PickedScenario()
	game := game.NewGame(s.title, pickedScenario, s.era,
		s.TimerEnabled, s.InitTime, s.TurnBonus, s.ReservoirTime, s.ActionBonus,
		s.version)
	game.NotNewGame = true
	game.Player1 = s.sides.Side(1).Player
	game.Player2 = s.sides.Side(2).Player
	next_scenario := "[store_next_scenario]\n" + game.String() + "[/store_next_scenario]\n"
	s.sendData([]byte(next_scenario))
	s.StartGame()
}

func (s *Server) MuteAll() {
	s.sendData(wml.EmptyTag("muteall").Bytes())
}

func (s *Server) LeaveGame() {
	s.sendData(wml.EmptyTag("leave_game").Bytes())
}

func (s *Server) Disconnect() {
	time.Sleep(time.Minute)
	s.disconnecting = true
}

func (s *Server) Listen() {
	for {
		if s.disconnecting == true {
			s.conn.Close()
			s.disconnecting = false
			break
		}
		data := wml.ParseData(string(s.receiveData()))
		if len(data) > 0 {
			fmt.Printf("Received: %q\n", data)
		}
	DataSwitch:
		switch {
		case data.Contains("name") && data.Contains("side") && s.sides.FreeSlots() > 0:
			name := data["name"].(string)
			side, _ := strconv.Atoi(data["side"].(string))
			if s.players.ContainsValue(name) && s.sides.HasSide(side) && !s.sides.HasPlayer(name) {
				s.ChangeSide(side, "insert", wml.Data{"current_player": name, "name": name, "player_id": name})
				s.sides.Side(side).Player = name
				if s.picking == true && s.sides.FreeSlots() == 0 {
					if s.scenarios.MustStart() == false {
						s.PickingMessage()
					} else {
						s.Message("The picked scenario is \"" +
							s.scenarios.PickedScenario().Name() +
							"\". Type \"ready\" to start.")
					}
				}
			} else if name != s.username && !s.observers.ContainsValue(name) {
				s.observers = append(s.observers, name)
			}
		case data.Contains("side_drop"):
			side_drop := data["side_drop"].(wml.Data)
			if side_drop.Contains("side_num") {
				side, _ := strconv.Atoi(side_drop["side_num"].(string))
				s.ChangeSide(side, "delete", wml.Data{"current_player": "x", "name": "x", "player_id": "x"})
				sideStruct := s.sides.Side(side)
				sideStruct.Ready = false
				sideStruct.Player = ""
			}
		case data.Contains("observer"):
			observer := data["observer"].(wml.Data)
			if observer.Contains("name") {
				name := observer["name"].(string)
				if name != s.username && !s.observers.ContainsValue(name) {
					s.observers = append(s.observers, name)
				}
			}
		case data.Contains("observer_quit"):
			observer_quit := data["observer_quit"].(wml.Data)
			if observer_quit.Contains("name") {
				name := observer_quit["name"].(string)
				if name != s.username && s.observers.ContainsValue(name) {
					s.observers.DeleteValue(name)
				}
			}
		case data.Contains("leave_game"):
			for _, v := range s.sides {
				v.Player = ""
				v.Ready = false
			}
			s.HostGame()
			for _, v := range s.sides {
				s.ChangeSide(v.Side, "insert", wml.Data{"color": v.Color})
			}
		case data.Contains("whisper"):
			whisper := data["whisper"].(wml.Data)
			if whisper.Contains("message") && whisper.Contains("receiver") && whisper.Contains("sender") {
				text := whisper["message"].(string)
				receiver := whisper["receiver"].(string)
				sender := whisper["sender"].(string)
				if receiver == s.username && s.admins.ContainsValue(sender) {
					// Commands that require an extra preparation
					{
						command := strings.SplitN(strings.TrimSpace(text), " ", 2)
						if command[0] == "players" && len(command) == 2 {
							players := types.StringList(strings.Split(command[1], ","))
							for i, v := range players {
								players[i] = strings.TrimSpace(v)
							}
							if players.Match("^[0-9A-Za-z_-]+$") == false {
								break DataSwitch
							}
							s.players = players.Unique()
							break DataSwitch
						}
					}
					// Ordinary commands
					command := strings.Fields(text)
					switch {
					case command[0] == "color" && len(command) == 3:
						side := types.ParseInt(command[1], -1)
						color := command[2]
						if s.sides.HasSide(side) && colors.ContainsValue(color) && !s.sides.HasColor(color) {
							s.ChangeSide(side, "insert", wml.Data{"color": color})
							s.sides.Side(side).Color = color
						}
					case command[0] == "slot" && len(command) == 3:
						side := types.ParseInt(command[1], -1)
						observer := command[2]
						if s.sides.HasSide(side) && s.observers.ContainsValue(observer) && !s.sides.HasPlayer(observer) {
							s.ChangeSide(side, "insert", wml.Data{"is_host": false, "is_local": false, "current_player": observer, "name": observer, "player_id": observer})
							player := s.sides.Side(side).Player
							s.observers.DeleteValue(observer)
							s.sides.Side(side).Player = observer
							s.observers = append(s.observers, player)
							if s.picking == true && s.sides.FreeSlots() == 0 {
								if s.scenarios.MustStart() == false {
									s.PickingMessage()
								} else {
									s.Message("The picked scenario is \"" +
										s.scenarios.PickedScenario().Name() +
										"\". Type \"ready\" to start.")
								}
							}

						}
						// Need to have an observers list
					case command[0] == "drop" && len(command) == 2:
						side := types.ParseInt(command[1], -1)
						// If the value is -1, so it's actually a nickname
						if side == -1 {
							side = s.sides.Find(command[1]).Side
						}
						if s.sides.HasSide(side) && s.sides.Side(side).Player != "" {
							player := s.sides.Side(side).Player
							s.ChangeSide(side, "delete", wml.Data{"current_player": "x", "name": "x", "player_id": "x"})
							s.sides.Side(side).Player = ""
							s.observers = append(s.observers, player)
						}
					case command[0] == "players" && len(command) == 1:
						s.Whisper(sender, "Player list: "+strings.Join(s.players, ", "))
					case command[0] == "admins" && len(command) == 1:
						s.Whisper(sender, "Admin list: "+strings.Join(s.admins, ", "))
					case command[0] == "stop" && len(command) == 1:
						s.Whisper(sender, "Logging out...")
						s.LeaveGame()
						s.Disconnect()
					case command[0] == "help" && len(command) == 1:
						s.Whisper(sender, "Command list:\n"+
							"color 1 {red,green,purple,orange,white,teal} - set up a side's color\n"+
							"slot 1 username - change the side's controller\n"+
							"drop 1 - remove player from a slot\n"+
							"drop username - the same but by username\n"+
							"players - display players allowed to join\n"+
							"players {username_1,username_2,...,username_N} - set up players allowed to join\n"+
							"admins - display admins list\n"+
							"stop - stop the bot instance\n"+
							"help - request command reference")
					}
				}
			}
		case data.Contains("message"):
			message := data["message"].(wml.Data)
			if message.Contains("message") && message.Contains("room") && message.Contains("sender") {
				text := message["message"].(string)
				room := message["room"].(string)
				sender := message["sender"].(string)
				if room == "this game" && s.sides.HasPlayer(sender) {
					side := s.sides.Find(sender)
					command := strings.Fields(text)
					switch {
					case s.picking == true && s.pickingPlayer != "" && command[0] == "pick" && len(command) == 2:
						if sender == s.pickingPlayer {
							scenarioNumber := types.ParseInt(command[1], -1)
							if scenarioNumber > 0 && scenarioNumber <= len(s.scenarios) {
								for i := range s.scenarios {
									if i != scenarioNumber-1 {
										s.scenarios[i].Skip = true
									}
								}
								s.Message("The scenario has been chosen. \"" + s.scenarios.PickedScenario().Name() + "\" is to be played.")
								if s.sides.MustStart() {
									s.StartGame()
									s.StoreNext()
									s.MuteAll()
									s.LeaveGame()
									s.Disconnect()
								}
							} else {
								s.Message("Incorrect scenario number, please try again.")
							}
						} else {
							s.Message("You aren't allowed to choose a scenario.")
						}
					case s.picking == true && s.pickingPlayer == "" && command[0] == "skip" && len(command) == 2:
						if s.sides.FreeSlots() == 0 {
							if s.lastSkip != sender {
								scenarioNumber := types.ParseInt(command[1], -1)
								if scenarioNumber > 0 && scenarioNumber <= len(s.scenarios) {
									if s.scenarios[scenarioNumber-1].Skip == false {
										s.scenarios[scenarioNumber-1].Skip = true
										s.lastSkip = sender
										if s.scenarios.PickedIndex() != -1 {
											s.Message("The scenario has been chosen. \"" + s.scenarios.PickedScenario().Name() + "\" is to be played.")
										} else {
											s.PickingMessage()
										}
										if s.sides.MustStart() && s.scenarios.MustStart() {
											s.StartGame()
											s.StoreNext()
											s.MuteAll()
											s.LeaveGame()
											s.Disconnect()
										}
									} else {
										s.Message("The chosen scenario is already skipped, choose another one.")
									}
								} else {
									s.Message("Incorrect scenario number, please try again.")
								}
							} else {
								s.Message("You've just skipped, please wait for your turn.")
							}
						} else {
							s.Message("Please wait until all the players are joined.")
						}
					case command[0] == "color" && len(command) == 2 &&
						colors.ContainsValue(command[1]) && !s.sides.HasColor(command[1]):
						s.ChangeSide(side.Side, "insert", wml.Data{"color": command[1]})
						side.Color = command[1]
					case command[0] == "ready" && len(command) == 1:
						if side.Ready == false {
							side.Ready = true
							s.Message("Player " + side.Player + " is ready to start the game")
							if s.picking == false && s.sides.MustStart() {
								s.StartGame()
								s.MuteAll()
								s.LeaveGame()
								s.Disconnect()
							}
							if s.picking == true && s.sides.MustStart() && s.scenarios.MustStart() {
								s.StartGame()
								s.StoreNext()
								s.MuteAll()
								s.LeaveGame()
								s.Disconnect()
							}
						}
					case len(command) == 2 && command[0] == "not" && command[1] == "ready":
						if side.Ready == true {
							side.Ready = false
							s.Message("Player " + side.Player + " isn't ready to start the game")
						}
					case command[0] == "help" && len(command) == 1:
						text := "Command list:\n"
						if s.picking == true && s.pickingPlayer != "" && s.scenarios.MustStart() == false {
							text += "pick 1 - pick a scenario when picking\n"
						}
						if s.picking == true && s.pickingPlayer == "" && s.scenarios.MustStart() == false {
							text += "skip 1 - skip a scenario when picking\n"
						}
						text += "color {red,green,purple,orange,white,teal} - set up a color\n" +
							"ready - ready to start the game\n" +
							"not ready - decline readiness\n" +
							"help - request command reference"
						s.Message(text)
					}
				}
			}
		}
		//time.Sleep(time.Millisecond * 100)
	}
}

func (s *Server) ChangeSide(side int, command string, data wml.Data) {
	s.sendData((&wml.Data{"scenario_diff": wml.Data{
		"change_child": wml.Data{"index": 0, "scenario": wml.Data{
			"change_child": wml.Data{"index": side - 1, "side": wml.Data{command: data}},
		}},
	}}).Bytes())
}

func (s *Server) Message(text string) {
	for _, v := range SplitMessage(wml.EscapeString(text)) {
		s.sendData((&wml.Data{"message": wml.Data{"message": v, "room": "this game", "sender": s.username}}).Bytes())
	}
}

func (s *Server) Whisper(receiver string, text string) {
	for _, v := range SplitMessage(wml.EscapeString(text)) {
		s.sendData((&wml.Data{"whisper": wml.Data{"sender": s.username, "receiver": receiver, "message": v}}).Bytes())
	}
}

func (s *Server) PickingMessage() {
	var scenarioList []string
	for i, v := range s.scenarios {
		var name string
		if v.Skip == false {
			name = v.Scenario.Name()
		} else {
			name = "_"
		}
		scenarioList = append(scenarioList, "["+strconv.Itoa(i+1)+"] "+name)
	}
	scenarioListString := ""
	rowSize := 4
	for i := 0; i < len(scenarioList); i += rowSize {
		var j int
		if i+rowSize < len(scenarioList) {
			j = i + rowSize
		} else {
			j = len(scenarioList)
		}
		scenarioListString += strings.Join(scenarioList[i:j], "     ") + "\n"
	}
	if s.pickingPlayer == "" {
		var allowedToChoose string
		if s.lastSkip == "" {
			allowedToChoose = "any side"
		} else {
			for _, v := range s.sides {
				if v.Player != s.lastSkip {
					allowedToChoose = v.Player
					break
				}
			}
		}
		s.Message("Please choose the most unwanted scenario:\n" +
			scenarioListString +
			"\nAllowed to choose: " + allowedToChoose +
			"\nCommand: \"skip 1\" (change \"1\" with an actual scenario number)")
	} else {
		allowedToChoose := s.pickingPlayer
		s.Message("Please choose the scenario to play:\n" +
			scenarioListString +
			"\nAllowed to choose: " + allowedToChoose +
			"\nCommand: \"pick 1\" (change \"1\" with an actual scenario number)")
	}
}

func (s *Server) Error() error {
	return s.err
}

func (s *Server) receiveData() []byte {
	buffer := s.read(4)
	if len(buffer) < 4 {
		return nil
	}
	size := int(binary.BigEndian.Uint32(buffer))
	reader, _ := gzip.NewReader(bytes.NewBuffer(s.read(size)))
	var result []byte
	if result, s.err = ioutil.ReadAll(reader); s.err != nil {
		return nil
	}
	if s.err = reader.Close(); s.err != nil {
		return nil
	}
	return result
}

func (s *Server) sendData(data []byte) {
	var b bytes.Buffer

	gz := gzip.NewWriter(&b)
	gz.Write([]byte(data))
	gz.Close()

	var length int = len(b.Bytes())
	s.conn.Write([]byte{0, 0, byte(length / 256), byte(length % 256)})
	s.conn.Write(b.Bytes())

}

func (s *Server) read(n int) []byte {
	result := []byte{}
	count := 0
	for count < n {
		buffer := make([]byte, n-count)
		var num int
		num, s.err = s.conn.Read(buffer)
		if s.err != nil {
			return nil
		}
		count += num
		result = append(result, buffer[:num]...)
	}
	return result
}
