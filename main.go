package main

import (
	"fmt"

	"github.com/dan6erbond/jamboree-api/pkg"
	"github.com/spf13/viper"
)

func main() {
	viper.SetConfigName("config")
	viper.SetConfigType("yaml")
	viper.AddConfigPath(".")

	viper.SetDefault("server.port", 5001)
	viper.SetDefault("server.host", "0.0.0.0")

	err := viper.ReadInConfig()
	if err != nil {
		panic(fmt.Errorf("fatal error config file: %w", err))
	}

	app := pkg.NewApp()
	app.Run()
}
