package main

import (
	"flag"
	"io"
	"log"
	"net/http"
	"net/url"
	"strings"
	"time"

	"github.com/nxadm/tail"
)

func main() {
	if err := run(); err != nil {
		log.Fatal(err)
	}
}

func run() error {
	baseURL := flag.String("server-url", "http://127.0.0.1:8080", "AlphaServer URL")
	file := flag.String("log-file", "/config/logs/openssh/current", "OpenSSH server log file")
	clientID := flag.String("client-id", "node-example", "AlphaClient client id")
	flag.Parse()

	t, err := tail.TailFile(*file, tail.Config{
		Follow:    true,
		MustExist: true,
		Location: &tail.SeekInfo{
			Whence: io.SeekEnd,
		},
	})
	if err != nil {
		return err
	}

	httpClient := http.Client{
		Timeout: 5 * time.Second,
	}

	for line := range t.Lines {
		if !strings.Contains(line.Text, "Accepted") && !strings.Contains(line.Text, "Failed") {
			continue
		}

		values := url.Values{
			"client-id": {*clientID},
		}
		req, err := http.NewRequest(http.MethodPost, *baseURL+"/increment", strings.NewReader(values.Encode()))
		if err != nil {
			return err
		}

		req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
		if _, err := httpClient.Do(req); err != nil {
			return err
		}
	}

	return nil
}
