package _dataloader

import (
	"context"

	"github.com/graph-gophers/dataloader/v6"
)

func (s GeneralDataloader) JobBatchFunc(ctx context.Context, keys dataloader.Keys) []*dataloader.Result {
	results := make([]*dataloader.Result, len(keys))
	// <You Can Start Making Things More Performant in Here>
	return results
}
