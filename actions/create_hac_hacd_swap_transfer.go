package actions

import (
	"encoding/hex"
	"fyne.io/fyne"
	"fyne.io/fyne/container"
	"fyne.io/fyne/widget"
	"github.com/hacash/core/account"
	"github.com/hacash/core/actions"
	"github.com/hacash/core/fields"
	"github.com/hacash/core/transactions"
	"github.com/hacash/pcwallet/widgets"
	"strconv"
	"time"
)

// curl http://192.168.2.108:3381/operatehex -X POST -d '00000001.....'

func AddOpenButtonOnMainOfCreateTransferHACswapHACD(box *fyne.Container, langChangeManager *widgets.LangChangeManager) {

	title := map[string]string{"en": "Create HAC ⇄ HACD swap transfer tx", "zh": "创建 HAC 与 HACD 互换转账交易"}

	button := langChangeManager.NewButton(title, func() {
		OpenWindowCreateTransferHACswapHACD(title, langChangeManager)
	})
	box.Add(button)
}

func OpenWindowCreateTransferHACswapHACD(title map[string]string, langChangeManager *widgets.LangChangeManager) fyne.Window {

	// 打开窗口测试
	testSize := fyne.Size{
		Width:  800,
		Height: 1000,
	}

	box := container.NewVBox()
	AddCanvasObjectCreateTransferHACswapHACD(title, box, langChangeManager)

	// 开启窗口
	return langChangeManager.NewWindowAndShow(title, &testSize, box)
}

