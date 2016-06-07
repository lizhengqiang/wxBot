package bot

import (
	"bytes"
	"encoding/json"
	"io/ioutil"
	"strings"
	"net/url"
	"github.com/google/go-querystring/query"
)

func (bot *WeixinBot) SimplePostJson(uri string, params interface{}) (b []byte, err error) {
	paramsBytes, paramsErr := json.Marshal(params)
	if paramsErr != nil {
		return nil, paramsErr
	}
	resp, err := bot.httpClient.Post(bot.getProperty(baseUri) + uri, "application/json", bytes.NewReader(paramsBytes))
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}
	return body, nil
}

func (bot *WeixinBot) PostJson(uri string, request interface{}, response interface{}) {
	paramsBytes, paramsErr := json.Marshal(request)
	if paramsErr != nil {
		return
	}
	var targetUrl = ""
	if strings.Contains(uri, `http://`) || strings.Contains(uri, `https://`) {
		targetUrl = uri
	} else {
		targetUrl = bot.getProperty(baseUri) + uri
	}
	resp, err := bot.httpClient.Post(targetUrl, "application/json", bytes.NewReader(paramsBytes))
	if err != nil {
		return
	}
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return
	}
	json.Unmarshal(body, response)
	return
}

func (bot *WeixinBot) GetJson(uri string, request interface{}, response interface{}) {

	var targetUrl = ""
	if strings.Contains(uri, `http://`) || strings.Contains(uri, `https://`) {
		targetUrl = uri
	} else {
		targetUrl = bot.getProperty(baseUri) + uri
	}
	// ..
	var params url.Values
	if request != nil {
		var paramsErr error
		params, paramsErr = query.Values(request)
		if paramsErr != nil {
			return
		}
		targetUrl = targetUrl + "?" + params.Encode()
	}
	resp, err := bot.httpClient.Get(targetUrl)
	if err != nil {
		return
	}
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return
	}
	json.Unmarshal(body, response)
	return
}
