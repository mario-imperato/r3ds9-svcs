package apicms_test

import (
	"github.com/GPA-Gruppo-Progetti-Avanzati-SRL/tpm-common/util"
	"github.com/GPA-Gruppo-Progetti-Avanzati-SRL/tpm-http-archive/hartracing"
	"github.com/GPA-Gruppo-Progetti-Avanzati-SRL/tpm-http-archive/hartracing/filetracer"
	"github.com/GPA-Gruppo-Progetti-Avanzati-SRL/tpm-http-client/restclient"
	"github.com/mario-imperato/r3ds9-svcs/clients"
	"github.com/mario-imperato/r3ds9-svcs/clients/apicms"
	"github.com/opentracing/opentracing-go"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	"github.com/stretchr/testify/require"
	jaegercfg "github.com/uber/jaeger-client-go/config"
	"github.com/uber/jaeger-lib/metrics"
	"io"
	"os"
	"testing"
)

var cliConfig = apicms.Config{
	Config: restclient.Config{
		RestTimeout:      0,
		SkipVerify:       false,
		Headers:          nil,
		TraceGroupName:   "r3ds9-client-apicms",
		TraceRequestName: "",
		RetryCount:       0,
		RetryWaitTime:    0,
		RetryMaxWaitTime: 0,
		RetryOnHttpError: nil,
	},
	Host: apicms.HostInfo{
		Scheme:   "http",
		HostName: "localhost",
		Port:     8082,
	},
}

var reqCtx = clients.ApiRequestContext{
	Dom:       "cvf",
	Ns:        "cvfchamp42",
	Lang:      "it",
	XAPIKey:   "r3ds9-api-key-01",
	RequestId: util.NewUUID(),
	Headers: []restclient.Header{
		{
			Name:  clients.SidHeaderName,
			Value: "123-stella",
		},
		{
			Name:  clients.UserHeaderNickName,
			Value: "root",
		},
		{
			Name:  clients.UserHeaderUserId,
			Value: "321-stella",
		},
	},
	Span:    nil,
	HarSpan: nil,
}

func TestMain(m *testing.M) {
	exitVal := m.Run()
	os.Exit(exitVal)
}

func TestApiCmsClient(t *testing.T) {

	log.Logger = log.Output(zerolog.ConsoleWriter{Out: os.Stderr})
	c, err := InitTracing(t)
	require.NoError(t, err)
	if c != nil {
		defer c.Close()
	}

	cli, err := apicms.NewClient(&cliConfig)
	require.NoError(t, err)
	defer cli.Close()

	err = cli.Liveness(reqCtx)
	require.NoError(t, err)

	err = cli.Hello(reqCtx)
	require.NoError(t, err)

}

const (
	JAEGER_SERVICE_NAME = "JAEGER_SERVICE_NAME"
)

func InitHarTracing(t *testing.T) (io.Closer, error) {
	trc, c, err := filetracer.NewTracer()
	if err != nil {
		return nil, err
	}

	hartracing.SetGlobalTracer(trc)
	return c, nil
}

func InitTracing(t *testing.T) (io.Closer, error) {

	if os.Getenv(JAEGER_SERVICE_NAME) == "" {
		t.Log("skipping jaeger config no vars in env.... (" + JAEGER_SERVICE_NAME + ")")
		return nil, nil
	}

	t.Log("initialize jaeger service " + os.Getenv(JAEGER_SERVICE_NAME))

	var tracer opentracing.Tracer
	var closer io.Closer

	jcfg, err := jaegercfg.FromEnv()
	if err != nil {
		log.Warn().Err(err).Msg("Unable to configure JAEGER from environment")
		return nil, err
	}

	tracer, closer, err = jcfg.NewTracer(
		jaegercfg.Logger(&jlogger{}),
		jaegercfg.Metrics(metrics.NullFactory),
	)
	if nil != err {
		log.Error().Err(err).Msg("Error in NewTracer")
		return nil, err
	}

	opentracing.SetGlobalTracer(tracer)

	return closer, nil
}

type jlogger struct{}

func (l *jlogger) Error(msg string) {
	log.Error().Msg("(jaeger) " + msg)
}

func (l *jlogger) Infof(msg string, args ...interface{}) {
	log.Info().Msgf("(jaeger) "+msg, args...)
}
