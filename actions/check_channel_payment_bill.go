package actions

import (
	"encoding/hex"
	"fyne.io/fyne"
	"fyne.io/fyne/container"
	"fyne.io/fyne/widget"
	"github.com/hacash/core/channel"
	"github.com/hacash/pcwallet/widgets"
	"strings"
)

// 检查通道链支付对账票据
func AddOpenButtonOnMainOfCheckChannelPaymentBill(box *fyne.Container, langChangeManager *widgets.LangChangeManager) {

	title := map[string]string{"en": "Check channel payment bill", "zh": "验证通道链支付对账票据"}

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

	txbodyshow := widget.NewEntry()
	txbodyshow.MultiLine = true
	txbodyshow.Wrapping = fyne.TextWrapBreak

	button1 := langChangeManager.NewButton(map[string]string{"en": "Check", "zh": "查看"}, func() {
		// 显示票据内容
		langChangeManager.SetText(txbodyshow, renderChannelPaymentBill(input1.Text))
	})

	// add item
	page.Add(input1)
	page.Add(button1)
	page.Add(txbodyshow)

	card := langChangeManager.NewCardSetTitle(title, page)
	box.Add(card)

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
	if channel.BillTypeCodeReconciliation == bill.TypeCode() {
		// 通道支付
		bobj := bill.(*channel.OffChainCrossNodeSimplePaymentReconciliationBill)
		bobj.TypeCode()

	} else if channel.BillTypeCodeSimplePay == bill.TypeCode() {
		// 通道对账
		bobj := bill.(*channel.OffChainFormPaymentChannelRealtimeReconciliation)
		bobj.TypeCode()

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

	// ret ok
	contents["en"] = en
	contents["zh"] = zh
	return contents
}
