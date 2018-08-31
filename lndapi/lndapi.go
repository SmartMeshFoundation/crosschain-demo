package lndapi

import (
	"context"
	"encoding/hex"
	"fmt"
	"os"
	"path/filepath"

	"github.com/btcsuite/btcutil"
	"github.com/ethereum/go-ethereum/common"
	"github.com/gogo/protobuf/jsonpb"
	"github.com/golang/protobuf/proto"
	"github.com/lightningnetwork/lnd/lncfg"
	"github.com/lightningnetwork/lnd/lnrpc"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
)

const (
	defaultTLSCertFilename  = "tls.cert"
	defaultMacaroonFilename = "admin.macaroon"
	defaultRPCPort          = "8002"
	defaultRPCHostPort      = "192.124.8" + defaultRPCPort
)

var (
	//Commit stores the current commit hash of this build. This should be
	//set using -ldflags during compilation.
	Commit string

	defaultLndDir       = btcutil.AppDataDir("lnd", false)
	defaultTLSCertPath  = filepath.Join(defaultLndDir, defaultTLSCertFilename)
	defaultMacaroonPath = filepath.Join(defaultLndDir, defaultMacaroonFilename)
)

func fatal(err error) {
	fmt.Fprintf(os.Stderr, "[lncli] %v\n", err)
	os.Exit(1)
}

type LndAPI struct {
	rpcServer string
	client    lnrpc.LightningClient
}

func NewLndAPI(rpcServer string) *LndAPI {
	l := &LndAPI{
		rpcServer: rpcServer,
		client:    newLightingClient(rpcServer),
	}
	return l
}
func newLightingClient(rpcServer string) lnrpc.LightningClient {
	creds, err := credentials.NewClientTLSFromFile(defaultTLSCertPath, "")
	if err != nil {
		fatal(err)
	}

	// Create a dial options array.
	opts := []grpc.DialOption{
		grpc.WithTransportCredentials(creds),
		grpc.WithDialer(lncfg.ClientAddressDialer("10009")),
	}
	conn, err := grpc.Dial(rpcServer, opts...)
	if err != nil {
		panic(err)
	}
	return lnrpc.NewLightningClient(conn)
}

func (l *LndAPI) ListInvoices() (invoices *lnrpc.ListInvoiceResponse, err error) {
	req := &lnrpc.ListInvoiceRequest{
		PendingOnly: false,
	}

	invoices, err = l.client.ListInvoices(context.Background(), req)
	if err != nil {
		return
	}
	return
}

func (l *LndAPI) ListChannels() (channels *lnrpc.ListChannelsResponse, err error) {
	req := &lnrpc.ListChannelsRequest{
		ActiveOnly: true,
		PublicOnly: true,
	}
	return l.client.ListChannels(context.Background(), req)
}
func (l *LndAPI) AddInvoice(amount int64, secret common.Hash) (invoice *lnrpc.AddInvoiceResponse, err error) {
	req := &lnrpc.Invoice{
		Memo:      "token swap",
		RPreimage: secret[:],
		Value:     amount,
	}
	return l.client.AddInvoice(context.Background(), req)
}
func (l *LndAPI) SendPayment(amount int64, receipt string, lockSecretHsh common.Hash, finalDeltaTimeout int) (resp *lnrpc.SendResponse, err error) {
	destNode, err := hex.DecodeString(receipt)
	if err != nil {
		return
	}
	if len(destNode) != 33 {
		err = fmt.Errorf("dest node pubkey must be exactly 33 bytes, is "+
			"instead: %v", len(destNode))
		return
	}
	req := &lnrpc.SendRequest{
		Dest:           destNode,
		Amt:            amount,
		FeeLimit:       nil,
		PaymentHash:    lockSecretHsh[:],
		FinalCltvDelta: int32(finalDeltaTimeout),
	}
	PrintRespJSON(req)
	return l.client.SendPaymentSync(context.Background(), req)
}
func PrintRespJSON(resp proto.Message) {
	jsonMarshaler := &jsonpb.Marshaler{
		EmitDefaults: true,
		Indent:       "    ",
	}

	jsonStr, err := jsonMarshaler.MarshalToString(resp)
	if err != nil {
		fmt.Println("unable to decode response: ", err)
		return
	}

	fmt.Println(jsonStr)
}
