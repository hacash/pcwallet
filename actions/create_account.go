package actions

import (
	"encoding/hex"
	"fyne.io/fyne"
	"fyne.io/fyne/container"
	"fyne.io/fyne/widget"
	"github.com/hacash/core/account"
	"strings"
)

func AddCanvasObjectCreateAccount(box *fyne.Container) {

	box.Add(widget.NewLabel("\n"))

	page := container.NewVBox()

	page.Add(widget.NewLabel("通过一个密码或者随机创建一个账户，强烈推荐随机创建账户！\n因为简单的密码将被人猜中你的私钥，导致你的代币丢失！\nThrough a password or create an account randomly, \nit is highly recommended to create an account randomly! \nBecause a simple password will be guessed your private key, \nresulting in the loss of your token!"))

	input := widget.NewEntry()
	input.MultiLine = true
	input.PlaceHolder = "这里输入密码 / input password"
	input.Wrapping = 2

	accshow := widget.NewEntry()
	accshow.MultiLine = true
	button1 := widget.NewButton("通过密码创建账户 / Create Account by Password", func() {
		if input.Text == "" {
			accshow.SetText("请输入密码 / Please input a password")
			return
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
