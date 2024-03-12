package tencent

import (
	"bufio"
	"encoding/json"
	"errors"
	"io"
	"net/http"

	utils "github.com/rehiy/one-llm/internal"
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
	reader := utils.NewEventStreamReader(bufio.NewReader(response.Body), 1024, emptyMessagesLimit)

	return &streamReader{
		reader:         reader,
		response:       response,
		errAccumulator: utils.NewErrorAccumulator(),
		unmarshaler:    &utils.JSONUnmarshaler{},
	}
}

func (stream *streamReader) Recv() (response ChatCompletionResponse, err error) {
	event, err := stream.reader.Recv()
	if err != nil {
		return
	}

	if stream.isFinished {
		err = io.EOF
		return
	}

	err = json.Unmarshal(event.Data, &response)
	if err != nil {
		return
	}

	if response.Error.Code > 0 {
		err = errors.New(response.Error.Message)
		return
	}

	if len(response.Choices) > 0 && response.Choices[0].FinishReason == "stop" {
		err = io.EOF
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
