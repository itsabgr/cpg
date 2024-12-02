package cpg

import (
	"context"
	"encoding/json"
	"github.com/itsabgr/ge"
	"runtime"
)

func ParseAssetsConfig(ctx context.Context, data []byte) (*Assets, error) {
	config := map[string]json.RawMessage{}

	if err := json.Unmarshal(data, &config); err != nil {
		return nil, err
	}
	var driverConfig struct {
		Type   string `json:"type"`
		Enable bool   `json:"enable"`
	}
	assets := NewAssets()
	for assetName, raw := range config {
		runtime.Gosched()
		if err := ctx.Err(); err != nil {
			return nil, err
		}
		driverConfig.Type = ""
		driverConfig.Enable = false
		if err := json.Unmarshal(raw, &driverConfig); err != nil {
			return nil, err
		}
		if false == driverConfig.Enable {
			continue
		}

		factory, factoryExists := assetFactories[driverConfig.Type]

		if false == factoryExists {
			return nil, ge.Detail(ge.New("asset factory not found"), ge.D{"asset": assetName, "factory": driverConfig.Type})
		}

		factoryConfig := factory.Config()

		if err := json.Unmarshal(raw, factoryConfig); err != nil {
			return nil, err
		}

		asset, err := factory.New(ctx, factoryConfig)
		if err != nil {
			return nil, ge.Wrap(ge.Detail(ge.New("failed to init asset"), ge.D{"asset": assetName, "factory": driverConfig.Type}), err)
		}
		ge.Assert(assets.Register(assetName, asset))
	}
	return assets, nil
}
