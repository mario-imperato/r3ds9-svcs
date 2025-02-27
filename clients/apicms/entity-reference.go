package apicms

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/GPA-Gruppo-Progetti-Avanzati-SRL/tpm-http-archive/har"
	"github.com/GPA-Gruppo-Progetti-Avanzati-SRL/tpm-http-client/restclient"
	"github.com/mario-imperato/r3ds9-svcs/clients"
	"github.com/rs/zerolog/log"
	"net/http"
	"net/url"
	"strings"
)

const (
	OpAddReference = "add"
	OpDelReference = "del"
)

type UpdateEntityRefsReq struct {
	OpType string          `json:"op,omitempty" bson:"op,omitempty"`
	FileId string          `json:"file-id,omitempty" bson:"file-id,omitempty"`
	Ref    EntityReference `json:"refs,omitempty" bson:"refs,omitempty"`
}

// EntityReference see https://github.com/mario-imperato/r3ds9-mongodb/blob/6b90ded13f0e24330ce898a0703dbb10fd477016/model/r3ds9-apicms/file/model.go#EntRefStruct
type EntityReference struct {
	Dom     string `json:"dom,omitempty" bson:"dom,omitempty"`
	Ns      string `json:"ns,omitempty" bson:"ns,omitempty"`
	EntType string `json:"entType,omitempty" bson:"entType,omitempty"`
	EntId   string `json:"entId,omitempty" bson:"entId,omitempty"`
}

func (tok *EntityReference) ToJSON() ([]byte, error) {
	return json.Marshal(tok)
}

func (c *Client) UpdateFileEntityReferences(reqCtx clients.ApiRequestContext, tokenRequest []UpdateEntityRefsReq) (*clients.ApiResponse, error) {
	const semLogContext = "cms-api-client::add-file-entity-reference"

	const ApiUrl = "/api/cms/:hostDomain/:hostNamespace/:hostLang/files/updateRefs"

	ep := c.entityReferenceUrl(ApiUrl, reqCtx.Dom, reqCtx.Ns, reqCtx.Lang, nil)

	b, err := json.Marshal(tokenRequest)
	if err != nil {
		return nil, clients.NewBadRequestError(clients.WithMessage(err.Error()))
	}

	req, err := c.client.NewRequest(http.MethodPut, ep, b, reqCtx.GetHeaders("application/json"), nil)
	if err != nil {
		return nil, clients.NewBadRequestError(clients.WithMessage(err.Error()))
	}

	harEntry, err := c.client.Execute(req,
		restclient.ExecutionWithOpName("api-cms-upd-entity-refs"),
		restclient.ExecutionWithRequestId(reqCtx.RequestId),
		restclient.ExecutionWithSpan(reqCtx.Span),
		restclient.ExecutionWithHarSpan(reqCtx.HarSpan))
	// c.harEntries = append(c.harEntries, harEntry)
	if err != nil {
		return nil, clients.NewInternalServerError(clients.WithMessage(err.Error()))
	}

	resp, err := DeserializeTAddEntityReferenceResponseBody(harEntry)
	return resp, err
}

func DeserializeTAddEntityReferenceResponseBody(resp *har.Entry) (*clients.ApiResponse, error) {

	const semLogContext = "cms-api-client::deserialize-add-file-entity-reference"
	if resp == nil || resp.Response == nil || resp.Response.Content == nil || resp.Response.Content.Data == nil {
		err := errors.New("cannot deserialize null response")
		log.Error().Err(err).Msg(semLogContext)
		return nil, clients.NewInternalServerError(clients.WithMessage(err.Error()))
	}

	var apiResponse clients.ApiResponse
	var err error
	switch resp.Response.Status {
	case http.StatusOK:
		apiResponse, err = clients.DeserializeApiResponseFromJson(resp.Response.Content.Data)
		if err != nil {
			return nil, clients.NewInternalServerError(clients.WithMessage(err.Error()))
		}

	default:
		apiResponse, err = clients.DeserializeApiResponseFromJson(resp.Response.Content.Data)
		if err != nil {
			return nil, clients.NewInternalServerError(clients.WithMessage(err.Error()))
		}
		err = &apiResponse
		return nil, err
	}

	return &apiResponse, nil
}

func (c *Client) entityReferenceUrl(apiPath string, domain, site, lang string, qParams []har.NameValuePair) string {
	var sb = strings.Builder{}
	sb.WriteString(c.host.Scheme)
	sb.WriteString("://")
	sb.WriteString(c.host.HostName)
	sb.WriteString(":")
	sb.WriteString(fmt.Sprint(c.host.Port))

	apiPath = strings.Replace(apiPath, HostDomainPathPlaceHolder, domain, 1)
	apiPath = strings.Replace(apiPath, HostNamespacePathPlaceHolder, site, 1)
	apiPath = strings.Replace(apiPath, HostLangPathPlaceHolder, lang, 1)
	sb.WriteString(apiPath)

	if len(qParams) > 0 {
		sb.WriteString("?")
		for i, qp := range qParams {
			if i > 0 {
				sb.WriteString("&")
			}
			sb.WriteString(qp.Name)
			sb.WriteString("=")
			sb.WriteString(url.QueryEscape(qp.Value))
		}
	}
	return sb.String()
}
