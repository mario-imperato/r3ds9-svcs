package linkedservices

import (
	"github.com/GPA-Gruppo-Progetti-Avanzati-SRL/tpm-aws-common/s3/awss3lks"
	"github.com/GPA-Gruppo-Progetti-Avanzati-SRL/tpm-kafka-common/kafkalks"
	"github.com/GPA-Gruppo-Progetti-Avanzati-SRL/tpm-mongo-common/mongolks"

	"github.com/rs/zerolog/log"
)

type ServiceRegistry struct {
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

	return nil
}
