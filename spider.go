package crawler

import (
	"golang.org/x/net/context"

	"github.com/zly-app/crawler/core"
)

var _ core.ISpider = (*Spider)(nil)

type Spider struct {
	core.ISpiderTool
}

func (s *Spider) Init(ctx context.Context) error {
	return nil
}

func (s *Spider) SubmitInitialSeed(ctx context.Context) error {
	return nil
}

func (s *Spider) Close(ctx context.Context) error {
	return nil
}
