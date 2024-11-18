package openai

import (
	"fmt"
	"strings"

	"github.com/songquanpeng/one-api/relay/channeltype"
	"github.com/songquanpeng/one-api/relay/model"
)

func ResponseText2Usage(responseText string, modeName string, promptTokens int) *model.Usage {
	usage := &model.Usage{}
	usage.PromptTokens = promptTokens
	usage.CompletionTokens = CountTokenText(responseText, modeName)
	usage.TotalTokens = usage.PromptTokens + usage.CompletionTokens
	return usage
}

func GetFullRequestURL(baseURL string, requestURL string, channelType int) string {
	fullRequestURL := fmt.Sprintf("%s%s", baseURL, requestURL)
	fmt.Println("openai GetFullRequestURL baseURL: ", baseURL, " requestURL: ", requestURL)
	fmt.Println("openai GetFullRequestURL requestURL: ", requestURL)
	fmt.Println("openai GetFullRequestURL fullRequestURL: ", fullRequestURL)
	fmt.Println("openai GetFullRequestURL channelType: ", channelType)
	if strings.HasPrefix(baseURL, "https://gateway.ai.cloudflare.com") {
		switch channelType {
		case channeltype.OpenAI:
			fullRequestURL = fmt.Sprintf("%s%s", baseURL, strings.TrimPrefix(requestURL, "/v1"))
		case channeltype.Azure:
			fullRequestURL = fmt.Sprintf("%s%s", baseURL, strings.TrimPrefix(requestURL, "/openai/deployments"))
		}
	}
	// 如果baseURL是千帆，则修改URL
	if strings.Contains(baseURL, "/ai_custom/v1/wenxinworkshop/") {
		if strings.Contains(baseURL, "AccessCode") {
			fullRequestURL = baseURL
		} else {
			accessCodePart := ""
			if strings.Contains(requestURL, "?AccessCode=") {
				accessCodePart = requestURL[strings.LastIndex(requestURL, "?AccessCode="):]
			}
			fullRequestURL = fmt.Sprintf("%s%s", baseURL, accessCodePart)
			// 打印
			fmt.Println("openai GetFullRequestURL accessCodePart: ", accessCodePart)
		}
	}

	// 如果baseURL是阿里，则修改URL（也可以和上面的百度代码合并）
	if strings.Contains(baseURL, "/compatible-mode/v1/") {
		if strings.Contains(baseURL, "AccessCode") {
			fullRequestURL = baseURL
		} else {
			accessCodePart := ""
			if strings.Contains(requestURL, "?AccessCode=") {
				accessCodePart = requestURL[strings.LastIndex(requestURL, "?AccessCode="):]
			}
			fullRequestURL = fmt.Sprintf("%s%s", baseURL, accessCodePart)
			// 打印
			fmt.Println("openai GetFullRequestURL accessCodePart: ", accessCodePart)
		}
	}

	fmt.Println("openai GetFullRequestURL final url: ", fullRequestURL)
	return fullRequestURL
}
