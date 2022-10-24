package Utils

import (
	"encoding/json"
	"fmt"
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
	r, err := c.Get(url)
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
