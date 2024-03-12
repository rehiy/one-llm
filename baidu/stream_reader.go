package baidu

import (
	"bufio"
	"errors"
	"fmt"
	"io"
	"net/http"

	utils "github.com/rehiy/one-llm/internal"
)

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
	event, err := stream.reader.Recv()
	if err != nil {
		return
	}

	if stream.isFinished {
		err = io.EOF
		return
	}

	if event.Data == nil {
		unmarshalErr := stream.unmarshaler.Unmarshal(event.Other, &response)
		if unmarshalErr != nil {
			return response, unmarshalErr
		}

		if response.ErrorCode > 0 {
			return response, fmt.Errorf("[%d][%s]", response.ErrorCode, response.ErrorMsg)
		}

		err = errors.New("data is empty")

		return
	}

	unmarshalErr := stream.unmarshaler.Unmarshal(event.Data, &response)
	if unmarshalErr != nil {
		return response, unmarshalErr
	}

	if response.IsEnd {
		stream.isFinished = true
		return
	}

	return
}

func (stream *streamReader) unmarshalError() (errResp *ErrorResponse) {
	errBytes := stream.errAccumulator.Bytes()
	if len(errBytes) == 0 {
		return
	}

	err := stream.unmarshaler.Unmarshal(errBytes, &errResp)
	if err != nil {
		errResp = nil
	}

	return
}

func (stream *streamReader) Close() {
	stream.response.Body.Close()
}
