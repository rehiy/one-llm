package google

import (
	"encoding/json"

	"github.com/rehiy/one-llm/internal/httpc"
)

type Client struct {
	ApiBaseUrl     string
	ApiVersion     string
	ApiKey         string
	Model          string
	SafetySettings []ChatCompletionSafetySetting
}

func NewClient(key string) *Client {

	return &Client{
		ApiBaseUrl: ApiBaseUrl,
		ApiVersion: ApiVersion,
		ApiKey:     key,
		Model:      "gemini-pro",
	}

}

func (c *Client) CreateChatCompletion(contents []ChatCompletionMessage) (ChatCompletionResponse, error) {

	var resp ChatCompletionResponse

	query := ChatCompletionRequest{
		Contents:       contents,
		SafetySettings: c.SafetySettings,
	}

	heaner := httpc.H{
		"Content-Type":   "application/json",
		"x-goog-api-key": c.ApiKey,
	}

	url := c.ApiBaseUrl + "/" + c.ApiVersion + "/models/" + c.Model + ":generateContent"
	response, err := httpc.JsonPost(url, query, heaner)
	if err != nil {
		return resp, err
	}

	err = json.Unmarshal(response, &resp)

	return resp, err

}

func (c *Client) CreateVisionCompletion(contents []ChatCompletionMessage) (ChatCompletionResponse, error) {

	var resp ChatCompletionResponse

	query := &ChatCompletionRequest{
		Contents:       contents,
		SafetySettings: c.SafetySettings,
	}

	heaner := httpc.H{
		"Content-Type":   "application/json",
		"x-goog-api-key": c.ApiKey,
	}

	url := c.ApiBaseUrl + "/" + c.ApiVersion + "/models/gemini-pro-vision:generateContent"
	response, err := httpc.JsonPost(url, query, heaner)
	if err != nil {
		return resp, err
	}

	err = json.Unmarshal(response, &resp)

	return resp, err

}
