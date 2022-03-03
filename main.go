package main

import (
	"bytes"
	"fmt"
	"go.clly.me/minic/sdk"
	"log"
	"net/http"
	"regexp"
	"sort"
	"strings"
)

func main() {
	fmt.Println("hello world")
	mux := serveMux()

	srv := sdk.ConfigureHTTP(mux)

	// TODO: Easier signal handling
	log.Fatal(srv.ListenAndServe())
}

func serveMux() *http.ServeMux {
	m := http.NewServeMux()
	m.HandleFunc("/", handler())
	return m
}

func handler() http.HandlerFunc {
	return func(writer http.ResponseWriter, request *http.Request) {
		if isRobot(request.UserAgent()) {
			remoteAddr := request.Header.Get("X-Forwarded-For")
			if remoteAddr == "" {
				remoteAddr = strings.Split(request.RemoteAddr, ":")[0]
			}
			fmt.Fprintln(writer, remoteAddr)
			return
		}
		resp := headersFromReq(request)
		fmt.Fprint(writer, resp)
	}
}

func headersFromReq(req *http.Request) string {
	b := bytes.NewBuffer(make([]byte, 0, 4096))
	fmt.Fprintln(b, "Headers:")
	headers := []string{}
	for k, v := range req.Header {
		headers = append(headers, fmt.Sprintf("\t%s: %s\n", k, v))
	}
	sort.Strings(headers)
	for _, v := range headers {
		b.WriteString(v)
	}
	return b.String()
}

var robotRegexp = regexp.MustCompile(`^curl/.*$`)

func isRobot(ua string) bool {
	return robotRegexp.MatchString(ua)
}
