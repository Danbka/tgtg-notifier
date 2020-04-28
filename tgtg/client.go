package tgtg

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/tidwall/gjson"
	"io/ioutil"
	"log"
	"net/http"
)

type Client struct {
	userId string
	accessToken string
	refreshToken string
}

type CheckerResult struct {
	storeName string
	availableItems int
}

func (c *Client) SetAccessToken(token string) {
	c.accessToken = token
}

func (c *Client) SetRefreshToken(token string) {
	c.refreshToken = token
}

func (c *Client) SetUserId(id string) {
	c.userId = id
}

func (c *Client) Auth(email, password string) error {
	requestBody, err := json.Marshal(map[string]string{
		"email": email,
		"password": password,
		"device_type": "UNKNOWN",
	})

	if err != nil {
		return err
	}

	request, err := http.NewRequest("POST", "https://apptoogoodtogo.com/api/auth/v1/loginByEmail", bytes.NewBuffer(requestBody))

	if err != nil {
		return err
	}

	request.Header.Set("Content-type", "application/json")
	request.Header.Set("Accept-Language", "ru-RU,ru;q=0.9,en-US;q=0.8,en;q=0.7")

	client := http.Client{}

	response, err := client.Do(request)
	if err != nil {
		log.Fatalln(err)
	}

	defer response.Body.Close()

	var result map[string]interface{}

	if json.NewDecoder(response.Body).Decode(&result) != nil {
		log.Fatalln(err)
	}

	startupData := result["startup_data"].(map[string]interface{})
	user := startupData["user"].(map[string]interface{})

	fmt.Println(user["user_id"])

	c.SetAccessToken(result["access_token"].(string))
	c.SetRefreshToken(result["refresh_token"].(string))
	c.SetUserId(user["user_id"].(string))

	return nil
}

func (c *Client) HasAvailableItem(storeId string) (bool, error) {
	requestBody, err := json.Marshal(map[string]string{
		"user_id": c.userId,
	})

	if err != nil {
		return false, err
	}

	request, err := http.NewRequest("POST", "https://apptoogoodtogo.com/api/store/v3/" + storeId, bytes.NewBuffer(requestBody))

	if err != nil {
		return false, err
	}

	request.Header.Set("Content-type", "application/json")
	request.Header.Set("Authorization", "Bearer " + c.accessToken)

	client := http.Client{}

	response, err := client.Do(request)
	if err != nil {
		return false, err
	}

	defer response.Body.Close()

	body, err := ioutil.ReadAll(response.Body)
	if err != nil {
		panic(err.Error())
	}

	bodyString := string(body)

	items := gjson.Get(bodyString, "items")

	for _, v := range items.Array() {
		if gjson.Get(v.String() , "items_available").Int() > 0 {
			return true, nil
		}
	}

	return false, nil
}