package aliyun

const ApiBaseUrl = "https://dashscope.aliyuncs.com"
const ApiVersion = "v1"

const ChatMessageRoleAssistant = "assistant"
const ChatMessageRoleSystem = "system"
const ChatMessageRoleUser = "user"

// request

type ChatCompletionRequest struct {
	Model      string                   `json:"model"`
	Input      ChatCompletionInput      `json:"input"`
	Parameters ChatCompletionParameters `json:"parameters"`
}

type ChatCompletionInput struct {
	Messages []ChatCompletionMessage `json:"messages"`
}

type ChatCompletionParameters struct {
	EnableSearch bool `json:"enable_search"`
}

type ChatCompletionMessage struct {
	Role    string `json:"role"`
	Content string `json:"content"`
}

// response

type ChatCompletionResponse struct {
	Code      string                `json:"code"`
	Message   string                `json:"message"`
	Output    ChatCompletioneOutput `json:"output"`
	Usage     ChatCompletionUsage   `json:"usage"`
	RequestID string                `json:"request_id"`
}

type ChatCompletioneOutput struct {
	Text         string `json:"text"`
	FinishReason string `json:"finish_reason"`
}

type ChatCompletionUsage struct {
	OutputTokens int `json:"output_tokens"`
	InputTokens  int `json:"input_tokens"`
}
