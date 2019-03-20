package server

import (
	logs "log"

	"github.com/alexflint/go-arg"
)

type Config struct {
	Port           int    `arg:"env:PORT"`
	TPLinkUsername string `arg:"env:TPLINK_USERNAME"`
	TPLinkPwd      string `arg:"env:TPLINK_PWD"`
	TPLinkUUID     string `arg:"env:TPLINK_UUID"`
}

func (c *Config) Parse() error {
	c.Port = 9009
	if err := arg.Parse(c); err != nil {
		return err
	}
	logs.Printf("loaded the following config: %+v", c)
	return nil
}
