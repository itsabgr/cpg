package eth

import (
	"context"
	"cpg/pkg/cpg"
	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/itsabgr/ge"
	"math/big"
	"time"
)

var _ cpg.AssetFactory = Factory{}

type Factory struct{}

type FactoryConfig struct {
	EthClientEndpoint  string   `json:"eth_client_endpoint"`
	TxGasLimit         uint64   `json:"tx_gas_limit"`
	MinDelaySeconds    uint16   `json:"min_delay_seconds"`
	ChainID            *big.Int `json:"chain_id"`
	MinAllowedAmount   *big.Int `json:"min_allowed_amount"`
	MaxAllowedGasPrice *big.Int `json:"max_allowed_gas_price"`
}

func (Factory) Name() string {
	return "eth"
}

func (Factory) Config() any {
	return &FactoryConfig{}
}

func (fac Factory) New(ctx context.Context, config any) (cpg.Asset, error) {
	conf := config.(*FactoryConfig)
	ethClient, err := ethclient.Dial(conf.EthClientEndpoint)
	if err != nil {
		return nil, ge.Wrap(ge.New("failed to dial eth endpoint"), err)
	}
	return New(ctx, Config{
		EthClient:          ethClient,
		TxGasLimit:         conf.TxGasLimit,
		ChainID:            conf.ChainID,
		MinAllowedAmount:   conf.MinAllowedAmount,
		MinDelay:           time.Duration(conf.MinDelaySeconds) * time.Second,
		MaxAllowedGasPrice: conf.MaxAllowedGasPrice,
	})
}
