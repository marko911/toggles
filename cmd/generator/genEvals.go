package main

import (
	"fmt"
	"math/rand"
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
	// var conn net.Conn
	fmt.Println("GEN EVALS STARTED")
	// client := &http.Client{
	// 	Transport: &http.Transport{
	// 		Dial: func(netw, addr string) (net.Conn, error) {
	// 			if conn != nil {
	// 				conn.Close()
	// 				conn = nil
	// 			}
	// 			netc, err := net.DialTimeout(netw, addr, 1*time.Second)
	// 			if err != nil {
	// 				return nil, err
	// 			}
	// 			conn = netc
	// 			return netc, nil
	// 		},
	// 	},
	// }
	// u, _ := url.Parse("http://localhost:8080/evals/flags/JDJhJDEwJDF0QkFYcGZYUC9OSHh1bFgybkNuT09wOXl0aWRZNzhsSXdpaUtVdFdyeFEuRmV2cFFPa3JX")
	// for i := 0; i < 10; i++ {
	// 	body := map[string]interface{}{
	// 		"flagKey": "on-off-bool",
	// 		"user": map[string]interface{}{
	// 			"key": fmt.Sprintf("%v@gmail.com", genString(8)),
	// 			"attributes": map[string]interface{}{
	// 				"groups": []string{"beta testers"},
	// 			},
	// 		}}
	// 	mars, _ := json.Marshal(body)
	// 	req, err := http.NewRequest("POST", u.String(), bytes.NewBuffer(mars))
	// 	if err != nil {
	// 		log.Println("creating post request failed:", err)
	// 	}
	// 	resp, err := client.Do(req)
	// 	if err != nil {
	// 		log.Println("Error getting response:", err)
	// 		continue
	// 	}
	// 	if resp.StatusCode != http.StatusOK && resp.StatusCode != http.StatusCreated {
	// 		log.Println("didnt work", resp.StatusCode)
	// 	}
	// }
}
