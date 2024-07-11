package model

type InvoiceStatus string

const (
	InvoiceStatusPend InvoiceStatus = "InvoiceStatusPend"
	InvoiceStatusDone InvoiceStatus = "InvoiceStatusDone"
)

func (sts InvoiceStatus) String() string {
	return string(sts)
}
