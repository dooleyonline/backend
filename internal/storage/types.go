package storage

import "net/http"

type PresignRequest struct {
	Method      string
	Bucket      string
	Key         string
	ContentType string
}

type PresignResponse struct {
	URL    string
	Header http.Header
}