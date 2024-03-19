# 大模型统一接口

各大厂的大语言模型非官方 SDK。主要特色是低依赖，支持以下模型：

- `aliyun` 阿里通义千问
- `baidu` 百度千帆（文心）
- `google` 谷歌 Gemini
- `minmax` MinMax
- `sense` 商汤日日新
- `tencent` 腾讯混元
- `xunfei` 讯飞星火
- `zhipu` 智谱 Ai

## 密钥申请

- 阿里 通义千问 <https://dashscope.console.aliyun.com/apiKey>
- 谷歌 Gemini <https://aistudio.google.com/app/apikey?hl=zh-cn>
- 讯飞 星火 <https://console.xfyun.cn/services/bm3>

## 使用方法

以 `百度千帆 Ernie` 为例，每个模型的引入和初始化方式都差不多，可以根据类型提示完成代码。

```go
package main

import (
    "context"
    "fmt"
    . "github.com/rehiy/one-llm/baidu"
)

func main() {
    client := NewClient("xxxx", "yyyy", true)
    stream, err := client.CreateChatCompletionStream(context.Background(), ChatCompletionRequest{
        Messages: []ChatCompletionMessage{
            {Content: "Hello!", Role: ChatMessageRoleUser},
        },
        Stream: true,
    })

    defer stream.Close()

    fmt.Println("Stream response: ")
    for {
        response, err := stream.Recv()
        if errors.Is(err, io.EOF) {
            fmt.Printf("\nStream finished: %d %s\n", response.ErrorCode, response.ErrorMsg)
            return
        }
        if err != nil {
            fmt.Printf("\nStream error: %v\n", err)
            return
        }
        fmt.Printf("error: %s\n", response.ErrorMsg)
        fmt.Printf("resp: %s\n", response.Result)
    }
}
```

## 其他说明

部分源码来自 [go-llm-api](https://github.com/liudding/go-llm-api)，感谢原作者的贡献。吃水不忘挖*坑*人，方便的话，请顺手赠送一个 Star。
