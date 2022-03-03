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
	"path/filepath"

	"github.com/Talis-Aspire-Developer-Community/examples/get-list-title-csv-input/golang/listapi"
	"golang.org/x/oauth2/clientcredentials"
)

func main() {
	// The list ID and tenant shortcode values are retrieved from the
	// parameters passed when running this tool
	var infilepath, outfilepath, tenant string
	flag.StringVar(&infilepath, "infile", "", "input csv file")
	flag.StringVar(&outfilepath, "outfile", "", "input csv file")
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

	if err := ensureOutFilepath(outfilepath); err != nil {
		log.Fatal(fmt.Errorf("could not ensure out filepath: %w", err))
	}

	// Open the 'out' file ready for writing
	out, err := os.Create(outfilepath)
	if err != nil {
		log.Fatal(fmt.Errorf("error opening csv file: %w", err))
	}
	defer out.Close()

	// Open the 'in' file ready for reading
	in, err := os.Open(infilepath)
	if err != nil {
		log.Fatal(fmt.Errorf("error opening csv file: %w", err))
	}
	defer in.Close()

	// Iterate over the csv rows until we reach the end of file or error
	if err := iterateAndWrite(in, out, tenant, c); err != nil {
		log.Fatal(err)
	}
}

func iterateAndWrite(in io.Reader, out io.Writer, tenant string, c *http.Client) error {
	// Create a new CSV reader
	r := csv.NewReader(in)
	r.FieldsPerRecord = 1
	r.TrimLeadingSpace = true
	// Create a new CSV writer
	w := csv.NewWriter(out)

	for {
		// Read each record from csv
		record, err := r.Read()
		if err == io.EOF {
			break
		}
		if err != nil {
			return fmt.Errorf("couldnt read csv row from input: %w", err)
		}
		listID := record[0]

		a := &listapi.Client{
			BaseURL:    "http://tenant-test.ac.uk",
			TenantCode: tenant,
			Client:     c,
		}
		// Get the list from the URL
		listResp, err := a.Get(listID)
		if err != nil {
			return fmt.Errorf("failed to get the list data: %w", err)
		}

		// Now we have the data in a struct, we can write the title
		// to the output file
		title := listResp.Data.Attr.Title
		outRow := []string{title}
		w.Write(outRow)
		w.Flush()
	}
	return nil
}

func ensureOutFilepath(f string) error {
	absPath, err := filepath.Abs(f)
	if err != nil {
		return err
	}
	d := filepath.Dir(absPath)
	fmt.Println(d)
	return os.MkdirAll(d, os.ModePerm)
}

func getClient(ID, secret, url string) *http.Client {
	c := clientcredentials.Config{
		ClientID:     ID,
		ClientSecret: secret,
		TokenURL:     url,
	}
	return c.Client(context.Background())
}
