package main

import (
    "strings"

    "github.com/sirupsen/logrus"
    "github.com/spf13/viper"

    "tinysyslog/internal/app/tinysyslogd"
)

func main() {
    viper.SetEnvPrefix("tinysyslog")
    replacer := strings.NewReplacer("-", "_")
    viper.SetEnvKeyReplacer(replacer)
    viper.AutomaticEnv()

    server := tinysyslogd.NewServer()
    err := server.Run()
    if err != nil {
        logrus.Fatalf("Error starting server: %v", err)
    }
}
