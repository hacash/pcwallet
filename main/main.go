package main

import (
	"fyne.io/fyne"
	"fyne.io/fyne/app"
	"fyne.io/fyne/container"
	"fyne.io/fyne/widget"
	"github.com/flopp/go-findfont"
	"github.com/hacash/pcwallet/actions"
	"net/url"
	"os"
	"strings"
)

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
	w := a.NewWindow("Hacash Offline PC Wallet / Hacash 离线电脑钱包")

	w.Resize(fyne.Size{
		Width:  1000,
		Height: 800,
	})

	objs := container.NewVBox(widget.NewLabel("欢迎使用 Hacash 离线安全钱包，\n本程序包含创建账户、生成 HAC、BTC 或 HACD 转账交易、开启关闭通道等与私钥安全相关的功能。\n无安全问题的查询余额、查询交易、提交签名后交易等等功能请使用在线钱包："))

	objs.Add(widget.NewLabel("Welcome to use hacash offline security wallet. \nThis program includes creating account, generating HAC, BTC or HACD \ntransfer transaction, opening and closing channel and \nother functions related to private key security. \nFor the functions of no security issues, such as balance inquiry, \ntransaction inquiry and transaction after submitting signature, \nplease use online Wallet:"))
	online_wallet_url, _ := url.Parse("https://wallet.hacash.org")
	objs.Add(widget.NewHyperlink("https://wallet.hacash.org", online_wallet_url))

	objs.Add(widget.NewLabel("如果你需要测试转账或者捐赠本钱包的开发者，请向以下地址转账：\nIf you need to test the transfer or donate the wallet to the developer, \nplease transfer to the following address:"))
	donate_url, _ := url.Parse("http://block.hacash.org/address/1K7vNCyurJXEPzFnauJE1t2S5tDdjkAaZY")
	objs.Add(widget.NewHyperlink("1K7vNCyurJXEPzFnauJE1t2S5tDdjkAaZY", donate_url))

	// 创建账户
	actions.AddCanvasObjectCreateAccount(objs)

	// 创建 HAC 普通转账交易
	actions.AddCanvasObjectCreateTransferHAC(objs)
	// 创建 BTC 普通转账交易

	objs.Add(widget.NewLabel("\n\n\n"))

	// 页面翻动
	scroll := container.NewVScroll(objs)

	w.SetContent(scroll)

	w.ShowAndRun()

	// 回退字体设置
	os.Unsetenv("FYNE_FONT")
}
