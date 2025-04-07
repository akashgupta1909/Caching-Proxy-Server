package main

import (
	"encoding/json"
	"io"
	"net/http"
	"net/url"
)

type redisClientHandler func(http.ResponseWriter, *http.Request, *RedisConfig)

func (redisClient *RedisConfig) redisClientWrapper(redisClientHandler redisClientHandler) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		redisClientHandler(w, r, redisClient)
	}
}

type ResponseCache struct {
	StatusCode int         `json:"status_code"`
	Header     http.Header `json:"header"`
	Body       string      `json:"body"`
}

func (config *Config) handleProxy(writer http.ResponseWriter, r *http.Request, redisClient *RedisConfig) {
	urlPath := r.URL.Path
	originURL := config.Origin.ResolveReference(&url.URL{Path: urlPath})

	cachedResponse, err := redisClient.GetHTTPResponse(urlPath)
	if err == nil {
		var cachedData ResponseCache
		err = json.Unmarshal(cachedResponse, &cachedData)
		if err != nil {
			http.Error(writer, "Error unmarshalling cached response", http.StatusInternalServerError)
			return
		}
		for key, values := range cachedData.Header {
			for _, value := range values {
				writer.Header().Add(key, value)
			}
		}
		writer.Header().Add("X-Cache", "HIT")
		writer.WriteHeader(cachedData.StatusCode)
		writer.Write([]byte(cachedData.Body))
		return
	}

	req, err := http.Get(originURL.String())
	if err != nil {
		http.Error(writer, "Error fetching the URL", http.StatusInternalServerError)
		return
	}
	defer req.Body.Close()
	bodyBytes, err := io.ReadAll(req.Body)
	if err != nil {
		http.Error(writer, "Error reading the response body", http.StatusInternalServerError)
		return
	}

	toStore := ResponseCache{
		StatusCode: req.StatusCode,
		Header:     req.Header,
		Body:       string(bodyBytes),
	}

	toStoreBytes, err := json.Marshal(toStore)
	if err != nil {
		http.Error(writer, "Error serializing the response for caching", http.StatusInternalServerError)
		return
	}
	err = redisClient.SaveHTTPResponse(urlPath, toStoreBytes)
	if err != nil {
		http.Error(writer, "Error storing the response in cache", http.StatusInternalServerError)
		return
	}

	for key, values := range toStore.Header {
		for _, value := range values {
			writer.Header().Add(key, value)
		}
	}

	writer.Header().Add("X-Cache", "MISS")
	writer.WriteHeader(toStore.StatusCode)
	writer.Write([]byte(toStore.Body))
}
