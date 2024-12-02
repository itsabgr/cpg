package eth

import (
	"bytes"
	"context"
	"cpg/pkg/cpg"
	"crypto/ecdsa"
	"errors"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/itsabgr/fak"
	"github.com/itsabgr/ge"
	"math/big"
	"time"
)

const SaltSize = 32

type Config struct {
	EthClient          *ethclient.Client
	MinDelay           time.Duration
	TxGasLimit         uint64
	ChainID            *big.Int
	MinAllowedAmount   *big.Int
	MaxAllowedGasPrice *big.Int
}

func New(ctx context.Context, config Config) (cpg.Asset, error) {
	if config.MinDelay < time.Second {
		panic(ge.New("too less min delay"))
	}
	if config.TxGasLimit <= 0 {
		panic(ge.New("zero gas limit"))
	}
	if config.MinAllowedAmount.Cmp(big.NewInt(1)) < 0 {
		panic(ge.New("non-positive min amount"))
	}
	if config.MaxAllowedGasPrice.Cmp(big.NewInt(0)) <= 0 {
		panic(ge.New("non-positive max fee"))
	}

	chainId, err := func() (*big.Int, error) {
		timeout, cancel := context.WithTimeout(ctx, time.Second*2)
		defer cancel()
		return config.EthClient.ChainID(timeout)
	}()

	if err != nil {
		return nil, ge.Wrap(ge.New("failed to get chain id"), err)
	}

	if chainId.Cmp(config.ChainID) != 0 {
		return nil, ge.Detail(ge.New("mismatched chain id"), ge.D{
			"endpoint": chainId.String(),
			"config":   config.ChainID.String(),
		})
	}

	return &asset{
		ethClient:          config.EthClient,
		txGasLimit:         *(&big.Int{}).SetUint64(config.TxGasLimit),
		minAllowedAmount:   *(&big.Int{}).Set(config.MinAllowedAmount),
		maxAllowedGasPrice: *(&big.Int{}).Set(config.MaxAllowedGasPrice),
		chainID:            *(&big.Int{}).Set(config.ChainID),
		info: cpg.AssetInfo{
			MinDelay:   config.MinDelay,
			SaltLength: SaltSize,
		},
	}, nil
}

type asset struct {
	ethClient          *ethclient.Client
	txGasLimit         big.Int
	minAllowedAmount   big.Int
	chainID            big.Int
	maxAllowedGasPrice big.Int
	info               cpg.AssetInfo
}

func (ass *asset) GetBalance(ctx context.Context, invoice *cpg.Invoice) (*big.Int, error) {
	walletAddress := common.HexToAddress(invoice.WalletAddress)
	timeout, cancel := context.WithTimeout(ctx, time.Second*5)
	defer cancel()
	balance, err := ass.ethClient.BalanceAt(timeout, walletAddress, nil)
	if err != nil {
		return nil, err
	}
	return balance, nil
}

func (ass *asset) PrepareInvoice(ctx context.Context, invoice *cpg.Invoice) error {

	if invoice.MinAmount.Cmp(&ass.minAllowedAmount) < 0 {
		return ge.New("too low amount")
	}

	if !validateAddress(invoice.Recipient) {
		return ge.New("invalid recipient")
	}

	if !validateAddress(invoice.Beneficiary) {
		return ge.New("invalid beneficiary")
	}

	salt := invoice.DecryptSalt()

	if len(salt) != SaltSize {
		return ge.New("invalid salt size")
	}

	walletPrivateKey := ge.Must(ecdsa.GenerateKey(crypto.S256(), bytes.NewReader(salt)))

	walletAddress := crypto.PubkeyToAddress(*walletPrivateKey.Public().(*ecdsa.PublicKey)).Hex()

	invoice.WalletAddress = walletAddress

	return nil

}

func (ass *asset) checkGasPrice(ctx context.Context) (*big.Int, error) {
	timeout, cancel := context.WithTimeout(ctx, time.Second*2)
	defer cancel()
	gasPrice, err := ass.ethClient.SuggestGasPrice(timeout)
	if err != nil {
		return nil, err
	}
	if gasPrice.Cmp(&ass.maxAllowedGasPrice) > 0 {
		return nil, errors.New("too high gas price")
	}
	return gasPrice, nil
}

func (ass *asset) TryFlush(ctx context.Context, invoice *cpg.Invoice) error {

	gasPrice, err := ass.checkGasPrice(ctx)
	if err != nil {
		return ge.Wrap(ge.New("failed to check gas price"), err)
	}

	walletBalance, err := ass.GetBalance(ctx, invoice)
	if err != nil {
		return ge.Wrap(ge.New("failed to get wallet balance"), err)
	}

	if walletBalance.Cmp(&ass.minAllowedAmount) < 0 {
		return ge.Wrap(ge.New("too less wallet balance"), err)
	}

	txFee := big.NewInt(0).Mul(gasPrice, &ass.txGasLimit)

	if txFee.Cmp(walletBalance) >= 0 {
		return errors.New("tx fee overcomes the wallet balance")
	}

	walletPendingNonce, err := ass.ethClient.PendingNonceAt(ctx, common.HexToAddress(invoice.WalletAddress))
	if err != nil {
		return ge.Wrap(ge.New("failed to get wallet pending nonce"), err)
	}

	salt := invoice.DecryptSalt()
	if len(salt) != SaltSize {
		return ge.New("invalid salt size")
	}

	walletPrivateKey := ge.Must(ecdsa.GenerateKey(crypto.S256(), bytes.NewReader(salt)))

	signedTx, err := types.SignTx(types.NewTx(&types.LegacyTx{
		Nonce:    walletPendingNonce,
		GasPrice: gasPrice,
		Gas:      ass.txGasLimit.Uint64(),
		To:       fak.Ptr(common.HexToAddress(invoice.Destination())),
		Value:    (&big.Int{}).Sub(walletBalance, txFee),
	}), types.NewEIP155Signer(&ass.chainID), walletPrivateKey)
	if err != nil {
		return ge.Wrap(ge.New("failed to sign tx"), err)
	}

	err = ass.ethClient.SendTransaction(ctx, signedTx)
	if err != nil {
		return ge.Wrap(ge.New("failed to send tx"), err)
	}

	return nil

}

func (ass *asset) Info() cpg.AssetInfo {
	return ass.info
}

func validateAddress(address string) bool {
	return len(address) == 42 && common.HexToAddress(address).String() == address
}
