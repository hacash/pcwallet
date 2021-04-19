package actions

import (
	"encoding/hex"
	"fmt"
	"fyne.io/fyne"
	"fyne.io/fyne/container"
	"fyne.io/fyne/widget"
	"github.com/hacash/core/actions"
	"github.com/hacash/core/transactions"
	"github.com/hacash/pcwallet/widgets"
	"strings"
	"time"
)

func AddCanvasObjectCheckTxContent(box *fyne.Container, langChangeManager *widgets.LangChangeManager) {

	box.Add(widget.NewLabel("\n\n"))

	page := container.NewVBox()

	page.Add(langChangeManager.NewTextWrapWordLabel(map[string]string{"en": "View the operation items, conditions and contents of a transaction, and check whether the signature is completed", "zh": "查看一笔交易具体包含的操作项目、条件和内容，以及检查是否完成签名"}))

	input1 := langChangeManager.NewEntrySetPlaceHolder(map[string]string{"en": "Transaction body hex string", "zh": "请输入交易体数据（txbody hex）"})
	input1.Wrapping = fyne.TextWrapWord
	input1.MultiLine = true

	txbodyshow := widget.NewEntry()
	txbodyshow.MultiLine = true
	txbodyshow.Wrapping = fyne.TextWrapBreak

	button1 := langChangeManager.NewButton(map[string]string{"en": "Check", "zh": "查看"}, func() {
		// 显示交易条款
		langChangeManager.SetText(txbodyshow, renderTxContent(input1.Text))
	})

	// add item
	page.Add(input1)
	page.Add(button1)
	page.Add(txbodyshow)

	card := langChangeManager.NewCardSetTitle(map[string]string{"en": "Check tx content", "zh": "验证交易内容"}, page)
	box.Add(card)

}

// 输出一笔交易的内容
func renderTxContent(txbodystr string) map[string]string {
	contents := map[string]string{"en": "", "zh": ""}
	txbodystr = strings.Trim(txbodystr, "\n ")
	// 检查
	if txbodystr == "" {
		return map[string]string{"en": "Please input the txbody ", "zh": "请输入交易体数据"}
	}
	txbody, e0 := hex.DecodeString(txbodystr)
	if e0 != nil {
		return map[string]string{"en": "txbody format error", "zh": "交易体数据格式错误"}
	}
	// 解析交易
	trs, _, e1 := transactions.ParseTransaction(txbody, 0)
	if e1 != nil {
		return map[string]string{"en": "txbody data error", "zh": "交易体数据错误"}
	}
	// 解析
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
		// 每一项条款
		/************************ Action_1_SimpleTransfer ************************/
		if a, ok := act.(*actions.Action_1_SimpleTransfer); ok {
			toaddr := a.ToAddress.ToReadable()
			amt := a.Amount.ToFinString()
			en += fmt.Sprintf("simple transfer: account <%s> transfers amount <%s> to account <%s>", mainaddr, amt, toaddr)
			zh += fmt.Sprintf("普通转账： 地址 <%s> 向地址 <%s> 转账 <%s>", mainaddr, toaddr, amt)
		}
		/************************ END ************************/
	}
	end1 := "\n\n]\n"
	en += end1
	zh += end1
	// 签名
	notsignedaddrnum := 0
	reqaddrs, _ := trs.RequestSignAddresses(nil, false)
	en += fmt.Sprintf("\nSignature accounts required (%d): {\n", len(reqaddrs))
	zh += fmt.Sprintf("\n交易签名检查 (%d): {\n", len(reqaddrs))
	for i, rqa := range reqaddrs {
		en_stat := "OK: completed the signature"
		zh_stat := "OK: 已完成签名"
		if ok, e1 := trs.VerifyTargetSign(rqa); !ok || e1 != nil {
			en_stat = "fail: not signed"
			zh_stat = "验证失败：未签名"
			notsignedaddrnum += 1
		}
		en += fmt.Sprintf("\n%d). %s <%s>", i+1, rqa.ToReadable(), en_stat)
		zh += fmt.Sprintf("\n%d). %s <%s>", i+1, rqa.ToReadable(), zh_stat)
	}
	end2 := "\n\n}\n"
	en += end2
	zh += end2

	if notsignedaddrnum == 0 {
		en += "\nSigned successfully: all signatures have been completed.\n"
		zh += "\n签名验证成功: 已全部完成签名。\n"
	} else {
		en += fmt.Sprintf("\nSignature verification failed: %d accounts not signed.\n", notsignedaddrnum)
		zh += fmt.Sprintf("\n签名验证失败: %d 个账户未签名。\n", notsignedaddrnum)
	}

	// ret ok
	contents["en"] = en
	contents["zh"] = zh
	return contents
}
