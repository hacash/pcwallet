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

func AddCanvasObjectCreateTransferHACD(box *fyne.Container) {

	box.Add(widget.NewLabel("\n\n\n\n"))

	page := container.NewVBox()

	page.Add(widgets.NewTextWrapWordLabel("创建一笔 HACD 批量转账交易。注意：交易手续费将额外扣除HAC；交易手续费建议不低于 0.002 枚；交易时间戳为选填，不填则默认取用当前时间。"))
	page.Add(widgets.NewTextWrapWordLabel("Creates a multiple HACD transaction. Note: the transaction fee will be use HAC deducted additionally; it is suggested that the transaction fee should not be less than 0.002 pieces; the transaction timestamp is optional, the current time will be used by default."))

	input1 := widget.NewEntry()
	input1.PlaceHolder = "这里输入逗号隔开的钻石字面值列表 / HACD name list split by comma"

	input2 := widget.NewEntry()
	input2.PlaceHolder = "这里输入钻石接收地址 / HACD Receive address"

	input3 := widget.NewEntry()
	input3.PlaceHolder = "这里输入钻石转出地址私钥或密码 / HACD Payment PrivateKey or Password"

	input4 := widget.NewEntry()
	input4.PlaceHolder = "这里输入交易手续费地址私钥或密码 / Tx Fee PrivateKey or Password"

	input5 := widget.NewEntry()
	input5.PlaceHolder = "这里输入交易手续费 - HAC（单位：枚 - :248） / Tx Fee - HAC (unit: 248)"

	input6 := widget.NewEntry()
	input6.PlaceHolder = "选填：交易时间戳 / Optional: Tx timestamp"

	txbodyshow := widget.NewEntry()
	txbodyshow.MultiLine = true
	txbodyshow.Wrapping = fyne.TextWrapBreak

	button1 := widget.NewButton("确认创建钻石交易 / Create HACD Tx", func() {
		if input1.Text == "" {
			txbodyshow.SetText("请输入钻石列表 / Please HACD names split by comma")
			return
		}
		diamondstrs := input1.Text
		toaddr, e1 := fields.CheckReadableAddress(input2.Text)
		if e1 != nil {
			txbodyshow.SetText("接收地址格式错误 / Receive address format error")
			return
		}
		payacc := account.GetAccountByPrivateKeyOrPassword(input3.Text)
		feeacc := account.GetAccountByPrivateKeyOrPassword(input4.Text)
		fee, e4 := fields.NewAmountFromMeiString(input5.Text)
		if e4 != nil {
			txbodyshow.SetText("交易手续费格式错误 / Tx Fee format error")
			return
		}
		if len(fee.Numeral) > 2 {
			txbodyshow.SetText("手续费位数过长 / Tx Fee  digits too long")
			return
		}
		usetime := time.Now().Unix()
		if len(input6.Text) > 0 {
			its, e1 := strconv.ParseInt(input6.Text, 10, 0)
			if e1 != nil {
				txbodyshow.SetText("时间戳格式错误 / Timestamp format error")
				return
			}
			usetime = its
		}

		// 创建交易
		tx, e0 := transactions.CreateOneTxOfOutfeeQuantityHACDTransfer(
			payacc, *toaddr, diamondstrs, feeacc, fee, usetime)
		if e0 != nil {
			txbodyshow.SetText("交易创建失败 / Transaction creation failed: \n\n" + e0.Error())
			return
		}
		if tx == nil {
			txbodyshow.SetText("交易创建失败 / Transaction creation failed")
			return
		}

		// 创建成功
		txbody, e3 := tx.Serialize()
		if e3 != nil {
			txbodyshow.SetText("交易创建失败 / Transaction creation failed")
			return
		}
		txbodyshow.SetText("HACD 转账交易创建成功！ / HACD Transfer transaction created successfully!" +
			"\n请复制下面 [交易体/txbody] 后面的内容去在线钱包提交交易 / Please copy the following [txbody] to submit transaction in online wallet:" +
			"\n\n[交易哈希/txhash] " + tx.Hash().ToHex() +
			"\n\n[交易体/txbody] " + hex.EncodeToString(txbody) +
			"\n\n[时间戳/timestamp] " + strconv.FormatInt(usetime, 10))
	})

	page.Add(input1)
	page.Add(input2)
	page.Add(input3)
	page.Add(input4)
	page.Add(input5)
	page.Add(input6)

	page.Add(button1)
	page.Add(txbodyshow)

	card := widget.NewCard("创建 HACD 转账交易 / Create HACD transfer tx", "", page)
	box.Add(card)

}
