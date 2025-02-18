package main

import (
	"ai/infrastructure"
	"flag"
	"fmt"
	"log"
	"os"

	// Update with your actual protobuf package path
	// Update with the correct path to your server package

	"github.com/spf13/viper"
)

var (
	version string
	build   string
	runEnv  string
)

func init() {
	// read running flag
	if len(os.Getenv("ENV")) != 0 {
		runEnv = os.Getenv("ENV")
	} else {
		flagEnv := flag.String("env", "dev", "A config file name without .env")
		flag.Parse()
		runEnv = *flagEnv
	}
	if err := LoadConfig(runEnv); err != nil {
		log.Fatalf("error while loading the env:\n %+v", err)
	}
}
func LoadConfig(configName string) (err error) {
	if len(configName) == 0 {
		configName = "dev"
	}
	viper.SetConfigName(configName) // name of config file (without extension)
	viper.SetConfigType("yaml")     // REQUIRED if the config file does not have the extension in the name
	viper.AddConfigPath("./config") // path to look for the config file in
	viper.AddConfigPath(".")        // optionally look for config in the working directory

	viper.AutomaticEnv()

	err = viper.ReadInConfig() // Find and read the config file
	if err != nil {            // Handle errors reading the config file
		panic(fmt.Errorf("fatal error config file: %w", err))
	}

	return
}
func main() {
	server, err := infrastructure.NewServer(version, build, runEnv)
	if err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}
	server.Run()
}
