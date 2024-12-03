// Code generated by ent, DO NOT EDIT.

package database

import (
	"context"
	"cpg/pkg/ent/database/invoice"
	"errors"
	"fmt"
	"math/big"
	"time"

	"entgo.io/ent/dialect/sql/sqlgraph"
	"entgo.io/ent/schema/field"
)

// InvoiceCreate is the builder for creating a Invoice entity.
type InvoiceCreate struct {
	config
	mutation *InvoiceMutation
	hooks    []Hook
}

// SetMinAmount sets the "min_amount" field.
func (ic *InvoiceCreate) SetMinAmount(b *big.Int) *InvoiceCreate {
	ic.mutation.SetMinAmount(b)
	return ic
}

// SetRecipient sets the "recipient" field.
func (ic *InvoiceCreate) SetRecipient(s string) *InvoiceCreate {
	ic.mutation.SetRecipient(s)
	return ic
}

// SetBeneficiary sets the "beneficiary" field.
func (ic *InvoiceCreate) SetBeneficiary(s string) *InvoiceCreate {
	ic.mutation.SetBeneficiary(s)
	return ic
}

// SetAsset sets the "asset" field.
func (ic *InvoiceCreate) SetAsset(s string) *InvoiceCreate {
	ic.mutation.SetAsset(s)
	return ic
}

// SetMetadata sets the "metadata" field.
func (ic *InvoiceCreate) SetMetadata(s string) *InvoiceCreate {
	ic.mutation.SetMetadata(s)
	return ic
}

// SetCreateAt sets the "create_at" field.
func (ic *InvoiceCreate) SetCreateAt(t time.Time) *InvoiceCreate {
	ic.mutation.SetCreateAt(t)
	return ic
}

// SetNillableCreateAt sets the "create_at" field if the given value is not nil.
func (ic *InvoiceCreate) SetNillableCreateAt(t *time.Time) *InvoiceCreate {
	if t != nil {
		ic.SetCreateAt(*t)
	}
	return ic
}

// SetDeadline sets the "deadline" field.
func (ic *InvoiceCreate) SetDeadline(t time.Time) *InvoiceCreate {
	ic.mutation.SetDeadline(t)
	return ic
}

// SetFillAt sets the "fill_at" field.
func (ic *InvoiceCreate) SetFillAt(t time.Time) *InvoiceCreate {
	ic.mutation.SetFillAt(t)
	return ic
}

// SetNillableFillAt sets the "fill_at" field if the given value is not nil.
func (ic *InvoiceCreate) SetNillableFillAt(t *time.Time) *InvoiceCreate {
	if t != nil {
		ic.SetFillAt(*t)
	}
	return ic
}

// SetLastCheckoutAt sets the "last_checkout_at" field.
func (ic *InvoiceCreate) SetLastCheckoutAt(t time.Time) *InvoiceCreate {
	ic.mutation.SetLastCheckoutAt(t)
	return ic
}

// SetNillableLastCheckoutAt sets the "last_checkout_at" field if the given value is not nil.
func (ic *InvoiceCreate) SetNillableLastCheckoutAt(t *time.Time) *InvoiceCreate {
	if t != nil {
		ic.SetLastCheckoutAt(*t)
	}
	return ic
}

// SetCancelAt sets the "cancel_at" field.
func (ic *InvoiceCreate) SetCancelAt(t time.Time) *InvoiceCreate {
	ic.mutation.SetCancelAt(t)
	return ic
}

// SetNillableCancelAt sets the "cancel_at" field if the given value is not nil.
func (ic *InvoiceCreate) SetNillableCancelAt(t *time.Time) *InvoiceCreate {
	if t != nil {
		ic.SetCancelAt(*t)
	}
	return ic
}

// SetWalletAddress sets the "wallet_address" field.
func (ic *InvoiceCreate) SetWalletAddress(s string) *InvoiceCreate {
	ic.mutation.SetWalletAddress(s)
	return ic
}

// SetEncryptedSalt sets the "encrypted_salt" field.
func (ic *InvoiceCreate) SetEncryptedSalt(b []byte) *InvoiceCreate {
	ic.mutation.SetEncryptedSalt(b)
	return ic
}

