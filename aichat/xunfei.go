package aichat

import (
	"context"
	"errors"
	"io"
	"strings"

	"github.com/rehiy/one-llm/xunfei"
)

func XunfeiText(ask string, llmc *UserConfig) (string, error) {

	keys := strings.Split(llmc.Secret, ",")
	if len(keys) != 3 {
		return "", errors.New("密钥格式错误")
	}

	model := "v3.5"
	if len(llmc.Model) > 1 {
		model = llmc.Model
	}

	// 初始化模型

	config := xunfei.DefaultConfig(keys[0], keys[1], keys[2])

	if len(llmc.Endpoint) > 1 {
		config.BaseURL = llmc.Endpoint
	}

	client := xunfei.NewClientWithConfig(config)

	req := xunfei.ChatCompletionRequest{
		MaxTokens: 2048,
		Messages:  []xunfei.ChatCompletionMessage{},
	}

	// 设置上下文

	if llmc.RoleContext != "" {
		req.Messages = []xunfei.ChatCompletionMessage{
			{Content: llmc.RoleContext, Role: xunfei.ChatMessageRoleUser},
			{Content: "OK", Role: xunfei.ChatMessageRoleAssistant},
		}
	}

	for _, msg := range llmc.MsgHistorys {
		role := msg.Role
		req.Messages = append(req.Messages, xunfei.ChatCompletionMessage{
			Content: msg.Content, Role: role,
		})
	}

	req.Messages = append(req.Messages, xunfei.ChatCompletionMessage{
		Content: ask, Role: xunfei.ChatMessageRoleUser,
	})

	// 请求模型接口

	stream, err := client.CreateChatCompletionStream(context.Background(), req, model)
	if err != nil {
		return "", err
	}

	defer stream.Close()

	reply := ""

	for {
		response, err := stream.Recv()
		if err != nil {
			if errors.Is(err, io.EOF) {
				break
			}
			return reply, err
		}
		if len(response.Payload.Choices.Text) > 0 {
			reply += response.Payload.Choices.Text[0].Content
		} else {
			reply += response.Payload.Choices.Content
		}
	}

	if reply == "" {
		return "", errors.New("未得到预期的结果，请稍后重试")
	}

	// 更新历史记录

	item1 := &MsgHistory{Content: ask, Role: "user"}
	item2 := &MsgHistory{Content: reply, Role: "assistant"}

	llmc.AddHistory(item1, item2)

	return item2.Content, nil

}
