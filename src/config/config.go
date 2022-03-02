package config

import (
	"fmt"
	"io/ioutil"

	"gopkg.in/yaml.v2"
)

type Config struct {
	MongoDB struct {
		URL string `yaml:"url"`
	} `yaml:"mongodb"`
	RabbitMQ struct {
		URL           string `yaml:"url"`
		CommandQueue  Queue  `yaml:"command_queue"`
		PostsQueue    Queue  `yaml:"posts_queue"`
		MessagesTopic Topic  `yaml:"messages_topic"`
	} `yaml:"rabbitmq"`
}

type Queue struct {
	Name string `yaml:"name"`
}

type Topic struct {
	Name string `yaml:"name"`
}

func GetConfig(fileName string) (Config, error) {
	yamlFile, err := ioutil.ReadFile(fileName)
	if err != nil {
		return Config{}, fmt.Errorf("Error reading config file '%s' | error: %s", fileName, err.Error())
	}

	var config Config
	err = yaml.Unmarshal(yamlFile, &config)
	if err != nil {
		return Config{}, fmt.Errorf("Error parsing config file '%s' | error: %s", fileName, err.Error())
	}

	return config, nil
}