// SetID sets the "id" field.
func (ic *InvoiceCreate) SetID(s string) *InvoiceCreate {
	ic.mutation.SetID(s)
	return ic
}

// Mutation returns the InvoiceMutation object of the builder.
func (ic *InvoiceCreate) Mutation() *InvoiceMutation {
	return ic.mutation
}

// Save creates the Invoice in the database.
func (ic *InvoiceCreate) Save(ctx context.Context) (*Invoice, error) {
	ic.defaults()
	return withHooks(ctx, ic.sqlSave, ic.mutation, ic.hooks)
}

// SaveX calls Save and panics if Save returns an error.
func (ic *InvoiceCreate) SaveX(ctx context.Context) *Invoice {
	v, err := ic.Save(ctx)
	if err != nil {
		panic(err)
	}
	return v
}

// Exec executes the query.
func (ic *InvoiceCreate) Exec(ctx context.Context) error {
	_, err := ic.Save(ctx)
	return err
}

// ExecX is like Exec, but panics if an error occurs.
func (ic *InvoiceCreate) ExecX(ctx context.Context) {
	if err := ic.Exec(ctx); err != nil {
		panic(err)
	}
}

// defaults sets the default values of the builder before save.
func (ic *InvoiceCreate) defaults() {
	if _, ok := ic.mutation.CreateAt(); !ok {
		v := invoice.DefaultCreateAt()
		ic.mutation.SetCreateAt(v)
	}
}

// check runs all checks and user-defined validators on the builder.
func (ic *InvoiceCreate) check() error {
	if _, ok := ic.mutation.MinAmount(); !ok {
		return &ValidationError{Name: "min_amount", err: errors.New(`database: missing required field "Invoice.min_amount"`)}
	}
	if v, ok := ic.mutation.MinAmount(); ok {
		if err := invoice.MinAmountValidator(v.String()); err != nil {
			return &ValidationError{Name: "min_amount", err: fmt.Errorf(`database: validator failed for field "Invoice.min_amount": %w`, err)}
		}
	}
	if _, ok := ic.mutation.Recipient(); !ok {
		return &ValidationError{Name: "recipient", err: errors.New(`database: missing required field "Invoice.recipient"`)}
	}
	if v, ok := ic.mutation.Recipient(); ok {
		if err := invoice.RecipientValidator(v); err != nil {
			return &ValidationError{Name: "recipient", err: fmt.Errorf(`database: validator failed for field "Invoice.recipient": %w`, err)}
		}
	}
	if _, ok := ic.mutation.Beneficiary(); !ok {
		return &ValidationError{Name: "beneficiary", err: errors.New(`database: missing required field "Invoice.beneficiary"`)}
	}
	if v, ok := ic.mutation.Beneficiary(); ok {
		if err := invoice.BeneficiaryValidator(v); err != nil {
			return &ValidationError{Name: "beneficiary", err: fmt.Errorf(`database: validator failed for field "Invoice.beneficiary": %w`, err)}
		}
	}
	if _, ok := ic.mutation.Asset(); !ok {
		return &ValidationError{Name: "asset", err: errors.New(`database: missing required field "Invoice.asset"`)}
	}
	if v, ok := ic.mutation.Asset(); ok {
		if err := invoice.AssetValidator(v); err != nil {
			return &ValidationError{Name: "asset", err: fmt.Errorf(`database: validator failed for field "Invoice.asset": %w`, err)}
		}
	}
	if _, ok := ic.mutation.Metadata(); !ok {
		return &ValidationError{Name: "metadata", err: errors.New(`database: missing required field "Invoice.metadata"`)}
	}
	if v, ok := ic.mutation.Metadata(); ok {
		if err := invoice.MetadataValidator(v); err != nil {
			return &ValidationError{Name: "metadata", err: fmt.Errorf(`database: validator failed for field "Invoice.metadata": %w`, err)}
		}
	}
	if _, ok := ic.mutation.CreateAt(); !ok {
		return &ValidationError{Name: "create_at", err: errors.New(`database: missing required field "Invoice.create_at"`)}
	}
	if _, ok := ic.mutation.Deadline(); !ok {
		return &ValidationError{Name: "deadline", err: errors.New(`database: missing required field "Invoice.deadline"`)}
	}
	if _, ok := ic.mutation.WalletAddress(); !ok {
		return &ValidationError{Name: "wallet_address", err: errors.New(`database: missing required field "Invoice.wallet_address"`)}
	}
	if v, ok := ic.mutation.WalletAddress(); ok {
		if err := invoice.WalletAddressValidator(v); err != nil {
			return &ValidationError{Name: "wallet_address", err: fmt.Errorf(`database: validator failed for field "Invoice.wallet_address": %w`, err)}
		}
	}
	if _, ok := ic.mutation.EncryptedSalt(); !ok {
		return &ValidationError{Name: "encrypted_salt", err: errors.New(`database: missing required field "Invoice.encrypted_salt"`)}
	}
	if v, ok := ic.mutation.EncryptedSalt(); ok {
		if err := invoice.EncryptedSaltValidator(v); err != nil {
			return &ValidationError{Name: "encrypted_salt", err: fmt.Errorf(`database: validator failed for field "Invoice.encrypted_salt": %w`, err)}
		}
	}
	if v, ok := ic.mutation.ID(); ok {
		if err := invoice.IDValidator(v); err != nil {
			return &ValidationError{Name: "id", err: fmt.Errorf(`database: validator failed for field "Invoice.id": %w`, err)}
		}
	}
	return nil
}

