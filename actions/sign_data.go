package actions

import (
	"encoding/hex"
	"fmt"
	fyne "fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/widget"
	"github.com/hacash/core/account"
	"github.com/hacash/core/fields"
	"github.com/hacash/pcwallet/widgets"
)

func AddOpenButtonOnMainOfSignData(box *fyne.Container, langChangeManager *widgets.LangChangeManager) {
	title := map[string]string{"en": "Sign some data", "zh": "签署数据"}

	button := langChangeManager.NewButton(title, func() {
		OpenWindowSignData(title, langChangeManager)
	})
	box.Add(button)
}

func OpenWindowSignData(title map[string]string, langChangeManager *widgets.LangChangeManager) fyne.Window {
	// 打开窗口测试
	testSize := fyne.Size{
		Width:  800,
		Height: 1000,
	}

	box := container.NewVBox()
	AddCanvasObjectSignData(title, box, langChangeManager)

	// 开启窗口
	return langChangeManager.NewWindowAndShow(title, &testSize, box)
}

func AddCanvasObjectSignData(title map[string]string, box *fyne.Container, langChangeManager *widgets.LangChangeManager) {
	box.Add(widget.NewLabel("\n\n"))

	page := container.NewVBox()

	page.Add(langChangeManager.NewTextWrapWordLabel(map[string]string{"en": "Use the private key or password to sign a transaction and display the signed transaction data", "zh": "使用私钥或者密码签署一条数据，显示签名后的证明数据"}))

	// 交易输入框
	input1 := langChangeManager.NewEntrySetPlaceHolder(map[string]string{"en": "data for sign", "zh": "请输入待签名数据"})
	input1.Wrapping = fyne.TextWrapWord
	input1.MultiLine = true

	// 密码或私钥
	input2 := langChangeManager.NewEntrySetPlaceHolder(map[string]string{"en": "Input private key or password", "zh": "这里输入私钥或密码"})

	// 结果显示
	txbodyshow := widget.NewEntry()
	txbodyshow.MultiLine = true
	txbodyshow.Wrapping = fyne.TextWrapBreak

	// 按钮
	button1 := langChangeManager.NewButton(map[string]string{"en": "Do sign", "zh": "开始签署数据"}, func() {
		// 执行签名
		if len(input1.Text) == 0 {
			langChangeManager.SetText(txbodyshow, map[string]string{"en": "Please input the data for sign ", "zh": "请输入待签名数据"})
			return
		}

		if len(input2.Text) == 0 {
			langChangeManager.SetText(txbodyshow, map[string]string{"en": "Please input private key or password", "zh": "请输入私钥或密码"})
			return
		}

		signacc := account.GetAccountByPrivateKeyOrPassword(input2.Text)
		sckdata := fields.CreateSignCheckData(input1.Text)
		e := sckdata.FillSign(signacc)
		if e != nil {
			langChangeManager.SetText(txbodyshow, map[string]string{"en": fmt.Sprintf("do sign error: %s", e.Error()), "zh": fmt.Sprintf("签名发生错误: %s", e.Error())})
			return
		}

		sckdts, e := sckdata.Serialize()
		if e != nil {
			langChangeManager.SetText(txbodyshow, map[string]string{"en": fmt.Sprintf("SignCheckData serialize error: %s", e.Error()), "zh": fmt.Sprintf("SignCheckData 序列化错误: %s", e.Error())})
			return
		}

		// 输出
		var en, zh string
		// 显示交易体
		showSignDataHexCall := func() {
			bodyhex := "\n\n-------- signed data hex start --------\n"
			bodyhex += hex.EncodeToString(sckdts)
			bodyhex += "\n-------- signed data hex  end  --------\n\n"
			en += fmt.Sprintf("\nSign data successfully! hex data is：%s", bodyhex)
			zh += fmt.Sprintf("\n签署数据成功！签名数据为：%s", bodyhex)
		}
		showSignDataHexCall()
		langChangeManager.SetText(txbodyshow, map[string]string{"en": en, "zh": zh})
	})

	// add item
	page.Add(input1)
	page.Add(input2)
	page.Add(button1)
	page.Add(txbodyshow)

	card := langChangeManager.NewCardSetTitle(title, page)
	box.Add(card)
}
