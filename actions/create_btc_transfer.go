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

func AddCanvasObjectCreateTransferBTC(box *fyne.Container) {

	box.Add(widget.NewLabel("\n\n\n\n"))

	page := container.NewVBox()

	page.Add(widgets.NewTextWrapWordLabel("创建一笔 BTC 普通转账交易。注意：转账数量为实到数额，交易手续费将额外扣除 HAC；交易手续费建议不低于 0.0001 枚；交易时间戳为选填，不填则默认取用当前时间。"))
	page.Add(widgets.NewTextWrapWordLabel("Creates a normal BTC transaction. Note: the amount of transfer is actual receive amount, the transaction fee will be use HAC deducted additionally; it is suggested that the transaction fee should not be less than 0.0001 pieces; the transaction timestamp is optional, the current time will be used by default."))
	page.Add(widgets.NewTextWrapWordLabel("[比特币单位/BTC unit]: 1 BTC = 100000000 SAT_satoshi"))

	input1 := widget.NewEntry()
	input1.PlaceHolder = "这里输入BTC付款地址 / BTC Payment address"

	input2 := widget.NewEntry()
	input2.PlaceHolder = "这里输入BTC接收地址 / BTC Receive address"

	input3 := widget.NewEntry()
	input3.PlaceHolder = "这里输入转账数量（单位：SAT_satoshi） / Transfer quantity (unit: SAT_satoshi)"

	input4 := widget.NewEntry()
	input4.PlaceHolder = "这里输入交易手续费私钥或密码 / BTC Payment PrivateKey or Password"

	input5 := widget.NewEntry()
	input5.PlaceHolder = "这里输入交易手续费私钥或密码 / Tx Fee PrivateKey or Password"

	input6 := widget.NewEntry()
	input6.PlaceHolder = "这里输入交易手续费 - HAC（单位：枚 - :248） / Tx Fee - HAC (unit: 248)"

	input7 := widget.NewEntry()
	input7.PlaceHolder = "选填：交易时间戳 / Optional: Tx timestamp"

	txbodyshow := widget.NewEntry()
	txbodyshow.MultiLine = true
	txbodyshow.Wrapping = fyne.TextWrapBreak

	button1 := widget.NewButton("确认创建 BTC 交易 / Create BTC Tx", func() {
		if input1.Text == "" {
			txbodyshow.SetText("请输入输入BTC付款地址 / Please input BTC Payment address")
			return
		}
		addr1, e1 := fields.CheckReadableAddress(input1.Text)
		if e1 != nil {
			txbodyshow.SetText("BTC付款地址格式错误 / BTC Payment address format error")
			return
		}
		addr2, e2 := fields.CheckReadableAddress(input2.Text)
		if e2 != nil {
			txbodyshow.SetText("BTC接收地址格式错误 / BTC Receive address format error")
			return
		}
		amount, e3 := strconv.ParseUint(input3.Text, 10, 0)
		if e3 != nil {
			txbodyshow.SetText("转账数量格式错误 / BTC Transfer quantity format error")
			return
		}

		payacc := account.GetAccountByPrivateKeyOrPassword(input4.Text)
		if bytes.Compare(payacc.Address, *addr1) != 0 {
			txbodyshow.SetText("私钥或密码不匹配付款地址 /\n The private key or password does not \nmatch the payment address")
			return
		}
		feeacc := account.GetAccountByPrivateKeyOrPassword(input5.Text)
		fee, e4 := fields.NewAmountFromString(input6.Text)
		if e4 != nil {
			txbodyshow.SetText("交易手续费格式错误 / Tx Fee format error")
			return
		}
		if len(fee.Numeral) > 2 {
			txbodyshow.SetText("手续费位数过长 / Tx Fee digits too long")
			return
		}
		usetime := time.Now().Unix()
		if len(input7.Text) > 0 {
			its, e1 := strconv.ParseInt(input7.Text, 10, 0)
			if e1 != nil {
				txbodyshow.SetText("时间戳格式错误 / Timestamp format error")
				return
			}
			usetime = its
		}

		// 创建交易
		tx, e0 := transactions.CreateOneTxOfBTCTransfer(payacc, *addr2, amount, feeacc, fee, usetime)
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
		txbodyshow.SetText("BTC 转账交易创建成功！ / BTC Transfer transaction created successfully!" +
			"\n请复制下面 [交易体/txbody] 后面的内容去在线钱包提交交易 /\n Please copy the following [txbody] \nto submit transaction in online wallet:" +
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

	card := widget.NewCard("创建 BTC 转账交易 / Create BTC transfer tx", "", page)
	box.Add(card)

}
