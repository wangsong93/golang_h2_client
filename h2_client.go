package main

import (
	"crypto/tls"
	"errors"
	"flag"
	"fmt"
	"io"
	"io/ioutil"
	"net"
	"net/http"
	"os"
	"strings"

	"golang.org/x/net/http2"
)

var (
	url      = flag.String("u", "https://www.baidu.com", "request url")
	method   = flag.String("m", http.MethodGet, "http method, e.g.: GET, POST, PUT, DELETE")
	header   = flag.String("h", "", "header use semicolon as separator, e.g.: User-Agent:Go Client;Accept:*/*")
	query    = flag.String("q", "", "query string, e.g.: a=b?c=d")
	body     = flag.String("b", "", "string body")
	useHttp2 = flag.Bool("http2", false, "has this flag:h2, no this flag:h2c")

	h2Client = &http.Client{
		Transport: &http2.Transport{
			DialTLS: func(network, addr string, cfg *tls.Config) (net.Conn, error) {
				cfg.InsecureSkipVerify = true // TODO skip or not
				return tls.Dial(network, addr, cfg)
			},
		},
	}
	h2cClient = &http.Client{
		Transport: &http2.Transport{
			AllowHTTP: true,
			DialTLS: func(network, addr string, cfg *tls.Config) (net.Conn, error) {
				return net.Dial(network, addr)
			},
		},
	}
)

func main() {
	flag.Parse()
	checkHttpMethod(*method)
	req, err := http.NewRequest(*method, *url+*query, strings.NewReader(*body))
	checkErrorAndExit("new request", err)
	var (
		resp *http.Response
	)
	if *useHttp2 {
		fmt.Println("use http2\n")
		resp, err = h2Client.Do(req)
	} else {
		fmt.Println("use http2 over cleartext\n")
		resp, err = h2cClient.Do(req)
	}
	checkErrorAndExit("do request", err)
	printResponse(resp)
}

func printResponse(resp *http.Response) {
	body, err := readBody(resp.Body)
	checkErrorAndExit("read response body", err)
	fmt.Printf(
		`Status:%s
StatusCode:%d
Proto:%s
ProtoMajor:%d
ProtoMinor:%d
Header:%v
Body:%s
ContentLength:%d
`, resp.Status, resp.StatusCode, resp.Proto, resp.ProtoMajor, resp.ProtoMinor, formatHeader(resp.Header), string(body), resp.ContentLength)
}

func checkErrorAndExit(msg string, err error) {
	if err != nil {
		fmt.Printf("msg:%s, err:%+v\n", msg, err)
		os.Exit(2)
	}
}

func checkHttpMethod(method string) string {
	switch method {
	case http.MethodGet, http.MethodPost, http.MethodPut, http.MethodDelete:
		return method
	default:
		checkErrorAndExit("method invalid", errors.New(method))
		return ""
	}
}

func formatHeader(header http.Header) string {
	if len(header) == 0 {
		return ""
	}
	buf := &strings.Builder{}
	idx := 0
	for k, v := range header {
		if idx > 0 {
			buf.WriteString("       ") // len("header:")
		}
		buf.WriteString(k + ": ")
		buf.WriteString(strings.Join(v, ","))
		idx++
		if idx != len(header) {
			buf.WriteByte('\n')
		}
	}
	return buf.String()
}

func readBody(body io.ReadCloser) ([]byte, error) {
	defer body.Close()
	return ioutil.ReadAll(body)
}
