// Code generated by ent, DO NOT EDIT.

package database

import (
	"cpg/pkg/ent/database/invoice"
	"cpg/pkg/ent/schema"
	"math/big"
	"time"

	"entgo.io/ent/schema/field"
)

// The init function reads all schema descriptors with runtime code
// (default values, validators, hooks and policies) and stitches it
// to their package variables.
func init() {
	invoiceFields := schema.Invoice{}.Fields()
	_ = invoiceFields
	// invoiceDescMinAmount is the schema descriptor for min_amount field.
	invoiceDescMinAmount := invoiceFields[1].Descriptor()
	invoice.ValueScanner.MinAmount = invoiceDescMinAmount.ValueScanner.(field.TypeValueScanner[*big.Int])
	// invoice.MinAmountValidator is a validator for the "min_amount" field. It is called by the builders before save.
	invoice.MinAmountValidator = invoiceDescMinAmount.Validators[0].(func(string) error)
	// invoiceDescRecipient is the schema descriptor for recipient field.
	invoiceDescRecipient := invoiceFields[2].Descriptor()
	// invoice.RecipientValidator is a validator for the "recipient" field. It is called by the builders before save.
	invoice.RecipientValidator = invoiceDescRecipient.Validators[0].(func(string) error)
	// invoiceDescBeneficiary is the schema descriptor for beneficiary field.
	invoiceDescBeneficiary := invoiceFields[3].Descriptor()
	// invoice.BeneficiaryValidator is a validator for the "beneficiary" field. It is called by the builders before save.
	invoice.BeneficiaryValidator = invoiceDescBeneficiary.Validators[0].(func(string) error)
	// invoiceDescAsset is the schema descriptor for asset field.
	invoiceDescAsset := invoiceFields[4].Descriptor()
	// invoice.AssetValidator is a validator for the "asset" field. It is called by the builders before save.
	invoice.AssetValidator = invoiceDescAsset.Validators[0].(func(string) error)
	// invoiceDescMetadata is the schema descriptor for metadata field.
	invoiceDescMetadata := invoiceFields[5].Descriptor()
	// invoice.MetadataValidator is a validator for the "metadata" field. It is called by the builders before save.
	invoice.MetadataValidator = invoiceDescMetadata.Validators[0].(func(string) error)
	// invoiceDescCreateAt is the schema descriptor for create_at field.
	invoiceDescCreateAt := invoiceFields[6].Descriptor()
	// invoice.DefaultCreateAt holds the default value on creation for the create_at field.
	invoice.DefaultCreateAt = invoiceDescCreateAt.Default.(func() time.Time)
	// invoiceDescWalletAddress is the schema descriptor for wallet_address field.
	invoiceDescWalletAddress := invoiceFields[10].Descriptor()
	// invoice.WalletAddressValidator is a validator for the "wallet_address" field. It is called by the builders before save.
	invoice.WalletAddressValidator = invoiceDescWalletAddress.Validators[0].(func(string) error)
	// invoiceDescEncryptedSalt is the schema descriptor for encrypted_salt field.
	invoiceDescEncryptedSalt := invoiceFields[11].Descriptor()
	// invoice.EncryptedSaltValidator is a validator for the "encrypted_salt" field. It is called by the builders before save.
	invoice.EncryptedSaltValidator = invoiceDescEncryptedSalt.Validators[0].(func([]byte) error)
	// invoiceDescID is the schema descriptor for id field.
	invoiceDescID := invoiceFields[0].Descriptor()
	// invoice.IDValidator is a validator for the "id" field. It is called by the builders before save.
	invoice.IDValidator = invoiceDescID.Validators[0].(func(string) error)
}
