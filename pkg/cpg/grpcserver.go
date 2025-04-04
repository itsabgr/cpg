package cpg

import (
	"context"
	"cpg/pkg/proto"
	"cpg/pkg/ratelimit"
	"github.com/golang/protobuf/ptypes/empty"
	"github.com/itsabgr/ge"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/durationpb"
	"google.golang.org/protobuf/types/known/timestamppb"
	"math/big"
	"time"
)

type grpcServer struct {
	rateLimit *ratelimit.RateLimit
	cpg       *CPG
	assets    map[string]*proto.AssetInfo
	proto.UnimplementedCPGServer
}

func NewGRPCServer(cpg *CPG, rateLimit *ratelimit.RateLimit) proto.CPGServer {

	serv := grpcServer{
		rateLimit: rateLimit,
		cpg:       cpg,
	}

	serv.assets = make(map[string]*proto.AssetInfo, len(cpg.assets.map_))

	for name, info := range cpg.assets.Infos() {
		serv.assets[name] = &proto.AssetInfo{
			MinDelay: durationpb.New(info.MinDelay),
		}
	}

	return serv

}

func (serv grpcServer) Ping(ctx context.Context, input *proto.PingInput) (*proto.PingOutput, error) {

	return &proto.PingOutput{
		Now:  timestamppb.Now(),
		Pong: input.GetPing(),
	}, nil

}

func (serv grpcServer) ListAssets(ctx context.Context, _ *empty.Empty) (*proto.ListAssetsOutput, error) {

	return &proto.ListAssetsOutput{
		Assets: serv.assets,
	}, nil

}

func (serv grpcServer) RecoverInvoice(ctx context.Context, input *proto.RecoverInvoiceInput) (*empty.Empty, error) {

	cancel, err := serv.rateLimitRequest(&ctx, time.Second*2, input)
	if err != nil {
		return nil, err
	}
	defer cancel()

	if ok, err := serv.rateLimit.Limit(ctx, input.GetInvoiceId(), time.Second*5); err != nil {
		return nil, err
	} else if !ok {
		return nil, status.Error(codes.ResourceExhausted, "busy")
	}

	err = serv.cpg.RecoverInvoice(ctx, RecoverInvoiceParams{
		InvoiceID:     input.GetInvoiceId(),
		InvoiceBackup: input.GetInvoiceBackup(),
	})

	if err != nil {
		return nil, err
	}

	return nil, nil

}

func (serv grpcServer) CreateInvoice(ctx context.Context, input *proto.CreateInvoiceInput) (*proto.CreateInvoiceOutput, error) {

	result, err := serv.cpg.CreateInvoice(ctx, CreateInvoiceParams{
		AssetName:    input.GetAssetName(),
		Metadata:     input.GetMetadata(),
		Recipient:    input.GetRecipient(),
		Beneficiary:  input.GetBeneficiary(),
		AutoCheckout: input.GetAutoCheckout(),
		MinAmount:    str2BigInt(input.GetMinAmount(), 10),
		Deadline:     input.GetDeadline().AsTime(),
	})
	if err != nil {
		return nil, err
	}

	return &proto.CreateInvoiceOutput{
		InvoiceId:     result.InvoiceID,
		InvoiceBackup: result.InvoiceBackup,
	}, nil

}

func (serv grpcServer) CancelInvoice(ctx context.Context, input *proto.CancelInvoiceInput) (*empty.Empty, error) {

	cancel, err := serv.rateLimitRequest(&ctx, time.Second*2, input)
	if err != nil {
		return nil, err
	}
	defer cancel()

	err = serv.cpg.CancelInvoice(ctx, CancelInvoiceParams{
		InvoiceID:     input.GetInvoiceId(),
		WalletAddress: input.GetWalletAddress(),
	})
	if err != nil {
		return nil, err
	}

	return nil, nil

}

