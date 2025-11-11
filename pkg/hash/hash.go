// Copyright (C) 2025 Alan Barbosa Lima.
//
// Almodon is licensed under the GNU General Public License
// version 3. You should have received a copy of the
// license, located in LICENSE, at the root of the source
// tree. If not, see <https://www.gnu.org/licenses/>.

// Package hash wraps functionality from the standard
// [golang.org/x/crypto/bcrypt] package.
package hash

import (
	"errors"

	"golang.org/x/crypto/bcrypt"
)

// Hash takes an arbitrarily long password and hashes it using the
// BCrypt algorithm.
//
// To compare a hash to its password, you MUST use the [Compare]
// function, hashing the password again and using == should yield a
// different hash, failing the equality test.
//
// It handles the passwords that are larger than 72 bytes, which
// BCrypt wouldn't accept, by XORing the content over itself in
// 72-byte chunks. This introduces chance for matching a passwords
// in a simplistic manner, because XOR is trivially reversible.
func Hash(password []byte) ([60]byte, error) {
	ingest := atMax72Bytes(password)

	digest, err := bcrypt.GenerateFromPassword(ingest, bcrypt.DefaultCost)
	if err != nil {
		return [60]byte{}, err
	}
	if len(digest) != 60 {
		return [60]byte{}, errors.New("hash: digest is not of length 60")
	}

	return [60]byte(digest), nil
}

// Compare compares a hash generated through [Hash] with a password,
// it returns true for a match, and false for not a match or an
// error.
func Compare(hash, password []byte) bool {
	err := bcrypt.CompareHashAndPassword(hash, atMax72Bytes(password))
	return err == nil
}

func atMax72Bytes(data []byte) []byte {
	const size = 72

	if len(data) <= size {
		return data
	}

	result := make([]byte, size)
	for i, d := range data {
		result[i%size] ^= d
	}

	return result
}
