package crawler

import (
	"github.com/zly-app/crawler/core"
)

var _ core.ISpider = (*Spider)(nil)

type Spider struct {
	core.ISpiderTool
}

func (s *Spider) Init(tool core.ISpiderTool) error {
	s.ISpiderTool = tool
	return nil
}

func (s *Spider) SubmitInitialSeed() error {
	return nil
}

func (s *Spider) Close() error {
	return nil
}
