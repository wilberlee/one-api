package baidu

import (
	"errors"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/songquanpeng/one-api/relay/adaptor"
	"github.com/songquanpeng/one-api/relay/adaptor/openai"
	"github.com/songquanpeng/one-api/relay/meta"
	"github.com/songquanpeng/one-api/relay/model"
	"github.com/songquanpeng/one-api/relay/relaymode"
	"io"
	"net/http"
	"strings"
)

type Adaptor struct {
}

func (a *Adaptor) Init(meta *meta.Meta) {

}

func (a *Adaptor) GetRequestURL(meta *meta.Meta) (string, error) {
	// https://cloud.baidu.com/doc/WENXINWORKSHOP/s/clntwmv7t
	if strings.Contains(meta.BaseURL, "/ai_custom/v1/wenxinworkshop/") {
		return openai.GetFullRequestURL(meta.BaseURL, meta.RequestURLPath, meta.ChannelType), nil
	}
	baseURL := "http://172.40.233.2:8687"
	suffix := "chat/"
	if strings.HasPrefix(meta.ActualModelName, "Embedding") {
		suffix = "embeddings/"
	}
	if strings.HasPrefix(meta.ActualModelName, "bge-large") {
		suffix = "embeddings/"
	}
	if strings.HasPrefix(meta.ActualModelName, "tao-8k") {
		suffix = "embeddings/"
	}
	switch meta.ActualModelName {
	case "Qianfan-70b":
		suffix += "qianfan-70b"
	case "LLama2-7b":
		suffix += "llama-2-7b"
	case "ChatGLM3-6b":
		suffix += "chatglm3-6b"
	case "ERNIE-4.0":
		suffix += "llama-2-7b"
	case "ERNIE-Bot-4":
		suffix += "qianfan-70b"
	case "ERNIE-Bot":
		suffix += "chatglm3-6b"
	case "ERNIE-Bot-turbo":
		suffix += "eb-instant"
	case "ERNIE-Speed":
		suffix += "ernie_speed"
	case "ERNIE-4.0-8K":
		suffix += "completions_pro"
	case "ERNIE-3.5-8K":
		suffix += "completions"
	case "ERNIE-3.5-8K-0205":
		suffix += "ernie-3.5-8k-0205"
	case "ERNIE-3.5-8K-1222":
		suffix += "ernie-3.5-8k-1222"
	case "ERNIE-Bot-8K":
		suffix += "ernie_bot_8k"
	case "ERNIE-3.5-4K-0205":
		suffix += "ernie-3.5-4k-0205"
	case "ERNIE-Speed-8K":
		suffix += "ernie_speed"
	case "ERNIE-Speed-128K":
		suffix += "ernie-speed-128k"
	case "ERNIE-Lite-8K-0922":
		suffix += "eb-instant"
	case "ERNIE-Lite-8K-0308":
		suffix += "ernie-lite-8k"
	default:
		suffix += strings.ToLower(meta.ActualModelName)
	}
	fullRequestURL := fmt.Sprintf("%s/rpc/2.0/ai_custom/v1/wenxinworkshop/%s", baseURL, suffix)

	var accessToken string
	var err error
	if accessToken, err = GetAccessToken(meta.APIKey); err != nil {
		return "", err
	}
	fullRequestURL += "?AccessCode=" + accessToken
	fmt.Println("url is:", fullRequestURL)
	return fullRequestURL, nil
}

func (a *Adaptor) SetupRequestHeader(c *gin.Context, req *http.Request, meta *meta.Meta) error {
	adaptor.SetupCommonRequestHeader(c, req, meta)
	return nil
}

func (a *Adaptor) ConvertRequest(c *gin.Context, relayMode int, request *model.GeneralOpenAIRequest) (any, error) {
	if request == nil {
		return nil, errors.New("request is nil")
	}
	fmt.Println("******** baidu convertRequest *******")
	switch relayMode {
	case relaymode.Embeddings:
		baiduEmbeddingRequest := ConvertEmbeddingRequest(*request)
		return baiduEmbeddingRequest, nil
	default:
		baiduRequest := ConvertRequest(*request)
		return baiduRequest, nil
	}
}

func (a *Adaptor) ConvertImageRequest(request *model.ImageRequest) (any, error) {
	if request == nil {
		return nil, errors.New("request is nil")
	}
	return request, nil
}

func (a *Adaptor) DoRequest(c *gin.Context, meta *meta.Meta, requestBody io.Reader) (*http.Response, error) {
	return adaptor.DoRequestHelper(a, c, meta, requestBody)
}

func (a *Adaptor) DoResponse(c *gin.Context, resp *http.Response, meta *meta.Meta) (usage *model.Usage, err *model.ErrorWithStatusCode) {
	if meta.IsStream {
		err, usage = StreamHandler(c, resp)
	} else {
		switch meta.Mode {
		case relaymode.Embeddings:
			err, usage = EmbeddingHandler(c, resp)
		default:
			err, usage = Handler(c, resp)
		}
	}
	return
}

func (a *Adaptor) GetModelList() []string {
	return ModelList
}

func (a *Adaptor) GetChannelName() string {
	return "baidu"
}
