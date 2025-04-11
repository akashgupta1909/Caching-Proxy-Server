package main

import (
	"fmt"
	"os"

	"github.com/joho/godotenv"
)

func main() {
	godotenv.Load()

	config, err := initCommand()
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	redisConfig, err := initRedisClient()
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	defer redisConfig.Client.Close()

	if config.CleanCache {
		if err := redisConfig.FlushAll(); err != nil {
			fmt.Println(err)
			os.Exit(1)
		}
		fmt.Println("Cache Cleaned")
		return
	}

	if err := config.initRouter(&redisConfig); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
