package auth

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
)

type DeviceCodeResponse struct {
	DeviceCode      string `json:"device_code"`
	UserCode        string `json:"user_code"`
	VerificationUri string `json:"verification_uri"`
	ExpiresIn       int    `json:"expires_in"`
	Interval        int    `json:"interval"`
}

func Login(w http.ResponseWriter, r *http.Request) error {
	clientId := os.Getenv("GH_CLIENT_ID")

	req, err := http.NewRequest(
		"POST",
		fmt.Sprintf("github.com/login/device/code?client_id=%s&scope=repo,read:user", clientId),
		nil,
	)

	if err != nil {
		return err
	}

	res, err := http.DefaultClient.Do(req)

	if err != nil {
		return err
	}

	bytes, err := io.ReadAll(res.Body)

	if err != nil {
		return err
	}

	deviceCodeResponse := DeviceCodeResponse{}

	if err := json.Unmarshal(bytes, &deviceCodeResponse); err != nil {
		return err
	}

	w.WriteHeader(200)
	w.Write([]byte(fmt.Sprintf("Enter the code %s at %s to finish authenticating. Once authenticated, you'll be able to make requests.", deviceCodeResponse.DeviceCode, deviceCodeResponse.VerificationUri)))

	Callback(w, r, deviceCodeResponse.DeviceCode, deviceCodeResponse.Interval)

	return nil
}
