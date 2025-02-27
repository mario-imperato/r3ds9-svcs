package apicms

import (
	"fmt"
	"github.com/GPA-Gruppo-Progetti-Avanzati-SRL/tpm-http-archive/har"
	"github.com/GPA-Gruppo-Progetti-Avanzati-SRL/tpm-http-client/restclient"
	"github.com/mario-imperato/r3ds9-svcs/clients"
	"github.com/rs/zerolog/log"
	"net/http"
	"net/url"
	"strings"
)

func (c *Client) Hello(reqCtx clients.ApiRequestContext) error {
	const semLogContext = "cms-api-client::hello"

	const ApiUrl = "/api/cms/:hostDomain/:hostSite/{:hostLang/hello"

	ep := c.helloUrl(ApiUrl, reqCtx.Dom, reqCtx.Ns, reqCtx.Lang, nil)

	req, err := c.client.NewRequest(http.MethodGet, ep, nil, reqCtx.GetHeaders("application/json"), nil)
	if err != nil {
		return clients.NewBadRequestError(clients.WithMessage(err.Error()))
	}

	harEntry, err := c.client.Execute(req,
		restclient.ExecutionWithOpName("api-cms-hello"),
		restclient.ExecutionWithRequestId(reqCtx.RequestId),
		restclient.ExecutionWithSpan(reqCtx.Span),
		restclient.ExecutionWithHarSpan(reqCtx.HarSpan))
	// c.harEntries = append(c.harEntries, harEntry)
	if err != nil {
		return clients.NewInternalServerError(clients.WithMessage(err.Error()))
	}

	log.Info().Interface("har", harEntry).Msg(semLogContext)
	return nil
}

func (c *Client) helloUrl(apiPath string, domain, site, lang string, qParams []har.NameValuePair) string {
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
