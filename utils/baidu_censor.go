package utils

import (
	"bytes"
	"encoding/json"
	"io/ioutil"
	"net/http"
	"net/url"
	"sync"
	"time"
)

var (
	baiduToken     string
	baiduTokenLock sync.Mutex
	baiduTokenTime time.Time
)

func getBaiduAccessToken(apiKey, secretKey string) (string, error) {
	baiduTokenLock.Lock()
	defer baiduTokenLock.Unlock()
	if baiduToken != "" && time.Since(baiduTokenTime) < time.Hour {
		return baiduToken, nil
	}
	resp, err := http.PostForm("https://aip.baidubce.com/oauth/2.0/token", url.Values{
		"grant_type":    {"client_credentials"},
		"client_id":     {apiKey},
		"client_secret": {secretKey},
	})
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()
	var result struct {
		AccessToken string `json:"access_token"`
	}
	json.NewDecoder(resp.Body).Decode(&result)
	baiduToken = result.AccessToken
	baiduTokenTime = time.Now()
	return baiduToken, nil
}

func BaiduTextCensor(text, apiKey, secretKey string) (conclusion string, err error) {
	token, err := getBaiduAccessToken(apiKey, secretKey)
	if err != nil {
		return "", err
	}
	api := "https://aip.baidubce.com/rest/2.0/solution/v1/text_censor/v2/user_defined?access_token=" + token
	data := url.Values{"text": {text}}
	req, _ := http.NewRequest("POST", api, bytes.NewBufferString(data.Encode()))
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()
	body, _ := ioutil.ReadAll(resp.Body)
	var result struct {
		Conclusion string `json:"conclusion"`
	}
	json.Unmarshal(body, &result)
	return result.Conclusion, nil
}
