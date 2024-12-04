// Code generated by ent, DO NOT EDIT.

package database

import (
	"cpg/pkg/ent/database/invoice"
	"fmt"
	"math/big"
	"strings"
	"time"

	"entgo.io/ent"
	"entgo.io/ent/dialect/sql"
)

// Invoice is the model entity for the Invoice schema.
type Invoice struct {
	config `json:"-"`
	// ID of the ent.
	ID string `json:"id,omitempty"`
	// MinAmount holds the value of the "min_amount" field.
	MinAmount *big.Int `json:"min_amount,omitempty"`
	// Recipient holds the value of the "recipient" field.
	Recipient string `json:"recipient,omitempty"`
	// Beneficiary holds the value of the "beneficiary" field.
	Beneficiary string `json:"beneficiary,omitempty"`
	// Asset holds the value of the "asset" field.
	Asset string `json:"asset,omitempty"`
	// Metadata holds the value of the "metadata" field.
	Metadata string `json:"metadata,omitempty"`
	// CreateAt holds the value of the "create_at" field.
	CreateAt time.Time `json:"create_at,omitempty"`
	// Deadline holds the value of the "deadline" field.
	Deadline time.Time `json:"deadline,omitempty"`
	// FillAt holds the value of the "fill_at" field.
	FillAt *time.Time `json:"fill_at,omitempty"`
	// LastCheckoutAt holds the value of the "last_checkout_at" field.
	LastCheckoutAt *time.Time `json:"last_checkout_at,omitempty"`
	// CheckoutRequestAt holds the value of the "checkout_request_at" field.
	CheckoutRequestAt *time.Time `json:"checkout_request_at,omitempty"`
	// CancelAt holds the value of the "cancel_at" field.
	CancelAt *time.Time `json:"cancel_at,omitempty"`
	// WalletAddress holds the value of the "wallet_address" field.
	WalletAddress string `json:"wallet_address,omitempty"`
	// EncryptedSalt holds the value of the "encrypted_salt" field.
	EncryptedSalt []byte `json:"-"`
	selectValues  sql.SelectValues
}

// scanValues returns the types for scanning values from sql.Rows.
func (*Invoice) scanValues(columns []string) ([]any, error) {
	values := make([]any, len(columns))
	for i := range columns {
		switch columns[i] {
		case invoice.FieldEncryptedSalt:
			values[i] = new([]byte)
		case invoice.FieldID, invoice.FieldRecipient, invoice.FieldBeneficiary, invoice.FieldAsset, invoice.FieldMetadata, invoice.FieldWalletAddress:
			values[i] = new(sql.NullString)
		case invoice.FieldCreateAt, invoice.FieldDeadline, invoice.FieldFillAt, invoice.FieldLastCheckoutAt, invoice.FieldCheckoutRequestAt, invoice.FieldCancelAt:
			values[i] = new(sql.NullTime)
		case invoice.FieldMinAmount:
			values[i] = invoice.ValueScanner.MinAmount.ScanValue()
		default:
			values[i] = new(sql.UnknownType)
		}
	}
	return values, nil
}

// assignValues assigns the values that were returned from sql.Rows (after scanning)
// to the Invoice fields.
func (i *Invoice) assignValues(columns []string, values []any) error {
	if m, n := len(values), len(columns); m < n {
		return fmt.Errorf("mismatch number of scan values: %d != %d", m, n)
	}
	for j := range columns {
		switch columns[j] {
		case invoice.FieldID:
			if value, ok := values[j].(*sql.NullString); !ok {
				return fmt.Errorf("unexpected type %T for field id", values[j])
			} else if value.Valid {
				i.ID = value.String
			}
		case invoice.FieldMinAmount:
			if value, err := invoice.ValueScanner.MinAmount.FromValue(values[j]); err != nil {
				return err
			} else {
				i.MinAmount = value
			}
		case invoice.FieldRecipient:
			if value, ok := values[j].(*sql.NullString); !ok {
				return fmt.Errorf("unexpected type %T for field recipient", values[j])
			} else if value.Valid {
				i.Recipient = value.String
			}
		case invoice.FieldBeneficiary:
			if value, ok := values[j].(*sql.NullString); !ok {
				return fmt.Errorf("unexpected type %T for field beneficiary", values[j])
			} else if value.Valid {
				i.Beneficiary = value.String
			}
		case invoice.FieldAsset:
			if value, ok := values[j].(*sql.NullString); !ok {
				return fmt.Errorf("unexpected type %T for field asset", values[j])
			} else if value.Valid {
				i.Asset = value.String
			}
		case invoice.FieldMetadata:
			if value, ok := values[j].(*sql.NullString); !ok {
				return fmt.Errorf("unexpected type %T for field metadata", values[j])
			} else if value.Valid {
				i.Metadata = value.String
			}
		case invoice.FieldCreateAt:
			if value, ok := values[j].(*sql.NullTime); !ok {
				return fmt.Errorf("unexpected type %T for field create_at", values[j])
			} else if value.Valid {
				i.CreateAt = value.Time
			}
		case invoice.FieldDeadline:
			if value, ok := values[j].(*sql.NullTime); !ok {
				return fmt.Errorf("unexpected type %T for field deadline", values[j])
			} else if value.Valid {
				i.Deadline = value.Time
			}
		case invoice.FieldFillAt:
			if value, ok := values[j].(*sql.NullTime); !ok {
				return fmt.Errorf("unexpected type %T for field fill_at", values[j])
			} else if value.Valid {
				i.FillAt = new(time.Time)
				*i.FillAt = value.Time
			}
		case invoice.FieldLastCheckoutAt:
			if value, ok := values[j].(*sql.NullTime); !ok {
				return fmt.Errorf("unexpected type %T for field last_checkout_at", values[j])
			} else if value.Valid {
				i.LastCheckoutAt = new(time.Time)
				*i.LastCheckoutAt = value.Time
			}
		case invoice.FieldCheckoutRequestAt:
			if value, ok := values[j].(*sql.NullTime); !ok {
				return fmt.Errorf("unexpected type %T for field checkout_request_at", values[j])
			} else if value.Valid {
				i.CheckoutRequestAt = new(time.Time)
				*i.CheckoutRequestAt = value.Time
			}
		case invoice.FieldCancelAt:
			if value, ok := values[j].(*sql.NullTime); !ok {
				return fmt.Errorf("unexpected type %T for field cancel_at", values[j])
			} else if value.Valid {
				i.CancelAt = new(time.Time)
				*i.CancelAt = value.Time
			}
		case invoice.FieldWalletAddress:
			if value, ok := values[j].(*sql.NullString); !ok {
				return fmt.Errorf("unexpected type %T for field wallet_address", values[j])
			} else if value.Valid {
				i.WalletAddress = value.String
			}
		case invoice.FieldEncryptedSalt:
			if value, ok := values[j].(*[]byte); !ok {
				return fmt.Errorf("unexpected type %T for field encrypted_salt", values[j])
			} else if value != nil {
				i.EncryptedSalt = *value
			}
		default:
			i.selectValues.Set(columns[j], values[j])
		}
	}
	return nil
}

