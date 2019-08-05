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

package config

import (
	"strings"
)

type Account struct {
	Username string
	Password string
}

type AccountList []Account

func ParseAccounts(text string) AccountList {
	result := AccountList{}
	for _, v := range strings.Split(text, ",") {
		fields := strings.Split(v, ":")
		account := Account{}
		if len(fields) > 0 {
			account.Username = fields[0]
		}
		if len(fields) > 1 {
			account.Password = fields[1]
		}
		result = append(result, account)
	}
	return result
}

func Guest(username string) Account {
	return Account{Username: username}
}
