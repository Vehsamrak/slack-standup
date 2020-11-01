package app

import (
	"fmt"
	"gopkg.in/yaml.v2"
	"io/ioutil"
	"log"
	"os"
)

type Config struct {
	Token string `yaml:"token"`
}

func (Config) Load() *Config {
	fmt.Println("Loading configuration file")

	file, err := os.Open("config.yml")
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