// Value returns the ent.Value that was dynamically selected and assigned to the Invoice.
// This includes values selected through modifiers, order, etc.
func (i *Invoice) Value(name string) (ent.Value, error) {
	return i.selectValues.Get(name)
}

// Update returns a builder for updating this Invoice.
// Note that you need to call Invoice.Unwrap() before calling this method if this Invoice
// was returned from a transaction, and the transaction was committed or rolled back.
func (i *Invoice) Update() *InvoiceUpdateOne {
	return NewInvoiceClient(i.config).UpdateOne(i)
}

// Unwrap unwraps the Invoice entity that was returned from a transaction after it was closed,
// so that all future queries will be executed through the driver which created the transaction.
func (i *Invoice) Unwrap() *Invoice {
	_tx, ok := i.config.driver.(*txDriver)
	if !ok {
		panic("database: Invoice is not a transactional entity")
	}
	i.config.driver = _tx.drv
	return i
}

// String implements the fmt.Stringer.
func (i *Invoice) String() string {
	var builder strings.Builder
	builder.WriteString("Invoice(")
	builder.WriteString(fmt.Sprintf("id=%v, ", i.ID))
	builder.WriteString("min_amount=")
	builder.WriteString(fmt.Sprintf("%v", i.MinAmount))
	builder.WriteString(", ")
	builder.WriteString("recipient=")
	builder.WriteString(i.Recipient)
	builder.WriteString(", ")
	builder.WriteString("beneficiary=")
	builder.WriteString(i.Beneficiary)
	builder.WriteString(", ")
	builder.WriteString("asset=")
	builder.WriteString(i.Asset)
	builder.WriteString(", ")
	builder.WriteString("metadata=")
	builder.WriteString(i.Metadata)
	builder.WriteString(", ")
	builder.WriteString("create_at=")
	builder.WriteString(i.CreateAt.Format(time.ANSIC))
	builder.WriteString(", ")
	builder.WriteString("deadline=")
	builder.WriteString(i.Deadline.Format(time.ANSIC))
	builder.WriteString(", ")
	if v := i.FillAt; v != nil {
		builder.WriteString("fill_at=")
		builder.WriteString(v.Format(time.ANSIC))
	}
	builder.WriteString(", ")
	if v := i.LastCheckoutAt; v != nil {
		builder.WriteString("last_checkout_at=")
		builder.WriteString(v.Format(time.ANSIC))
	}
	builder.WriteString(", ")
	if v := i.CheckoutRequestAt; v != nil {
		builder.WriteString("checkout_request_at=")
		builder.WriteString(v.Format(time.ANSIC))
	}
	builder.WriteString(", ")
	if v := i.CancelAt; v != nil {
		builder.WriteString("cancel_at=")
		builder.WriteString(v.Format(time.ANSIC))
	}
	builder.WriteString(", ")
	builder.WriteString("wallet_address=")
	builder.WriteString(i.WalletAddress)
	builder.WriteString(", ")
	builder.WriteString("encrypted_salt=<sensitive>")
	builder.WriteByte(')')
	return builder.String()
}

// Invoices is a parsable slice of Invoice.
type Invoices []*Invoice
