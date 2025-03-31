package main

import (
	"bytes"
	"fmt"
	"time"

	vegeta "github.com/tsenart/vegeta/v12/lib"
)

func main() {
	// Define target
	targeter := vegeta.NewStaticTargeter(
		vegeta.Target{
			Method: "POST",
			URL:    "http://localhost:8001/v1/api/deposit/test",
			Header: map[string][]string{
				"Authorization": {"test"},
				"Content-Type":  {"application/json"},
			},
			Body: []byte(`{
				"routingKey": "user.registered.id1",
				"messageID": "msg-1",
				"email": "id1@email.com",
				"userID": "1",
				"hashKey": "batch-123"
			}`),
		},
		vegeta.Target{
			Method: "POST",
			URL:    "http://localhost:8001/v1/api/deposit/test",
			Header: map[string][]string{
				"Authorization": {"test"},
				"Content-Type":  {"application/json"},
			},
			Body: []byte(`{
				"routingKey": "user.registered.id1",
				"messageID": "msg-2",
				"email": "id2@email.com",
				"userID": "2",
				"hashKey": "batch-123"
			}`),
		},
		// vegeta.Target{
		// 	Method: "POST",
		// 	URL:    "http://localhost:8001/v1/api/deposit/test",
		// 	Header: map[string][]string{
		// 		"Authorization": {"test"},
		// 		"Content-Type":  {"application/json"},
		// 	},
		// 	Body: []byte(`{
		// 		"routingKey": "user.registered.id2",
		// 		"messageID": "msg-3",
		// 		"email": "id3@email.com",
		// 		"userID": "3",
		// 		"hashKey": "batch-124"
		// 	}`),
		// },
		// vegeta.Target{
		// 	Method: "POST",
		// 	URL:    "http://localhost:8001/v1/api/deposit/test",
		// 	Header: map[string][]string{
		// 		"Authorization": {"test"},
		// 		"Content-Type":  {"application/json"},
		// 	},
		// 	Body: []byte(`{
		// 		"routingKey": "user.registered.id2",
		// 		"messageID": "msg-4",
		// 		"email": "id4@email.com",
		// 		"userID": "4",
		// 		"hashKey": "batch-124"
		// 	}`),
		// },
	)

	// Define attack rate (2 requests per second)
	rate := vegeta.Rate{Freq: 2, Per: time.Second}
	duration := 1 * time.Second
	attacker := vegeta.NewAttacker()

	// Execute attack
	results := attacker.Attack(targeter, rate, duration, "API Load Test")
	var metrics vegeta.Metrics
	for res := range results {
		metrics.Add(res)
	}
	metrics.Close()

	// Print report
	var report bytes.Buffer
	vegeta.NewTextReporter(&metrics).Report(&report)
	fmt.Println(report.String())
}