func AddCanvasObjectCreateTransferHACswapHACD(title map[string]string, box *fyne.Container, langChangeManager *widgets.LangChangeManager) {

	page := container.NewVBox()

	page.Add(langChangeManager.NewTextWrapWordLabel(map[string]string{"en": "Create an atomic swap transaction between HAC and HACD. The transaction is an atomic transaction. It can only succeed at the same time and will not let one of the transfer parties fail.", "zh": "创建一笔 HAC 与 HACD 的原子互换转账交易，交易是原子事务，只可能同时成功，不会让其中转账一方失败。"}))

	input1 := langChangeManager.NewEntrySetPlaceHolder(map[string]string{"en": "HACD Payment (HAC Received) Address or Password or PrivateKey", "zh": "钻石转出（HAC收款）地址、密码或私钥"})
	input2 := langChangeManager.NewEntrySetPlaceHolder(map[string]string{"en": "HACD name list split by comma", "zh": "输入逗号隔开的钻石字面值列表"})
	input2.Wrapping = fyne.TextWrapWord
	input2.MultiLine = true
	input3 := langChangeManager.NewEntrySetPlaceHolder(map[string]string{"en": "HAC Payment (HACD Received) Address or Password or PrivateKey", "zh": "输入HAC转出（钻石收款）地址、密码或私钥"})
	input4 := langChangeManager.NewEntrySetPlaceHolder(map[string]string{"en": "HAC payment quantity - HAC (unit: 248) or use 'ㄜ1:248' format, Example: '0.25' or 'ㄜ25:246'", "zh": "输入HAC支付数量- HAC（单位：枚 - :248）也可直接使用 'ㄜ1:248' 格式， 例如：'0.25' or 'ㄜ25:246'"})
	input5 := langChangeManager.NewEntrySetPlaceHolder(map[string]string{"en": "Tx Fee Address or Password or PrivateKey", "zh": "输入支付交易手续费地址、密码或私钥"})
	input6 := langChangeManager.NewEntrySetPlaceHolder(map[string]string{"en": "Tx Fee - HAC (unit: 248) or use 'ㄜ1:248' format, Example: '0.25' or 'ㄜ25:246'", "zh": "输入交易手续费 - HAC（单位：枚 - :248）也可直接使用 'ㄜ1:248' 格式， 例如：'0.25' or 'ㄜ25:246'"})
	input6.SetText("ㄜ8:244")
	input7 := langChangeManager.NewEntrySetPlaceHolder(map[string]string{"en": "Optional: Tx timestamp", "zh": "选填：交易时间戳"})

	txbodyshow := widget.NewEntry()
	txbodyshow.MultiLine = true
	txbodyshow.Wrapping = fyne.TextWrapBreak

	button1 := langChangeManager.NewButton(map[string]string{"en": "Create HAC & HACD swap transfer Tx", "zh": "确认创建 HAC 和 HACD 原子互换交易"}, func() {
		var e error = nil
		if input2.Text == "" {
			langChangeManager.SetText(txbodyshow, map[string]string{"en": "Please HACD names split by comma", "zh": "请输入钻石列表"})
			return
		}
		//diamonds, e := transactions.CreateHACDlistBySplitCommaFromString(input2.Text)
		diamonds := fields.DiamondListMaxLen200{}
		e = diamonds.ParseHACDlistBySplitCommaFromString(input2.Text)
		if e != nil {
			en := "HACD names split by comma error: " + e.Error()
			zh := "钻石名称列表错误: " + e.Error()
			langChangeManager.SetText(txbodyshow, map[string]string{"en": en, "zh": zh})
			return
		}
		payHACDacc := account.GetAccountByPrivateKeyOrPassword(input1.Text)
		payHACDaddr, payHACDacc := parseAccountFromAddressOrPasswordOrPrivateKey(input1.Text)
		payHACaddr, payHACacc := parseAccountFromAddressOrPasswordOrPrivateKey(input3.Text)
		hacAmt, e := fields.NewAmountFromString(input4.Text)
		if e != nil {
			langChangeManager.SetText(txbodyshow, map[string]string{"en": "HAC payment quantity", "zh": "HAC数量格式错误"})
			return
		}
		feeaddr, feeacc := parseAccountFromAddressOrPasswordOrPrivateKey(input5.Text)
		fee, e4 := fields.NewAmountFromString(input6.Text)
		if e4 != nil {
			langChangeManager.SetText(txbodyshow, map[string]string{"en": "Tx Fee format error", "zh": "交易手续费格式错误"})
			return
		}
		if len(fee.Numeral) > 2 {
			langChangeManager.SetText(txbodyshow, map[string]string{"en": "Tx Fee digits too long", "zh": "手续费位数过长"})
			return
		}
		usetime := time.Now().Unix()
		if len(input7.Text) > 0 {
			its, e1 := strconv.ParseInt(input7.Text, 10, 0)
			if e1 != nil {
				langChangeManager.SetText(txbodyshow, map[string]string{"en": "Timestamp format error", "zh": "时间戳格式错误"})
				return
			}
			usetime = its
		}

		if feeaddr.NotEqual(*payHACaddr) {
			langChangeManager.SetText(txbodyshow, map[string]string{"en": "Pay fee address and HAC transfer out address must be the same", "zh": "手续费支付地址与HAC转出地址必须相同"})
			return
		}

		// 创建交易
		tx, e0 := transactions.NewEmptyTransaction_2_Simple(*feeaddr)
		if e0 != nil {
			langChangeManager.SetText(txbodyshow, map[string]string{"en": "Transaction creation failed: \n\n" + e0.Error(), "zh": "交易创建失败: \n\n" + e0.Error()})
			return
		}
		if tx == nil {
			langChangeManager.SetText(txbodyshow, map[string]string{"en": "Transaction creation failed", "zh": "交易创建失败"})
			return
		}
		tx.Timestamp = fields.BlockTxTimestamp(usetime)
		tx.Fee = *fee
		// HACD 转账
		hacdact := &actions.Action_6_OutfeeQuantityDiamondTransfer{
			FromAddress: *payHACDaddr,
			ToAddress:   *payHACaddr,
			DiamondList: diamonds,
		}
		// 添加 HAC 支付
		hacact := actions.NewAction_1_SimpleToTransfer(*payHACDaddr, hacAmt)
		// 添加antions
		tx.AppendAction(hacdact)
		tx.AppendAction(hacact)

		// 判断签名 sign 私钥签名
		if payHACDacc != nil {
			tx.FillTargetSign(payHACDacc)
		}
		if payHACacc != nil {
			tx.FillTargetSign(payHACacc)
		}
		if feeacc != nil {
			tx.FillTargetSign(feeacc)
		}
		// 创建成功
		txbody, e3 := tx.Serialize()
		if e3 != nil {
			langChangeManager.SetText(txbodyshow, map[string]string{"en": "Transaction creation failed", "zh": "交易创建失败"})
			return
		}

		txbodyhex := "\n\n-------- txbody hex start --------\n"
		txbodyhex += hex.EncodeToString(txbody)
		txbodyhex += "\n-------- txbody hex  end  --------\n\n"

		resEn := "HAC & HACD swap transaction created successfully!" +
			"\nPlease copy the following [txbody] to sign the tx then submit transaction in online wallet:" +
			"\n\n[txhash] " + tx.Hash().ToHex() +
			"\n\n[timestamp] " + strconv.FormatInt(usetime, 10) +
			"\n\n[txbody] " + txbodyhex
		resZh := "HAC和HACD原子互换交易创建成功！" +
			"\n请复制下面 [交易体/txbody] 内容，先完成签名操作，然后去在线钱包提交交易:" +
			"\n\n[交易哈希/txhash] " + tx.Hash().ToHex() +
			"\n\n[时间戳/timestamp] " + strconv.FormatInt(usetime, 10) +
			"\n\n[交易体/txbody] " + txbodyhex

		// 签名检查
		checkTxSignStatus(tx, &resEn, &resZh)

		langChangeManager.SetText(txbodyshow, map[string]string{"en": resEn, "zh": resZh})

	})

	page.Add(input1)
	page.Add(input2)
	page.Add(input3)
	page.Add(input4)
	page.Add(input5)
	page.Add(input6)
	page.Add(input7)

	page.Add(button1)
	page.Add(txbodyshow)

	card := langChangeManager.NewCardSetTitle(title, page)
	box.Add(card)

}
