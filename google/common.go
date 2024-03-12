package google

const ApiBaseUrl = "https://generativelanguage.googleapis.com"
const ApiVersion = "v1beta"

const ChatMessageRoleAssistant = "model"
const ChatMessageRoleSystem = "user"
const ChatMessageRoleUser = "user"

// chat request

type ChatCompletionRequest struct {
	Contents       []ChatCompletionMessage       `json:"contents"`
	SafetySettings []ChatCompletionSafetySetting `json:"safetySettings"`
}

type ChatCompletionMessage struct {
	Parts []ChatCompletionMessagePart `json:"parts"`
	Role  string                      `json:"role"`
}

type ChatCompletionMessagePart struct {
	Text       string                    `json:"text,omitempty"`
	InlineData *ChatCompletionInlineData `json:"inline_data,omitempty"`
}

type ChatCompletionInlineData struct {
	MimeType string `json:"mime_type"`
	Data     string `json:"data"`
}

type ChatCompletionSafetySetting struct {
	Category  string `json:"category"`
	Threshold string `json:"threshold"`
}

// chat response

type ChatCompletionResponse struct {
	Candidates     []ChatCompletionCandidate `json:"candidates"`
	Error          ChatCompletionError       `json:"error"`
	PromptFeedback struct {
		BlockReason   string                       `json:"blockReason"`
		SafetyRatings []ChatCompletionSafetyRating `json:"safetyRatings"`
	} `json:"promptFeedback"`
}

type ChatCompletionCandidate struct {
	Content       ChatCompletionMessage        `json:"content"`
	FinishReason  string                       `json:"finishReason"`
	Index         int                          `json:"index"`
	SafetyRatings []ChatCompletionSafetyRating `json:"safetyRatings"`
}

type ChatCompletionSafetyRating struct {
	Category    string `json:"category"`
	Probability string `json:"probability"`
}

type ChatCompletionError struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
	Status  string `json:"status"`
}
