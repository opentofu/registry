package http

import (
	"github.com/hashicorp/go-retryablehttp"
	"net/http"
)

func GetHttpRetryClient() *http.Client {
	retryClient := retryablehttp.NewClient()
	retryClient.RetryMax = 10
	retryClient.Logger = nil

	return retryClient.StandardClient()
}
