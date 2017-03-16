package main

import (
	"encoding/json"
	"io/ioutil"
	"net/http"
	"testing"
	"time"
)

func TestDeserialization(t *testing.T) {
	api_url := "https://redisq.zkillboard.com/listen.php"
	client := http.Client{
		Timeout: time.Duration(20 * time.Second),
	}
	resp, err := client.Get(api_url)
	if err != nil {
		t.Skip("Error returned from zKill, please run test later.")
	}

	var kill zKill
	defer resp.Body.Close()
	response_bytes, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		t.Skip("Failed to read bytes, please run test later.")
	}

	err = json.Unmarshal(response_bytes, &kill)

	if kill.Payload == nil {
		t.Skip("zKill returned no kill info, please run test later.")
	}
	if kill.Payload.KillId == 0 {
		t.Error("JSON did not deserialize, KillId was 0.")
	}
	if kill.Payload.Zkb.Value <= 0 {
		t.Error("JSON did not deserialize, Zkill value was 0.")
	}
}
