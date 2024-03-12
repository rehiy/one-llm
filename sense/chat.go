package sense

import (
	"context"
	"errors"
	"net/http"
)

const chatCompletionsSuffix = "/chat-completions"

var (
	ErrChatCompletionInvalidModel       = errors.New("this model is not supported with this method, please use CreateCompletion client method instead") //nolint:lll
	ErrChatCompletionStreamNotSupported = errors.New("streaming is not supported with this method, please use CreateChatCompletionStream")              //nolint:lll
)

// CreateChatCompletion — API call to Create a completion for the chat message.
func (c *Client) CreateChatCompletion(
	ctx context.Context,
	request ChatCompletionRequest,
	args ...any,
) (response ChatCompletionResponse, err error) {
	if request.Stream {
		err = ErrChatCompletionStreamNotSupported
		return
	}
	if request.Model == "" && len(args) > 0 {
		m, ok := args[0].(string)
		if ok {
			request.Model = m
		}
	}

	req, err := c.newRequest(ctx, http.MethodPost, c.config.BaseURL+chatCompletionsSuffix, withBody(request))
	if err != nil {
		return
	}

	err = c.sendRequest(req, &response)
	return
}
