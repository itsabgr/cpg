package schema

import (
	"entgo.io/ent"
	"entgo.io/ent/schema/field"
	"math/big"
	"time"
)

type Invoice struct {
	ent.Schema
}

func (Invoice) Fields() []ent.Field {
	return []ent.Field{
		field.String("id").Unique().NotEmpty().Immutable(),
		field.String("min_amount").GoType(&big.Int{}).ValueScanner(field.TextValueScanner[*big.Int]{}).NotEmpty().Immutable(),
		field.String("recipient").NotEmpty().Immutable(),
		field.String("beneficiary").NotEmpty().Immutable(),
		field.String("asset").NotEmpty().Immutable(),
		field.String("metadata").MaxLen(256).Immutable(),
		field.Time("create_at").Default(time.Now).Immutable(),
		field.Time("deadline").Immutable(),
		field.Time("fill_at").Optional().Nillable(),
		field.Time("last_checkout_at").Optional().Nillable(),
		field.Time("checkout_request_at").Optional().Nillable(),
		field.Time("cancel_at").Optional().Nillable(),
		field.String("wallet_address").Unique().NotEmpty().Immutable(),
		field.Bytes("encrypted_salt").Sensitive().Unique().NotEmpty().Immutable(),
	}
}
