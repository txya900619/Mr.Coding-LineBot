package mrcoding

import (
	"bytes"
	"encoding/json"
	"io/ioutil"
	"log"
	"net/http"
)

//Send request to backend to create chatroom and get chatroom id
func (bot *Bot) createChatroomAndGetID(userID, name string) string {
	client := &http.Client{}
	reqBody := map[string]string{"lineChatroomUserID": userID, "name": name}
	jsonReqBody, _ := json.Marshal(reqBody)
	req, err := http.NewRequest(http.MethodPost, "https://mrcoding.org/api/chatrooms", bytes.NewBuffer(jsonReqBody))
	if err != nil {
		log.Fatalf("newReq, err: %v", err)
	}

	req.Header.Set("Authorization", bot.backendToken)
	req.Header.Set("Content-Type", "application/json")
	resp, err := client.Do(req)
	if err != nil {
		log.Fatalf("reqDo, err: %v", err)
	}

	//Backend err
	if !(resp.StatusCode == 200 || resp.StatusCode == 201) {
		log.Fatal("statusCode err")

	}

	body, err := ioutil.ReadAll(resp.Body)
	defer resp.Body.Close()
	if err != nil {
		log.Fatalf("readAll, err: %v", err)
	}

	result := make(map[string]interface{})
	err = json.Unmarshal([]byte(body), &result)
	if err != nil {
		log.Fatalf("jsonUnmarshal, err: %v", err)
	}
	return result["_id"].(string)
}
