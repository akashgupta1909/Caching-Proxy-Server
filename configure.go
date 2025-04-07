package main

import (
	"flag"
	"fmt"
	"net/url"
)

type Config struct {
	Port       int
	Origin     *url.URL
	CleanCache bool
}

func initCommand() (Config, error) {
	var config Config

	flag.Usage = func() {
		fmt.Println("Usage: go run main.go [flags]")
		fmt.Println("cache-proxy [flags]")
		fmt.Println("Flags:")
		fmt.Println("  -h, --help            Show this help message and exit")
		fmt.Println("  -p, --port <port>     Port to run the server on (default: 8080)")
		fmt.Println("  -o, --origin <url>    Origin to make proxy requests to")
		fmt.Println("  -c, --clean-cache     Clean the cache before starting the server")
	}

	portLongFlag := flag.Int("port", 0, "Port to run the server on")
	originLongFlag := flag.String("origin", "", "Origin to make proxy requests to")
	cleanCacheLongFlag := flag.Bool("clean-cache", false, "Clean the cache before starting the server")

	portShortFlag := flag.Int("p", 0, "Port to run the server on")
	originShortFlag := flag.String("o", "", "Origin to make proxy requests to")
	cleanCacheShortFlag := flag.Bool("c", false, "Clean the cache before starting the server")

	flag.Parse()

	cleanCache := *cleanCacheLongFlag
	if !*cleanCacheLongFlag && *cleanCacheShortFlag {
		cleanCache = *cleanCacheShortFlag
	}

	if cleanCache {
		config = Config{
			CleanCache: true,
		}
		return config, nil
	}

	origin := *originLongFlag
	if origin == "" {
		origin = *originShortFlag
	}
	if origin == "" {
		return config, fmt.Errorf("error: --origin flag is required")
	}

	port := *portLongFlag
	if port == 0 {
		port = *portShortFlag
	}
	if port == 0 {
		port = 8080
	}

	originURL, err := url.Parse(origin)
	if err != nil {
		return config, fmt.Errorf("error: --origin must be a valid URL: %v", err)
	}
	if originURL.Scheme == "" {
		return config, fmt.Errorf("error: --origin must be a valid URL with scheme (http or https)")
	}
	config = Config{
		Port:       port,
		Origin:     originURL,
		CleanCache: cleanCache,
	}

	fmt.Printf("Arguments parsed: port=%d, origin=%s, clean-cache=%t\n", config.Port, config.Origin, config.CleanCache)
	return config, nil
}