func (serv grpcServer) RequestCheckout(ctx context.Context, input *proto.RequestCheckoutInput) (*empty.Empty, error) {

	cancel, err := serv.rateLimitRequest(&ctx, time.Second*1, input)
	if err != nil {
		return nil, err
	}
	defer cancel()

	err = serv.cpg.RequestCheckout(ctx, RequestCheckoutParams{
		InvoiceID: input.GetInvoiceId(),
	})
	if err != nil {
		return nil, err
	}

	return nil, nil

}

func (serv grpcServer) GetInvoice(ctx context.Context, input *proto.GetInvoiceInput) (*proto.GetInvoiceOutput, error) {

	result, err := serv.cpg.GetInvoice(ctx, GetInvoiceParams{
		InvoiceID: input.GetInvoiceId(),
	})
	if err != nil {
		return nil, err
	}

	return &proto.GetInvoiceOutput{
		MinAmount:         result.MinAmount.Text(10),
		Recipient:         result.Recipient,
		Beneficiary:       result.Beneficiary,
		Asset:             result.Asset,
		Metadata:          result.Metadata,
		CreateAt:          timestamppb.New(result.CreateAt),
		Deadline:          timestamppb.New(result.Deadline),
		FillAt:            optionalTime2timestamp(result.FillAt),
		CancelAt:          optionalTime2timestamp(result.CancelAt),
		CheckoutRequestAt: optionalTime2timestamp(result.CheckoutRequestAt),
		LastCheckoutAt:    optionalTime2timestamp(result.LastCheckoutAt),
		AutoCheckout:      result.AutoCheckout,
		WalletAddress:     result.WalletAddress,
		Status:            proto.InvoiceStatus(result.Status),
	}, nil

}

func (serv grpcServer) CheckInvoice(ctx context.Context, input *proto.CheckInvoiceInput) (*proto.CheckInvoiceOutput, error) {

	if input.GetInvoiceId() != "" {
		cancel, err := serv.rateLimitRequest(&ctx, time.Second*5, input)
		if err != nil {
			return nil, err
		}
		defer cancel()
	}

	result, err := serv.cpg.CheckInvoice(ctx, CheckInvoiceParams{
		InvoiceID:     input.GetInvoiceId(),
		WalletAddress: input.GetWalletAddress(),
	})
	if err != nil {
		return nil, err
	}

	return &proto.CheckInvoiceOutput{
		InvoiceStatus: proto.InvoiceStatus(result.InvoiceStatus),
	}, nil

}

func (serv grpcServer) TryCheckoutInvoice(ctx context.Context, input *proto.TryCheckoutInvoiceInput) (*empty.Empty, error) {

	cancel, err := serv.rateLimitRequest(&ctx, time.Second*10, input)
	if err != nil {
		return nil, err
	}
	defer cancel()

	err = serv.cpg.TryCheckoutInvoice(ctx, TryCheckoutInvoiceParams{
		InvoiceID: input.GetInvoiceId(),
	})

	if err != nil {
		return nil, err
	}

	return nil, nil

}

func optionalTime2timestamp(t *time.Time) *timestamppb.Timestamp {
	if t == nil {
		return nil
	}
	return timestamppb.New(*t)
}

func str2BigInt(str string, base int) *big.Int {
	n := &big.Int{}
	if n2, ok := n.SetString(str, base); ok {
		n = n2
	} else {
		n = big.NewInt(0)
	}
	return n
}

func (serv grpcServer) rateLimitRequest(ctx *context.Context, duration time.Duration, input interface{ GetInvoiceId() string }) (context.CancelFunc, error) {

	ge.Assert(duration > 0)

	if ok, err := serv.rateLimit.Limit(*ctx, input.GetInvoiceId(), duration+time.Second); err != nil {
		return nil, err
	} else if !ok {
		return nil, status.Error(codes.ResourceExhausted, "invoice is busy")
	}

	timeoutCtx, cancelCtx := context.WithTimeout(*ctx, duration)

	*ctx = timeoutCtx

	return cancelCtx, nil

}
