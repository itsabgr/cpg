package crypto

import (
	"crypto/rand"
	"github.com/itsabgr/ge"
	"io"
)

func ReadN(src io.Reader, n int) []byte {
	if src == nil {
		src = rand.Reader
	}
	b := make([]byte, n)
	ge.Assert(ge.Must(io.ReadFull(src, b)) == n)
	return b
}
