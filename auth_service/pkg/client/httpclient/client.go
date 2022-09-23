package httpclient

import (
	"net/http"
	"time"
)

func GetClient() http.Client {
	return http.Client{Timeout: 10 * time.Second}
}
