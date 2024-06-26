package http

import (
	"fmt"
	"log"
	"net/http"
	"regexp"
	"strconv"
	"strings"
	"url-shortener/internal/service"
)

type UrlHandler struct {
	urlService *service.UrlService
}

func NewUrlHandler(urlService *service.UrlService) *UrlHandler {
	return &UrlHandler{urlService}
}

func (handler *UrlHandler) Get(writer http.ResponseWriter, request *http.Request) {
	requestURIValid, err := regexp.MatchString(`/s/[A-z0-9]{8}$`, request.RequestURI)
	if !requestURIValid {
		writer.WriteHeader(http.StatusBadRequest)
		log.Println("invalid request uri")
		return
	}
	if err != nil {
		writer.WriteHeader(http.StatusBadRequest)
		log.Println(err.Error())
		return
	}

	splittedRequestURI := strings.Split(request.RequestURI, "/")
	hash := splittedRequestURI[len(splittedRequestURI)-1]

	url, err := (*handler.urlService).GetUrl(hash)
	if err != nil {
		writer.WriteHeader(http.StatusNotFound)
		log.Println(err.Error())
		return
	}

	http.Redirect(writer, request, url, http.StatusFound)
	fmt.Fprintf(writer, "status: %s; url: %s", strconv.Itoa(http.StatusFound), url)
}

func (handler *UrlHandler) Set(writer http.ResponseWriter, request *http.Request) {
	requestURIValid := strings.Contains(request.RequestURI, "/a/?url=")
	if !requestURIValid {
		writer.WriteHeader(http.StatusBadRequest)
		log.Println("invalid request uri")
		return
	}

	requestUrl := request.URL.Query().Get("url")
	if requestUrl == "" {
		writer.WriteHeader(http.StatusBadRequest)
		log.Println("query parameter 'url' is missing")
		return
	}

	hash, err := (*handler.urlService).SetUrl(requestUrl)
	if err != nil {
		writer.WriteHeader(http.StatusInternalServerError)
		log.Println(err.Error())
		return
	}

	writer.WriteHeader(http.StatusOK)
	fmt.Fprintf(writer, "status: %s; hash: %s", strconv.Itoa(http.StatusOK), hash)
}
