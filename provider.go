package coin

import "context"

type Group string

const (
	BalanceGroup   Group = "balance"
	FlaggerGroup   Group = "flagger"
	ThresholdGroup Group = "threshold"
)

type Broker interface {
	Execute(ctx context.Context, groupName Group) error
}
