package centrifugo

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"server/pkg/app/service"
)

func NewCentrifugoClient() service.CentrifugoClient {
	return &centrifugoClient{}
}

type centrifugoClient struct{}

func (c *centrifugoClient) Publish(channel string, data interface{}) error {
	url := "http://centrifugo:8000/api/publish"

	payload := map[string]interface{}{
		"channel": channel,
		"data":    data,
	}

	jsonData, err := json.Marshal(payload)
	if err != nil {
		return fmt.Errorf("error marshaling payload: %v", err)
	}

	req, err := http.NewRequest("POST", url, bytes.NewBuffer(jsonData))
	if err != nil {
		return fmt.Errorf("error creating request: %v", err)
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "_salt") // Замените на ваш API ключ

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return fmt.Errorf("error making request: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("unexpected status code: %d", resp.StatusCode)
	}

	return nil
}
