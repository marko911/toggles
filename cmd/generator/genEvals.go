package main

import (
	"bytes"
	"encoding/json"
	"log"
	"math/rand"
	"net"
	"net/http"
	"net/url"
	"time"
)

const charset = "abcdefghijklmnopqrstuvwxyz" +
	"ABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"

var seededRand *rand.Rand = rand.New(
	rand.NewSource(time.Now().UnixNano()))

func stringWithCharset(length int, charset string) string {
	b := make([]byte, length)
	for i := range b {
		b[i] = charset[seededRand.Intn(len(charset))]
	}
	return string(b)
}

func genString(length int) string {
	return stringWithCharset(length, charset)
}

// Used to generate evaluations on a flag
func main() {
	var conn net.Conn

	client := &http.Client{
		Transport: &http.Transport{
			Dial: func(netw, addr string) (net.Conn, error) {
				if conn != nil {
					conn.Close()
					conn = nil
				}
				netc, err := net.DialTimeout(netw, addr, 1*time.Second)
				if err != nil {
					return nil, err
				}
				conn = netc
				return netc, nil
			},
		},
	}
	u, _ := url.Parse("http://localhost:8080/evals/flags/JDJhJDEwJGZneG8xYzdNUmdhZUdOMUdCY3AvWnU0RlhnbEFRQWdSVmRrZGk4bzNmQTdWU2MzQ0RlbFl1")
	for i := 0; i < 2; i++ {
		body := map[string]interface{}{
			"flagKey": "fetch-flags-was",
			"user": map[string]interface{}{
				"key": genString(8),
			}}
		mars, _ := json.Marshal(body)
		req, err := http.NewRequest("POST", u.String(), bytes.NewBuffer(mars))
		if err != nil {
			log.Println("creating post request failed:", err)
		}
		resp, err := client.Do(req)
		if err != nil {
			log.Println("Error getting response:", err)
			continue
		}
		if resp.StatusCode != http.StatusOK && resp.StatusCode != http.StatusCreated {
			log.Println("didnt work", resp.StatusCode)
		}
	}
}
