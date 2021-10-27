package actions

import (
	"encoding/hex"
	"fmt"
	"fyne.io/fyne"
	"fyne.io/fyne/container"
	"fyne.io/fyne/widget"
	"github.com/hacash/core/account"
	"github.com/hacash/core/actions"
	"github.com/hacash/core/channel"
	"github.com/hacash/core/fields"
	"github.com/hacash/core/stores"
	"github.com/hacash/core/transactions"
	"github.com/hacash/pcwallet/widgets"
	"strconv"
	"strings"
	"time"
)

// 检查通道链支付对账票据
func AddOpenButtonOnMainOfCheckChannelPaymentBill(box *fyne.Container, langChangeManager *widgets.LangChangeManager) {

	title := map[string]string{"en": "Check channel bill and create arbitration tx ", "zh": "检查通道链票据及创建仲裁交易"}

	button := langChangeManager.NewButton(title, func() {
		OpenWindowCheckChannelPaymentBill(title, langChangeManager)
	})
	box.Add(button)
}

func OpenWindowCheckChannelPaymentBill(title map[string]string, langChangeManager *widgets.LangChangeManager) fyne.Window {

	// 打开窗口测试
	testSize := fyne.Size{
		Width:  800,
		Height: 1000,
	}

	box := container.NewVBox()
	AddCanvasObjectCheckChannelPaymentBill(title, box, langChangeManager)

	// 开启窗口
	return langChangeManager.NewWindowAndShow(title, &testSize, box)
}

func AddCanvasObjectCheckChannelPaymentBill(title map[string]string, box *fyne.Container, langChangeManager *widgets.LangChangeManager) {

	page := container.NewVBox()

	page.Add(langChangeManager.NewTextWrapWordLabel(map[string]string{"en": "View the operation items, conditions and contents of a transaction, and check whether the signature is completed", "zh": "查看一笔交易具体包含的操作项目、条件和内容，以及检查是否完成签名"}))

	input1 := langChangeManager.NewEntrySetPlaceHolder(map[string]string{"en": "Reconciliation or payment bill hex string", "zh": "请输入对账或支付票据数据（bill hex）"})
	input1.Wrapping = fyne.TextWrapWord
	input1.MultiLine = true
	//input1.SetText("02cdfb81f2c55e814b03e1a33653666bc3000000010000000000000024f402275ef40226c200000012336ca7aad576b58d25fb9f0eab8b1e05966272009d7d95e7e9997a3355e6a4d04e1be8adf0fb95320061763fbe021750656f71814e9eccb2fbc692b38c2f64a86651f39c404711218c160701f0de7bb482e574d4c89f3d2226a9eb0551178996125a7fd28464afd403f93a89bfdd13a966ad49b69624820673e90e44d7eedf4ac4d199f8010c2f6a8aef2cbb200c038d8dd5292df3a0f4e8d9dea43a3177844edf34de685b40ec14f748dbb1319ea1b7d7be1bab3a4db4467b59fe425e3f3ffdc298dbeb34bb8c458bc6dd71560d513da31a8b0113e6b2f84d50653581b12239137412897ddec27438a3d40af962ac")

	// 创建仲裁交易
	input2 := langChangeManager.NewEntrySetPlaceHolder(map[string]string{"en": "Arbitration tx main address Password or PrivateKey", "zh": "申请仲裁的账户密码或私钥"})
	input3 := langChangeManager.NewEntrySetPlaceHolder(map[string]string{"en": "Arbitration tx fee - HAC (unit: 248) or use 'ㄜ1:248' format, Example: '0.25' or 'ㄜ25:246'", "zh": "输入交易手续费 - HAC（单位：枚 - :248）也可直接使用 'ㄜ1:248' 格式， 例如：'0.25' or 'ㄜ25:246'"})

	input4 := langChangeManager.NewEntrySetPlaceHolder(map[string]string{"en": "Target channel ID", "zh": "目标通道ID"})

	input5 := langChangeManager.NewEntrySetPlaceHolder(map[string]string{"en": "Optional: Tx timestamp", "zh": "选填：交易时间戳"})

	txbodyshow := widget.NewEntry()
	txbodyshow.MultiLine = true
	txbodyshow.Wrapping = fyne.TextWrapBreak

	getUseTime := func() (int64, bool) {
		usetime := time.Now().Unix()
		if len(input5.Text) > 0 {
			its, e1 := strconv.ParseInt(input5.Text, 10, 0)
			if e1 != nil {
				langChangeManager.SetText(txbodyshow, map[string]string{"en": "Timestamp format error", "zh": "时间戳格式错误"})
				return 0, false
			}
			usetime = its
		}
		return usetime, true
	}

	button1 := langChangeManager.NewButton(map[string]string{"en": "Check", "zh": "查看票据内容"}, func() {
		// 显示票据内容
		langChangeManager.SetText(txbodyshow, renderChannelPaymentBill(input1.Text))
	})
	button2 := langChangeManager.NewButton(map[string]string{"en": "Create Submit Bill Arbitration Tx", "zh": "创建提交票据仲裁交易"}, func() {
		// 创建仲裁交易
		usetime, ok := getUseTime()
		if !ok {
			return
		}
		langChangeManager.SetText(txbodyshow, renderCreateArbitrationTx(input1.Text, input2.Text, input3.Text, usetime))
	})
	button3 := langChangeManager.NewButton(map[string]string{"en": "Create No Bill Arbitration Close Channel Tx", "zh": "创建无票据单方面关闭通道交易"}, func() {
		// 创建无票据单方面关闭通道交易
		usetime, ok := getUseTime()
		if !ok {
			return
		}
		langChangeManager.SetText(txbodyshow, renderCreateNoBillCloseTx(input4.Text, input2.Text, input3.Text, usetime))
	})
	button4 := langChangeManager.NewButton(map[string]string{"en": "Create end arbitration period and close the channel tx", "zh": "创建结束仲裁期并关闭通道交易"}, func() {
		// 创建结束仲裁期并关闭通道交易
		usetime, ok := getUseTime()
		if !ok {
			return
		}
		langChangeManager.SetText(txbodyshow, renderCreateFinishEndArbitrationTx(input4.Text, input2.Text, input3.Text, usetime))
	})

	// 票据
	page.Add(input1)

	// 创建仲裁交易
	page.Add(input2)
	page.Add(input3)
	page.Add(input4)
	page.Add(input5)

	page.Add(button1)
	page.Add(button2)
	page.Add(button3)
	page.Add(button4)

	// 显示
	page.Add(txbodyshow)

	card := langChangeManager.NewCardSetTitle(title, page)
	box.Add(card)

}

