package crawler

import (
	"github.com/zly-app/crawler/core"
)

var _ core.ISpider = (*Spider)(nil)

type Spider struct {
	core.ISpiderTool
}

func (s *Spider) Init() error {
	return nil
}

func (s *Spider) SubmitInitialSeed() error {
	return nil
}

func (s *Spider) Close() error {
	return nil
}
