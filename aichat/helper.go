package aichat

import (
	"encoding/base64"
	"net/http"
	"os"
)

func Text(msg string, llmc *UserConfig) string {

	var err error
	var res string

	// 调用接口生成文本
	switch llmc.Provider {
	case "aliyun":
		res, err = AliyunText(msg, llmc)
	case "baidu":
		res, err = BaiduText(msg, llmc)
	case "google":
		res, err = GoogleText(msg, llmc)
	case "openai":
		res, err = OpenaiText(msg, llmc)
	case "tencent":
		res, err = TencentText(msg, llmc)
	case "xunfei":
		res, err = XunfeiText(msg, llmc)
	case "":
		res = "当前模型已失效，请重新选择"
	default:
		res = "暂不支持此模型"
	}

	// 返回结果
	if err != nil {
		return err.Error()
	}
	return res

}

func Vison(msg, img string, llmc *UserConfig) string {

	var err error
	var res string

	// 调用接口生成文本
	switch llmc.Provider {
	case "google":
		res, err = GoogleVison(msg, img, llmc)
	case "":
		res = "当前模型已失效，请重新选择"
	default:
		res = "当前模型不支持分析图片"
	}

	// 返回结果
	if err != nil {
		return err.Error()
	}
	return res

}

// 用户配置

type MsgHistory struct {
	Role    string `json:"role"`    // user,assistant
	Content string `json:"content"` // 消息内容
}

type UserConfig struct {
	Family   string `json:"family"`   // 族类描述
	Provider string `json:"provider"` // 供应商
	Endpoint string `json:"endpoint"` // 接口地址
	Model    string `json:"model"`    // 模型
	// 密钥格式
	// 科大讯飞 APP-ID,API-KEY,API-SECRET
	// 文心一言 API-KEY,API-SECRET
	// 腾讯混元 APP-ID,API-KEY,API-SECRET
	// 阿里百炼（通义千问） APP-ID,AGENT-KEY,ACCESS_KEY_ID,ACCESS_KEY_SECRET
	// 其他服务商 API-KEY
	Secret        string        `json:"secret"`          // 密钥
	RoleContext   string        `json:"role_context"`    // 角色设定
	MsgHistorys   []*MsgHistory `json:"msg_historys"`    // 消息历史记录
	MsgHistoryMax int           `json:"msg_history_max"` // 消息记录最大条数
}

func (u *UserConfig) AddHistory(items ...*MsgHistory) {

	if len(u.MsgHistorys) >= u.MsgHistoryMax {
		u.MsgHistorys = u.MsgHistorys[len(items):]
	}

	u.MsgHistorys = append(u.MsgHistorys, items...)

}

func (u *UserConfig) ResetHistory() {

	u.MsgHistorys = []*MsgHistory{}

}

// 读取图片

func ReadImage(img string) (string, string) {

	fileContent, err := os.ReadFile(img)
	if err != nil {
		return "", ""
	}

	base64String := base64.StdEncoding.EncodeToString(fileContent)
	mimeType := http.DetectContentType(fileContent)

	return base64String, mimeType

}