// 无票据单方面关闭通道
func renderCreateFinishEndArbitrationTx(cidstr, mainpassword, txfee string, usetime int64) map[string]string {
	contents := map[string]string{"en": "", "zh": ""}
	cidstr = strings.Trim(cidstr, "\n ")
	mainpassword = strings.Trim(mainpassword, "\n ")
	txfee = strings.Trim(txfee, "\n ")
	// 检查票据
	if cidstr == "" {
		return map[string]string{"en": "Please input channel id", "zh": "请输入通道ID"}
	}
	if mainpassword == "" {
		return map[string]string{"en": "Please input main private key", "zh": "请输入主账户私钥或密码"}
	}
	if txfee == "" {
		return map[string]string{"en": "Please input tx fee", "zh": "请输入交易手续费"}
	}
	// 检查通道id
	channelID, e := hex.DecodeString(cidstr)
	if e != nil || len(channelID) != stores.ChannelIdLength {
		return map[string]string{"en": "channel id format error", "zh": "通道ID格式错误"}
	}
	// 检查交易费
	fee, e := fields.NewAmountFromString(txfee)
	if e != nil {
		return map[string]string{"en": "tx fee format error", "zh": "手续费格式错误"}
	}
	mainacc := account.GetAccountByPrivateKeyOrPassword(mainpassword)
	// 按类型创建交易
	tx, e := transactions.NewEmptyTransaction_2_Simple(mainacc.Address)
	if e != nil {
		return map[string]string{"en": "create tx error", "zh": "创建交易失败"}
	}
	tx.Timestamp = fields.BlockTxTimestamp(usetime) // 时间戳
	tx.SetFee(fee)
	tx.AppendAction(&actions.Action_27_ClosePaymentChannelByClaimDistribution{
		ChannelId: channelID,
	})
	// 签名
	tx.FillTargetSign(mainacc)

	// 显示
	txhex, e := tx.Serialize()
	if e != nil {
		return map[string]string{"en": "Serialize tx error", "zh": "序列化交易失败"}
	}

	showcon := "Create FinishEndArbitrationTx Create Successfully!\n\n---- tx body start ----\n" +
		hex.EncodeToString(txhex) +
		"\n---- tx body end ----\n"

	contents["en"] = showcon
	contents["zh"] = showcon

	return contents

}

