package provider

import (
	"context"
	"log"

	"github.com/lovoo/goka"

	coin "github.com/hinha/coinbit_test"
)

type broker struct {
	listeners []string
}

func NewBroker(brokers []string) coin.Broker {
	return &broker{
		listeners: brokers,
	}
}

func (c *broker) Execute(ctx context.Context, groupName coin.Group) error {
	var group *goka.GroupGraph

	// state manager
	if groupName == coin.BalanceGroup {
		log.Println("starting balance group")
		group = balanceGroup
	} else if groupName == coin.FlaggerGroup {
		log.Println("starting flag group")
		group = flagGroup
	} else if groupName == coin.ThresholdGroup {
		log.Println("starting detector group")
		group = thresholdGroup
	} else {
		panic("group is not define")
	}

	return runProcessor(ctx, c.listeners, group)
}

func runProcessor(ctx context.Context, brokers []string, group *goka.GroupGraph) error {
	processor, err := goka.NewProcessor(brokers, group)
	if err != nil {
		return err
	}
	return processor.Run(ctx)
}
