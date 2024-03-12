package xunfei

const (
	ChatMessageRoleUser      = "user"
	ChatMessageRoleAssistant = "assistant"

	ModelV1   = "v1"
	ModelV2   = "v2"
	ModelV3   = "v3"
	ModelV3_5 = "v3.5"

	parameterChatDomainGeneral     = "general"
	parameterChatDomainGeneralV2   = "generalv2"
	parameterChatDomainGeneralV3   = "generalv3"
	parameterChatDomainGeneralV3_5 = "generalv3.5"
)

type ChatCompletionMessage struct {
	Role    string `json:"role"`
	Content string `json:"content"`
}

type ChatCompletionStreamRequest struct {
	Header    ChatCompletionStreamRequestHeader    `json:"header"`
	Parameter ChatCompletionStreamRequestParameter `json:"parameter"`
	Payload   ChatCompletionStreamRequestPayload   `json:"payload"`
}

type ChatCompletionStreamRequestHeader struct {
	AppId string `json:"app_id"`
	Uid   string `json:"uid,omitempty"`
}

type ChatCompletionStreamRequestParameter struct {
	Chat ChatCompletionStreamRequestParameterChat `json:"chat"`
}

type ChatCompletionStreamRequestPayload struct {
	Message ChatCompletionStreamRequestPayloadMessage `json:"message"`
}

type ChatCompletionStreamRequestPayloadMessage struct {
	Text []ChatCompletionMessage `json:"text"`
}

type ChatCompletionStreamRequestParameterChat struct {
	Domain      string  `json:"domain"`
	Temperature float64 `json:"temperature,omitempty"`
	TopK        int     `json:"top_k,omitempty"`
	MaxTokens   int     `json:"max_tokens,omitempty"`
	Auditing    string  `json:"auditing,omitempty"`
}

type ChatCompletionStreamResponse struct {
	Header  ChatCompletionStreamResponseHeader  `json:"header"`
	Payload ChatCompletionStreamResponsePayload `json:"payload"`
}

type ChatCompletionStreamResponseHeader struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
	Sid     string `json:"sid"`
	Status  int    `json:"status"` // 会话状态，取值为[0,1,2]；0代表首次结果；1代表中间结果；2代表最后一个结果
}

type ChatCompletionStreamResponsePayload struct {
	Choices ChatCompletionStreamResponsePayloadChoices `json:"choices"`
	Usage   Usage                                      `json:"usage"`
}

type ChatCompletionStreamResponsePayloadChoices struct {
	Status int `json:"status"` // 文本响应状态，取值为[0,1,2]; 0代表首个文本结果；1代表中间文本结果；2代表最后一个文本结果
	Seq    int `json:"seq"`    // 返回的数据序号，取值为[0,9999999]

	ChatCompletionMessage
	Text []ChatCompletionMessage `json:"text"` // 讯飞文档中有两种格式，这里兼容处理
}

// ChatCompletionRequest represents a request structure for chat completion API.
type ChatCompletionRequest struct {
	Messages    []ChatCompletionMessage `json:"messages"`
	Temperature float64                 `json:"temperature,omitempty"`
	TopK        int                     `json:"top_k,omitempty"`
	MaxTokens   int                     `json:"max_tokens,omitempty"`
	Auditing    string                  `json:"auditing,omitempty"`
	Uid         string                  `json:"uid,omitempty"`
}

type Usage struct {
	PromptTokens     int `json:"prompt_tokens"`
	CompletionTokens int `json:"completion_tokens"`
	TotalTokens      int `json:"total_tokens"`
}
