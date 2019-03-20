package datahandler

import (
	"bitbucket.org/emotech/common/golang/logs"
	"github.com/Sultan-IH/wind/plug"
)

type DataChannel chan []byte

type DataHandler struct {
	directory map[string]DataChannel
	plugs     map[string]plug.Plug
}

func (dh *DataHandler) StartTransmission(ID string) bool {

	_, ok := dh.directory[ID]
	if ok {
		logs.Printf("transmission already exists for channel with ID [%v], ignoring ...", ID)
		return false
	}
	// make a go routine that d
	plug, ok := dh.plugs[ID]
	if !ok {
		logs.Printf("start transmission: plug with ID %s not found", ID)
		return false
	}

	dh.directory[ID] = make(DataChannel)
	logs.Printf("started transmission channel with ID [%s]", ID)

	go controlVentilatorPWM(dh.directory[ID], plug)
	return true
}
func (dh *DataHandler) EndTransmission(ID string) {
	channel, ok := dh.directory[ID]
	if !ok {
		return
	}
	close(channel)
	delete(dh.directory, ID)
	logs.Printf("closed transmission channel with ID [%s]", ID)
}

func (dh *DataHandler) RecordData(ID string, data []byte) bool {
	channel, ok := dh.directory[ID]
	if !ok {
		logs.Printf("no transmission registered at ", ID)
		return false
	}
	// so send operation on a channel is actually a blocking operation
	// which is why we use select
	select {
	case channel <- data:
		logs.Printf("sent %v on channel [%s]", string(data), ID)
	default:
		logs.Printf("no receivers detected, value ignored")
	}
	return true
}

func (dh *DataHandler) ReceiveData(ID string) DataChannel {
	channel, ok := dh.directory[ID]
	if !ok {
		logs.Printf("cant find channels with this ID: %v", ID)
		return nil
	}

	return channel
}

func NewDataHandler(plugs []plug.Plug) *DataHandler {
	dh := &DataHandler{
		directory: make(map[string]DataChannel),
		plugs:     make(map[string]plug.Plug),
	}
	for _, plug := range plugs {
		dh.plugs[plug.Alias] = plug
	}
	return dh
}
