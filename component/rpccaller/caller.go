package rpccaller

import (
	appctx "github.com/hoangtk0100/app-context"
	"github.com/spf13/pflag"
)

type itemAPICaller struct {
	id         string
	serviceURL string
	logger     appctx.Logger
}

func NewItemAPICaller(id string) *itemAPICaller {
	return &itemAPICaller{id: id}
}

func (c *itemAPICaller) ID() string {
	return c.id
}

func (c *itemAPICaller) InitFlags() {
	pflag.StringVar(&c.serviceURL, "item-service-url", "http://localhost:9091", "URL of item service")
}

func (c *itemAPICaller) Run(ac appctx.AppContext) error {
	c.logger = ac.Logger("api.item")
	return nil
}

func (c *itemAPICaller) Stop() error {
	return nil
}

func (c *itemAPICaller) GetServiceURL() string {
	return c.serviceURL
}
