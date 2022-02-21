package main

import (
	"context"
	"encoding/json"
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"

	"golang.org/x/oauth2/clientcredentials"
)

// These three structs define the data response we are expecting
// to recieve in the response to the GetList call.
// The json body will looks something like:
// {
//     "data": {
//         "attributes": {
//             "title": "The list title",
//		   }
//     }
// }
type GetListResponse struct {
	Data GetListResponseData `json:"data"`
}

type GetListResponseData struct {
	Attr GetListResponseDataAttributes `json:"attributes"`
}

type GetListResponseDataAttributes struct {
	Title string `json:"title"`
}

// main is the point of entry for this tool
func main() {
	// The list ID and tenant shortcode values are retrieved from the
	// parameters passed when running this tool
	var listID, tenant string
	flag.StringVar(&listID, "id", "", "list ID")
	flag.StringVar(&tenant, "tenant", "", "the tenant shortcode")
	flag.Parse()

	// ID and secret come from local environment variables
	// Set in your terminal eg:
	//      export ACTIVE_TALIS_PERSONA_ID="my-id"
	personaID := os.Getenv("ACTIVE_TALIS_PERSONA_ID")
	personaSecret := os.Getenv("ACTIVE_TALIS_PERSONA_SECRET")
	personaURL := "https://users.talis.com/oauth/tokens"

	// Get an http client which can make GET calls to
	// a URL.  This client will be able to automatically
	// add in authorization tokens based on your credentials,
	// and will automatically refresh the token if it expires
	c := getClient(personaID, personaSecret, personaURL)

	// Build the URL for the API Get Lists endpoint, based on the list ID and tenant shortcode
	url := fmt.Sprintf("https://rl.talis.com/3/%s/lists/%s", tenant, listID)

	l, err := getListFromURL(url, c)
	if err != nil {
		// If something went wrong getting the list data
		// then quit out and print the error to the screen
		log.Fatal(fmt.Errorf("faild to get the list data: %w", err))
	}

	// Now we have the data correctly unmarshalled, we can print out the title
	title := l.Data.Attr.Title
	fmt.Println(title)
}

// getListFromURL will take a url and an http client and will return the list response
func getListFromURL(url string, c *http.Client) (*GetListResponse, error) {
	// Call the built URL to get the list details
	resp, err := c.Get(url)
	if err != nil {
		// If something went wrong retrieving the list data
		// then quit out and print the error to the screen
		return nil, fmt.Errorf("get call to retrieve the list failed: %w", err)
	}
	// We must make sure the response body is closed when this
	// function ends
	defer resp.Body.Close()

	// Check that we got a 200 OK status response from the API
	if resp.StatusCode != http.StatusOK {
		// If we did not get an "OK" status code in the response
		// then quit out and print the error to the screen
		return nil, fmt.Errorf("status code was: %s", resp.Status)
	}

	b, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		// If something went wrong reading the body
		// then quit out and print the error to the screen
		return nil, fmt.Errorf("failed reading the body response: %w", err)
	}

	// Unmarshal (convert) the json response into a GetListResponse struct
	// that we can easily work with in golang
	// Note that the variable is a pointer to a struct
	// (ie. it points to the memory location of the struct rather than
	// being the struct itself)
	l := &GetListResponse{}
	if err = json.Unmarshal(b, l); err != nil {
		// If something went wrong unmarshalling
		// then quit out and print the error to the screen
		return nil, fmt.Errorf("failed to unmarshal json response: %w", err)
	}

	// All is good, so return list with a 'nil' error
	return l, nil
}

// getClient will return an http client with the oauth
// credentials baked in
func getClient(ID, secret, url string) *http.Client {
	c := clientcredentials.Config{
		ClientID:     ID,
		ClientSecret: secret,
		TokenURL:     url,
	}
	return c.Client(context.Background())
}
