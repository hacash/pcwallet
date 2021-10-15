package actions

import (
	"encoding/hex"
	"fmt"
	"fyne.io/fyne"
	"fyne.io/fyne/container"
	"fyne.io/fyne/widget"
	"github.com/hacash/core/channel"
	"github.com/hacash/pcwallet/widgets"
	"strings"
	"time"
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
	//input1.SetText("01cdfb81f2c55e814b03e1a33653666bc300000001000000000000000e02f701010012336ca7aad576b58d25fb9f0eab8b1e05966272009d7d95e7e9997a3355e6a4d04e1be8adf0fb9532f1039828f6f10399040a00616831af4d7bbb0407d8681d0d86d1e91e00167908003230e909db0810e670a369c6868274136dad12ef005c110abc683fcdfa40027dab15f9f064a7a36bf7009d7d95e7e9997a3355e6a4d04e1be8adf0fb95320012336ca7aad576b58d25fb9f0eab8b1e0596627200f393e23a42bc3431eccfbd0b829929cb9abd20ad00a536637453340a1c1d8ca3f1dec0b44c5728d7a3009cba1cb8f332141964668ea906a38f339f1bbcca001ecf9afca1c31fdeacc2091acc91c2dc5ef28a7905ea0e26d85f5de5c1802e7c2c647be340401df498c67aed8796c5b952e4531be370d5ded9068ce40636ffd719679205916edf600936cf3e5464e47af3c32b3224ebb233e9cac026f8b6f582c52b667a9802ecc09b6ac3123cb5809cea7d39ee136a8c4e611c4328702450366532863bd1b3163759e4d7ba1914b0bc2bd5b0c44561aace63973f5b39d6a4ac97926b9de73c32a423e5ddc2d83c7443b4b087afc4e1e8f19791b2398b5383d11ba02594c41f02a83451e8167768731e89d04a89f24f9df9d5e5351b337ef3e564892590721fcd0f2b9b4dd35863b6f51eb35256a77df41a665b0a9b570b7b8a9755f5a6289e451e0150bf98751ac32e6c4f6a387899220a072f771662c16c3c76a5082de67c23038d8dd5292df3a0f4e8d9dea43a3177844edf34de685b40ec14f748dbb1319ea1ae7fd62805eec76939eebfac556117a9c95cc7ee8fb53dd6d0ff85dca8cc46b124447f9328a9b1f93305a3ee27823e8738636de1346a1127bbb83f95ca4e1ff1021750656f71814e9eccb2fbc692b38c2f64a86651f39c404711218c160701f0de02c28eae823731fbde9ca6d79288de6e31cf6a5fa6440a40a41a8902348a7853161bd533abdd36c63f558701512ea42a99aa5b4467f8f2d7b2b3efebcc33a74e03cb6f7785cabfb38e39f30c35a7a87c845d604c908d402ee232b98f3b4be8c007d982c70236d1fd39ae51211b75967f0788f43059ff96fab2ea48eab53b89bd7571307ca210013cb77a0156327e0bd87c583f3287c53eca46ade0de2df35733e80344ee89b8a8720daf111f3370d127919db2075c2bed1dbf1bfd3ef020f72276f864efefb4208e38b057a9f40e21bda5a3fb7b6429d3ed70babd062650ec27c6cb2704ec5f9cd0dd60819e95bbad3c392a4f1403ba882c033ed761efbb5ef79d570285a2fe9808a81e92aae51b80052478f76096a8eff337d0a6a4d875e4eea6f6e843da23e3c0c2a662a7934fc9b22cfd3ebcb918278c6013e0c7bd9e0ba983ec8252713801392aae0f2f24e4f83cd7ef5878475863494c2d4bf9907508a06a215503f1e4fbfb53b19229bc9cb67f3756f7d40fae6e512457adf3c3b0468b837c2526aab2c5909568cdd60c05ae36b31f9eddc78994b359c3debb07976b7dbd9c406f292e44ef719a8327835206490a3080aa6bf6a8f16f7f1965bda6c2ce77504043")

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
			bill.GetRightAddress().ToReadable(),
			bill.GetRightBalance().ToFinString(),
			bill.GetTimestamp(), tts,
		)
	}
	en += renderFmt("Bill type: %d\n" +
		"Channel ID: %s\n" +
		"Reuse version: %d\n" +
		"Reconciliation serial number: %d\n" +
		"Left address: %s\n" +
		"Left balance: %s\n" +
		"Right address: %s\n" +
		"Right balance: %s\n" +
		"Timestamp: %d (%s)\n")
	zh += renderFmt("票据类型: %d\n" +
		"通道ID: %s\n" +
		"重用版本: %d\n" +
		"对账流水号: %d\n" +
		"左侧地址: %s\n" +
		"左侧余额: %s\n" +
		"右侧地址: %s\n" +
		"右侧余额: %s\n" +
		"时间戳: %d (%s)\n")

	if channel.BillTypeCodeSimplePay == bill.TypeCode() {
		// 通道支付
		bobj := bill.(*channel.OffChainCrossNodeSimplePaymentReconciliationBill)
		dt := bobj.ChannelChainTransferData
		adinfo := fmt.Sprintf("\nChannel count: %d\n"+
			"Sign address count: %d\n"+
			"Order hash: %s\n",
			dt.ChannelCount,
			dt.MustSignCount,
			dt.OrderNoteHashHalfChecker.ToHex())
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
