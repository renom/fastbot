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
	"fmt"
	"strconv"
	_ "unsafe"

	_ "golang.org/x/crypto/bcrypt"
)

var itoa []byte = []byte("./0123456789ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz")

func Sum(password string, salt string) (string, error) {
	p := []byte(password)
	s := []byte(salt)
	if isMd5(salt) {
		salt1 := append([]byte{}, s[4:4+8]...)
		salt2 := append([]byte{}, s[12:12+8]...)
		count := 1 << uint(bytes.IndexByte(itoa, salt[3]))
		sum := md5Hash(p, salt1, count)
		sum = md5Hash(sum, salt2, 1<<10)
		return string(sum), nil
	} else if isBcrypt(salt) {
		salt1 := append([]byte{}, s[:29]...)
		salt2 := append([]byte{}, s[29:]...)
		hash, err := bcryptHash(p, salt1)
		if err != nil {
			return "", err
		}
		sum := md5Hash(hash, salt2, 1<<10)
		return string(sum), nil
	}
	return "", fmt.Errorf("Unknown encryption algorithm")
}

func md5Hash(password []byte, salt []byte, count int) []byte {
	hash := md5.Sum(append(salt, password...))

	for {
		hash = md5.Sum(append(hash[:], password...))

		count--
		if count == 0 {
			break
		}
	}

	return md5Encode(hash[:], 16)
}

func md5Encode(text []byte, count int) []byte {
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

func bcryptHash(password []byte, salt []byte) ([]byte, error) {
	cost, err := strconv.Atoi(string(salt[4:6]))
	if err != nil {
		return nil, err
	}
	err = crypto_checkCost(cost)
	if err != nil {
		return nil, err
	}

	hash, err := crypto_bcrypt(password, cost, salt[7:])
	if err != nil {
		return nil, err
	}
	return append(salt[:], hash[:]...), nil
}

func isMd5(salt string) bool {
	return salt[0:3] == "$H$"
}

func isBcrypt(salt string) bool {
	return salt[0:4] == "$2a$" ||
		salt[0:4] == "$2b$" ||
		salt[0:4] == "$2x$" ||
		salt[0:4] == "$2y$"
}

//go:linkname crypto_bcrypt golang.org/x/crypto/bcrypt.bcrypt
func crypto_bcrypt(password []byte, cost int, salt []byte) ([]byte, error)

//go:linkname crypto_checkCost golang.org/x/crypto/bcrypt.checkCost
func crypto_checkCost(cost int) error
