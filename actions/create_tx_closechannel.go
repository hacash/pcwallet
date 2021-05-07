package actions

import (
	"encoding/hex"
	"fyne.io/fyne"
	"fyne.io/fyne/container"
	"fyne.io/fyne/widget"
	"github.com/hacash/core/actions"
	"github.com/hacash/core/fields"
	"github.com/hacash/core/transactions"
	"github.com/hacash/pcwallet/widgets"
	"strconv"
	"strings"
	"time"
)

func AddCanvasObjectCreateTxCloseChannel(box *fyne.Container, langChangeManager *widgets.LangChangeManager) {

	box.Add(widget.NewLabel("\n\n"))

	page := container.NewVBox()

	page.Add(langChangeManager.NewTextWrapWordLabel(map[string]string{"en": "", "zh": ""}))

	input1 := langChangeManager.NewEntrySetPlaceHolder(map[string]string{"en": "Channel id", "zh": "通道ID"})
	input2 := langChangeManager.NewEntrySetPlaceHolder(map[string]string{"en": "Fee address or password or private key", "zh": "手续费支付地址或者密码私钥"})
	input3 := langChangeManager.NewEntrySetPlaceHolder(map[string]string{"en": "Fee amount", "zh": "手续费支付数额"})
	input3.SetText("ㄜ1:244") // 默认手续费
	input4 := langChangeManager.NewEntrySetPlaceHolder(map[string]string{"en": "Optional: Left password or private key", "zh": "选填：左侧账户的密码或私钥"})
	input5 := langChangeManager.NewEntrySetPlaceHolder(map[string]string{"en": "Optional: Right password or private key", "zh": "选填：右侧账户的密码或私钥"})
	input6 := langChangeManager.NewEntrySetPlaceHolder(map[string]string{"en": "Optional: Tx timestamp", "zh": "选填：交易时间戳"})

	txbodyshow := widget.NewEntry()
	txbodyshow.MultiLine = true
	txbodyshow.Wrapping = fyne.TextWrapBreak

	button1 := langChangeManager.NewButton(map[string]string{"en": "Create close channel Tx", "zh": "确认创建关闭通道的交易"}, func() {

		if input1.Text == "" || input2.Text == "" || input3.Text == "" {
			langChangeManager.SetText(txbodyshow, map[string]string{"en": "Please finish fields", "zh": "请完善输入内容"})
			return
		}

		// 通道id
		channelId := make([]byte, 16)
		idbts, e1 := hex.DecodeString(strings.Trim(input1.Text, "\n "))
		if e1 != nil || len(idbts) != 16 {
			langChangeManager.SetText(txbodyshow, map[string]string{"en": "Channel id format error", "zh": "通道ID格式错误"})
			return
		}
		channelId = idbts
		// 手续费地址和数额
		fee_addr, fee_acc := parseAccountFromAddressOrPasswordOrPrivateKey(input2.Text)
		fee_amt, e6 := fields.NewAmountFromString(strings.Trim(input3.Text, "\n "))
		if e6 != nil {
			langChangeManager.SetText(txbodyshow, map[string]string{"en": "Fee amount format error", "zh": "手续费数额格式错误"})
			return
		}

		// 左地址和右地址的签名
		_, acc1 := parseAccountFromAddressOrPasswordOrPrivateKey(input4.Text)
		_, acc2 := parseAccountFromAddressOrPasswordOrPrivateKey(input5.Text)

		// 交易时间
		usetime := time.Now().Unix()
		if len(input6.Text) > 0 {
			its, e1 := strconv.ParseInt((strings.Trim(input6.Text, "\n ")), 10, 0)
			if e1 != nil {
				langChangeManager.SetText(txbodyshow, map[string]string{"en": "Timestamp format error", "zh": "时间戳格式错误"})
				return
			}
			usetime = its
		}

		// 创建交易
		var newTrs = transactions.Transaction_2_Simple{
			Timestamp: fields.VarUint5(usetime),
			Address:   *fee_addr,
			Fee:       *fee_amt,
		}

		// action
		opcAct := actions.Action_3_ClosePaymentChannel{
			ChannelId: channelId,
		}

		// 添加 action
		newTrs.AppendAction(&opcAct)

		// 签名
		if acc1 != nil {
			newTrs.FillTargetSign(acc1)
		}
		if acc2 != nil {
			newTrs.FillTargetSign(acc2)
		}
		if fee_acc != nil {
			newTrs.FillTargetSign(fee_acc)
		}

		// 创建成功
		txbody, e3 := newTrs.Serialize()
		if e3 != nil {
			langChangeManager.SetText(txbodyshow, map[string]string{"en": "Transaction creation failed", "zh": "交易创建失败"})
			return
		}
		txbodyhex := "\n\n-------- signed txbody hex start --------\n"
		txbodyhex += hex.EncodeToString(txbody)
		txbodyhex += "\n-------- signed txbody hex  end  --------\n\n"

		resEn := "Close channel transaction created successfully!" +
			"\nPlease copy the following [txbody] to sign the tx then submit transaction in online wallet:" +
			"\n\n[txhash] " + newTrs.Hash().ToHex() +
			"\n\n[timestamp] " + strconv.FormatInt(usetime, 10) +
			"\n\n[txbody] " + txbodyhex
		resZh := "关闭通道交易创建成功！" +
			"\n请复制下面 [交易体/txbody] 内容，先完成签名操作，然后去在线钱包提交交易:" +
			"\n\n[交易哈希/txhash] " + newTrs.Hash().ToHex() +
			"\n\n[时间戳/timestamp] " + strconv.FormatInt(usetime, 10) +
			"\n\n[交易体/txbody] " + txbodyhex

		// 签名检查
		checkTxSignStatus(&newTrs, &resEn, &resZh)

		langChangeManager.SetText(txbodyshow, map[string]string{"en": resEn, "zh": resZh})
	})

	page.Add(input1)
	page.Add(input2)
	page.Add(input3)
	page.Add(input4)
	page.Add(input5)
	page.Add(input6)

	page.Add(button1)
	page.Add(txbodyshow)

	card := langChangeManager.NewCardSetTitle(map[string]string{"en": "Create close channel tx", "zh": "创建关闭通道的交易"}, page)
	box.Add(card)

}
