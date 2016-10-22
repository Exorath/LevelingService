package main

import (
	"errors"
	"golang.org/x/net/context"
	"net/http"
)







func main() {
	ctx := context.Background()
	svc := levelingService{}


	http.ListenAndServe(":8080", MakeHTTPHandler(ctx, svc))
}

// ErrEmpty is returned when an input string is empty.
var ErrEmpty = errors.New("empty string")