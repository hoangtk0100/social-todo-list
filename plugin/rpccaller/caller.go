package rpccaller

import (
	appctx "github.com/hoangtk0100/app-context"
	"github.com/spf13/pflag"
)

type apiItemCaller struct {
	id         string
	serviceURL string
	logger     appctx.Logger
}

func NewApiItemCaller(id string) *apiItemCaller {
	return &apiItemCaller{id: id}
}

func (c *apiItemCaller) ID() string {
	return c.id
}

func (c *apiItemCaller) InitFlags() {
	pflag.StringVar(&c.serviceURL, "item-service-url", "http://localhost:9091", "URL of item service")
}

func (c *apiItemCaller) Run(ac appctx.AppContext) error {
	c.logger = ac.Logger("api.item")
	return nil
}

func (c *apiItemCaller) Stop() error {
	return nil
}

func (c *apiItemCaller) GetServiceURL() string {
	return c.serviceURL
}
