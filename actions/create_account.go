package actions

import (
	"encoding/hex"
	"fyne.io/fyne"
	"fyne.io/fyne/container"
	"fyne.io/fyne/widget"
	"github.com/hacash/core/account"
	"github.com/hacash/pcwallet/widgets"
	"strings"
)

func AddCanvasObjectCreateAccount(box *fyne.Container) {

	box.Add(widget.NewLabel("\n\n\n\n"))

	page := container.NewVBox()

	page.Add(widgets.NewTextWrapWordLabel("通过一个密码或者随机创建一个账户，强烈推荐随机创建账户！因为简单的密码将被人猜中你的私钥，导致你的代币丢失！密码仅支持字母大小写、数字和特殊符号，不支持空格、中文或其他字符。"))
	page.Add(widgets.NewTextWrapWordLabel("Through a password or create an account randomly, it is highly recommended to create an account randomly! Because a simple password will be guessed your private key, resulting in the loss of your token! Passwords only support upper and lower case letters, numbers and special symbols, and do not support spaces, Chinese or other characters."))

	input := widget.NewEntry()
	input.PlaceHolder = "这里输入密码 / input password"
	input.Wrapping = 2

	accshow := widget.NewMultiLineEntry()
	button1 := widget.NewButton("通过密码创建账户 / Create Account by Password", func() {
		if input.Text == "" {
			accshow.SetText("请输入密码 / Please input a password")
			return
		}
		// 密码合法性
		for _, v := range input.Text {
			if v < 33 || v > 126 {
				accshow.SetText("密码内含有不支持的字符 / The password contains unsupported characters")
				return
			}
		}

		// 通过密码创建账户
		accobj := account.CreateAccountByPassword(input.Text)
		showAccount(accshow, accobj)

	})

	button2 := widget.NewButton("随机创建账户 / Create Account Randomly", func() {

		// 随机创建账户
		accobj := account.CreateNewRandomAccount()
		showAccount(accshow, accobj)

	})

	page.Add(input)
	page.Add(button1)
	page.Add(button2)
	page.Add(accshow)

	card := widget.NewCard("创建账户 / Create Account", "", page)
	box.Add(card)

}

func showAccount(text *widget.Entry, acc *account.Account) {

	text.SetText("创建成功 / Created successfully :" +
		"\n\n[地址/Address] " + acc.AddressReadable +
		"\n\n[公钥/PublicKey] " + strings.ToUpper(hex.EncodeToString(acc.PublicKey)) +
		"\n\n[私钥/PrivateKey] " + strings.ToUpper(hex.EncodeToString(acc.PrivateKey)) +
		"\n")
}
