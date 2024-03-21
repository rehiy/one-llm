package aichat

import (
	"errors"

	"github.com/rehiy/one-llm/aliyun"
)

func AliyunText(ask string, llmc *UserConfig) (string, error) {

	// 初始化模型

	client := aliyun.NewClient(llmc.Secret)

	if len(llmc.Model) > 1 {
		client.Model = llmc.Model
	}

	if len(llmc.Endpoint) > 1 {
		client.ApiBaseUrl = llmc.Endpoint
	}

	req := []aliyun.ChatCompletionMessage{}

	// 设置上下文

	if llmc.RoleContext != "" {
		req = []aliyun.ChatCompletionMessage{
			{Content: llmc.RoleContext, Role: aliyun.ChatMessageRoleSystem},
		}
	}

	for _, msg := range llmc.MsgHistorys {
		role := msg.Role
		req = append(req, aliyun.ChatCompletionMessage{
			Content: msg.Content, Role: role,
		})
	}

	req = append(req, aliyun.ChatCompletionMessage{
		Content: ask, Role: aliyun.ChatMessageRoleUser,
	})

	// 请求模型接口

	resp, err := client.CreateChatCompletion(req)
	if err != nil {
		return "", err
	}

	if resp.Message != "" {
		return "", errors.New(resp.Message)
	}

	if resp.Output.Text == "" {
		return "", errors.New(resp.Output.FinishReason)
	}

	// 更新历史记录

	item1 := &MsgHistory{Content: ask, Role: "user"}
	item2 := &MsgHistory{Content: resp.Output.Text, Role: "assistant"}

	llmc.AddHistory(item1, item2)

	return item2.Content, nil

}
