package main

import (
	"bitbucket.org/emotech/common/golang/logs"
	"github.com/Sultan-IH/wind/server"
)

func main() {
	s, err := server.NewServer()

	if err != nil {
		logs.Panicf("error starting server: %v", err)
	}
	s.Run()

	// err := plug.GetToken("ksula0155@gmail.com", "Kasa123Darkside", "f60c4c04-4f72-47d7-a150-ed5b48def372")
	// if err != nil {
	// 	return
	// }
	// devices, err := plug.GetDeviceList()
	// if err != nil {
	// 	return
	// }
	// logs.Printf("devices: %+v", devices)
	// plug := devices[0]

	// plug.TurnOFF()
	// time.Sleep(time.Second * 3)
	// plug.TurnON()
}
