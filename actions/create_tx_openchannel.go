package actions

import (
	"crypto/rand"
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

func AddOpenButtonOnMainOfCreateTxOpenChannel(box *fyne.Container, langChangeManager *widgets.LangChangeManager) {

	title := map[string]string{"en": "Create open channel tx", "zh": "创建开启通道的交易"}

	button := langChangeManager.NewButton(title, func() {
		OpenWindowCreateTxOpenChannel(title, langChangeManager)
	})
	box.Add(button)
}

func OpenWindowCreateTxOpenChannel(title map[string]string, langChangeManager *widgets.LangChangeManager) fyne.Window {

	// 打开窗口测试
	testSize := fyne.Size{
		Width:  800,
		Height: 1100,
	}

	box := container.NewVBox()
	AddCanvasObjectCreateTxOpenChannel(title, box, langChangeManager)

	// 开启窗口
	return langChangeManager.NewWindowAndShow(title, &testSize, box)
}

func AddCanvasObjectCreateTxOpenChannel(title map[string]string, box *fyne.Container, langChangeManager *widgets.LangChangeManager) {

	page := container.NewVBox()

	page.Add(langChangeManager.NewTextWrapWordLabel(map[string]string{"en": "", "zh": ""}))

	input1 := langChangeManager.NewEntrySetPlaceHolder(map[string]string{"en": "Left address or password or private key", "zh": "左侧账户地址或密码私钥"})
	input2 := langChangeManager.NewEntrySetPlaceHolder(map[string]string{"en": "Left amount", "zh": "左侧存入金额"})
	input3 := langChangeManager.NewEntrySetPlaceHolder(map[string]string{"en": "Right address or password or private key", "zh": "右侧账户地址或密码私钥"})
	input4 := langChangeManager.NewEntrySetPlaceHolder(map[string]string{"en": "Right amount", "zh": "右侧存入金额"})
	input5 := langChangeManager.NewEntrySetPlaceHolder(map[string]string{"en": "Fee address or password or private key", "zh": "手续费支付地址或密码私钥"})
	input6 := langChangeManager.NewEntrySetPlaceHolder(map[string]string{"en": "Fee amount", "zh": "手续费支付数额"})
	input6.SetText("ㄜ8:244") // 默认手续费
	input7 := langChangeManager.NewEntrySetPlaceHolder(map[string]string{"en": "Optional: channel id", "zh": "选填：通道ID"})
	input8 := langChangeManager.NewEntrySetPlaceHolder(map[string]string{"en": "Optional: Tx timestamp", "zh": "选填：交易时间戳"})

	txbodyshow := widget.NewEntry()
	txbodyshow.MultiLine = true
	txbodyshow.Wrapping = fyne.TextWrapBreak

	button1 := langChangeManager.NewButton(map[string]string{"en": "Create open channel Tx", "zh": "确认创建开启通道的交易"}, func() {

		if input1.Text == "" || input2.Text == "" || input3.Text == "" || input4.Text == "" || input5.Text == "" || input6.Text == "" {
			langChangeManager.SetText(txbodyshow, map[string]string{"en": "Please finish fields", "zh": "请完善输入内容"})
			return
		}
		addr1, acc1 := parseAccountFromAddressOrPasswordOrPrivateKey(input1.Text)
		addr2, acc2 := parseAccountFromAddressOrPasswordOrPrivateKey(input3.Text)
		amount1, e3 := fields.NewAmountFromString(strings.Trim(input2.Text, "\n "))
		if e3 != nil {
			langChangeManager.SetText(txbodyshow, map[string]string{"en": "Left amount format error", "zh": "左侧数额格式错误"})
			return
		}
		amount2, e4 := fields.NewAmountFromString(strings.Trim(input4.Text, "\n "))
		if e4 != nil {
			langChangeManager.SetText(txbodyshow, map[string]string{"en": "Right amount format error", "zh": "右侧数额格式错误"})
			return
		}
		fee_addr, fee_acc := parseAccountFromAddressOrPasswordOrPrivateKey(input5.Text)
		fee_amt, e6 := fields.NewAmountFromString(strings.Trim(input6.Text, "\n "))
		if e6 != nil {
			langChangeManager.SetText(txbodyshow, map[string]string{"en": "Fee amount format error", "zh": "手续费数额格式错误"})
			return
		}

		// 通道id
		channelId := make([]byte, 16)
		if input7.Text == "" {
			// 随机创建id
			rand.Read(channelId)
			if channelId[0] == 0 {
				channelId[0] = 255
			}
			if channelId[15] == 0 {
				channelId[15] = 255
			}
		} else {
			idbts, e1 := hex.DecodeString(strings.Trim(input7.Text, "\n "))
			if e1 != nil || len(idbts) != 16 {
				langChangeManager.SetText(txbodyshow, map[string]string{"en": "Channel id format error", "zh": "通道ID格式错误"})
				return
			}
			channelId = idbts
		}

		// 交易时间
		usetime := time.Now().Unix()
		if len(input8.Text) > 0 {
			its, e1 := strconv.ParseInt((strings.Trim(input8.Text, "\n ")), 10, 0)
			if e1 != nil {
				langChangeManager.SetText(txbodyshow, map[string]string{"en": "Timestamp format error", "zh": "时间戳格式错误"})
				return
			}
			usetime = its
		}

		// 创建交易
		var newTrs = transactions.Transaction_2_Simple{
			Timestamp:   fields.VarUint5(usetime),
			MainAddress: *fee_addr,
			Fee:         *fee_amt,
		}

		// action
		opcAct := actions.Action_2_OpenPaymentChannel{
			ChannelId:    channelId,
			LeftAddress:  *addr1,
			LeftAmount:   *amount1,
			RightAddress: *addr2,
			RightAmount:  *amount2,
		}

		// 添加 action
		newTrs.AppendAction(&opcAct)

		// 判断签名 sign 私钥签名
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

		resEn := "Open channel transaction created successfully!" +
			"\nPlease copy the following [txbody] to sign the tx then submit transaction in online wallet:" +
			"\n\n[txhash] " + newTrs.Hash().ToHex() +
			"\n\n[timestamp] " + strconv.FormatInt(usetime, 10) +
			"\n\n[channel id] " + hex.EncodeToString(channelId) +
			"\n\n[txbody] " + txbodyhex
		resZh := "开启通道交易创建成功！" +
			"\n请复制下面 [交易体/txbody] 内容，先完成签名操作，然后去在线钱包提交交易:" +
			"\n\n[交易哈希/txhash] " + newTrs.Hash().ToHex() +
			"\n\n[时间戳/timestamp] " + strconv.FormatInt(usetime, 10) +
			"\n\n[通道ID/channel id] " + hex.EncodeToString(channelId) +
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
	page.Add(input7)
	page.Add(input8)

	page.Add(button1)
	page.Add(txbodyshow)

	card := langChangeManager.NewCardSetTitle(title, page)
	box.Add(card)

}
