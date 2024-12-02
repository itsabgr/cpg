package cpg

import (
	"context"
	"github.com/itsabgr/ge"
	"iter"
	"math/big"
	"strings"
	"sync"
	"time"
	"unicode/utf8"
)

var assetFactories = map[string]AssetFactory{}

type AssetFactory interface {
	Name() string
	Config() any
	New(ctx context.Context, config any) (Asset, error)
}

type AssetInfo struct {
	MinDelay   time.Duration
	SaltLength int
}

type Asset interface {
	Info() AssetInfo
	PrepareInvoice(ctx context.Context, invoice *Invoice) error
	GetBalance(ctx context.Context, invoice *Invoice) (*big.Int, error)
	TryFlush(ctx context.Context, invoice *Invoice) error
}

type Assets struct {
	_    sync.Mutex
	map_ map[string]Asset
}

func NewAssets() *Assets {
	return &Assets{
		map_: make(map[string]Asset),
	}
}

func (a *Assets) Register(name string, asset Asset) bool {
	ge.Assert(asset != nil, ge.New("nil asset"))
	ge.Assert(name != "", ge.New("empty asset name"))
	ge.Assert(utf8.ValidString(name), ge.New("invalid utf8 ascii name"))
	ge.Assert(strings.ToLower(name) == name, ge.New("asset name should be lower case"))

	if a.map_ == nil {
		a.map_ = map[string]Asset{}
	}

	if a.map_[name] != nil {
		return false
	}

	a.map_[name] = asset
	return true
}

func (a *Assets) Get(name string) Asset {
	return a.map_[name]
}

func (a *Assets) Count() int {
	return len(a.map_)
}

func (a *Assets) Infos() iter.Seq2[string, AssetInfo] {
	return func(yield func(string, AssetInfo) bool) {
		for n, p := range a.map_ {
			if !yield(n, p.Info()) {
				return
			}
		}
	}
}

func AssetFactories() iter.Seq2[string, AssetFactory] {
	return func(yield func(string, AssetFactory) bool) {
		for name, factory := range assetFactories {
			if !yield(name, factory) {
				return
			}
		}
	}
}

func RegisterAssetFactory(factory AssetFactory) {

	ge.Assert(factory != nil, ge.New("try to register nil asset factory"))

	name := factory.Name()

	ge.Assert(name != "", ge.New("empty asset factory name"))

	ge.Assert(assetFactories[name] == nil, ge.Detail(ge.New("duplicate asset factory name"), ge.D{"name": name}))

	assetFactories[name] = factory

}
