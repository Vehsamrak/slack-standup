package config

import (
	"fmt"
	"gopkg.in/yaml.v2"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
)

type Config struct {
	Token        string   `yaml:"token"`
	ApiUrl       string   `yaml:"api_url"`
	Port         int      `yaml:"port"`
	ChannelNames []string `yaml:"channel_names"`
}

func (Config) Load() *Config {
	dir, err := filepath.Abs(filepath.Dir(os.Args[0]))
	if err != nil {
		log.Fatal(err)
	}

	fmt.Printf("Loading configuration file: %s/config.yml", dir)

	file, err := os.Open(dir + "/config.yml")
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()

	fileContent, _ := ioutil.ReadAll(file)

	config := &Config{}

	err = yaml.Unmarshal(fileContent, &config)
	if err != nil {
		panic(err)
	}

	fmt.Printf("Configuration file loaded: %#v\n", config)

	return config
}
