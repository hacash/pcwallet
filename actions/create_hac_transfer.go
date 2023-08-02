package actions

import (
	"bytes"
	"encoding/hex"
	fyne "fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/widget"
	"github.com/hacash/core/account"
	"github.com/hacash/core/fields"
	"github.com/hacash/core/transactions"
	"github.com/hacash/pcwallet/widgets"
	"strconv"
	"time"
)

func AddOpenButtonOnMainOfCreateTransferHAC(box *fyne.Container, langChangeManager *widgets.LangChangeManager) {
	title := map[string]string{"en": "Create HAC simple transfer tx", "zh": "创建 HAC 普通转账交易"}

	button := langChangeManager.NewButton(title, func() {
		OpenWindowCreateTransferHAC(title, langChangeManager)
	})
	box.Add(button)
}

func OpenWindowCreateTransferHAC(title map[string]string, langChangeManager *widgets.LangChangeManager) fyne.Window {
	// 打开窗口测试
	testSize := fyne.Size{
		Width:  800,
		Height: 700,
	}

	box := container.NewVBox()
	AddCanvasObjectCreateTransferHAC(title, box, langChangeManager)

	// 开启窗口
	return langChangeManager.NewWindowAndShow(title, &testSize, box)
}

func AddCanvasObjectCreateTransferHAC(title map[string]string, box *fyne.Container, langChangeManager *widgets.LangChangeManager) {
	page := container.NewVBox()

	page.Add(langChangeManager.NewTextWrapWordLabel(map[string]string{"en": "Creates a normal HAC transaction. Note: the amount of transfer is actual receive amount, the transaction fee will be deducted additionally; it is suggested that the transaction fee should not be less than 0.0001 pieces; the transaction timestamp is optional, the current time will be used by default.", "zh": "创建一笔HAC普通转账交易。注意：转账数量为实到数额，交易手续费将额外扣除；交易手续费建议不低于 0.0001 枚；交易时间戳为选填，不填则默认取用当前时间。"}))

	input1 := langChangeManager.NewEntrySetPlaceHolder(map[string]string{"en": "Payment address", "zh": "输入付款地址"})
	input2 := langChangeManager.NewEntrySetPlaceHolder(map[string]string{"en": "Receive address", "zh": "输入接收地址"})
	input3 := langChangeManager.NewEntrySetPlaceHolder(map[string]string{"en": "Transfer quantity (unit: 248) or use 'ㄜ1:248' format, Example: '12.25' or 'ㄜ1225:246'", "zh": "输入转账数量 - HAC（单位：枚 - :248）也可直接使用 ㄜ1:248 格式， 例如：'12.25' or 'ㄜ1225:246'"})
	input4 := langChangeManager.NewEntrySetPlaceHolder(map[string]string{"en": "Tx Fee - HAC (unit: 248) or use 'ㄜ1:248' format, Example: '0.25' or 'ㄜ25:246'", "zh": "输入交易手续费 - HAC（单位：枚 - :248）也可直接使用 'ㄜ1:248' 格式， 例如：'0.25' or 'ㄜ25:246'"})
	input5 := langChangeManager.NewEntrySetPlaceHolder(map[string]string{"en": "PrivateKey or Password", "zh": "输入私钥或密码"})
	input6 := langChangeManager.NewEntrySetPlaceHolder(map[string]string{"en": "Optional: Tx timestamp", "zh": "选填：交易时间戳"})
	txbodyshow := widget.NewEntry()
	txbodyshow.MultiLine = true
	txbodyshow.Wrapping = fyne.TextWrapBreak

	button1 := langChangeManager.NewButton(map[string]string{"en": "Create HAC transfer Tx", "zh": "确认创建 HAC 转账交易"}, func() {
		if input1.Text == "" {
			langChangeManager.SetText(txbodyshow, map[string]string{"en": "Please input Payment address", "zh": "请输入付款地址"})
			return
		}

		addr1, e1 := fields.CheckReadableAddress(input1.Text)
		if e1 != nil {
			langChangeManager.SetText(txbodyshow, map[string]string{"en": "Payment address format error", "zh": "付款地址格式错误"})
			return
		}

		addr2, e2 := fields.CheckReadableAddress(input2.Text)
		if e2 != nil {
			langChangeManager.SetText(txbodyshow, map[string]string{"en": "Receive address format error", "zh": "接收地址格式错误"})
			return
		}

		amount, e3 := fields.NewAmountFromString(input3.Text)
		if e3 != nil {
			langChangeManager.SetText(txbodyshow, map[string]string{"en": "Transfer quantity format error", "zh": "转账数量格式错误"})
			return
		}

		if len(amount.Numeral) > 6 {
			langChangeManager.SetText(txbodyshow, map[string]string{"en": "Transfer quantity digits too long", "zh": "转账数量位数过长"})
			return
		}

		fee, e4 := fields.NewAmountFromString(input4.Text)
		if e4 != nil {
			langChangeManager.SetText(txbodyshow, map[string]string{"en": "Tx Fee format error", "zh": "交易手续费格式错误"})
			return
		}

		if len(fee.Numeral) > 2 {
			langChangeManager.SetText(txbodyshow, map[string]string{"en": "Tx Fee digits too long", "zh": "手续费位数过长"})
			return
		}

		payacc := account.GetAccountByPrivateKeyOrPassword(input5.Text)
		if bytes.Compare(payacc.Address, *addr1) != 0 {
			langChangeManager.SetText(txbodyshow, map[string]string{"en": "The private key or password does not \nmatch the payment address", "zh": "私钥或密码不匹配付款地址"})
			return
		}

		usetime := time.Now().Unix()
		if len(input6.Text) > 0 {
			its, e1 := strconv.ParseInt(input6.Text, 10, 0)
			if e1 != nil {
				langChangeManager.SetText(txbodyshow, map[string]string{"en": "Timestamp format error", "zh": "时间戳格式错误"})
				return
			}
			usetime = its
		}

		// 创建交易
		tx := transactions.CreateOneTxOfSimpleTransfer(
			payacc, *addr2, amount, fee, usetime)
		if tx == nil {
			langChangeManager.SetText(txbodyshow, map[string]string{"en": "Transaction creation failed", "zh": "交易创建失败"})
			return
		}

		// for fork or test chain ID
		MaybeForTransactionAddCheckChainID(tx)

		tx.FillTargetSign(payacc)

		// 创建成功
		txbody, e3 := tx.Serialize()
		if e3 != nil {
			langChangeManager.SetText(txbodyshow, map[string]string{"en": "Transaction creation failed", "zh": "交易创建失败"})
			return
		}

		resEn := "HAC Transfer transaction created successfully!" +
			"\nPlease copy the following [txbody] to submit transaction in online wallet:" +
			"\n\n[txhash] " + tx.Hash().ToHex() +
			"\n\n[txbody] " + hex.EncodeToString(txbody) +
			"\n\n[timestamp] " + strconv.FormatInt(usetime, 10)
		resZh := "HAC 转账交易创建成功！" +
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
