package provider

import (
	"log"
	"time"

	"github.com/lovoo/goka"

	coin "github.com/hinha/coinbit_test"
	"github.com/hinha/coinbit_test/pb"
)

var (
	ThresholdGroupTable = goka.GroupTable(goka.Group(coin.ThresholdGroup))
	thresholdGroup      = goka.DefineGroup(goka.Group(coin.ThresholdGroup),
		goka.Input(coin.DepositStream, new(coin.DepositEncoder), detector),
		goka.Output(coin.FlagWalletStream, new(coin.FlagEventEncoder)),
		goka.Persist(new(coin.CounterEncoder)),
	)

	rollingPeriod int64 = 120
	maxAmount           = 10000
)

func detector(ctx goka.Context, msg interface{}) {
	var counter *pb.Counter
	if val := ctx.Value(); val != nil {
		counter = val.(*pb.Counter)
	} else {
		counter = new(pb.Counter)
	}

	m := msg.(*pb.Deposit)
	counter.Received += m.Amount
	log.Println(m, counter)

	if counter.RollingUpPeriodStart == 0 {
		counter.RollingUpPeriodStart = time.Now().Unix()
	} else {
		if time.Now().Unix()-counter.RollingUpPeriodStart > rollingPeriod {
			counter.RollingUpPeriodStart = 0
			counter.Received = 0
		}
	}
	ctx.SetValue(counter)

	if counter.Received >= float64(maxAmount) && counter.RollingUpPeriodStart != rollingPeriod {
		ctx.Emit(coin.FlagWalletStream, ctx.Key(), &pb.FlagEvent{FlagRemoved: false, RollingUpPeriodStart: counter.RollingUpPeriodStart})
	} else {
		ctx.Emit(coin.FlagWalletStream, ctx.Key(), &pb.FlagEvent{FlagRemoved: true})
	}
}
