package sense

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	neturl "net/url"

	utils "github.com/rehiy/one-llm/internal"
)

type Client struct {
	config ClientConfig

	requestBuilder    utils.RequestBuilder
	createFormBuilder func(io.Writer) utils.FormBuilder

	authToken AuthToken
}

// NewClient creates new Sense AI API client.
func NewClient(token string) *Client {
	config := DefaultConfig(token)
	return NewClientWithConfig(config)
}

// NewClientWithConfig creates new Sense Nova API client for specified config.
func NewClientWithConfig(config ClientConfig) *Client {
	return &Client{
		config:         config,
		requestBuilder: utils.NewRequestBuilder(),
		createFormBuilder: func(body io.Writer) utils.FormBuilder {
			return utils.NewFormBuilder(body)
		},
	}
}

type requestOptions struct {
	body   any
	query  map[string]string
	header http.Header
}

type requestOption func(*requestOptions)

func withBody(body any) requestOption {
	return func(args *requestOptions) {
		args.body = body
	}
}

func withQuery(query map[string]string) requestOption {
	return func(args *requestOptions) {
		for k, v := range query {
			args.query[k] = v
		}
	}
}

func withContentType(contentType string) requestOption {
	return func(args *requestOptions) {
		args.header.Set("Content-Type", contentType)
	}
}

func (c *Client) newRequest(ctx context.Context, method, url string, setters ...requestOption) (*http.Request, error) {
	// Default Options
	args := &requestOptions{
		body:   nil,
		query:  map[string]string{},
		header: make(http.Header),
	}
	for _, setter := range setters {
		setter(args)
	}

	if args.query != nil {
		baseURL, _ := neturl.Parse(url)
		params := neturl.Values{}
		for k, v := range args.query {
			params.Add(k, v)
		}
		baseURL.RawQuery = params.Encode()
		url = baseURL.String()
	}

	req, err := c.requestBuilder.Build(ctx, method, url, args.body, args.header)
	if err != nil {
		return nil, err
	}
	c.setCommonHeaders(req)
	return req, nil
}

func (c *Client) sendRequest(req *http.Request, v any) error {
	req.Header.Set("Accept", "application/json; charset=utf-8")

	// Check whether Content-Type is already set, Upload Files API requires
	// Content-Type == multipart/form-data
	contentType := req.Header.Get("Content-Type")
	if contentType == "" {
		req.Header.Set("Content-Type", "application/json; charset=utf-8")
	}

	res, err := c.config.HTTPClient.Do(req)
	if err != nil {
		return err
	}

	defer res.Body.Close()

	if isFailureStatusCode(res) {
		return c.handleErrorResp(res)
	}

	return decodeResponse(res.Body, v)
}

func (c *Client) sendRequestRaw(req *http.Request) (body io.ReadCloser, err error) {
	resp, err := c.config.HTTPClient.Do(req)
	if err != nil {
		return
	}

	if isFailureStatusCode(resp) {
		err = c.handleErrorResp(resp)
		return
	}
	return resp.Body, nil
}

func sendRequestStream[T streamable](client *Client, req *http.Request) (*streamReader, error) {
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Accept", "text/event-stream")
	req.Header.Set("Cache-Control", "no-cache")
	req.Header.Set("Connection", "keep-alive")

	resp, err := client.config.HTTPClient.Do(req) //nolint:bodyclose // body is closed in stream.Close()
	if err != nil {
		return nil, err
	}
	if isFailureStatusCode(resp) {
		//return nil, client.handleErrorResp(resp)
	}
	return newStreamReader(resp, client.config.EmptyMessagesLimit), nil
}

func (c *Client) setCommonHeaders(req *http.Request) {
	req.Header.Set("Authorization", fmt.Sprintf("%s", c.config.authToken))
}

func isFailureStatusCode(resp *http.Response) bool {
	return resp.StatusCode < http.StatusOK || resp.StatusCode >= http.StatusBadRequest
}

func decodeResponse(body io.Reader, v any) error {
	if v == nil {
		return nil
	}

	if result, ok := v.(*string); ok {
		return decodeString(body, result)
	}
	return json.NewDecoder(body).Decode(v)
}

func decodeString(body io.Reader, output *string) error {
	b, err := io.ReadAll(body)
	if err != nil {
		return err
	}
	*output = string(b)
	return nil
}

// fullURL returns full URL for request.
func (c *Client) fullURL(model string) string {
	urlSuffix := chatCompletionsSuffix
	if model != "" {
		urlSuffix = "/chat" + model
	}

	return fmt.Sprintf("%s%s", c.config.BaseURL, urlSuffix)
}

func (c *Client) handleErrorResp(resp *http.Response) error {
	var errRes ErrorResponse
	//fmt.Printf("%s", resp.Body)

	err := json.NewDecoder(resp.Body).Decode(&errRes)
	if err != nil || errRes.Error == nil {
		reqErr := &RequestError{
			HTTPStatusCode: resp.StatusCode,
			Err:            err,
		}
		if errRes.Error != nil {
			reqErr.Err = errRes.Error
		}
		return reqErr
	}

	//errRes.Error.HTTPStatusCode = resp.StatusCode
	return errRes.Error
}
