package actions

import (
	"encoding/hex"
	fyne "fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/dialog"
	"fyne.io/fyne/v2/widget"
	"github.com/hacash/core/account"
	"github.com/hacash/pcwallet/widgets"
	"strings"
)

func AddOpenButtonOnMainOfCreateAccount(box *fyne.Container, langChangeManager *widgets.LangChangeManager) {
	title := map[string]string{"en": "Create Account", "zh": "创建账户"}

	button := langChangeManager.NewButton(title, func() {
		OpenWindowCreateAccount(title, langChangeManager)
	})
	box.Add(button)
}

var createAccountWindow fyne.Window

func OpenWindowCreateAccount(title map[string]string, langChangeManager *widgets.LangChangeManager) fyne.Window {
	// 打开窗口测试
	testSize := fyne.Size{
		Width:  800,
		Height: 500,
	}

	box := container.NewVBox()
	AddCanvasObjectCreateAccount(title, box, langChangeManager)

	// 开启窗口
	createAccountWindow = langChangeManager.NewWindowAndShow(title, &testSize, box)
	return createAccountWindow
}

func AddCanvasObjectCreateAccount(title map[string]string, box *fyne.Container, langChangeManager *widgets.LangChangeManager) {
	page := container.NewVBox()

	page.Add(langChangeManager.NewTextWrapWordLabel(map[string]string{"en": "Through a password or create an account randomly, it is highly recommended to create an account randomly! Because a simple password will be guessed your private key, resulting in the loss of your token! Passwords only support upper and lower case letters, numbers and special symbols, and do not support spaces, Chinese or other characters.", "zh": "通过一个密码或者随机创建一个账户，强烈推荐随机创建账户！因为简单的密码将被人猜中你的私钥，导致你的代币丢失！密码仅支持字母大小写、数字和特殊符号，不支持空格、中文或其他字符。"}))

	input := langChangeManager.NewEntrySetPlaceHolder(map[string]string{"en": "Input password", "zh": "这里输入密码或私钥"})

	accshow := widget.NewMultiLineEntry()
	button1 := langChangeManager.NewButton(map[string]string{"en": "Create Account by Password", "zh": "通过密码或私钥创建账户"}, func() {
		if input.Text == "" {
			langChangeManager.SetText(accshow, map[string]string{"en": "Please input a password", "zh": "请输入密码或私钥"})
			return
		}
		// 密码合法性
		for _, v := range input.Text {
			if v < 33 || v > 126 {
				langChangeManager.SetText(accshow, map[string]string{"en": "The password contains unsupported characters", "zh": "密码内含有不支持的字符"})
				return
			}
		}

		// 通过密码创建账户
		accobj := account.GetAccountByPrivateKeyOrPassword(input.Text)
		showAccount(langChangeManager, accshow, accobj)

	})

	button2 := langChangeManager.NewButton(map[string]string{"en": "Create Account Randomly", "zh": "随机创建账户"}, func() {

		// 随机创建账户
		accobj := account.CreateNewRandomAccount()
		showAccount(langChangeManager, accshow, accobj)

	})

	var eabpswd *widget.Check
	var pswdcnf *dialog.ConfirmDialog
	var create_pswdcnf = func() *dialog.ConfirmDialog {
		if pswdcnf == nil {
			pswdcnf = dialog.NewConfirm("Security Warning", "Creating by password will make \naccount security issues \n(your assets may be stolen), \nare you sure you are aware of the risks?", func(ok bool) {
				if ok {
					input.Show()
					button1.Show()
					eabpswd.Hide()
				} else {
					eabpswd.SetChecked(false)
					//pswdcnf.Hide()
				}
			}, createAccountWindow)
		}
		return pswdcnf
	}
	eabpswd = langChangeManager.NewCheck(map[string]string{"en": "Enable password creation", "zh": "启用密码创建模式"}, func(ok bool) {
		pswdcnf = create_pswdcnf()
		pswdcnf.Show()
	})

	page.Add(input)
	page.Add(button1)
	input.Hide()
	button1.Hide()
	//fmt.Println(input.Text, button1.Text)
	page.Add(button2)
	page.Add(accshow)

	card := langChangeManager.NewCardSetTitle(title, page)
	box.Add(card)
	page.Add(eabpswd)
}

func showAccount(langChangeManager *widgets.LangChangeManager, text *widget.Entry, acc *account.Account) {
	en := "Created successfully :" +
		"\n\n[Address] " + acc.AddressReadable +
		"\n\n[PublicKey] " + strings.ToUpper(hex.EncodeToString(acc.PublicKey)) +
		"\n\n[PrivateKey] " + strings.ToUpper(hex.EncodeToString(acc.PrivateKey)) +
		"\n"
	zh := "创建成功:" +
		"\n\n[地址] " + acc.AddressReadable +
		"\n\n[公钥] " + strings.ToUpper(hex.EncodeToString(acc.PublicKey)) +
		"\n\n[私钥] " + strings.ToUpper(hex.EncodeToString(acc.PrivateKey)) +
		"\n"

	langChangeManager.SetText(text, map[string]string{"en": en, "zh": zh})
}
