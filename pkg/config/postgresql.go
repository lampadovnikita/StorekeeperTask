package config

import (
	"io/ioutil"

	"gopkg.in/yaml.v2"
)

type PGConfig struct {
	Username string `yaml:"username"`
	Password string `yaml:"password"`
	Host     string `yaml:"host"`
	Port     string `yaml:"port"`
	Database string `yaml:"database"`
}

func GetPGConfig() (pgcfg *PGConfig, err error) {
	// TODO: посмотреть, что можно сделать с путём к файлу
	yamlFile, err := ioutil.ReadFile("../config/postgresql.yaml")
	if err != nil {
		return nil, err
	}

	pgcfg = &PGConfig{}

	err = yaml.Unmarshal(yamlFile, pgcfg)
	if err != nil {
		return nil, err
	}

	return pgcfg, nil
}
