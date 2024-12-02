package cpg

import (
	"cpg/pkg/crypto"
	"encoding/json"
	"github.com/itsabgr/ge"
	"math/big"
	"sync"
	"time"
)

//go:generate go run golang.org/x/tools/cmd/stringer@latest -type=InvoiceStatus

type InvoiceStatus int

const (
	InvoiceStatusPending  InvoiceStatus = 1
	InvoiceStatusFilled   InvoiceStatus = 2
	InvoiceStatusExpired  InvoiceStatus = 3
	InvoiceStatusCanceled InvoiceStatus = 4
)

type Invoice struct {
	_             sync.Mutex
	ID            string
	MinAmount     big.Int
	Recipient     string
	Beneficiary   string
	Asset         string
	Metadata      string
	CreateAt      time.Time
	Deadline      time.Time
	FillAt        *time.Time
	CancelAt      *time.Time
	WalletAddress string
	EncryptedSalt []byte
	saltKeyring   *crypto.KeyRing
}

func (inv *Invoice) DecryptSalt() []byte {
	if inv.EncryptedSalt == nil {
		return nil
	}
	key, _ := inv.saltKeyring.Decrypt(inv.EncryptedSalt)
	return key
}

func (inv *Invoice) Destination() string {
	switch inv.Status() {
	case InvoiceStatusExpired, InvoiceStatusCanceled:
		return inv.Beneficiary
	case InvoiceStatusFilled, InvoiceStatusPending:
		return inv.Recipient
	default:
		panic(ge.UNREACHABLE)
	}
}

func (inv *Invoice) Status() InvoiceStatus {
	if inv.FillAt == nil {
		if inv.CancelAt == nil {
			if inv.Deadline.After(time.Now()) {
				return InvoiceStatusPending
			} else {
				return InvoiceStatusExpired
			}
		} else {
			return InvoiceStatusCanceled
		}
	} else {
		return InvoiceStatusFilled
	}
}

func (inv *Invoice) pack() []byte {
	return ge.Must(json.Marshal(inv))
}

func unpackInvoice(data []byte) (*Invoice, error) {
	inv := &Invoice{}
	return inv, json.Unmarshal(data, inv)
}

func (inv *Invoice) Encrypt(keyring *crypto.KeyRing) []byte {
	return keyring.Encrypt(inv.pack(), (*[24]byte)(crypto.ReadN(nil, 24)))
}

func DecryptInvoice(keyring *crypto.KeyRing, encrypted []byte) (*Invoice, error) {
	invData, _ := keyring.Decrypt(encrypted)
	return unpackInvoice(invData)
}
