package provider

import (
	"github.com/lovoo/goka"

	coin "github.com/hinha/coinbit_test"
	"github.com/hinha/coinbit_test/pb"
)

var (
	BalanceGroupTable = goka.GroupTable(goka.Group(coin.BalanceGroup))
	balanceGroup      = goka.DefineGroup(goka.Group(coin.BalanceGroup),
		goka.Input(coin.DepositStream, new(coin.DepositEncoder), collect),
		goka.Persist(new(coin.DepositListEncoder)))
)

// collect call every message from DepositStream
func collect(ctx goka.Context, message interface{}) {
	ml := &pb.DepositHistory{}
	if v := ctx.Value(); v != nil {
		ml = v.(*pb.DepositHistory)
	}

	msg := message.(*pb.Deposit)

	ml.WalletId = msg.WalletId
	ml.Deposits = append(ml.Deposits, msg)

	ctx.SetValue(ml)
}
