package transact

import (
	"context"
	"eticket-api/pkg/gotann"
)

type Transactor interface {
	Execute(ctx context.Context, fn func(tx gotann.Transaction) error) error
	ExecuteWithRetry(ctx context.Context, fn func(tx gotann.Transaction) error) error
	ExecuteReadOnly(ctx context.Context, fn func(tx gotann.Transaction) error) error
}
