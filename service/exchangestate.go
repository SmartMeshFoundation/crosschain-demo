package service

import (
	"bytes"

	"math/big"

	"time"

	"log"

	"encoding/json"

	"fmt"

	"errors"

	"github.com/SmartMeshFoundation/Atmosphere/lndapi"
	"github.com/SmartMeshFoundation/Atmosphere/smapi"
	"github.com/SmartMeshFoundation/Atmosphere/utils"
	"github.com/ethereum/go-ethereum/common"
	"github.com/lightningnetwork/lnd/lnrpc"
)

// Partition : Exchange交易参与者
type Partition struct {
	SmAddress  common.Address
	LndAddress string // 即闪电网络公钥
}

// ExchangeState :
// 核心交易状态,包含一笔完成exchange交易
type ExchangeState struct {
	/*
		交易数据
	*/
	SmSender       *Partition     // 雷电网络上转出资金的参与者
	SmTokenAddress common.Address // 雷电网络上交易的token地址
	SmAmount       *big.Int       // 雷电网络上流通的金额

	LndSender *Partition // 闪电网络上转出资金的参与者
	LndAmount *big.Int   // 闪电网络上流通的金额

	LockSecretHash common.Hash // 密码hash,同时也是一笔交易的唯一ID,我们认为一个密码不能同时进行多笔exchange交易
	Secret         common.Hash // 密码

	/*
		handler
	*/
	SmAPI  *smapi.SmAPI
	LndAPI *lndapi.LndAPI
}

// RegisterExchangeStateBySmSender :
// 在雷电网络上转出资产的参与者注册
func RegisterExchangeStateBySmSender(partnerSmAddress, smTokenAddress common.Address, smAmount *big.Int, lndAmount *big.Int, secret common.Hash) (state *ExchangeState, err error) {
	state = new(ExchangeState)
	state.SmSender, state.LndSender = new(Partition), new(Partition)
	state.LndAPI = LndAPI
	state.SmAPI = SmAPI
	// 1. 注册信息
	// 自己的sm地址从本地节点获取
	state.SmSender.SmAddress = state.SmAPI.AccountAddress
	state.LndSender.SmAddress = partnerSmAddress
	state.SmTokenAddress = smTokenAddress
	state.SmAmount = smAmount
	state.LndAmount = lndAmount
	state.Secret = secret
	state.LockSecretHash = utils.ShaSecret(secret.Bytes())
	log.Println("ExchangeState register done ...")
	// 2. 异步调用SmAPI发送MediatedTransfer
	state.sendTransferOnSmartraiden()
	log.Println("ExchangeState send transfer on smartraiden done ...")
	// 3. 发送完成后SmSender调用AddInvoice发送paymentReq
	err = state.sendAddInvoice()
	if err != nil {
		log.Println("ExchangeState send paymentReq on lnd FAIL !!!")
		return
	}
	log.Println("ExchangeState send paymentReq on lnd done ...")
	// 3. 启动轮询,调用LndAPI查询接收到的Payment
	log.Println("ExchangeState start waiting for payment on lnd ...")
	err = state.waitingTransferOnLnd()
	if err != nil {
		log.Println("ExchangeState start waiting for payment on lnd FAIL !!!")
		return
	}
	return
}

func (es *ExchangeState) sendTransferOnSmartraiden() {
	es.SmAPI.SendTransferWithSecretAsync(es.LndSender.SmAddress, es.SmTokenAddress, es.SmAmount, es.Secret)
}

func (es *ExchangeState) sendAddInvoice() (err error) {
	resp, err := es.LndAPI.AddInvoice(es.LndAmount.Int64(), es.Secret)
	log.Println("add invoice on lnd ...")
	lndapi.PrintRespJSON(resp)
	if err != nil {
		// 交易取消???
		return
	}
	return
}

