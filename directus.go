package directus

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"time"
)

// Client interface for interacting with Client API
type Client struct {
	accessToken string
	httpClient  http.Client
	url         url.URL
}

type ICollectionItem interface {
	GetCollectionFields() string
	GetCollectionName() string
	GetId() int
	// SetFieldsFromCollectionItem([]byte) error
	SetId(int)
}

type CollectionItem struct {
	Id          int        `json:"id,omitempty"`
	DateCreated *time.Time `json:"date_created,omitempty"`
	DateUpdated *time.Time `json:"date_updated,omitempty"`
}

// Returns a new initialized client interface for interacting with Directus API.
//   - baseUrl: directus API location.
//   - accessToken represents the directus user used for integrations.
func NewClient(baseUrl string, accessToken string) (Client, error) {
	dc := new(Client)

	u, err := url.Parse(baseUrl)
	if err != nil {
		return *dc, err
	}

	dc.accessToken = accessToken
	dc.httpClient = http.Client{}
	dc.url = *u

	return *dc, nil
}

// Adds a new item to directus collection.
func (dc *Client) CreateItem(item ICollectionItem) (ICollectionItem, error) {
	url := dc.url
	url.Path = fmt.Sprintf("/items/%s", item.GetCollectionName())

	queryParams := url.Query()
	queryParams.Add("fields", item.GetCollectionFields())
	url.RawQuery = queryParams.Encode()

	dataBytes, err := json.Marshal(item)
	if err != nil {
		return nil, err
	}

	req, err := http.NewRequest("POST", url.String(), bytes.NewBuffer(dataBytes))
	if err != nil {
		return nil, err
	}

	body, err := dc.sendRequest(req, 5, 0)
	if err != nil {
		return nil, err
	}

	// Remove data wrapper
	body = body[8:]
	body = body[:len(body)-1]

	err = json.Unmarshal(body, &item)

	// t := reflect.TypeOf(item)
	// UnmarshalCollectionItem[t](body)

	// err = item.SetFieldsFromCollectionItem(body)
	// if err != nil {
	// 	return nil, err
	// }

	return item, err
}

func (dc *Client) FindItem(item ICollectionItem, filter string) (any, error) {
	data, err := dc.FindItems(item, filter)
	if err != nil {
		return nil, err
	}

	var items []any
	err = json.Unmarshal(data, &items)
	if err != nil {
		return nil, err
	}

	if len(items) == 0 {
		return item, nil
	}

	itemBytes, err := json.Marshal(items[0])
	if err != nil {
		return nil, err
	}

	if err = json.Unmarshal(itemBytes, &item); err != nil {
		return nil, err
	}

	return item, nil
}

func (dc *Client) FindItemId(item ICollectionItem, filter string) (int, error) {
	url := dc.url

	url.Path = fmt.Sprintf("/items/%s", item.GetCollectionName())

	queryParams := url.Query()
	queryParams.Add("fields", "id")
	if filter != "" {
		queryParams.Add("filter", filter)
	}

	url.RawQuery = queryParams.Encode()

	req, err := http.NewRequest("GET", url.String(), nil)
	if err != nil {
		return 0, err
	}

	body, err := dc.sendRequest(req, 5, 0)
	if err != nil {
		return 0, err
	}

	// Remove data wrapper
	body = body[8:]
	body = body[:len(body)-1]

	var items []CollectionItem
	err = json.Unmarshal(body, &items)
	if err != nil {
		return 0, err
	}

	if len(items) == 0 {
		return 0, nil
	}

	return items[0].Id, nil
}

// Queries items from directus collection.
func (dc *Client) FindItems(item ICollectionItem, filter string) ([]byte, error) {
	url := dc.url

	url.Path = fmt.Sprintf("/items/%s", item.GetCollectionName())

	queryParams := url.Query()
	if item.GetCollectionFields() != "" {
		queryParams.Add("fields", item.GetCollectionFields())
	}
	if filter != "" {
		queryParams.Add("filter", filter)
	}

	url.RawQuery = queryParams.Encode()

	req, err := http.NewRequest("GET", url.String(), nil)
	if err != nil {
		return nil, err
	}

	body, err := dc.sendRequest(req, 5, 0)
	if err != nil {
		return nil, err
	}

	// Remove data wrapper
	body = body[8:]
	body = body[:len(body)-1]

	return body, nil
}

