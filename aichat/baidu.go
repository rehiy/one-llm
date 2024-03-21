package aichat

import (
	"context"
	"errors"
	"strings"

	"github.com/rehiy/one-llm/baidu"
)

func BaiduText(ask string, llmc *UserConfig) (string, error) {

	keys := strings.Split(llmc.Secret, ",")
	if len(keys) != 2 {
		return "", errors.New("密钥格式错误")
	}

	model := "completions_pro"
	if len(llmc.Model) > 1 {
		model = llmc.Model
	}

	// 初始化模型

	config := baidu.DefaultConfig(keys[0], keys[1], true)

	if len(llmc.Endpoint) > 1 {
		config.BaseURL = llmc.Endpoint
	}

	client := baidu.NewClientWithConfig(config)

	req := baidu.ChatCompletionRequest{
		Messages: []baidu.ChatCompletionMessage{},
	}

	// 设置上下文

	if llmc.RoleContext != "" {
		req.Messages = []baidu.ChatCompletionMessage{
			{Content: llmc.RoleContext, Role: baidu.ChatMessageRoleUser},
			{Content: "OK", Role: baidu.ChatMessageRoleAssistant},
		}
	}

	for _, msg := range llmc.MsgHistorys {
		role := msg.Role
		req.Messages = append(req.Messages, baidu.ChatCompletionMessage{
			Content: msg.Content, Role: role,
		})
	}

	req.Messages = append(req.Messages, baidu.ChatCompletionMessage{
		Content: ask, Role: baidu.ChatMessageRoleUser,
	})

	// 请求模型接口

	res, err := client.CreateChatCompletion(context.Background(), req, model)
	if err != nil {
		return "", err
	}

	if res.Result == "" {
		return "", errors.New("未得到预期的结果，请稍后重试")
	}

	// 更新历史记录

	item1 := &MsgHistory{Content: ask, Role: "user"}
	item2 := &MsgHistory{Content: res.Result, Role: "assistant"}

	llmc.AddHistory(item1, item2)

	return item2.Content, nil

}
