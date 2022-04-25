package main

import (
	"bytes"
	"fmt"
	"log"
	"net/http"
	"regexp"
	"sort"
	"strings"

	"go.clly.me/minic/sdk"
)

func main() {
	fmt.Println("hello world")

	mux := sdk.NewMux()
	handlers(mux)
	srv := sdk.ConfigureHTTP(mux)

	// TODO: Easier signal handling
	log.Fatal(srv.ListenAndServe())
}

func handlers(m *http.ServeMux) {
	if m == nil {
		panic("serve mux is nil")
	}
	m.HandleFunc("/", handler())

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