// 无票据单方面关闭通道
func renderCreateNoBillCloseTx(cidstr, mainpassword, txfee string, usetime int64) map[string]string {
	contents := map[string]string{"en": "", "zh": ""}
	cidstr = strings.Trim(cidstr, "\n ")
	mainpassword = strings.Trim(mainpassword, "\n ")
	txfee = strings.Trim(txfee, "\n ")
	// 检查票据
	if cidstr == "" {
		return map[string]string{"en": "Please input channel id", "zh": "请输入通道ID"}
	}
	if mainpassword == "" {
		return map[string]string{"en": "Please input main private key", "zh": "请输入主账户私钥或密码"}
	}
	if txfee == "" {
		return map[string]string{"en": "Please input tx fee", "zh": "请输入交易手续费"}
	}
	// 检查通道id
	channelID, e := hex.DecodeString(cidstr)
	if e != nil || len(channelID) != stores.ChannelIdLength {
		return map[string]string{"en": "channel id format error", "zh": "通道ID格式错误"}
	}
	// 检查交易费
	fee, e := fields.NewAmountFromString(txfee)
	if e != nil {
		return map[string]string{"en": "tx fee format error", "zh": "手续费格式错误"}
	}
	mainacc := account.GetAccountByPrivateKeyOrPassword(mainpassword)
	// 按类型创建交易
	tx, e := transactions.NewEmptyTransaction_2_Simple(mainacc.Address)
	if e != nil {
		return map[string]string{"en": "create tx error", "zh": "创建交易失败"}
	}
	tx.Timestamp = fields.BlockTxTimestamp(usetime) // 时间戳
	tx.SetFee(fee)
	tx.AppendAction(&actions.Action_22_UnilateralClosePaymentChannelByNothing{
		ChannelId:          channelID,
		AssertCloseAddress: mainacc.Address,
	})
	// 签名
	tx.FillTargetSign(mainacc)

	// 显示
	txhex, e := tx.Serialize()
	if e != nil {
		return map[string]string{"en": "Serialize tx error", "zh": "序列化交易失败"}
	}

	showcon := "Create Tx Create Successfully!\n\n---- tx body start ----\n" +
		hex.EncodeToString(txhex) +
		"\n---- tx body end ----\n"

	contents["en"] = showcon
	contents["zh"] = showcon

	return contents

}

