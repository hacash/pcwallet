package actions

import (
	"encoding/hex"
	"fmt"
	"fyne.io/fyne"
	"fyne.io/fyne/container"
	"fyne.io/fyne/widget"
	"github.com/hacash/core/account"
	"github.com/hacash/core/fields"
	"github.com/hacash/core/interfaces"
	"github.com/hacash/core/transactions"
	"github.com/hacash/pcwallet/widgets"
	"strings"
)

func AddOpenButtonOnMainOfSignTx(box *fyne.Container, langChangeManager *widgets.LangChangeManager) {

	title := map[string]string{"en": "Sign the tx", "zh": "签署交易"}

	button := langChangeManager.NewButton(title, func() {
		OpenWindowSignTx(title, langChangeManager)
	})
	box.Add(button)
}

func OpenWindowSignTx(title map[string]string, langChangeManager *widgets.LangChangeManager) fyne.Window {

	// 打开窗口测试
	testSize := fyne.Size{
		Width:  800,
		Height: 1000,
	}

	box := container.NewVBox()
	AddCanvasObjectSignTx(title, box, langChangeManager)

	// 开启窗口
	return langChangeManager.NewWindowAndShow(title, &testSize, box)
}

func AddCanvasObjectSignTx(title map[string]string, box *fyne.Container, langChangeManager *widgets.LangChangeManager) {

	box.Add(widget.NewLabel("\n\n"))

	page := container.NewVBox()

	page.Add(langChangeManager.NewTextWrapWordLabel(map[string]string{"en": "Use the private key or password to sign a transaction and display the signed transaction data", "zh": "使用私钥或者密码签署一笔交易，显示签名后的交易数据"}))

	// 交易输入框
	input1 := langChangeManager.NewEntrySetPlaceHolder(map[string]string{"en": "Transaction body hex string", "zh": "请输入交易体数据（txbody hex）"})
	input1.Wrapping = fyne.TextWrapWord
	input1.MultiLine = true

	// 密码或私钥
	input2 := langChangeManager.NewEntrySetPlaceHolder(map[string]string{"en": "Input private key or password", "zh": "这里输入私钥或密码"})

	// 结果显示
	txbodyshow := widget.NewEntry()
	txbodyshow.MultiLine = true
	txbodyshow.Wrapping = fyne.TextWrapBreak

	// 签名 forceAddSign = 是否强制附加签名
	doSignCallFunc := func(cleanAdditionalSignature bool, forceAddSign bool) {

		signacc := account.GetAccountByPrivateKeyOrPassword(input2.Text)
		signaddr := fields.Address(signacc.Address)
		// 解析交易
		txbodystr := input1.Text
		txbodystr = strings.Trim(txbodystr, "\n ")
		// 检查
		if txbodystr == "" {
			langChangeManager.SetText(txbodyshow, map[string]string{"en": "Please input the txbody ", "zh": "请输入交易体数据"})
			return
		}
		txbody, e0 := hex.DecodeString(txbodystr)
		if e0 != nil {
			langChangeManager.SetText(txbodyshow, map[string]string{"en": "txbody format error", "zh": "交易体数据格式错误"})
			return
		}
		// 解析交易
		trs, _, e1 := transactions.ParseTransaction(txbody, 0)
		if e1 != nil {
			langChangeManager.SetText(txbodyshow, map[string]string{"en": "txbody data error", "zh": "交易体数据错误"})
			return
		}
		// 输出
		var en, zh string
		// 显示交易体
		showTxBodyHexCall := func() {
			newtxbody, _ := trs.Serialize()
			txbodyhex := "\n\n-------- signed txbody hex start --------\n"
			txbodyhex += hex.EncodeToString(newtxbody)
			txbodyhex += "\n-------- signed txbody hex  end  --------\n\n"
			en += fmt.Sprintf("\nTxbody hex：%s", txbodyhex)
			zh += fmt.Sprintf("\n交易体数据 (txbody hex)：%s", txbodyhex)
		}

		// 清除附加签名
		if cleanAdditionalSignature {
			reqaddrmaps := make(map[string]bool)
			reqaddrs, _ := trs.RequestSignAddresses(nil, false)
			for _, v := range reqaddrs {
				reqaddrmaps[string(v)] = true
			}
			newsigns := make([]fields.Sign, 0)
			for _, v := range trs.GetSigns() {
				addr := account.NewAddressFromPublicKeyV0(v.PublicKey)
				if _, hav := reqaddrmaps[string(addr)]; hav {
					newsigns = append(newsigns, v)
				}
			}
			// 重设
			trs.SetSigns(newsigns)
			// 打印
			en = fmt.Sprintf("Clean additional signature successfully!\n")
			zh = fmt.Sprintf("清除所有附加签名成功！\n")
			// 清除成功，显示交易体
			showTxBodyHexCall()
			// 显示签名检查
			checkTxSignStatus(trs, &en, &zh)
			langChangeManager.SetText(txbodyshow, map[string]string{"en": en, "zh": zh})
			return
		}

		// 检查私钥面膜
		if input2.Text == "" {
			langChangeManager.SetText(txbodyshow, map[string]string{"en": "Please input private key or password", "zh": "请输入私钥或密码"})
			return
		}
		doSign := func() {
			// 签署
			err := trs.FillTargetSign(signacc)
			if err != nil {
				en = fmt.Sprintf("sign error: %s\n\n%s\n\n", err.Error(), en)
				zh = fmt.Sprintf("签名发生错误：%s\n\n%s\n\n", err.Error(), zh)
				checkTxSignStatus(trs, &en, &zh) // 显示签名检查
				langChangeManager.SetText(txbodyshow, map[string]string{"en": en, "zh": zh})
				return
			}
			// 签名成功，显示交易体
			en += fmt.Sprintf("Transaction signed successfully!\nSign account: <%s>\n", signacc.AddressReadable)
			zh += fmt.Sprintf("交易签名成功!\n签署账户：<%s>\n", signacc.AddressReadable)
			showTxBodyHexCall()
			// 加上签名检查
			checkTxSignStatus(trs, &en, &zh) // 显示签名检查
			langChangeManager.SetText(txbodyshow, map[string]string{"en": en, "zh": zh})
			return
		}

		// 选择执行
		if forceAddSign {
			// 强制执行签名
			doSign()

		} else {
			// 待签名地址 按需签名
			reqaddrs, _ := trs.RequestSignAddresses(nil, false)
			for _, v := range reqaddrs {
				if signaddr.Equal(v) {
					doSign()
					return
				}
			}
			// 账户不匹配
			en = "This transaction does not require this private key or password to sign\n\n"
			zh = "本交易无需此私钥或密码签署\n\n"
			checkTxSignStatus(trs, &en, &zh) // 显示签名检查
			langChangeManager.SetText(txbodyshow, map[string]string{"en": en, "zh": zh})
			return
		}

	}

	// 按钮
	button1 := langChangeManager.NewButton(map[string]string{"en": "Do sign on demand", "zh": "执行按需签名"}, func() {
		doSignCallFunc(false, false)
	})
	button2 := langChangeManager.NewButton(map[string]string{"en": "Do force sign additionally", "zh": "强制附加签名"}, func() {
		doSignCallFunc(false, true)
	})
	button3 := langChangeManager.NewButton(map[string]string{"en": "Clean all additional signature", "zh": "清除所有附加签名"}, func() {
		doSignCallFunc(true, false)
	})

	// add item
	page.Add(input1)
	page.Add(input2)
	page.Add(button1)
	page.Add(button2)
	page.Add(button3)
	page.Add(txbodyshow)

	card := langChangeManager.NewCardSetTitle(title, page)
	box.Add(card)

}

