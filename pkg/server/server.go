package server

import (
	"context"
	"github.com/google/uuid"
	"github.com/itsabgr/cpg"
	"github.com/itsabgr/cpg/pkg/proto"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/durationpb"
	"google.golang.org/protobuf/types/known/emptypb"
	"math/big"
)

var _ proto.CPGServer = &Server{}

type Server struct {
	proto.UnimplementedCPGServer
	impl *cpg.CPG
}

func (srv *Server) ListAssets(context.Context, *emptypb.Empty) (*proto.ListAssetsResponse, error) {
	assets := srv.impl.Assets()
	response := &proto.ListAssetsResponse{Assets: make(map[string]*proto.ListAssetsResponse_Asset, len(assets))}
	for _, asset := range assets {
		response.Assets[asset.Name] = &proto.ListAssetsResponse_Asset{
			Title: asset.Name,
			Delay: durationpb.New(asset.Delay),
		}
	}
	return response, nil
}

func (srv *Server) CreateInvoice(ctx context.Context, request *proto.CreateInvoiceRequest) (*proto.CreateInvoiceResponse, error) {
	amount, ok := (&big.Int{}).SetString(request.GetAmount(), 10)
	if !ok {
		return nil, status.Error(codes.InvalidArgument, "failed to parse amount")
	}
	invoiceId, err := srv.impl.CreateInvoice(ctx, &cpg.CreateInvoiceInput{
		Asset:     request.GetAsset(),
		Recipient: request.GetRecipient(),
		Amount:    amount,
	})
	if err != nil {
		return nil, err
	}
	return &proto.CreateInvoiceResponse{Uuid: invoiceId.String()}, nil
}

func (srv *Server) GetInvoice(ctx context.Context, request *proto.GetInvoiceRequest) (*proto.GetInvoiceResponse, error) {
	invoiceId, err := uuid.Parse(request.GetUuid())
	if err != nil {
		return nil, status.Error(codes.InvalidArgument, err.Error())
	}
	invoice, err := srv.impl.GetInvoice(ctx, invoiceId)
	if err != nil {
		return nil, err
	}
	if invoice == nil {
		return nil, status.Error(codes.NotFound, "invoice not found")
	}
	return &proto.GetInvoiceResponse{
		Status:        invoice.Status.String(),
		WalletAddress: invoice.Wallet,
		Request: &proto.CreateInvoiceRequest{
			Asset:     invoice.Asset,
			Recipient: invoice.Recipient,
			Amount:    invoice.Amount.Text(10),
		},
	}, nil
}
func (srv *Server) CheckWallet(ctx context.Context, request *proto.CheckWalletRequest) (*emptypb.Empty, error) {
	err := srv.impl.EnqueueCheckWallet(ctx, request.GetWalletAddress(), request.GetAsset())
	return nil, err
}
func (srv *Server) Checkout(ctx context.Context, request *proto.CheckoutRequest) (*emptypb.Empty, error) {
	err := srv.impl.EnqueueCheckWallet(ctx, request.GetInvoiceId(), request.GetAsset())
	return nil, err
}
func New(impl *cpg.CPG) *Server {
	return &Server{
		impl: impl,
	}
}