// 提交票据仲裁交易
func renderCreateArbitrationTx(billbody, mainpassword, txfee string, usetime int64) map[string]string {
	contents := map[string]string{"en": "", "zh": ""}
	billbody = strings.Trim(billbody, "\n ")
	mainpassword = strings.Trim(mainpassword, "\n ")
	txfee = strings.Trim(txfee, "\n ")
	// 检查票据
	if billbody == "" {
		return map[string]string{"en": "Please input the bill data", "zh": "请输入支付票据数据"}
	}
	if mainpassword == "" {
		return map[string]string{"en": "Please input main private key", "zh": "请输入主账户私钥或密码"}
	}
	if txfee == "" {
		return map[string]string{"en": "Please input tx fee", "zh": "请输入交易手续费"}
	}
	txbody, e0 := hex.DecodeString(billbody)
	if e0 != nil {
		return map[string]string{"en": "bill data format error", "zh": "支付票据数据格式错误"}
	}
	// 解析支付票据
	bill, _, e1 := channel.ParseReconciliationBalanceBillByPrefixTypeCode(txbody, 0)
	if e1 != nil {
		return map[string]string{"en": "bill data error", "zh": "支付票据数据错误"}
	}
	// 检查交易费
	fee, e := fields.NewAmountFromString(txfee)
	if e != nil {
		return map[string]string{"en": "tx fee format error", "zh": "手续费格式错误"}
	}
	mainacc := account.GetAccountByPrivateKeyOrPassword(mainpassword)
	// 按类型创建交易
	tx, e := transactions.NewEmptyTransaction_2_Simple(mainacc.Address)
	if e != nil {
		return map[string]string{"en": "create tx error", "zh": "创建交易失败"}
	}
	tx.Timestamp = fields.BlockTxTimestamp(usetime) // 时间戳
	tx.SetFee(fee)
	// 票据类型
	if bill.TypeCode() == channel.BillTypeCodeReconciliation {
		billobj := bill.(*channel.OffChainFormPaymentChannelRealtimeReconciliation)
		oncbill := billobj.ConvertToOnChain()
		tx.AppendAction(&actions.Action_23_UnilateralCloseOrRespondChallengePaymentChannelByRealtimeReconciliation{
			AssertAddress:  mainacc.Address,
			Reconciliation: *oncbill,
		})

	} else if bill.TypeCode() == channel.BillTypeCodeSimplePay {
		billobj := bill.(*channel.OffChainCrossNodeSimplePaymentReconciliationBill)
		tx.AppendAction(&actions.Action_24_UnilateralCloseOrRespondChallengePaymentChannelByChannelChainTransferBody{
			AssertAddress:                       mainacc.Address,
			ChannelChainTransferData:            billobj.ChannelChainTransferData,
			ChannelChainTransferTargetProveBody: billobj.ChannelChainTransferTargetProveBody,
		})
	}
	// 签名
	tx.FillTargetSign(mainacc)

	// 显示
	txhex, e := tx.Serialize()
	if e != nil {
		return map[string]string{"en": "Serialize tx error", "zh": "序列化交易失败"}
	}

	showcon := "Submit Bill Arbitration Tx Create Successfully!\n\n---- tx body start ----\n" +
		hex.EncodeToString(txhex) +
		"\n---- tx body end ----\n"

	contents["en"] = showcon
	contents["zh"] = showcon

	return contents
}

