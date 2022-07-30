package coin

import (
	"github.com/lovoo/goka"
	"google.golang.org/protobuf/proto"

	"github.com/hinha/coinbit_test/pb"
)

var (
	// reference: https://github.com/lovoo/goka/blob/master/examples/3-messaging/collector/collector.go
	// encoder protobuf
	DepositStream    goka.Stream = "deposit"
	FlagWalletStream goka.Stream = "flag_wallet"
)

type DepositEncoder struct{}

func (c *DepositEncoder) Encode(value interface{}) ([]byte, error) {
	return proto.Marshal(value.(*pb.Deposit))
}
func (c *DepositEncoder) Decode(data []byte) (interface{}, error) {
	var m pb.Deposit
	return &m, proto.Unmarshal(data, &m)
}

type DepositListEncoder struct{}

func (c *DepositListEncoder) Encode(value interface{}) ([]byte, error) {
	return proto.Marshal(value.(*pb.DepositHistory))
}
func (c *DepositListEncoder) Decode(data []byte) (interface{}, error) {
	var m pb.DepositHistory
	return &m, proto.Unmarshal(data, &m)
}

type FlagEventEncoder struct{}

func (c *FlagEventEncoder) Encode(value interface{}) ([]byte, error) {
	return proto.Marshal(value.(*pb.FlagEvent))
}
func (c *FlagEventEncoder) Decode(data []byte) (interface{}, error) {
	var m pb.FlagEvent
	return &m, proto.Unmarshal(data, &m)
}

type FlagValueEncoder struct{}

func (c *FlagValueEncoder) Encode(value interface{}) ([]byte, error) {
	return proto.Marshal(value.(*pb.FlagValue))
}
func (c *FlagValueEncoder) Decode(data []byte) (interface{}, error) {
	var m pb.FlagValue
	return &m, proto.Unmarshal(data, &m)
}

type CounterEncoder struct{}

func (c *CounterEncoder) Encode(value interface{}) ([]byte, error) {
	return proto.Marshal(value.(*pb.Counter))
}

func (c *CounterEncoder) Decode(data []byte) (interface{}, error) {
	var m pb.Counter
	return &m, proto.Unmarshal(data, &m)
}
