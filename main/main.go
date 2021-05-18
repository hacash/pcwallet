package main

import (
	"fyne.io/fyne"
	"fyne.io/fyne/app"
	"fyne.io/fyne/container"
	"fyne.io/fyne/widget"
	"github.com/flopp/go-findfont"
	"github.com/hacash/pcwallet/actions"
	"github.com/hacash/pcwallet/widgets"
	"net/url"
	"os"
	"strings"
)

/*

go build -o test/pcwallet pcwallet/main/main.go  && ./test/pcwallet


*/

func init() {

	// 中文字体支持
	fontPaths := findfont.List()
	for _, path := range fontPaths {
		if strings.Contains(path, "simkai.ttf") ||
			strings.Contains(path, "simhei.ttf") ||
			strings.Contains(path, "simsun.ttf") {
			os.Setenv("FYNE_FONT", path)
			break
		}
	}

}

func main() {

	a := app.New()
	w := a.NewWindow("Hacash Offline PC Wallet")
	windowSize := fyne.Size{
		Width:  500,
		Height: 800,
	}
	w.Resize(windowSize)

	objs := container.NewVBox()

	langChangeManager := widgets.NewLangChangeManager(a)

	langchange := widget.NewRadioGroup([]string{"English", "简体中文"}, func(s string) {
		if s == "English" {
			langChangeManager.ChangeLangByName("en")
			w.SetTitle("Hacash Offline PC Wallet")
		} else {
			langChangeManager.ChangeLangByName("zh")
			w.SetTitle("Hacash 离线电脑钱包")
		}
	})
	langchange.Horizontal = true
	langchange.Selected = "English"
	objs.Add(langchange)

	// label

	lb1 := langChangeManager.NewTextWrapWordLabel(map[string]string{"en": "Welcome to use hacash offline security wallet. This program includes creating account, generating HAC, BTC or HACD transfer transaction, opening and closing channel and other functions related to private key security. For the functions of no security issues, such as balance inquiry, transaction inquiry and transaction after submitting signature, please use online Wallet:", "zh": "欢迎使用 Hacash 离线安全钱包，本程序包含创建账户、生成 HAC、BTC 或 HACD 转账交易、开启关闭通道等与私钥安全相关的功能。无安全问题的查询余额、查询交易、提交签名后交易等等功能请使用在线钱包："})
	objs.Add(lb1)

	online_wallet_url, _ := url.Parse("https://wallet.hacash.org")
	objs.Add(widget.NewHyperlink("https://wallet.hacash.org", online_wallet_url))

	donate_address := "1QDc1twwVy3acuftAv3GuNnKwxopYi9VLb"

	lb3 := langChangeManager.NewTextWrapWordLabel(map[string]string{"en": "If you need to test the transfer or donate the wallet to the developer, please transfer to the following address:", "zh": "如果你需要测试转账或者捐赠本钱包的开发者，请向以下地址转账："})
	objs.Add(lb3)
	//donate_url, _ := url.Parse("https://explorer.hacash.org/address/" + donate_address)
	//objs.Add(widget.NewHyperlink(donate_address, donate_url))
	// 可复制输入框
	donate_address_input := widget.NewEntry()
	donate_address_input.Disable()
	donate_address_input.SetText(donate_address)
	objs.Add(donate_address_input)

	//
	prettl := langChangeManager.NewTextWrapWordLabel(map[string]string{"en": "\nPlease click the button below to use the wallet function:\n", "zh": "\n请点击以下按钮使用钱包功能:\n"})
	objs.Add(prettl)

	// 创建账户
	actions.AddOpenButtonOnMainOfCreateAccount(objs, langChangeManager)
	// 创建 HAC 普通转账交易
	actions.AddOpenButtonOnMainOfCreateTransferHAC(objs, langChangeManager)
	// 创建 HACD 转账交易
	actions.AddOpenButtonOnMainOfCreateTransferHACD(objs, langChangeManager)
	// 创建 BTC 转账交易
	actions.AddOpenButtonOnMainOfCreateTransferBTC(objs, langChangeManager)
	// 创建开启、关闭通道交易
	actions.AddOpenButtonOnMainOfCreateTxOpenChannel(objs, langChangeManager)
	actions.AddOpenButtonOnMainOfCreateTxCloseChannel(objs, langChangeManager)
	// 签署交易
	actions.AddOpenButtonOnMainOfSignTx(objs, langChangeManager)
	// 检查一笔交易的结构和内容
	actions.AddOpenButtonOnMainOfCheckTxContents(objs, langChangeManager)

	objs.Add(widget.NewLabel("\n\n"))

	// 页面翻动
	scroll := container.NewVScroll(objs)

	w.SetContent(scroll)

	w.Show()

	a.Run()

	// 回退字体设置
	os.Unsetenv("FYNE_FONT")

}
