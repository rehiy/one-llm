package minmax

import (
	"bufio"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"

	utils "github.com/rehiy/one-llm/internal"
)

var (
	errorPrefix = []byte(`{"error`)
)

type streamable interface {
	ChatCompletionResponse
}

type streamReader struct {
	isFinished bool

	reader         *utils.EventStreamReader
	response       *http.Response
	errAccumulator utils.ErrorAccumulator
	unmarshaler    utils.Unmarshaler
}

func newStreamReader(response *http.Response, emptyMessagesLimit uint) *streamReader {
	reader := utils.NewEventStreamReader(bufio.NewReader(response.Body), 4096, emptyMessagesLimit)

	return &streamReader{
		reader:         reader,
		response:       response,
		errAccumulator: utils.NewErrorAccumulator(),
		unmarshaler:    &utils.JSONUnmarshaler{},
	}
}

func (stream *streamReader) Recv() (response ChatCompletionResponse, err error) {
	if stream.isFinished {
		err = io.EOF
		return
	}

	event, err := stream.reader.Recv()
	if err != nil {
		return
	}

	if event.Data != nil {
		err = json.Unmarshal(event.Data, &response)
		if err != nil {
			return
		}
	} else {
		err = json.Unmarshal(event.Raw, &response)
		if err != nil {
			return
		}
	}

	if response.BaseResp.StatusCode != 0 {
		err = fmt.Errorf("[%d] %s", response.BaseResp.StatusCode, response.BaseResp.StatusMsg)
		return
	}

	if len(response.Choices) == 0 {
		err = errors.New("empty content")
		return
	}

	if response.Choices[0].FinishReason == "stop" {
		stream.isFinished = true
		return
	}

	if response.Choices[0].FinishReason == "length" {
		err = errors.New("tokens_to_generate limit exceeded")
		return
	}

	if response.Choices[0].FinishReason == "max-output" {
		err = errors.New("model max tokens limit exceeded")
		return
	}

	return
}

func (stream *streamReader) Close() {
	err := stream.response.Body.Close()
	if err != nil {
		return
	}
}
