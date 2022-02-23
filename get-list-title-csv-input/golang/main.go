package main

import (
	"context"
	"encoding/csv"
	"flag"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"

	"github.com/Talis-Aspire-Developer-Community/examples/get-list-title-csv-input/golang/listapi"
	"golang.org/x/oauth2/clientcredentials"
)

func main() {
	// The list ID and tenant shortcode values are retrieved from the
	// parameters passed when running this tool
	var filename, tenant string
	flag.StringVar(&filename, "file", "", "input csv file")
	flag.StringVar(&tenant, "tenant", "", "the tenant shortcode")
	flag.Parse()

	// ID and secret come from local environment variables
	// Set in your terminal eg:
	//      export ACTIVE_TALIS_PERSONA_ID="my-id"
	personaID := os.Getenv("ACTIVE_TALIS_PERSONA_ID")
	personaSecret := os.Getenv("ACTIVE_TALIS_PERSONA_SECRET")
	personaURL := "https://users.talis.com/oauth/tokens"

	// Get an http client which can make GET calls to
	// a URL.  This client will manage authorization tokens,
	// and will automatically refresh the token if it expires
	c := getClient(personaID, personaSecret, personaURL)

	// Open the file ready for reading
	file, err := os.Open(filename)
	if err != nil {
		log.Fatal(fmt.Errorf("error opening csv file: %w", err))
	}
	defer file.Close()

	// Create a new CSV reader
	r := csv.NewReader(file)
	r.FieldsPerRecord = 1
	r.TrimLeadingSpace = true
	// Iterate over the csv rows until we reach the end of file or error
	for {
		// Read each record from csv
		record, err := r.Read()
		if err == io.EOF {
			break
		}
		if err != nil {
			log.Fatal(err)
		}
		listID := record[0]

		// Get the list from the URL
		listResp, err := listapi.Get(tenant, listID, c)
		if err != nil {
			log.Fatal(fmt.Errorf("faild to get the list data: %w", err))
		}

		// Now we have the data in a struct, we can print out the title
		title := listResp.Data.Attr.Title
		fmt.Println(title)
	}
}

func getClient(ID, secret, url string) *http.Client {
	c := clientcredentials.Config{
		ClientID:     ID,
		ClientSecret: secret,
		TokenURL:     url,
	}
	return c.Client(context.Background())
}