// 输出签名检查
func checkTxSignStatus(trs interfaces.Transaction, en *string, zh *string) {
	// 签名
	notsignedaddrnum := 0
	reqaddrmaps := make(map[string]bool)
	reqaddrs, _ := trs.RequestSignAddresses(nil, false)
	*en += fmt.Sprintf("\nSignature accounts required (%d): {\n", len(reqaddrs))
	*zh += fmt.Sprintf("\n交易必要签名检查 (%d): {\n", len(reqaddrs))
	for i, rqa := range reqaddrs {
		reqaddrmaps[string(rqa)] = true // 标记
		en_stat := "OK: completed the signature"
		zh_stat := "OK: 已完成签名"
		rqas := []fields.Address{rqa}
		if ok, e1 := trs.VerifyTargetSigns(rqas); !ok || e1 != nil {
			en_stat = "fail: not signed"
			zh_stat = "验证失败：未签名"
			notsignedaddrnum += 1
		}
		*en += fmt.Sprintf("\n%d). %s <%s>", i+1, rqa.ToReadable(), en_stat)
		*zh += fmt.Sprintf("\n%d). %s <%s>", i+1, rqa.ToReadable(), zh_stat)
	}
	end2 := "\n\n}\n"
	*en += end2
	*zh += end2

	// 检查是否有附带签名
	allsigns := trs.GetSigns()
	forceAppendSignAddrs := []fields.Address{}
	for _, v := range allsigns {
		addr := account.NewAddressFromPublicKeyV0(v.PublicKey)
		if _, has := reqaddrmaps[string(addr)]; !has {
			// 一个强制附加签名
			forceAppendSignAddrs = append(forceAppendSignAddrs, addr)
		}
	}
	fasac := len(forceAppendSignAddrs)
	if fasac > 0 {
		// 显示强制附加签名
		*en += fmt.Sprintf("\nAdditional signature accounts (%d): {\n", fasac)
		*zh += fmt.Sprintf("\n附加签名账户 (%d): {\n", fasac)
		for i, v := range forceAppendSignAddrs {
			*en += fmt.Sprintf("\n%d). %s", i+1, v.ToReadable())
			*zh += fmt.Sprintf("\n%d). %s", i+1, v.ToReadable())
		}
		*en += end2
		*zh += end2
	}

	if notsignedaddrnum == 0 {
		*en += "\nSigned successfully: all signatures have been completed.\n"
		*zh += "\n签名验证成功: 已全部完成签名。\n"
	} else {
		*en += fmt.Sprintf("\nSignature verification failed: %d accounts not signed.\n", notsignedaddrnum)
		*zh += fmt.Sprintf("\n签名验证失败: %d 个账户未签名。\n", notsignedaddrnum)
	}
}
