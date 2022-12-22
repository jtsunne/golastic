package Utils

import (
	"encoding/json"
	"fmt"
	"github.com/tidwall/pretty"
	"net"
	"net/http"
	"net/url"
	"os"
	"strings"
	"time"
)

var (
	c = &http.Client{Timeout: 10 * time.Second}
)

func GetJson(url string, target interface{}) error {
	req, _ := http.NewRequest(http.MethodGet, url, nil)
	req.Header.Set("Content-Type", "application/json; charset=UTF-8")
	r, err := c.Do(req)
	if err != nil {
		return err
	}
	defer r.Body.Close()

	return json.NewDecoder(r.Body).Decode(target)
}

func PostJson(url string, body string, target interface{}) error {
	req, _ := http.NewRequest(http.MethodPost, url, strings.NewReader(body))
	req.Header.Set("Content-Type", "application/json; charset=UTF-8")
	r, err := c.Do(req)
	if err != nil {
		return err
	}
	defer r.Body.Close()

	return json.NewDecoder(r.Body).Decode(target)
}

func ParseEsUrl(u string) string {
	schema := "http"
	port := "9200"
	domain := ""
	s, err := url.Parse(u)
	if s.Host == "" {
		s, err = url.Parse("//" + u)
	}
	if err != nil {
		fmt.Println("Can't parse ESURL. Exiting...")
		os.Exit(1)
	}
	if s.Scheme == "https" {
		fmt.Println("GoLastic can't work with https yet. Exiting...")
		os.Exit(1)
	}
	if s.Port() != "" {
		port = s.Port()
	}
	if strings.Contains(s.Host, ":") {
		domain, port, _ = net.SplitHostPort(s.Host)
	} else {
		domain = s.Host
	}
	return fmt.Sprintf("%s://%s:%s", schema, domain, port)
}

func PrettyJson(s string) string {
	return string(pretty.Pretty([]byte(s)))
}

func ColorizeJson(s string) string {
	var r string
	open := false
	for _, chr := range strings.Split(s, "") {
		if chr == "\"" && !open {
			open = true
			r = r + "[darkred]\""
		} else if chr == "\"" && open {
			open = false
			r = r + "\"[white]"
		} else {
			r = r + chr
		}
	}
	sa := strings.Split(r, "")
	r = ""
	for _, chr := range sa {
		switch chr {
		case "{":
			r = r + "[yellow]{[white]"
		case "}":
			r = r + "[yellow]}[white]"
		default:
			r = r + chr
		}
	}

	return r
}