func (ic *InvoiceCreate) sqlSave(ctx context.Context) (*Invoice, error) {
	if err := ic.check(); err != nil {
		return nil, err
	}
	_node, _spec, err := ic.createSpec()
	if err != nil {
		return nil, err
	}
	if err := sqlgraph.CreateNode(ctx, ic.driver, _spec); err != nil {
		if sqlgraph.IsConstraintError(err) {
			err = &ConstraintError{msg: err.Error(), wrap: err}
		}
		return nil, err
	}
	if _spec.ID.Value != nil {
		if id, ok := _spec.ID.Value.(string); ok {
			_node.ID = id
		} else {
			return nil, fmt.Errorf("unexpected Invoice.ID type: %T", _spec.ID.Value)
		}
	}
	ic.mutation.id = &_node.ID
	ic.mutation.done = true
	return _node, nil
}

func (ic *InvoiceCreate) createSpec() (*Invoice, *sqlgraph.CreateSpec, error) {
	var (
		_node = &Invoice{config: ic.config}
		_spec = sqlgraph.NewCreateSpec(invoice.Table, sqlgraph.NewFieldSpec(invoice.FieldID, field.TypeString))
	)
	if id, ok := ic.mutation.ID(); ok {
		_node.ID = id
		_spec.ID.Value = id
	}
	if value, ok := ic.mutation.MinAmount(); ok {
		vv, err := invoice.ValueScanner.MinAmount.Value(value)
		if err != nil {
			return nil, nil, err
		}
		_spec.SetField(invoice.FieldMinAmount, field.TypeString, vv)
		_node.MinAmount = value
	}
	if value, ok := ic.mutation.Recipient(); ok {
		_spec.SetField(invoice.FieldRecipient, field.TypeString, value)
		_node.Recipient = value
	}
	if value, ok := ic.mutation.Beneficiary(); ok {
		_spec.SetField(invoice.FieldBeneficiary, field.TypeString, value)
		_node.Beneficiary = value
	}
	if value, ok := ic.mutation.Asset(); ok {
		_spec.SetField(invoice.FieldAsset, field.TypeString, value)
		_node.Asset = value
	}
	if value, ok := ic.mutation.Metadata(); ok {
		_spec.SetField(invoice.FieldMetadata, field.TypeString, value)
		_node.Metadata = value
	}
	if value, ok := ic.mutation.CreateAt(); ok {
		_spec.SetField(invoice.FieldCreateAt, field.TypeTime, value)
		_node.CreateAt = value
	}
	if value, ok := ic.mutation.Deadline(); ok {
		_spec.SetField(invoice.FieldDeadline, field.TypeTime, value)
		_node.Deadline = value
	}
	if value, ok := ic.mutation.FillAt(); ok {
		_spec.SetField(invoice.FieldFillAt, field.TypeTime, value)
		_node.FillAt = &value
	}
	if value, ok := ic.mutation.LastCheckoutAt(); ok {
		_spec.SetField(invoice.FieldLastCheckoutAt, field.TypeTime, value)
		_node.LastCheckoutAt = &value
	}
	if value, ok := ic.mutation.CancelAt(); ok {
		_spec.SetField(invoice.FieldCancelAt, field.TypeTime, value)
		_node.CancelAt = &value
	}
	if value, ok := ic.mutation.WalletAddress(); ok {
		_spec.SetField(invoice.FieldWalletAddress, field.TypeString, value)
		_node.WalletAddress = value
	}
	if value, ok := ic.mutation.EncryptedSalt(); ok {
		_spec.SetField(invoice.FieldEncryptedSalt, field.TypeBytes, value)
		_node.EncryptedSalt = value
	}
	return _node, _spec, nil
}

