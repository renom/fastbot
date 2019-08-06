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

package scenario

import (
	"io/ioutil"
	"path/filepath"
	"strings"

	"github.com/renom/fastbot/config"
)

type Scenario struct {
	name    string
	body    string
	path    string
	defines []string
}

func NewScenario(name string, body string, defines []string) Scenario {
	return Scenario{name: name, body: body, defines: defines}
}

func FromPath(path string, defines []string) Scenario {
	filename := filepath.Base(path)
	extension := filepath.Ext(filename)
	name := strings.ReplaceAll(filename[0:len(filename)-len(extension)], "_", " ")
	return Scenario{name: name, path: path, defines: defines}
}

func (s *Scenario) Name() string {
	return s.name
}

func (s *Scenario) Body() string {
	result, _ := ioutil.ReadFile(s.path)
	return string(result)
}

func (s *Scenario) Path() string {
	if s.path == "" && s.body != "" {
		path := config.TmpDir + "/" + strings.ReplaceAll(s.name, " ", "_") + ".cfg"
		ioutil.WriteFile(path, []byte(s.body), 0644)
		s.path = path
	}
	return s.path
}

func (s *Scenario) Defines() []string {
	return s.defines
}