func (dc *Client) GetItem(item ICollectionItem) (any, error) {
	url := dc.url

	url.Path = fmt.Sprintf("/items/%s/%d", item.GetCollectionName(), item.GetId())

	queryParams := url.Query()
	if item.GetCollectionFields() != "" {
		queryParams.Add("fields", item.GetCollectionFields())
	}

	url.RawQuery = queryParams.Encode()

	req, err := http.NewRequest("GET", url.String(), nil)
	if err != nil {
		return nil, err
	}

	body, err := dc.sendRequest(req, 5, 0)
	if err != nil {
		return nil, err
	}

	// Remove data wrapper
	body = body[8:]
	body = body[:len(body)-1]

	// item.SetFieldsFromCollectionItem(body)

	err = json.Unmarshal(body, &item)

	return item, err
}

// Get raw path with query parameters and returns raw response data
func (dc *Client) GetPath(path string, queryParams url.Values) (response []byte, err error) {
	url := dc.url
	url.Path = path

	if queryParams != nil {
		url.RawQuery = queryParams.Encode()
	}

	req, err := http.NewRequest("GET", url.String(), nil)
	if err != nil {
		return
	}

	body, err := dc.sendRequest(req, 5, 0)
	if err != nil {
		return
	}

	// Remove data wrapper
	body = body[8:]
	body = body[:len(body)-1]

	response = body

	return
}

// Updates an existing item in directus collection.
func (dc *Client) UpdateItem(item ICollectionItem) (any, error) {
	url := dc.url
	url.Path = fmt.Sprintf("/items/%s/%d", item.GetCollectionName(), item.GetId())
	item.SetId(0)

	queryParams := url.Query()
	queryParams.Add("fields", item.GetCollectionFields())
	url.RawQuery = queryParams.Encode()

	dataBytes, err := json.Marshal(item)
	if err != nil {
		return nil, err
	}

	req, err := http.NewRequest("PATCH", url.String(), bytes.NewBuffer(dataBytes))
	if err != nil {
		return nil, err
	}

	body, err := dc.sendRequest(req, 5, 0)
	if err != nil {
		return nil, err
	}

	// Remove data wrapper
	body = body[8:]
	body = body[:len(body)-1]

	// item.SetFieldsFromCollectionItem(body)

	// fmt.Println(string(body))

	err = json.Unmarshal(body, &item)

	return item, err
}

// Upserts an existing item in directus collection.
func (dc *Client) UpsertItem(item ICollectionItem) (any, error) {
	if item.GetId() == 0 {
		return dc.CreateItem(item)
	} else {
		return dc.UpdateItem(item)
	}
}

// Executes HTTP request with Directus API
func (dc *Client) sendRequest(request *http.Request, maxRetries int, retryCounter int) ([]byte, error) {
	var data = []byte{}

	if maxRetries == retryCounter {
		return data, fmt.Errorf("directus api max attempts reached")
	}

	request.Header.Set("Authorization", fmt.Sprintf("Bearer %s", dc.accessToken))
	request.Header.Set("Content-Type", "application/json")

	resp, err := dc.httpClient.Do(request)
	if err != nil {
		fmt.Println(("error sending request to directus API"))
		fmt.Println(err.Error())
		return data, err
	}

	if resp.StatusCode != 200 {
		fmt.Printf("Directus response status: %d\n", resp.StatusCode)

		if resp.StatusCode == 500 {
			return dc.sendRequest(request, maxRetries, (retryCounter + 1))
		}
	}

	defer resp.Body.Close()

	return ioutil.ReadAll(resp.Body)
}
