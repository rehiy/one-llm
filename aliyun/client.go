package aliyun

import (
	"encoding/json"

	"github.com/rehiy/one-llm/internal/httpc"
)

type Client struct {
	ApiBaseUrl string
	ApiVersion string
	ApiKey     string
	Model      string
	Parameters ChatCompletionParameters
}

func NewClient(key string) *Client {

	return &Client{
		ApiBaseUrl: ApiBaseUrl,
		ApiVersion: ApiVersion,
		ApiKey:     key,
		Model:      "qwen-max",
		Parameters: ChatCompletionParameters{EnableSearch: true},
	}

}

func (c *Client) CreateChatCompletion(messages []ChatCompletionMessage) (ChatCompletionResponse, error) {

	var resp ChatCompletionResponse

	query := ChatCompletionRequest{
		Model:      c.Model,
		Input:      ChatCompletionInput{messages},
		Parameters: c.Parameters,
	}

	heaner := httpc.H{
		"Content-Type":  "application/json",
		"Authorization": "Bearer " + c.ApiKey,
	}

	url := c.ApiBaseUrl + "/api/" + c.ApiVersion + "/services/aigc/text-generation/generation"
	response, err := httpc.JsonPost(url, query, heaner)
	if err != nil {
		return resp, err
	}

	err = json.Unmarshal(response, &resp)

	return resp, err

}
