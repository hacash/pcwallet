package actions

import (
	"encoding/hex"
	"fmt"
	"fyne.io/fyne"
	"fyne.io/fyne/container"
	"fyne.io/fyne/widget"
	"github.com/hacash/core/actions"
	"github.com/hacash/core/interfaces"
	"github.com/hacash/core/transactions"
	"github.com/hacash/pcwallet/widgets"
	"strings"
	"time"
)

func AddOpenButtonOnMainOfCheckTxContents(box *fyne.Container, langChangeManager *widgets.LangChangeManager) {

	title := map[string]string{"en": "Check tx content", "zh": "验证交易内容"}

	button := langChangeManager.NewButton(title, func() {
		OpenWindowCheckTxContents(title, langChangeManager)
	})
	box.Add(button)
}

func OpenWindowCheckTxContents(title map[string]string, langChangeManager *widgets.LangChangeManager) fyne.Window {

	// 打开窗口测试
	testSize := fyne.Size{
		Width:  800,
		Height: 1000,
	}

	box := container.NewVBox()
	AddCanvasObjectCheckTxContents(title, box, langChangeManager)

	// 开启窗口
	return langChangeManager.NewWindowAndShow(title, &testSize, box)
}

func AddCanvasObjectCheckTxContents(title map[string]string, box *fyne.Container, langChangeManager *widgets.LangChangeManager) {

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

	card := langChangeManager.NewCardSetTitle(title, page)
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
}

