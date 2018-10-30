package lndapi

import (
	"bytes"
	"fmt"
	"testing"

	"github.com/SmartMeshFoundation/Photon/utils"
	"github.com/ethereum/go-ethereum/common"
)

func TestGetConnection(t *testing.T) {
	newLightingClient("localhost:10001")
}

func TestListInvoices(t *testing.T) {
	l := NewLndAPI("localhost:10001")
	invoices, err := l.ListInvoices()
	if err != nil {
		t.Error(err)
	}
	PrintRespJSON(invoices)
}

func TestLndAPI_AddInvoice(t *testing.T) {
	l := NewLndAPI("localhost:10001")
	invoice, err := l.AddInvoice(1000, utils.NewRandomHash())
	if err != nil {
		t.Error(err)
	}
	PrintRespJSON(invoice)
}

func TestLndAPI_ListChannels(t *testing.T) {
	l := NewLndAPI("localhost:10001")
	channels, err := l.ListChannels()
	if err != nil {
		t.Error(err)
	}
	PrintRespJSON(channels)
}

func TestLndAPI_SendPayment(t *testing.T) {
	alice := NewLndAPI("localhost:10001")
	bob := NewLndAPI("localhost:10002")
	secret := utils.NewRandomHash()
	invoice, err := alice.AddInvoice(1000, secret)
	if err != nil {
		t.Error(err)
		return
	}
	PrintRespJSON(invoice)
	var secretHash common.Hash
	copy(secretHash[:], invoice.RHash)
	resp, err := bob.SendPayment(1000, "022d982f42d8e68d847c26f5feac9ef34ef4bbbb70c6c0221367faf9b398f6a84f", secretHash, 200)
	if err != nil {
		t.Error(err)
		return
	}
	PrintRespJSON(resp)
	if len(resp.PaymentError) != 0 {
		t.Error(resp.PaymentError)
		return
	}
	if bytes.Compare(resp.PaymentPreimage, secret[:]) != 0 {
		t.Error(fmt.Sprintf("not equal"))
		return
	}
	invoices, err := alice.ListInvoices()
	PrintRespJSON(invoices)
	found := false
	for _, i := range invoices.Invoices {
		if bytes.Compare(i.RHash, secretHash[:]) == 0 {
			found = true
			if !i.Settled {
				t.Error("should settled")
				return
			}
		}
	}
	if !found {
		t.Error("not found")
	}
}

func TestLndAPI_ChannelBalance(t *testing.T) {
	l := NewLndAPI("localhost:10001")
	resp, err := l.ChannelBalance()
	if err != nil {
		t.Error(err)
	}
	PrintRespJSON(resp)
}
