package listapi

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
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
type GetResponse struct {
	Data GetResponseData `json:"data"`
}

type GetResponseData struct {
	Attr GetResponseDataAttributes `json:"attributes"`
}

type GetResponseDataAttributes struct {
	Title string `json:"title"`
}

type Client struct {
	BaseURL    string
	TenantCode string
	Client     *http.Client
}

// Get the list api response
func (a *Client) Get(listID string) (*GetResponse, error) {
	// Build the URL we want to call
	url := fmt.Sprintf("%s/3/%s/lists/%s", a.BaseURL, a.TenantCode, listID)

	// Call the built URL to get the list details
	resp, err := a.Client.Get(url)
	if err != nil {
		return nil, fmt.Errorf("get call to retrieve the list failed: %w", err)
	}
	defer resp.Body.Close()

	// Check that we got a 200 OK status response from the API
	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("status code was: %s", resp.Status)
	}

	// Read the body into a byte array
	b, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed reading the body response: %w", err)
	}

	// Unmarshal (convert) the byte array into a GetListResponse struct
	// that we can easily work with in golang
	// Note that the variable is a pointer to a struct
	// (ie. it points to the memory location of the struct rather than
	// being the struct itself)
	l := &GetResponse{}
	if err = json.Unmarshal(b, l); err != nil {
		return nil, fmt.Errorf("failed to unmarshal json response: %w", err)
	}
	return l, nil
}