// InvoiceCreateBulk is the builder for creating many Invoice entities in bulk.
type InvoiceCreateBulk struct {
	config
	err      error
	builders []*InvoiceCreate
}

// Save creates the Invoice entities in the database.
func (icb *InvoiceCreateBulk) Save(ctx context.Context) ([]*Invoice, error) {
	if icb.err != nil {
		return nil, icb.err
	}
	specs := make([]*sqlgraph.CreateSpec, len(icb.builders))
	nodes := make([]*Invoice, len(icb.builders))
	mutators := make([]Mutator, len(icb.builders))
	for i := range icb.builders {
		func(i int, root context.Context) {
			builder := icb.builders[i]
			builder.defaults()
			var mut Mutator = MutateFunc(func(ctx context.Context, m Mutation) (Value, error) {
				mutation, ok := m.(*InvoiceMutation)
				if !ok {
					return nil, fmt.Errorf("unexpected mutation type %T", m)
				}
				if err := builder.check(); err != nil {
					return nil, err
				}
				builder.mutation = mutation
				var err error
				nodes[i], specs[i], err = builder.createSpec()
				if err != nil {
					return nil, err
				}
				if i < len(mutators)-1 {
					_, err = mutators[i+1].Mutate(root, icb.builders[i+1].mutation)
				} else {
					spec := &sqlgraph.BatchCreateSpec{Nodes: specs}
					// Invoke the actual operation on the latest mutation in the chain.
					if err = sqlgraph.BatchCreate(ctx, icb.driver, spec); err != nil {
						if sqlgraph.IsConstraintError(err) {
							err = &ConstraintError{msg: err.Error(), wrap: err}
						}
					}
				}
				if err != nil {
					return nil, err
				}
				mutation.id = &nodes[i].ID
				mutation.done = true
				return nodes[i], nil
			})
			for i := len(builder.hooks) - 1; i >= 0; i-- {
				mut = builder.hooks[i](mut)
			}
			mutators[i] = mut
		}(i, ctx)
	}
	if len(mutators) > 0 {
		if _, err := mutators[0].Mutate(ctx, icb.builders[0].mutation); err != nil {
			return nil, err
		}
	}
	return nodes, nil
}

// SaveX is like Save, but panics if an error occurs.
func (icb *InvoiceCreateBulk) SaveX(ctx context.Context) []*Invoice {
	v, err := icb.Save(ctx)
	if err != nil {
		panic(err)
	}
	return v
}

// Exec executes the query.
func (icb *InvoiceCreateBulk) Exec(ctx context.Context) error {
	_, err := icb.Save(ctx)
	return err
}

// ExecX is like Exec, but panics if an error occurs.
func (icb *InvoiceCreateBulk) ExecX(ctx context.Context) {
	if err := icb.Exec(ctx); err != nil {
		panic(err)
	}
}
