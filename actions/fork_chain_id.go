package actions

import (
	"github.com/hacash/core/actions"
	"github.com/hacash/core/fields"
	"github.com/hacash/core/transactions"
)

var SetCheckChainID uint64 = 0

func MaybeForTransactionAddCheckChainID(trs *transactions.Transaction_2_Simple) {
	trs.AddAction(&actions.Action_30_SupportDistinguishForkChainID{
		CheckChainID: fields.VarUint8(SetCheckChainID),
	})
}
