package main

import (
	"io"
	"io/ioutil"
	"net/http"
)

func sendGetRequest(url string) ([]byte, error) {

	// Create a new request using http.NewRequest()
	req, err := http.NewRequest("GET", url, nil)

	if err != nil {
		return nil, err
	}

	// Send the request and retrieve the response using http.DefaultClient.Do()
	resp, err := http.DefaultClient.Do(req)

	if err != nil {
		return nil, err
	}

	// Read the response body using ioutil.ReadAll()
	body, err := ioutil.ReadAll(resp.Body)

	if err != nil {
		return nil, err
	}

	// Return the response body and no error to indicate success
	return body, nil
}

func sendPostRequest(url string, bodyRequest io.Reader) ([]byte, error) {

	// Create a new POST request with the given body as parameter.
	req, err := http.NewRequest("POST", url, bodyRequest)

	if err != nil { // Handle error if any.
		panic(err) // Panic if there is an error.
	}

	// Send the request and get the response from the server.
	resp, err := http.DefaultClient.Do(req)

	if err != nil { // Handle error if any.
		panic(err) // Panic if there is an error.
	}

	defer resp.Body.Close() // Close the response body when finished with it.

	// Read the response body and store it in a variable for later use.
	body, err := ioutil.ReadAll(resp.Body)

	if err != nil { // Handle error if any.
		panic(err) // Panic if there is an error.
	}

	println(string(body)) // Print out the response body as a string for debugging purposes.
	return body, err
}
