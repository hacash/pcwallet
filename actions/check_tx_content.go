package actions

import (
	"encoding/hex"
	"fmt"
	fyne "fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/widget"
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
	acts := trs.GetActionList()
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

	// 格式化
	fmtEnZh := func(enstr, zhstr string, params ...interface{}) {
		en += fmt.Sprintf(enstr, params...)
		zh += fmt.Sprintf(zhstr, params...)
	}

	// 每一项条款
	if a, ok := act.(*actions.Action_1_SimpleToTransfer); ok {
		/**************** Action_1_SimpleToTransfer ****************/
		toaddr := a.ToAddress.ToReadable()
		amt := a.Amount.ToFinString()
		fmtEnZh("Simple transfer: Account <%s> transfers to account <%s> for amount <%s>",
			"普通 HAC 转账： 地址 <%s> 向地址 <%s> 转账 <%s>",
			mainaddr, toaddr, amt)
	} else if a, ok := act.(*actions.Action_13_FromTransfer); ok {
		/**************** Action_13_FromTransfer ****************/
		fromaddr := a.FromAddress.ToReadable()
		amt := a.Amount.ToFinString()
		fmtEnZh("HAC From transfer: Account <%s> transfers to account <%s> for amount <%s>",
			"HAC From 转账： 地址 <%s> 向地址 <%s> 转账 <%s>",
			fromaddr, mainaddr, amt)
	} else if a, ok := act.(*actions.Action_14_FromToTransfer); ok {
		/**************** Action_14_FromToTransfer ****************/
		fromaddr := a.FromAddress.ToReadable()
		toaddr := a.ToAddress.ToReadable()
		amt := a.Amount.ToFinString()
		fmtEnZh("HAC From -> To transfer: Account <%s> transfers to account <%s> for amount <%s>",
			"HAC From -> To 转账： 地址 <%s> 向地址 <%s> 转账 <%s>",
			fromaddr, toaddr, amt)
	} else if a, ok := act.(*actions.Action_2_OpenPaymentChannel); ok {
		/**************** Action_2_OpenPaymentChannel *************/
		cid := a.ChannelId.ToHex()
		addr1 := a.LeftAddress.ToReadable()
		amt1 := a.LeftAmount.ToFinString()
		addr2 := a.RightAddress.ToReadable()
		amt2 := a.RightAmount.ToFinString()
		fmtEnZh("Open payment channel: Open ID = <%s> channel, account <%s> and <%s> respective deposit <%s> and <%s> into channel",
			"开启支付通道： 开启 ID 为 <%s> 的通道，账户 <%s> 和 <%s> 分别各自向通道内存入 <%s> 及 <%s>",
			cid, addr1, addr2, amt1, amt2)
	} else if a, ok := act.(*actions.Action_3_ClosePaymentChannel); ok {

		/**************** Action_3_ClosePaymentChannel ************/
		fmtEnZh("Close payment channel: Close ID = <%s> channel, Balance allocation is the same as deposit",
			"关闭支付通道： 关闭 ID 为 <%s> 的通道，余额分配与存入相同",
			a.ChannelId.ToHex())

	} else if a, ok := act.(*actions.Action_21_ClosePaymentChannelBySetupOnlyLeftAmount); ok {

		/**************** Action_21_ClosePaymentChannelBySetupOnlyLeftAmount ************/
		fmtEnZh("Close payment channel: Close ID = <%s> channel, Left Amount: %sHAC %dSAT, The remaining part is automatically allocated on the right.",
			"关闭支付通道： 关闭 ID 为 <%s> 的通道，设定余额分配为：左 %sHAC %dSAT ，右侧自动分配剩余部分",
			a.ChannelId.ToHex(), a.LeftAmount.ToFinString(), a.LeftSatoshi.GetRealSatoshi())

	} else if a, ok := act.(*actions.Action_4_DiamondCreate); ok {

		/**************** Action_4_DiamondCreate *******************/
		fmtEnZh("Mint block diamond: name <%s>, number <%d>, miner account <%s>", "铸造区块钻石： 字面值 <%s>, 序号 <%d>, 矿工账户 <%s>", string(a.Diamond), a.Number, a.Address.ToReadable())

	} else if a, ok := act.(*actions.Action_5_DiamondTransfer); ok {

		/**************** Action_5_DiamondTransfer *****************/
		fmtEnZh("Transfer diamond: name <%s>, collection account <%s>",
			"区块钻石转账： 字面值 <%s>, 收取账户 <%s>",
			a.ToAddress.ToReadable(), string(a.Diamond))

	} else if a, ok := act.(*actions.Action_6_OutfeeQuantityDiamondTransfer); ok {

		/**************** Action_5_DiamondTransfer *****************/
		fromaddr := a.FromAddress.ToReadable()
		toaddr := a.ToAddress.ToReadable()
		dianames := a.GetDiamondNamesSplitByComma() // 名称列表
		fmtEnZh("Batch transfer diamonds: account <%s> transfer %d diamonds to account <%s> names is <%s>",
			"区块钻石批量转账： 账户 <%s> 向账户 <%s> 转移字面值为 <%s> 的 %d 枚钻石",
			fromaddr, toaddr, dianames, a.DiamondList.Count)

	} else if a, ok := act.(*actions.Action_7_SatoshiGenesis); ok {

		/**************** Action_7_SatoshiGenesis *****************/
		fmtEnZh("Bitcoin genesis move: <%d> bitcoin move by account <%s>",
			"比特币单向转移： <%d> 枚比特币被账户 <%s> 转移进来",
			uint32(a.BitcoinQuantity), a.OriginAddress.ToReadable())

	} else if a, ok := act.(*actions.Action_8_SimpleSatoshiTransfer); ok {

		/**************** Action_8_SimpleSatoshiTransfer ***********/
		fmtEnZh("Satoshi simple transfer: Account <%s> transfers to account <%s> with amount <%d> SAT",
			"比特币普通转账： 地址 <%s> 向地址 <%s> 转账 <%d> SAT", mainaddr, a.ToAddress.ToReadable(), uint64(a.Amount))

	} else if a, ok := act.(*actions.Action_9_LockblsCreate); ok {

		/**************** Action_9_LockblsCreate ***********/
		lid := hex.EncodeToString(a.LockblsId)
		payaddr := a.PaymentAddress.ToReadable()
		gotaddr := a.MasterAddress.ToReadable()
		hei1 := uint64(a.EffectBlockHeight)
		num1 := uint64(a.LinearBlockNumber)
		ttl1 := a.TotalStockAmount.ToFinString()
		rls1 := a.LinearReleaseAmount.ToFinString()
		fmtEnZh("Create linear lock release HAC contract: lock_id <%s>, lock total amount <%s>, deduction account <%s>, income account <%s>, effect block height <%d>, linear block number<%d>, linear release amount <%s>",
			"创建 HAC 线性锁仓释放合约： 锁仓ID <%s>, 总共锁入金额 <%s>, 锁仓扣款账户 <%s>, 提取权益账户 <%s>, 生效区块高度 <%d>, 线性提取区块间隔 <%d>, 单次可提取金额 <%s>",
			lid, ttl1, payaddr, gotaddr, hei1, num1, rls1)

	} else if a, ok := act.(*actions.Action_10_LockblsRelease); ok {

		/**************** Action_10_LockblsRelease ***********/
		fmtEnZh("Release the locked HAC: Release lock_id = <%s> amount <%s>",
			"释放锁仓的货币： 释放 lock_id = <%s> 的线性锁仓 HAC, 释放数额为 <%s>",
			a.LockblsId.ToHex(), a.ReleaseAmount.ToFinString())

	} else if a, ok := act.(*actions.Action_11_FromToSatoshiTransfer); ok {

		/**************** Action_11_FromToSatoshiTransfer ***********/
		fmtEnZh("Satoshi transfer: Account <%s> transfers to account <%s> with amount <%d> SAT ",
			"比特币普通转账： 地址 <%s> 向地址 <%s> 转账 <%d> SAT",
			a.FromAddress.ToReadable(), a.ToAddress.ToReadable(), uint64(a.Amount))

	} else if a, ok := act.(*actions.Action_28_FromSatoshiTransfer); ok {

		/**************** Action_28_FromSatoshiTransfer ***********/
		fmtEnZh("Satoshi transfer: Account <%s> transfers to account <%s> with amount <%d> SAT ",
			"比特币转账： 地址 <%s> 向地址 <%s> 转账 <%d> SAT",
			a.FromAddress.ToReadable(), mainaddr, uint64(a.Amount))

	} else if a, ok := act.(*actions.Action_22_UnilateralClosePaymentChannelByNothing); ok {

		/**************** Action_22_UnilateralClosePaymentChannelByNothing ***********/
		fmtEnZh("Close the channel <%s> unilaterally without a bill, assert address: %s",
			"无票据单方面关闭通道 <%s> ，主张地址：%s",
			a.ChannelId.ToHex(), a.AssertCloseAddress.ToReadable())

	} else if a, ok := act.(*actions.Action_23_UnilateralCloseOrRespondChallengePaymentChannelByRealtimeReconciliation); ok {

		/**************** Action_23_UnilateralCloseOrRespondChallengePaymentChannelByRealtimeReconciliation ***********/
		rec := a.Reconciliation
		billsignck := "[fail]"
		if e := a.Reconciliation.CheckAddressAndSign(rec.LeftSign.GetAddress(), rec.RightSign.GetAddress()); e == nil {
			billsignck = "[OK]"
		}
		ens := "Submit realtime reconciliation, close channels or respond to challenges, ReuseVersion: %d, BillAutoNumber: %d, ChannelId: <%s>, AssertAddress: %s, Bill signature check: %s"
		zhs := "提交实时对账单，关闭通道或回应挑战，重用版本：%d, 账单序列号：%d，通道ID：<%s> ，主张地址：%s，票据签名检查：%s"
		fmtEnZh(ens, zhs, rec.ReuseVersion, rec.BillAutoNumber, rec.ChannelId.ToHex(), a.AssertAddress.ToReadable(), billsignck)

	} else if a, ok := act.(*actions.Action_24_UnilateralCloseOrRespondChallengePaymentChannelByChannelChainTransferBody); ok {

		/**************** Action_24_UnilateralCloseOrRespondChallengePaymentChannelByChannelChainTransferBody ***********/
		rec := a.ChannelChainTransferTargetProveBody
		billsignck := "[fail]"
		if e := rec.CheckAddressAndSign(rec.LeftAddress, rec.RightAddress); e == nil {
			billsignck = "[OK]"
		}
		ens := "Submit channel transfer bill, close channels or respond to challenges, ReuseVersion: %d, BillAutoNumber: %d, ChannelId: <%s>, AssertAddress: %s, Bill signature check: %s"
		zhs := "提交通道支付票据，关闭通道或回应挑战，重用版本：%d, 账单序列号：%d，通道ID：<%s> ，主张地址：%s, 票据签名检查：%s"
		fmtEnZh(ens, zhs, rec.ReuseVersion, rec.BillAutoNumber, rec.ChannelId.ToHex(), a.AssertAddress.ToReadable(), billsignck)

	} else if a, ok := act.(*actions.Action_27_ClosePaymentChannelByClaimDistribution); ok {
		/**************** Action_27_ClosePaymentChannelByClaimDistribution ***********/
		fmtEnZh("Over arbitration lock period, the balance is allocated according to the existing claims and close channel %s",
			"仲裁锁定期限到达，按既有主张分配余额并关闭通道 <%s>",
			a.ChannelId.ToHex())

	} else if a, ok := act.(*actions.Action_30_SupportDistinguishForkChainID); ok {
		/**************** Action_30_SupportDistinguishForkChainID ***********/
		fmtEnZh("Set this transaction to only be valid for test chain or fork chain ID = %d",
			"将此交易设定为仅为测试链或分叉链 ID = %d 有效", uint64(a.CheckChainID))

	} else {

		/************************ Other ************************/
		en += fmt.Sprintf("Other action <action_kind: %d>", actId)
		zh += fmt.Sprintf("其他内容 <action_kind:%d>", actId)
	}
	/************************ END ************************/

	// 返回
	return en, zh
}
