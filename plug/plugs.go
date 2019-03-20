package plug

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"

	"bitbucket.org/emotech/common/golang/logs"
)

type Plug struct {
	// this id is not the same one as used for transmission, this is a device ID
	ID string `json:"deviceId"`
	// this is the transmission ID, this is set in the kasa app
	Alias           string  `json:"alias"`
	VentilatorState float32 `json:"-"`
}

type controlPlugParams struct {
	DeviceID    string `json:"deviceId"`
	RequestData string `json:"requestData"`
}

type controlPlugResponse struct {
	ErrorCode int `json:"error_code"`
}

func (p *Plug) setState(state int) error {
	logs.Printf("[ PLUG:%s ] set state to %d", p.Alias, state)

	params := controlPlugParams{
		DeviceID:    p.ID,
		RequestData: fmt.Sprintf("{\"system\":{\"set_relay_state\":{\"state\":%d}}}", state),
	}
	reqBody := requestBody{
		Method: "passthrough",
		Params: params,
	}
	bodyBytes, err := json.Marshal(reqBody)
	if err != nil {
		logs.Printf("error marshalling register request body: %v", err)
		return err
	}
	reader := bytes.NewReader(bodyBytes)
	url := fmt.Sprintf("https://eu-wap.tplinkcloud.com/?token=%s", token)
	resp, err := http.Post(url, "application/json", reader)
	respBytes, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		logs.Printf("error reading response body: %v", err)
		return err
	}
	result := &controlPlugResponse{}
	if err := json.Unmarshal(respBytes, result); err != nil {
		logs.Printf("error unmarshalling resp bytes: %v", err)
		return err
	}
	if result.ErrorCode != 0 {
		logs.Printf("unexpecteds response from server: %v", err)
		return err
	}
	return nil
}

func (p *Plug) TurnON() error {
	return p.setState(1)

}
func (p *Plug) TurnOFF() error {
	return p.setState(0)
}