func (es *ExchangeState) waitingTransferOnLnd() (err error) {
	var invoices *lnrpc.ListInvoiceResponse
	for {
		invoices, err = es.LndAPI.ListInvoices()
		if err != nil {
			// 重试还是取消交易 ???
			continue
		}
		var invoice *lnrpc.Invoice
		for _, i := range invoices.Invoices {
			if bytes.Compare(i.RHash, es.LockSecretHash.Bytes()) == 0 {
				invoice = i
			}
		}
		// 找到
		if invoice != nil {
			if invoice.Settled {
				log.Println("Receive transfer on lnd ...")
				lndapi.PrintRespJSON(invoice)
				// 已经settle
				// 校验金额
				if invoice.Value != es.LndAmount.Int64() {
					// 金额不匹配,交易失败,通知上层
					return errors.New("lnd_amount does't match")
				}
				// 允许在雷电网络上发送密码,并结束轮询
				err = es.SmAPI.AllowRevealSecret(es.SmTokenAddress, es.LockSecretHash)
				if err != nil {
					// 这里出错,无关紧要,LndSender已经拿到密码了,可以自己注册
					log.Println("send secret on smartraiden failed !!!")
				}
				// 交易成功,结束轮询,通知上层
				return
			} else if invoice.CreationDate+invoice.Expiry <= int64(time.Now().Second()) {
				// 已经过期,不可能再成功,结束轮询,通知上层
				return errors.New("Receive transfer unsettled on lnd and already expire, must be fail")
			} else {
				// 继续等待交易settle,什么都不用做
			}
		} else {
			fmt.Println("waiting for invoice on lnd ...")
		}
		// 轮询间隔1秒
		time.Sleep(5 * time.Second)
	}
}

// RegisterExchangeStateByLndSender :
// 在闪电网络上转出资产的参与者注册
func RegisterExchangeStateByLndSender(targetLndAddress string, smTokenAddress common.Address, smAmount *big.Int, lndAmount *big.Int, lockSecretHash common.Hash) (state *ExchangeState, err error) {
	state = new(ExchangeState)
	state.SmSender, state.LndSender = new(Partition), new(Partition)
	state.LndAPI = LndAPI
	state.SmAPI = SmAPI
	// 1. 注册信息
	// 自己的sm地址从本地节点获取
	state.LndSender.SmAddress = state.SmAPI.AccountAddress
	state.SmSender.LndAddress = targetLndAddress
	state.SmTokenAddress = smTokenAddress
	state.SmAmount = smAmount
	state.LndAmount = lndAmount
	state.LockSecretHash = lockSecretHash
	// 2. 启动轮询,等待在smartraiden上收到锁为lockSecretHash的交易
	log.Println("ExchangeState start waiting for transfer on smartraiden ...")
	err = state.waitingTransferOnSmartraiden()
	return
}

func (es *ExchangeState) waitingTransferOnSmartraiden() (err error) {
	type SmTransferDataResponse struct {
		Initiator      string   `json:"initiator_address"`
		Target         string   `json:"target_address"`
		Token          string   `json:"token_address"`
		Amount         *big.Int `json:"amount"`
		Secret         string   `json:"secret"`
		LockSecretHash string   `json:"lock_secret_hash"`
		Expiration     int64    `json:"expiration"`
		Fee            *big.Int `json:"fee"`
		IsDirect       bool     `json:"is_direct"`
	}
	period := time.Second * 5
	var jsonStr string
	for {
		jsonStr, err = es.SmAPI.GetUnfinishedReceivedTransfer(es.SmTokenAddress, es.LockSecretHash)
		if err != nil || jsonStr == "null" {
			//出错或没查到,继续
			log.Println("waiting for mediated transfer on smartraiden...")
			time.Sleep(period)
			continue
		}
		log.Printf("recevice mediated transfer on smartraiden : %s\n", jsonStr)
		var smTransferData SmTransferDataResponse
		err = json.Unmarshal([]byte(jsonStr), &smTransferData)
		if err != nil {
			log.Println("decode fail")
			return errors.New("mediated transfer on smartraiden decode failed")
		}
		// 校验金额及target是否为自己
		if smTransferData.Amount.Cmp(es.SmAmount) != 0 || smTransferData.Target != es.LndSender.SmAddress.String() {
			// 交易失败,通知上层
			return errors.New("mediated transfer data does't match")
		}
		// 在lnd上发送交易
		var resp *lnrpc.SendResponse
		resp, err = es.LndAPI.SendPayment(es.LndAmount.Int64(), es.SmSender.LndAddress, es.LockSecretHash, 200)
		log.Println("send payment on lnd...")
		lndapi.PrintRespJSON(resp)
		if err != nil {
			// 重试还是失败 ???
			log.Println("send payment on lnd err", err)
		}
		if len(resp.PaymentError) != 0 {
			// 发送失败
			log.Println("send payment on lnd err")
		}
		fmt.Println("send payment on lnd success")
		return
	}
}
