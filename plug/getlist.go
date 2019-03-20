package plug

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"

	"bitbucket.org/emotech/common/golang/logs"
)

type getListResponse struct {
	ErrorCode int           `json:"error_code"`
	Result    getListResult `json:"result"`
}
type getListResult struct { //assumes that all devices you have are plugs
	DeviceList []Plug `json:"deviceList"`
}

func GetDeviceList() ([]Plug, error) {
	reqBody := requestBody{
		Method: "getDeviceList",
	}
	bodyBytes, err := json.Marshal(reqBody)
	if err != nil {
		logs.Printf("error marshalling register request body: %v", err)
		return nil, err
	}
	reader := bytes.NewReader(bodyBytes)
	url := fmt.Sprintf("https://wap.tplinkcloud.com?token=%s", token)
	resp, err := http.Post(url, "application/json", reader)
	if err != nil {
		logs.Printf("error making post request to get device list: %v", err)
		return nil, err
	}
	respBytes, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		logs.Printf("get device list error: cant read body response: %v", err)
		return nil, err
	}
	result := &getListResponse{}
	if err := json.Unmarshal(respBytes, result); err != nil {
		logs.Printf("get device list error: unmarshalling resp bytes: %v", err)
		return nil, err
	}
	if result.ErrorCode != 0 {
		logs.Printf("get device list error: unexpected response from server: %v", err)
		return nil, err
	}
	return result.Result.DeviceList, nil
}