// 输出票据的内容
func renderChannelPaymentBill(txbodystr string) map[string]string {
	contents := map[string]string{"en": "", "zh": ""}
	txbodystr = strings.Trim(txbodystr, "\n ")
	// 检查
	if txbodystr == "" {
		return map[string]string{"en": "Please input the bill data", "zh": "请输入支付票据数据"}
	}
	txbody, e0 := hex.DecodeString(txbodystr)
	if e0 != nil {
		return map[string]string{"en": "bill data format error", "zh": "支付票据数据格式错误"}
	}
	// 解析支付票据
	bill, _, e1 := channel.ParseReconciliationBalanceBillByPrefixTypeCode(txbody, 0)
	if e1 != nil {
		return map[string]string{"en": "bill data error", "zh": "支付票据数据错误"}
	}
	// 显示内容
	var en, zh string
	// 基础内容
	tts := time.Unix(int64(bill.GetTimestamp()), 0).Format("2006-01-02 15:04:05")
	renderFmt := func(fmtstr string) string {
		return fmt.Sprintf(
			fmtstr,
			bill.TypeCode(),
			bill.GetChannelId().ToHex(),
			bill.GetReuseVersion(),
			bill.GetAutoNumber(),
			bill.GetLeftAddress().ToReadable(),
			bill.GetLeftBalance().ToFinString(),
			bill.GetLeftSatoshi(),
			bill.GetRightAddress().ToReadable(),
			bill.GetRightBalance().ToFinString(),
			bill.GetRightSatoshi(),
			bill.GetTimestamp(), tts,
		)
	}
	en += renderFmt("Bill type: %d\n" +
		"Channel ID: %s\n" +
		"Reuse version: %d\n" +
		"Reconciliation serial number: %d\n" +
		"Left address: %s\n" +
		"Left balance: %s\n" +
		"Left satoshi: %d\n" +
		"Right address: %s\n" +
		"Right balance: %s\n" +
		"Right satoshi: %d\n" +
		"Timestamp: %d (%s)\n")
	zh += renderFmt("票据类型: %d\n" +
		"通道ID: %s\n" +
		"重用版本: %d\n" +
		"对账流水号: %d\n" +
		"左侧地址: %s\n" +
		"左侧HAC余额: %s\n" +
		"左侧SAT余额: %d\n" +
		"右侧地址: %s\n" +
		"右侧HAC余额: %s\n" +
		"右侧SAT余额: %d\n" +
		"时间戳: %d (%s)\n")

	if channel.BillTypeCodeSimplePay == bill.TypeCode() {
		// 通道支付
		bobj := bill.(*channel.OffChainCrossNodeSimplePaymentReconciliationBill)
		dt := bobj.ChannelChainTransferData
		prbd := bobj.ChannelChainTransferTargetProveBody
		signdts := make([]string, dt.MustSignCount)
		for i, v := range dt.MustSigns {
			signdts[i] = fmt.Sprintf("    %s: %s",
				v.GetAddress().ToReadable(), v.Signature.ToHex())
		}
		paydrct := "left to right"
		payamt := prbd.PayAmount.ToFinString()
		pd := uint8(prbd.PayDirection)
		if pd == channel.ChannelTransferDirectionHacashRightToLeft ||
			pd == channel.ChannelTransferDirectionSatoshiRightToLeft {
			paydrct = "right to left"
		}
		if pd >= channel.ChannelTransferDirectionSatoshiLeftToRight {
			payamt = fmt.Sprintf("SAT %d", prbd.PaySatoshi.GetRealSatoshi())
		}
		adinfo := fmt.Sprintf("\nLast pay: %s %s\n"+
			"\nChannel count: %d\n"+
			"Sign address count: %d\n"+
			"Order hash: %s\n"+
			"Addresses & Signs: {\n%s\n}",
			paydrct, payamt,
			dt.ChannelCount,
			dt.MustSignCount,
			dt.OrderNoteHashHalfChecker.ToHex(),
			strings.Join(signdts, "\n"))
		en += adinfo
		zh += adinfo

	} else if channel.BillTypeCodeReconciliation == bill.TypeCode() {
		// 通道对账
		bobj := bill.(*channel.OffChainFormPaymentChannelRealtimeReconciliation)

		adinfo := fmt.Sprintf("\nLeft sign data: {%s, %s}\n"+
			"Right sign data: {%s, %s}\n",
			bobj.LeftSign.PublicKey.ToHex(), bobj.LeftSign.Signature.ToHex(),
			bobj.RightSign.PublicKey.ToHex(), bobj.RightSign.Signature.ToHex())
		en += adinfo
		zh += adinfo

	} else {
		return map[string]string{"en": "bill data type error", "zh": "支付票据类型错误"}
	}

	/*/ 解析
	var en, zh string
	var mainaddr = trs.GetAddress().ToReadable()
	acts := trs.GetActions()
	txtimestamp := int64(trs.GetTimestamp())
	txtimestr := time.Unix(txtimestamp, 0).Format("2006-01-02 15:04:05")
	// en
	en += "Tx hash: <" + trs.Hash().ToHex() + ">"
	en += "\nFee account: <" + mainaddr + "> pay fee <" + trs.GetFee().ToFinString() + ">"
	en += fmt.Sprintf("\nTimestamp: <%d> (%s)", txtimestamp, txtimestr)
	en += fmt.Sprintf("\n\nTx actions (%d): [\n", len(acts))
	// zh
	zh += "交易哈希: <" + trs.Hash().ToHex() + ">"
	zh += "\n手续费账户: <" + mainaddr + "> 支付手续费 <" + trs.GetFee().ToFinString() + ">"
	zh += fmt.Sprintf("\n时间戳: <%d> (%s)", txtimestamp, txtimestr)
	zh += fmt.Sprintf("\n\n交易包含内容 (%d): [\n", len(acts))
	// loop actions
	for i, act := range acts {
		nod := fmt.Sprintf("\n%d). ", i+1)
		en += nod
		zh += nod
		// 解析每一项 action 的描述
		l1, l2 := renderTxActionDescribe(mainaddr, act)
		en += l1
		zh += l2
	}
	end1 := "\n\n]\n"
	en += end1
	zh += end1
	// 签名检查
	checkTxSignStatus(trs, &en, &zh)
	// ret ok
	contents["en"] = en
	contents["zh"] = zh
	return contents
	*/
	sgckr := "\nAll signatures checked successfully."
	if e := bill.VerifySignature(); e != nil {
		sgckr = "\nSignature check failed."
	}
	en += sgckr
	zh += sgckr

	// ret ok
	contents["en"] = en
	contents["zh"] = zh
	return contents
}
