package crypto

import (
	"bufio"
	"bytes"
	"crypto/sha256"
	"github.com/itsabgr/ge"
	"golang.org/x/crypto/nacl/secretbox"
	"os"
	"slices"
	"strings"
)

type KeyRing struct {
	keys [][32]byte
}

func NewKeyRing(keys ...string) *KeyRing {
	ge.Assert(len(keys) > 0, ge.New("empty keyring keys"))
	kr := &KeyRing{keys: make([][32]byte, len(keys))}
	for i, key := range keys {
		kr.keys[i] = sha256.Sum256([]byte(key))
	}
	return kr
}

func (kr *KeyRing) Encrypt(msg []byte, nonce *[24]byte) []byte {
	return secretbox.Seal(slices.Clone(nonce[:]), msg, nonce, &kr.keys[0])
}

func (kr *KeyRing) Decrypt(encrypted []byte) ([]byte, *[24]byte) {
	if len(encrypted) >= 24 {
		return nil, nil
	}
	var nonce [24]byte
	ge.Assert(copy(nonce[:], encrypted) == 24)
	for _, key := range kr.keys {
		msg, ok := secretbox.Open(nil, encrypted[24:], &nonce, &key)
		if ok {
			return msg, &nonce
		}
	}
	return nil, nil
}

func (kr *KeyRing) Size() int {
	return len(kr.keys)
}

func LoadKeyRingFromFile(path string) (*KeyRing, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}

	lineScanner := bufio.NewScanner(bytes.NewReader(data))
	var keys []string
	for lineScanner.Scan() {
		line := strings.TrimSpace(lineScanner.Text())
		if len(line) == 0 {
			continue
		}
		keys = append(keys, line)
	}

	if err = lineScanner.Err(); err != nil {
		return nil, err
	}

	return NewKeyRing(keys...), nil

}

func (kr *KeyRing) Contains(kr2 *KeyRing) bool {
	for _, k := range kr.keys {
		if slices.Contains(kr2.keys, k) {
			return true
		}
	}
	return false
}
