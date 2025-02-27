package linkedservices

import (
	"github.com/GPA-Gruppo-Progetti-Avanzati-SRL/tpm-aws-common/s3/awss3lks"
	"github.com/GPA-Gruppo-Progetti-Avanzati-SRL/tpm-http-client/restclient"
	"github.com/GPA-Gruppo-Progetti-Avanzati-SRL/tpm-kafka-common/kafkalks"
	"github.com/GPA-Gruppo-Progetti-Avanzati-SRL/tpm-mongo-common/mongolks"
	"github.com/mario-imperato/r3ds9-svcs/clients/apicms"

	"github.com/rs/zerolog/log"
)

type ServiceRegistry struct {
	ApiCmsLks *apicms.LinkedService
}

var registry ServiceRegistry

func InitRegistry(cfg *Config) error {

	registry = ServiceRegistry{}
	log.Info().Msg("initialize services registry")

	_, err := mongolks.Initialize(cfg.Mongo)
	if err != nil {
		return err
	}

	_, err = awss3lks.Initialize(cfg.S3)
	if err != nil {
		return err
	}

	_, err = kafkalks.Initialize(cfg.Kafka)
	if err != nil {
		return err
	}

	err = initializeApiCmsClientLinkedService(cfg.CmsClientCfg)
	if err != nil {
		return err
	}

	return nil
}

/*
 * TokensApiClient Initialization
 */

func initializeApiCmsClientLinkedService(cfg *apicms.Config) error {
	const semLogContext = "service-registry::initialize-cms-api-client-provider"
	log.Info().Msg(semLogContext)
	if cfg != nil {
		lks, err := apicms.NewInstance(cfg)
		if err != nil {
			return err
		}

		registry.ApiCmsLks = lks
	}

	return nil
}

func NewApiCmsClient(opts ...restclient.Option) (*apicms.Client, error) {
	return registry.ApiCmsLks.NewClient(opts...)
}
