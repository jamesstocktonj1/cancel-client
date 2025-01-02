package main

import (
	"context"
	"log"
	"net/http"
	"time"
)

const (
	CancelCount = 1
	NormalCount = 1

	// Request Config
	Method = "GET"
	Url    = "http://localhost:8080/"
)

func main() {

	// Cause instances to hang, one by one
	// sendCancelRequest()

	// Count number of instances running
	countNumberInstances()
}

func sendCancelRequest() {
	// HTTP Client configured to drop connection after response headers
	cancelClient := &http.Client{
		Transport: &http.Transport{
			ResponseHeaderTimeout: 2 * time.Microsecond,
		},
	}

	// Perform <CancelCount> of cancelled requests
	for i := 0; i < CancelCount; i++ {
		err := doRequest(cancelClient)
		if err != nil {
			log.Printf("Cancel Request: %s\n", err.Error())
		} else {
			log.Println("Cancel Request: success")
		}
	}

	time.Sleep(time.Second)

	// Regular HTTP Client
	normalClient := http.DefaultClient

	// Perform <NormalCount> of normal requests
	for i := 0; i < NormalCount; i++ {
		err := doRequest(normalClient)
		if err != nil {
			log.Printf("Normal Request: %s\n", err.Error())
		} else {
			log.Println("Normal Request: success")
		}
		time.Sleep(time.Second)
	}
}

func countNumberInstances() {
	cancelClient := &http.Client{
		Transport: &http.Transport{
			ResponseHeaderTimeout: 2 * time.Microsecond,
		},
	}
	normalClient := http.DefaultClient

	// Itterate until normal request returns error
	for i := 1; i < 100; i++ {
		// Perform cancelled request
		err := doRequest(cancelClient)
		if err != nil {
			log.Printf("Cancel Request %d: %s\n", i, err.Error())
		} else {
			log.Println("Cancel Request %d: success", i)
		}

		time.Sleep(10 * time.Millisecond)

		// Check if instances are still running
		err = doRequest(normalClient)
		if err != nil {
			log.Printf("Normal Request %d: %s\n", i, err.Error())
			log.Printf("\nComponent Instances: %d\n", i)
			return
		}
	}
}

func doRequest(client *http.Client) error {
	req, err := http.NewRequest(Method, Url, nil)
	if err != nil {
		return err
	}

	// Regular timeout context
	ctx, cancel := context.WithTimeout(req.Context(), time.Second)
	defer cancel()
	req = req.WithContext(ctx)

	// Send request
	_, err = client.Do(req)
	if err != nil {
		return err
	}
	return nil
}
