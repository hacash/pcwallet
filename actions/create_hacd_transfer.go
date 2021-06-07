package actions

import (
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

func AddOpenButtonOnMainOfCreateTransferHACD(box *fyne.Container, langChangeManager *widgets.LangChangeManager) {

	title := map[string]string{"en": "Create HACD transfer tx", "zh": "创建 HACD 转账交易"}

	button := langChangeManager.NewButton(title, func() {
		OpenWindowCreateTransferHACD(title, langChangeManager)
	})
	box.Add(button)
}

func OpenWindowCreateTransferHACD(title map[string]string, langChangeManager *widgets.LangChangeManager) fyne.Window {

	// 打开窗口测试
	testSize := fyne.Size{
		Width:  800,
		Height: 800,
	}

	box := container.NewVBox()
	AddCanvasObjectCreateTransferHACD(title, box, langChangeManager)

	// 开启窗口
	return langChangeManager.NewWindowAndShow(title, &testSize, box)
}

func AddCanvasObjectCreateTransferHACD(title map[string]string, box *fyne.Container, langChangeManager *widgets.LangChangeManager) {

	page := container.NewVBox()

	page.Add(langChangeManager.NewTextWrapWordLabel(map[string]string{"en": "Creates a multiple HACD transaction. Note: the transaction fee will be use HAC deducted additionally; it is suggested that the transaction fee should not be less than 0.002 pieces; the transaction timestamp is optional, the current time will be used by default.", "zh": "创建一笔 HACD 批量转账交易。注意：交易手续费将额外扣除HAC；交易手续费建议不低于 0.002 枚；交易时间戳为选填，不填则默认取用当前时间。"}))

	input1 := langChangeManager.NewEntrySetPlaceHolder(map[string]string{"en": "HACD name list split by comma", "zh": "输入逗号隔开的钻石字面值列表"})
	input1.Wrapping = fyne.TextWrapWord
	input1.MultiLine = true

	input2 := langChangeManager.NewEntrySetPlaceHolder(map[string]string{"en": "HACD Receive address", "zh": "输入钻石接收地址"})
	input3 := langChangeManager.NewEntrySetPlaceHolder(map[string]string{"en": "HACD Payment PrivateKey or Password", "zh": "输入钻石转出地址私钥或密码"})
	input4 := langChangeManager.NewEntrySetPlaceHolder(map[string]string{"en": "Tx Fee PrivateKey or Password", "zh": "输入交易手续费地址私钥或密码"})
	input5 := langChangeManager.NewEntrySetPlaceHolder(map[string]string{"en": "Tx Fee - HAC (unit: 248) or use 'ㄜ1:248' format, Example: '0.25' or 'ㄜ25:246'", "zh": "输入交易手续费 - HAC（单位：枚 - :248）也可直接使用 'ㄜ1:248' 格式， 例如：'0.25' or 'ㄜ25:246'"})
	input6 := langChangeManager.NewEntrySetPlaceHolder(map[string]string{"en": "Optional: Tx timestamp", "zh": "选填：交易时间戳"})

	txbodyshow := widget.NewEntry()
	txbodyshow.MultiLine = true
	txbodyshow.Wrapping = fyne.TextWrapBreak

	button1 := langChangeManager.NewButton(map[string]string{"en": "Create HACD transfer Tx", "zh": "确认创建 HACD 钻石转账交易"}, func() {
		if input1.Text == "" {
			langChangeManager.SetText(txbodyshow, map[string]string{"en": "Please HACD names split by comma", "zh": "请输入钻石列表"})
			return
		}
		diamondstrs := input1.Text
		toaddr, e1 := fields.CheckReadableAddress(input2.Text)
		if e1 != nil {
			langChangeManager.SetText(txbodyshow, map[string]string{"en": "Receive address format error", "zh": "接收地址格式错误"})
			return
		}
		payacc := account.GetAccountByPrivateKeyOrPassword(input3.Text)
		feeacc := account.GetAccountByPrivateKeyOrPassword(input4.Text)
		fee, e4 := fields.NewAmountFromStringUnsafe(input5.Text)
		if e4 != nil {
			langChangeManager.SetText(txbodyshow, map[string]string{"en": "Tx Fee format error", "zh": "交易手续费格式错误"})
			return
		}
		if len(fee.Numeral) > 2 {
			langChangeManager.SetText(txbodyshow, map[string]string{"en": "Tx Fee digits too long", "zh": "手续费位数过长"})
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
		tx, e0 := transactions.CreateOneTxOfOutfeeQuantityHACDTransfer(
			payacc, *toaddr, diamondstrs, feeacc, fee, usetime)
		if e0 != nil {
			langChangeManager.SetText(txbodyshow, map[string]string{"en": "Transaction creation failed: \n\n" + e0.Error(), "zh": "交易创建失败: \n\n" + e0.Error()})
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

		resEn := "HACD Transfer transaction created successfully!" +
			"\nPlease copy the following [txbody] to submit transaction in online wallet:" +
			"\n\n[txhash] " + tx.Hash().ToHex() +
			"\n\n[txbody] " + hex.EncodeToString(txbody) +
			"\n\n[timestamp] " + strconv.FormatInt(usetime, 10)
		resZh := "HACD 转账交易创建成功！" +
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
