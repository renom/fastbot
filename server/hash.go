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
	"crypto/md5"
)

var itoa []byte = []byte("./0123456789ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz")

func Sum(password string, salt string) string {
	p := []byte(password)
	s := []byte(salt)
	salt1 := append([]byte{}, s[4:4+8]...)
	salt2 := append([]byte{}, s[12:12+8]...)
	count := 1 << uint(bytes.IndexByte(itoa, salt[3]))
	sum := hash(p, salt1, count)
	sum = hash(sum, salt2, 1<<10)
	return string(sum)
}

func hash(password []byte, salt []byte, count int) []byte {
	hash := md5.Sum(append(salt, password...))

	for {
		hash = md5.Sum(append(hash[:], password...))

		count--
		if count == 0 {
			break
		}
	}

	return encode(hash[:], 16)
}

func encode(text []byte, count int) []byte {
	result := []byte{}
	i := 0
	for {
		value := int(text[i])
		result = append(result, itoa[value&0x3f])
		i++
		if i < count {
			value |= int(text[i]) << 8
		}
		result = append(result, itoa[(value>>6)&0x3f])
		if i >= count {
			break
		}
		i++
		if i < count {
			value |= int(text[i]) << 16
		}
		result = append(result, itoa[(value>>12)&0x3f])
		if i >= count {
			break
		}
		result = append(result, itoa[(value>>18)&0x3f])
		i++
		if i >= count {
			break
		}
	}
	return result
}
