package main

import (
	"context"
	"cpg/internal/assets/eth"
	"cpg/internal/daemon"
	"cpg/pkg/cpg"
	"cpg/pkg/crypto"
	"cpg/pkg/ent/database"
	"cpg/pkg/proto"
	"github.com/itsabgr/ge"
	_ "github.com/joho/godotenv/autoload"
	_ "github.com/lib/pq"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
	"log/slog"
	"net"
	"os"
)

type env struct {
	AssetsConfig  string `env:"ASSETS_CONFIG,notEmpty"`
	SaltKeyring   string `env:"SALT_KEYRING,notEmpty"`
	BackupKeyring string `env:"BACKUP_KEYRING,notEmpty"`
	GRPCServer    string `env:"GRPC_SERVER,notEmpty"`
	PostgresURI   string `env:"PG_URI,notEmpty"`
}

func main() {
	daemon.Run(func(ctx context.Context, config env) {
		defer slog.Info("bye")

		assets := prepareAssets(ctx, config.AssetsConfig)

		ge.Assert(assets.Count() > 0, ge.New("no asset loaded"))

		slog.Debug("assets loaded", slog.Int("count", assets.Count()))

		ln := ge.Must(net.Listen("tcp", config.GRPCServer))
		defer func() { _ = ln.Close() }()

		slog.Debug("tcp server listening", slog.String("addr", ln.Addr().String()))

		dbClient := ge.Must(database.Open("postgres", config.PostgresURI))

		defer func() { _ = dbClient.Close() }()

		ge.Throw(dbClient.Schema.Create(ctx))

		slog.Debug("connected to postgres")
		saltKR := ge.Must(crypto.LoadKeyRingFromFile(config.SaltKeyring))
		backupKR := ge.Must(crypto.LoadKeyRingFromFile(config.BackupKeyring))

		if saltKR.Contains(backupKR) {
			slog.Warn("keyrings have at least one match key")
		}
		slog.Debug("keyring loaded", slog.Int("salt", saltKR.Size()), slog.Int("backup", backupKR.Size()))

		grpcServer := grpc.NewServer()

		proto.RegisterCPGServer(grpcServer, cpg.NewGRPCServer(cpg.NewCPG(assets, cpg.NewDB(dbClient), saltKR, backupKR)))

		if daemon.Debug() {
			reflection.Register(grpcServer)
			slog.Debug("grpc reflection enabled")
		}

		go func() {
			<-ctx.Done()
			slog.Info("stopping grpc server")
			grpcServer.GracefulStop()
		}()

		slog.Info("start serving")
		ge.Throw(grpcServer.Serve(ln))

	})
}

func prepareAssets(ctx context.Context, configFilePath string) *cpg.Assets {
	cpg.RegisterAssetFactory(eth.Factory{})
	configData := ge.Must(os.ReadFile(configFilePath))
	return ge.Must(cpg.ParseAssetsConfig(ctx, configData))
}
