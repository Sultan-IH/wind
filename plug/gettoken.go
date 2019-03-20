package plug

import (
	"bytes"
	"encoding/json"
	"io/ioutil"
	"net/http"

	logs "log"
)

var token string

type requestBody struct {
	Method string      `json:"method"`
	Params interface{} `json:"params,omitempty"`
}
type authenticationParams struct {
	AppType  string `json:"appType"`
	Username string `json:"cloudUserName"`
	Pwd      string `json:"cloudPassword"`
	Uuid     string `json:"terminalUUID"`
}

type authenticationResponse struct {
	ErrorCode int                  `json:"error_code"`
	Result    authenticationResult `json:"result"`
}
type authenticationResult struct {
	AccountId string `json:"accountId"`
	RegTime   string `json:"regTime"`
	Email     string `json:"email"`
	Token     string `json:"token"`
}

func GetToken(username, pwd, uuid string) error {
	params := authenticationParams{
		AppType:  "Kasa_Android",
		Username: username,
		Pwd:      pwd,
		Uuid:     uuid,
	}
	body := requestBody{
		Method: "login",
		Params: params,
	}
	bodyBytes, err := json.Marshal(body)
	if err != nil {
		logs.Printf("error marshalling register request body: %v", err)
		return err
	}
	logs.Printf("body: %s", string(bodyBytes))
	reader := bytes.NewReader(bodyBytes)
	resp, err := http.Post("https://wap.tplinkcloud.com", "application/json", reader)
	if err != nil {
		logs.Printf("error making post request to register: %v", err)
		return err
	}
	logs.Printf("get token response status: %d", resp.StatusCode)
	respBytes, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		logs.Printf("error reading response body: %v", err)
		return err
	}
	result := &authenticationResponse{}
	if err := json.Unmarshal(respBytes, result); err != nil {
		logs.Printf("error unmarshalling resp bytes: %v", err)
		return err
	}
	if result.ErrorCode != 0 {
		logs.Printf("unexpecteds response from server: %v", err)
		return err
	}
	token = result.Result.Token
	return nil
}
