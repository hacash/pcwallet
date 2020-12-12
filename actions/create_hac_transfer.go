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
	"strconv"
	"time"
)

func AddCanvasObjectCreateTransferHAC(box *fyne.Container) {

	box.Add(widget.NewLabel("\n"))

	page := container.NewVBox()

	page.Add(widget.NewLabel("创建一笔HAC普通转账交易。注意：转账数量为实到数额，交易手续费将额外扣除；\n交易手续费建议不低于 0.0001 枚；交易时间戳为选填，不填则默认取用当前时间。\nCreates a normal HAC transaction. Note: the amount of transfer is \nactual receive amount, the transaction fee will be deducted additionally; \nit is suggested that the transaction fee should not be less than 0.001 pieces; \nthe transaction timestamp is optional, the current time will be used by default."))

	input1 := widget.NewEntry()
	input1.PlaceHolder = "这里输入付款地址 / Payment address"

	input2 := widget.NewEntry()
	input2.PlaceHolder = "这里输入接收地址 / Receive address"

	input3 := widget.NewEntry()
	input3.PlaceHolder = "这里输入转账数量（单位：枚 - :248） / Transfer quantity (unit: 248)"

	input4 := widget.NewEntry()
	input4.PlaceHolder = "这里输入交易手续费（单位：枚 - :248） / Tx Fee (unit: 248)"

	input5 := widget.NewEntry()
	input5.PlaceHolder = "这里输入私钥或密码 / PrivateKey or Password"

	input6 := widget.NewEntry()
	input6.PlaceHolder = "选填：交易时间戳 / Optional: Tx timestamp"

	txbodyshow := widget.NewEntry()
	txbodyshow.MultiLine = true
	txbodyshow.Wrapping = fyne.TextWrapBreak

	button1 := widget.NewButton("确认创建交易 / Create Tx", func() {
		if input1.Text == "" {
			txbodyshow.SetText("请输入输入付款地址 / Please input Payment address")
			return
		}
		addr1, e1 := fields.CheckReadableAddress(input1.Text)
		if e1 != nil {
			txbodyshow.SetText("付款地址格式错误 / Payment address format error")
			return
		}
		addr2, e2 := fields.CheckReadableAddress(input2.Text)
		if e2 != nil {
			txbodyshow.SetText("接收地址格式错误 / Receive address format error")
			return
		}
		amount, e3 := fields.NewAmountFromMeiString(input3.Text)
		if e3 != nil {
			txbodyshow.SetText("转账数量格式错误 / Transfer quantity format error")
			return
		}
		if len(amount.Numeral) > 6 {
			txbodyshow.SetText("转账数量位数过长 / Transfer quantity digits too long")
			return
		}
		fee, e4 := fields.NewAmountFromMeiString(input4.Text)
		if e4 != nil {
			txbodyshow.SetText("交易手续费格式错误 / Tx Fee format error")
			return
		}
		if len(fee.Numeral) > 2 {
			txbodyshow.SetText("手续费位数过长 / Tx Fee  digits too long")
			return
		}
		payacc := account.GetAccountByPrivateKeyOrPassword(input5.Text)
		if bytes.Compare(payacc.Address, *addr1) != 0 {
			txbodyshow.SetText("私钥或密码不匹配付款地址 /\n The private key or password does not \nmatch the payment address")
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
		tx := transactions.CreateOneTxOfSimpleTransfer(
			payacc, *addr2, amount, fee, usetime)

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
		txbodyshow.SetText("转账交易创建成功！ / ransfer transaction created successfully!" +
			"\n请复制下面 [交易体/txbody] 后面的内容去在线钱包提交交易 /\n Please copy the following [txbody] \nto submit transaction in online wallet:" +
			"\n\n[交易哈希/txhash] " + tx.Hash().ToHex() +
			"\n\n[交易体/txbody] " + hex.EncodeToString(txbody))

	})

	page.Add(input1)
	page.Add(input2)
	page.Add(input3)
	page.Add(input4)
	page.Add(input5)
	page.Add(input6)

	page.Add(button1)
	page.Add(txbodyshow)

	card := widget.NewCard("创建HAC普通转账交易 / Create HAC simple transfer tx", "", page)
	box.Add(card)

}
