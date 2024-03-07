package processor

import (
	"context"
	"github.com/mymmrac/telego"
)

type Fetcher interface {
	GetChan() <-chan telego.Update
}

type Processor interface {
	Process(ctx context.Context, u telego.Update) error
}