// 解析 action 的描述
func renderTxActionDescribe(mainaddr string, act interfaces.Action) (string, string) {
	var en, zh string
	var actId = act.Kind()
	// 每一项条款
	if a, ok := act.(*actions.Action_1_SimpleToTransfer); ok {

		/**************** Action_1_SimpleToTransfer ****************/
		toaddr := a.ToAddress.ToReadable()
		amt := a.Amount.ToFinString()
		en += fmt.Sprintf("Simple transfer: Account <%s> transfers amount <%s> to account <%s>", mainaddr, amt, toaddr)
		zh += fmt.Sprintf("普通 HAC 转账： 地址 <%s> 向地址 <%s> 转账 <%s>", mainaddr, toaddr, amt)

	} else if a, ok := act.(*actions.Action_13_FromTransfer); ok {

		/**************** Action_13_FromTransfer ****************/
		fromaddr := a.FromAddress.ToReadable()
		amt := a.Amount.ToFinString()
		en += fmt.Sprintf("HAC From transfer: Account <%s> transfers amount <%s> to account <%s>", fromaddr, amt, mainaddr)
		zh += fmt.Sprintf("HAC From 转账： 地址 <%s> 向地址 <%s> 转账 <%s>", fromaddr, mainaddr, amt)

	} else if a, ok := act.(*actions.Action_14_FromToTransfer); ok {

		/**************** Action_14_FromToTransfer ****************/
		fromaddr := a.FromAddress.ToReadable()
		toaddr := a.ToAddress.ToReadable()
		amt := a.Amount.ToFinString()
		en += fmt.Sprintf("HAC From -> To transfer: Account <%s> transfers amount <%s> to account <%s>", fromaddr, amt, toaddr)
		zh += fmt.Sprintf("HAC From -> To 转账： 地址 <%s> 向地址 <%s> 转账 <%s>", fromaddr, toaddr, amt)

	} else if a, ok := act.(*actions.Action_2_OpenPaymentChannel); ok {

		/**************** Action_2_OpenPaymentChannel *************/
		cid := hex.EncodeToString(a.ChannelId)
		addr1 := a.LeftAddress.ToReadable()
		amt1 := a.LeftAmount.ToFinString()
		addr2 := a.RightAddress.ToReadable()
		amt2 := a.RightAmount.ToFinString()
		en += fmt.Sprintf("Open payment channel: Open ID = <%s> channel, account <%s> and <%s> respective deposit <%s> and <%s> into channel", cid, addr1, addr2, amt1, amt2)
		zh += fmt.Sprintf("开启支付通道： 开启 ID 为 <%s> 的通道，账户 <%s> 和 <%s> 分别各自向通道内存入 <%s> 及 <%s>", cid, addr1, addr2, amt1, amt2)

	} else if a, ok := act.(*actions.Action_3_ClosePaymentChannel); ok {

		/**************** Action_3_ClosePaymentChannel ************/
		cid := hex.EncodeToString(a.ChannelId)
		en += fmt.Sprintf("Close payment channel: Close ID = <%s> channel", cid)
		zh += fmt.Sprintf("关闭支付通道： 关闭 ID 为 <%s> 的通道", cid)

	} else if a, ok := act.(*actions.Action_4_DiamondCreate); ok {

		/**************** Action_4_DiamondCreate *******************/
		dianame := string(a.Diamond)
		rwdaddr := a.Address.ToReadable()
		en += fmt.Sprintf("Mint block diamond: name <%s>, number <%d>, miner account <%s>", dianame, a.Number, rwdaddr)
		zh += fmt.Sprintf("铸造区块钻石： 字面值 <%s>, 序号 <%d>, 矿工账户 <%s>", dianame, a.Number, rwdaddr)

	} else if a, ok := act.(*actions.Action_5_DiamondTransfer); ok {

		/**************** Action_5_DiamondTransfer *****************/
		toaddr := a.ToAddress.ToReadable()
		dianame := string(a.Diamond)
		en += fmt.Sprintf("Transfer diamond: name <%s>, collection account <%s>", dianame, toaddr)
		zh += fmt.Sprintf("区块钻石转账： 字面值 <%s>, 收取账户 <%s>", dianame, toaddr)

	} else if a, ok := act.(*actions.Action_6_OutfeeQuantityDiamondTransfer); ok {

		/**************** Action_5_DiamondTransfer *****************/
		fromaddr := a.FromAddress.ToReadable()
		toaddr := a.ToAddress.ToReadable()
		dianames := a.GetDiamondNamesSplitByComma() // 名称列表
		en += fmt.Sprintf("Batch transfer diamonds: account <%s> transfer %d diamonds to account <%s> names is <%s>", fromaddr, a.DiamondCount, toaddr, dianames)
		zh += fmt.Sprintf("区块钻石批量转账： 账户 <%s> 向账户 <%s> 转移字面值为 <%s> 的 %d 枚钻石", fromaddr, toaddr, dianames, a.DiamondCount)

	} else if a, ok := act.(*actions.Action_7_SatoshiGenesis); ok {

		/**************** Action_7_SatoshiGenesis *****************/
		btn := uint32(a.BitcoinQuantity)
		addr := a.OriginAddress.ToReadable()
		en += fmt.Sprintf("Bitcoin genesis move: <%d> bitcoin move by account <%s>", btn, addr)
		zh += fmt.Sprintf("比特币单向转移： <%d> 枚比特币被账户 <%s> 转移进来", btn, addr)

	} else if a, ok := act.(*actions.Action_8_SimpleSatoshiTransfer); ok {

		/**************** Action_8_SimpleSatoshiTransfer ***********/
		toaddr := a.ToAddress.ToReadable()
		satamt := uint64(a.Amount)
		en += fmt.Sprintf("Satoshi simple transfer: Account <%s> transfers amount <%s> SAT to account <%s>", mainaddr, satamt, toaddr)
		zh += fmt.Sprintf("比特币普通转账： 地址 <%s> 向地址 <%s> 转账 <%s> SAT", mainaddr, toaddr, satamt)

	} else if a, ok := act.(*actions.Action_9_LockblsCreate); ok {

		/**************** Action_9_LockblsCreate ***********/
		lid := hex.EncodeToString(a.LockblsId)
		payaddr := a.PaymentAddress.ToReadable()
		gotaddr := a.MasterAddress.ToReadable()
		hei1 := uint64(a.EffectBlockHeight)
		num1 := uint64(a.LinearBlockNumber)
		ttl1 := a.TotalStockAmount.ToFinString()
		rls1 := a.LinearReleaseAmount.ToFinString()
		en += fmt.Sprintf("Create linear lock release HAC contract: lock_id <%s>, lock total amount <%s>, deduction account <%s>, income account <%s>, effect block height <%d>, linear block number<%d>, linear release amount <%s>", lid, ttl1, payaddr, gotaddr, hei1, num1, rls1)
		zh += fmt.Sprintf("创建 HAC 线性锁仓释放合约： 锁仓ID <%s>, 总共锁入金额 <%s>, 锁仓扣款账户 <%s>, 提取权益账户 <%s>, 生效区块高度 <%d>, 线性提取区块间隔 <%d>, 单次可提取金额 <%s>", lid, ttl1, payaddr, gotaddr, hei1, num1, rls1)

	} else if a, ok := act.(*actions.Action_10_LockblsRelease); ok {

		/**************** Action_10_LockblsRelease ***********/
		lid := hex.EncodeToString(a.LockblsId)
		amt := a.ReleaseAmount.ToFinString()
		en += fmt.Sprintf("Release the locked HAC: Release lock_id = <%s> amount <%s>", lid, amt)
		zh += fmt.Sprintf("释放锁仓的货币： 释放 lock_id = <%s> 的线性锁仓 HAC, 释放数额为 <%s>", lid, amt)

	} else if a, ok := act.(*actions.Action_11_FromToSatoshiTransfer); ok {

		/**************** Action_11_FromToSatoshiTransfer ***********/
		fromaddr := a.FromAddress.ToReadable()
		toaddr := a.FromAddress.ToReadable()
		satamt := uint64(a.Amount)
		en += fmt.Sprintf("Satoshi simple transfer: Account <%s> transfers amount <%s> SAT to account <%s>", fromaddr, satamt, toaddr)
		zh += fmt.Sprintf("比特币普通转账： 地址 <%s> 向地址 <%s> 转账 <%s> SAT", fromaddr, toaddr, satamt)

	} else {

		/************************ Other ************************/
		en += fmt.Sprintf("Other action <action_kind: %d>", actId)
		zh += fmt.Sprintf("其他内容 <action_kind:%d>", actId)
	}
	/************************ END ************************/

	// 返回
	return en, zh
}
