package clients

import (
	"github.com/GPA-Gruppo-Progetti-Avanzati-SRL/tpm-common/util"
	"github.com/GPA-Gruppo-Progetti-Avanzati-SRL/tpm-http-archive/har"
	"github.com/GPA-Gruppo-Progetti-Avanzati-SRL/tpm-http-archive/hartracing"
	"github.com/GPA-Gruppo-Progetti-Avanzati-SRL/tpm-http-client/restclient"
	"github.com/opentracing/opentracing-go"
)

type ApiRequestContext struct {
	Dom       string              `json:"dom,omitempty" bson:"dom,omitempty"`
	Ns        string              `json:"ns,omitempty" bson:"ns,omitempty"`
	Lang      string              `json:"lang,omitempty" bson:"lang,omitempty"`
	XAPIKey   string              `yaml:"x-api-key,omitempty" mapstructure:"x-api-key,omitempty" json:"x-api-key,omitempty"`
	RequestId string              `yaml:"request-id,omitempty" mapstructure:"request-id,omitempty" json:"request-id,omitempty"`
	Headers   []restclient.Header `mapstructure:"headers,omitempty" json:"headers,omitempty" yaml:"headers,omitempty"`
	Span      opentracing.Span    `yaml:"-" mapstructure:"-" json:"-"`
	HarSpan   hartracing.Span     `yaml:"-" mapstructure:"-" json:"-"`
}

const (
	ApiKeyHeaderName      = "X-R3ds9-Api-Key"
	SidHeaderName         = "X-R3ds9-Sid"
	UserHeaderNickName    = "X-R3ds9-Nickname"
	UserHeaderUserId      = "X-R3ds9-Userid"
	RequestIdHeaderName   = "Request-Id"
	ContentTypeHeaderName = "Content-Type"
)

var headers2Propagate = []string{
	SidHeaderName,
	UserHeaderNickName,
	UserHeaderUserId,
}

func (arc *ApiRequestContext) GetHeaders(ct string) []har.NameValuePair {

	const semLogContext = "tpm-tokens-client::get-headers"
	var nvp []har.NameValuePair

	if arc.RequestId == "" {
		arc.RequestId = util.NewObjectId().String()
	}
	nvp = append(nvp, har.NameValuePair{Name: RequestIdHeaderName, Value: arc.RequestId})

	if arc.XAPIKey != "" {
		nvp = append(nvp, har.NameValuePair{Name: ApiKeyHeaderName, Value: arc.XAPIKey})
	}

	if ct != "" {
		nvp = append(nvp, har.NameValuePair{Name: ContentTypeHeaderName, Value: ct})
	}

	for _, h := range arc.Headers {
		nvp = append(nvp, har.NameValuePair{Name: h.Name, Value: h.Value})
	}
	return nvp
}

func NewApiRequestContext(domain, site, lang string, apiKey string, contextHeaders map[string][]string, span opentracing.Span, harSpan hartracing.Span) ApiRequestContext {

	var hds []restclient.Header
	if len(contextHeaders) > 0 {
		for _, n := range headers2Propagate {
			if h, ok := contextHeaders[n]; ok {
				hds = append(hds, restclient.Header{
					Name:  n,
					Value: h[0],
				})
			}
		}
	}

	ctx := ApiRequestContext{
		Dom:       domain,
		Ns:        site,
		Lang:      lang,
		XAPIKey:   apiKey,
		RequestId: util.NewUUID(),
		Headers:   hds,
		Span:      span,
		HarSpan:   harSpan,
	}

	return ctx
}
