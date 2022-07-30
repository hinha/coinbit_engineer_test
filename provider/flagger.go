package provider

import (
	"log"

	"github.com/lovoo/goka"

	coin "github.com/hinha/coinbit_test"
	"github.com/hinha/coinbit_test/pb"
)

var (
	FlagGroupTable = goka.GroupTable(goka.Group(coin.FlaggerGroup))
	flagGroup      = goka.DefineGroup(
		goka.Group(coin.FlaggerGroup),
		goka.Input(coin.FlagWalletStream, new(coin.FlagEventEncoder), flag),
		goka.Persist(new(coin.FlagValueEncoder)),
	)
)

// flag processing one or more deposit within a single rolling-up
func flag(ctx goka.Context, msg interface{}) {
	log.Println(msg)
	var s *pb.FlagValue
	if v := ctx.Value(); v == nil {
		s = new(pb.FlagValue)
	} else {
		s = v.(*pb.FlagValue)
	}

	flagEvent := msg.(*pb.FlagEvent)
	if flagEvent.FlagRemoved {
		s.Flagged = false
		s.RollingUpPeriodStart = 0
	} else {
		s.Flagged = true
		s.RollingUpPeriodStart = flagEvent.RollingUpPeriodStart
	}
	ctx.SetValue(s)
}
