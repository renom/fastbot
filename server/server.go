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

	"github.com/renom/fastbot/types"
	"github.com/renom/fastbot/wml"
)

type Server struct {
	hostname  string
	port      uint16
	version   string
	username  string
	password  string
	title     string
	game      []byte
	admins    types.StringList
	players   types.StringList
	observers types.StringList
	factions  []wml.Data
	timeout   time.Duration
	err       error
	conn      net.Conn
	sides     SideList
}

var colors = types.StringList{"red", "blue", "green", "purple", "black", "brown", "orange", "white", "teal"}

func NewServer(hostname string, port uint16, version string, username string,
	password string, title string, game []byte, admins types.StringList, players types.StringList, timeout time.Duration) Server {
	s := Server{
		hostname: hostname,
		port:     port,
		version:  version,
		username: username,
		password: password,
		title:    title,
		game:     game,
		admins:   admins,
		players:  players,
		timeout:  timeout,
	}
	s.sides = SideList{&Side{Side: 1, Color: "red"}, &Side{Side: 2, Color: "blue"}}

	r, _ := regexp.Compile(`(?U)\[multiplayer_side\](.*\n)*[\t ]*\[/multiplayer_side\]`)
	factions := r.FindAll(s.game, -1)
	rData, _ := regexp.Compile(`(?U)[\t ]*[0-9a-z_]+[\t ]*=[\t ]*_?"[^"](.|\n)*` + `([^"]"[\t\n ]*\+[\t\n ]*_?"[^"])+` +
		`(.|\n)*[^"]"\n`)
	s.factions = []wml.Data{}
	for i, v := range factions {
		factions[i] = rData.ReplaceAll(v, []byte(""))
		factionData := wml.ParseTag(string(factions[i])).Data
		if !factionData.Contains("random_faction") || factionData["random_faction"] == false {
			s.factions = append(s.factions, factionData)
		}
	}
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
	if data := s.receiveData(); bytes.Equal(data, wml.EmptyTag("mustlogin").Bytes()) {
		s.sendData((&wml.Tag{"login", wml.Data{"selective_ping": "1", "username": s.username}}).Bytes())
	} else {
		return errors.New("Expects the server to require a log in step, but it doesn't.")
	}
	return nil
}

func (s *Server) HostGame() {
	s.sendData((&wml.Tag{"create_game", wml.Data{"name": s.title, "password": ""}}).Bytes())
	s.sendData(s.game)
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
	faction1 := s.factions[index]
	faction2 := append(s.factions[:index], s.factions[index+1:]...)[rand.Int31n(5)]
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
	s.sendData(wml.EmptyTag("leave_game").Bytes())
}

func (s *Server) Listen() {
	for {
		data := wml.ParseData(string(s.receiveData()))
		fmt.Printf("Received: %q\n", data)
	DataSwitch:
		switch {
		case data.Contains("name") && data.Contains("side") && s.sides.FreeSlots() > 0:
			name := data["name"].(string)
			side, _ := strconv.Atoi(data["side"].(string))
			if s.players.ContainsValue(name) && s.sides.HasSide(side) && !s.sides.HasPlayer(name) {
				s.ChangeSide(side, "insert", wml.Data{"current_player": name, "name": name, "player_id": name})
				s.sides.Side(side).Player = name
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
						if s.sides.HasSide(side) && s.observers.ContainsValue(observer) {
							s.ChangeSide(side, "insert", wml.Data{"is_host": false, "is_local": false, "current_player": observer, "name": observer, "player_id": observer})
							s.sides.Side(side).Player = observer
						}
						// Need to have an observers list
					case command[0] == "drop" && len(command) == 2:
						side := types.ParseInt(command[1], -1)
						// If the value is -1, so it's actually a nickname
						if side == -1 {
							side = s.sides.Find(command[1]).Side
						}
						if s.sides.HasSide(side) && s.sides.Side(side).Player != "" {
							s.ChangeSide(side, "delete", wml.Data{"current_player": "x", "name": "x", "player_id": "x"})
						}
					case command[0] == "players" && len(command) == 1:
						s.Whisper(sender, "Player list: "+strings.Join(s.players, ", "))
					case command[0] == "help" && len(command) == 1:
						s.Whisper(sender, "Command list:\n"+
							"color 1 {red,blue,green,purple,black,brown,orange,white,teal} - set up a side's color\n"+
							"slot 1 username - change the side's controller\n"+
							"drop 1 - remove player from a slot\n"+
							"drop username - the same but by username\n"+
							"players - display players allowed to join\n"+
							"players {username_1,username_2,...,username_N} - set up players allowed to join\n"+
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
					case command[0] == "color" && len(command) == 2 &&
						colors.ContainsValue(command[1]) && !s.sides.HasColor(command[1]):
						s.ChangeSide(side.Side, "insert", wml.Data{"color": command[1]})
						side.Color = command[1]
					case command[0] == "ready" && len(command) == 1:
						if side.Ready == false {
							side.Ready = true
							s.Message("Player " + side.Player + " is ready to start the game")
							if s.sides.MustStart() {
								s.StartGame()
							}
						}
					case len(command) == 2 && command[0] == "not" && command[1] == "ready":
						if side.Ready == true {
							side.Ready = false
							s.Message("Player " + side.Player + " isn't ready to start the game")
						}
					case command[0] == "help" && len(command) == 1:
						s.Message("Command list:\n" +
							"color {red,blue,green,purple,black,brown,orange,white,teal} - set up a color\n" +
							"ready - ready to start the game\n" +
							"not ready - decline readiness\n" +
							"help - request command reference")
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
	buffer := make([]byte, n)
	_, s.err = s.conn.Read(buffer)
	if s.err != nil {
		return nil
	}
	return buffer
}
