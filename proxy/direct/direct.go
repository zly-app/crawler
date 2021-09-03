package direct

import (
	"net/http"

	zapp_core "github.com/zly-app/zapp/core"

	"github.com/zly-app/crawler/core"
)

type DirectProxy struct{}

func (d *DirectProxy) EnableProxy(transport *http.Transport)  {}
func (d *DirectProxy) DisableProxy(transport *http.Transport) {}
func (d *DirectProxy) ChangeProxy(transport *http.Transport)  {}
func (d *DirectProxy) Close() error                           { return nil }

func NewDirectProxy(app zapp_core.IApp) core.IProxy {
	return new(DirectProxy)
}
