package daemon

import (
	"context"
	"github.com/caarlos0/env/v11"
	"github.com/itsabgr/ge"
	"github.com/itsabgr/ge/plot"
	"io"
	"os"
	"os/signal"
)

func Warn(err error) {
	if debug {
		ge.Must(os.Stderr.Write(plot.Tree(err).Bytes()))
	}
}

func Run[C any](main func(ctx context.Context, config C)) {
	rec := ge.Try(func() {
		ctx, cancel := signal.NotifyContext(context.Background(), os.Kill, os.Interrupt)
		defer cancel()
		var conf C
		ge.Throw(env.Parse(&conf))
		main(ctx, conf)
	})
	switch e := rec.(type) {
	case nil:
	case string:
		ge.Must(io.WriteString(os.Stderr, e))
	case error:
		ge.Must(os.Stderr.Write(plot.Tree(e.(error)).Bytes()))
	default:
		panic(e)
	}
}
