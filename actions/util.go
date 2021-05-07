package actions

import (
	"github.com/hacash/core/account"
	"github.com/hacash/core/fields"
	"strings"
)

func parseAccountFromAddressOrPasswordOrPrivateKey(stuff string) (*fields.Address, *account.Account) {
	stuff = strings.Trim(stuff, "\n ")
	var e error = nil
	var addr *fields.Address = nil
	var acc *account.Account = nil
	// 判断是否为地址
	if addr, e = fields.CheckReadableAddress(stuff); e == nil {
		return addr, acc
	}
	// 判断是否为私钥
	acc = account.GetAccountByPrivateKeyOrPassword(stuff)
	a := fields.Address(acc.Address)
	return &a, acc
}
