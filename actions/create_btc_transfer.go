package actions

import (
	"bytes"
	"encoding/hex"
	"fyne.io/fyne"
	"fyne.io/fyne/container"
	"fyne.io/fyne/widget"
	"github.com/hacash/core/account"
	"github.com/hacash/core/fields"
	"github.com/hacash/core/transactions"
	"github.com/hacash/pcwallet/widgets"
	"strconv"
	"time"
)

func AddOpenButtonOnMainOfCreateTransferBTC(box *fyne.Container, langChangeManager *widgets.LangChangeManager) {

	title := map[string]string{"en": "Create BTC transfer tx", "zh": "创建 BTC 转账交易"}

	button := langChangeManager.NewButton(title, func() {
		OpenWindowCreateTransferBTC(title, langChangeManager)
	})
	box.Add(button)
}

func OpenWindowCreateTransferBTC(title map[string]string, langChangeManager *widgets.LangChangeManager) fyne.Window {

	// 打开窗口测试
	testSize := fyne.Size{
		Width:  800,
		Height: 700,
	}

	box := container.NewVBox()
	AddCanvasObjectCreateTransferBTC(title, box, langChangeManager)

	// 开启窗口
	return langChangeManager.NewWindowAndShow(title, &testSize, box)
}

func AddCanvasObjectCreateTransferBTC(title map[string]string, box *fyne.Container, langChangeManager *widgets.LangChangeManager) {

	page := container.NewVBox()

	page.Add(langChangeManager.NewTextWrapWordLabel(map[string]string{"en": "Creates a normal BTC transaction. Note: the amount of transfer is actual receive amount, the transaction fee will be use HAC deducted additionally; it is suggested that the transaction fee should not be less than 0.0001 pieces; the transaction timestamp is optional, the current time will be used by default. BTC unit:", "zh": "创建一笔 BTC 普通转账交易。注意：转账数量为实到数额，交易手续费将额外扣除 HAC；交易手续费建议不低于 0.0001 枚；交易时间戳为选填，不填则默认取用当前时间。比特币单位:"}))
	page.Add(widgets.NewTextWrapWordLabel("1 BTC = 100000000 SAT_satoshi"))

	input1 := langChangeManager.NewEntrySetPlaceHolder(map[string]string{"en": "BTC Payment address", "zh": "输入 BTC 付款地址"})
	input2 := langChangeManager.NewEntrySetPlaceHolder(map[string]string{"en": "BTC Receive address", "zh": "输入 BTC 接收地址"})
	input3 := langChangeManager.NewEntrySetPlaceHolder(map[string]string{"en": "Transfer quantity (unit: SAT satoshi)", "zh": "输入转账数量（单位：SAT satoshi）"})
	input4 := langChangeManager.NewEntrySetPlaceHolder(map[string]string{"en": "Payment PrivateKey or Password", "zh": "输入付款私钥或密码"})
	input5 := langChangeManager.NewEntrySetPlaceHolder(map[string]string{"en": "Tx Fee PrivateKey or Password", "zh": "输入交易手续费私钥或密码"})
	input6 := langChangeManager.NewEntrySetPlaceHolder(map[string]string{"en": "Tx Fee - HAC (unit: 248) or use 'ㄜ1:248' format", "zh": "输入交易手续费 - HAC（单位：枚 - :248）也可直接使用 ㄜ1:248 格式"})
	input7 := langChangeManager.NewEntrySetPlaceHolder(map[string]string{"en": "Optional: Tx timestamp", "zh": "选填：交易时间戳"})

	txbodyshow := widget.NewEntry()
	txbodyshow.MultiLine = true
	txbodyshow.Wrapping = fyne.TextWrapBreak

	button1 := langChangeManager.NewButton(map[string]string{"en": "Create BTC transfer Tx", "zh": "确认创建 BTC 转账交易"}, func() {
		if input1.Text == "" {
			langChangeManager.SetText(txbodyshow, map[string]string{"en": "Please input BTC Payment address", "zh": "请输入 BTC 付款地址"})
			return
		}
		addr1, e1 := fields.CheckReadableAddress(input1.Text)
		if e1 != nil {
			langChangeManager.SetText(txbodyshow, map[string]string{"en": "BTC Payment address format error", "zh": "BTC 付款地址格式错误"})
			return
		}
		addr2, e2 := fields.CheckReadableAddress(input2.Text)
		if e2 != nil {
			langChangeManager.SetText(txbodyshow, map[string]string{"en": "BTC Receive address format error", "zh": "BTC 接收地址格式错误"})
			return
		}
		amount, e3 := strconv.ParseUint(input3.Text, 10, 0)
		if e3 != nil {
			langChangeManager.SetText(txbodyshow, map[string]string{"en": "BTC Transfer quantity format error", "zh": "转账数量格式错误"})
			return
		}

		payacc := account.GetAccountByPrivateKeyOrPassword(input4.Text)
		if bytes.Compare(payacc.Address, *addr1) != 0 {
			langChangeManager.SetText(txbodyshow, map[string]string{"en": "The private key or password does not \nmatch the payment address", "zh": "私钥或密码不匹配付款地址"})
			return
		}
		feeacc := account.GetAccountByPrivateKeyOrPassword(input5.Text)
		fee, e4 := fields.NewAmountFromString(input6.Text)
		if e4 != nil {
			langChangeManager.SetText(txbodyshow, map[string]string{"en": "Tx Fee format error", "zh": "交易手续费格式错误"})
			txbodyshow.SetText(" / ")
			return
		}
		if len(fee.Numeral) > 2 {
			langChangeManager.SetText(txbodyshow, map[string]string{"en": "Tx Fee digits too long", "zh": "手续费数值位数过长"})
			txbodyshow.SetText(" / ")
			return
		}
		usetime := time.Now().Unix()
		if len(input7.Text) > 0 {
			its, e1 := strconv.ParseInt(input7.Text, 10, 0)
			if e1 != nil {
				langChangeManager.SetText(txbodyshow, map[string]string{"en": "Timestamp format error", "zh": "时间戳格式错误"})
				return
			}
			usetime = its
		}

		// 创建交易
		tx, e0 := transactions.CreateOneTxOfBTCTransfer(payacc, *addr2, amount, feeacc, fee, usetime)
		if e0 != nil {
			langChangeManager.SetText(txbodyshow, map[string]string{"en": "Transaction creation failed: \n\n" + e0.Error(), "zh": "交易创建失败\n\n" + e0.Error()})
			return
		}

		if tx == nil {
			langChangeManager.SetText(txbodyshow, map[string]string{"en": "Transaction creation failed", "zh": "交易创建失败"})
			return
		}

		// 创建成功
		txbody, e3 := tx.Serialize()
		if e3 != nil {
			langChangeManager.SetText(txbodyshow, map[string]string{"en": "Transaction creation failed", "zh": "交易创建失败"})
			return
		}

		resEn := "BTC Transfer transaction created successfully!" +
			"\n Please copy the following [txbody] \nto submit transaction in online wallet:" +
			"\n\n[txhash] " + tx.Hash().ToHex() +
			"\n\n[txbody] " + hex.EncodeToString(txbody) +
			"\n\n[timestamp] " + strconv.FormatInt(usetime, 10)
		resZh := "BTC 转账交易创建成功!" +
			"\n请复制下面 [交易体/txbody] 后面的内容去在线钱包提交交易:" +
			"\n\n[交易哈希/txhash] " + tx.Hash().ToHex() +
			"\n\n[交易体/txbody] " + hex.EncodeToString(txbody) +
			"\n\n[时间戳/timestamp] " + strconv.FormatInt(usetime, 10)

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

	card := langChangeManager.NewCardSetTitle(title, page)
	box.Add(card)

}
