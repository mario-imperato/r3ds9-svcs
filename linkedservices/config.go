package linkedservices

import (
	"github.com/GPA-Gruppo-Progetti-Avanzati-SRL/tpm-aws-common/s3/awss3lks"
	"github.com/GPA-Gruppo-Progetti-Avanzati-SRL/tpm-kafka-common/kafkalks"
	"github.com/GPA-Gruppo-Progetti-Avanzati-SRL/tpm-mongo-common/mongolks"
)

type Config struct {
	Kafka []kafkalks.Config `mapstructure:"kafka,omitempty" json:"kafka,omitempty" yaml:"kafka,omitempty"`
	Mongo []mongolks.Config `mapstructure:"mongo,omitempty" json:"mongo,omitempty" yaml:"mongo,omitempty"`
	S3    []awss3lks.Config `mapstructure:"aws-s3,omitempty" json:"aws-s3,omitempty" yaml:"aws-s3,omitempty"`
}

func (c *Config) PostProcess() error {

	var err error
	for i, _ := range c.Kafka {
		err = c.Kafka[i].PostProcess()
	}
	if err != nil {
		return err
	}

	return nil
}
