package contracts

import "context"

type Deployer interface {
	DeployEsdtSafeContract(ctx context.Context, contractLocation string) error
	IsInterfaceNil() bool
}
