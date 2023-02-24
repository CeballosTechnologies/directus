package directus

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"time"
)

type Date time.Time

func (d *Date) MarshalJSON() ([]byte, error) {
	return json.Marshal(time.Time(*d).Format("2006-02-01"))
}

func (d *Date) UnmarshalJSON(data []byte) error {
	data = data[1:]
	data = data[:len(data)-1]

	t, err := time.Parse("2006-02-01", string(data))
	if err != nil {
		return err
	}
	*d = Date(t)
	return nil
}

// Client interface for interacting with Client API
type Client struct {
	accessToken string
	httpClient  http.Client
	url         url.URL
}

type IdObject struct {
	Id int `json:"id"`
}

type ICollectionItem interface {
	GetCollectionFields() string
	GetCollectionName() string
	GetId() int
	MarshalJSON() ([]byte, error)
	SetId(int)
}

type CollectionItem struct {
	Id          int        `json:"id,omitempty"`
	DateCreated *time.Time `json:"date_created,omitempty"`
	DateUpdated *time.Time `json:"date_updated,omitempty"`
}

type ISingletonItem interface {
	GetCollectionFields() string
	GetCollectionName() string
}

// NewClient returns a new initialized client interface for interacting with Directus API.
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

// CreateItem adds a new item to directus collection.
func (dc *Client) CreateItem(item ICollectionItem) (ICollectionItem, error) {
	u := dc.url
	u.Path = fmt.Sprintf("/items/%s", item.GetCollectionName())

	queryParams := u.Query()
	queryParams.Add("fields", item.GetCollectionFields())
	u.RawQuery = queryParams.Encode()

	dataBytes, err := SerializeItem(item) // json.Marshal(item)
	if err != nil {
		return nil, err
	}

	req, err := http.NewRequest("POST", u.String(), bytes.NewBuffer(dataBytes))
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

func (dc *Client) FindItem(item ICollectionItem, filter string) (ICollectionItem, error) {
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
	u := dc.url

	u.Path = fmt.Sprintf("/items/%s", item.GetCollectionName())

	queryParams := u.Query()
	queryParams.Add("fields", "id")
	if filter != "" {
		queryParams.Add("filter", filter)
	}

	u.RawQuery = queryParams.Encode()

	req, err := http.NewRequest("GET", u.String(), nil)
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

// FindItems queries items from directus collection.
func (dc *Client) FindItems(item ICollectionItem, filter string) ([]byte, error) {
	u := dc.url

	u.Path = fmt.Sprintf("/items/%s", item.GetCollectionName())

	queryParams := u.Query()
	if item.GetCollectionFields() != "" {
		queryParams.Add("fields", item.GetCollectionFields())
	}
	if filter != "" {
		queryParams.Add("filter", filter)
	}

	u.RawQuery = queryParams.Encode()

	req, err := http.NewRequest("GET", u.String(), nil)
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

func (dc *Client) GetItem(item ICollectionItem) (ICollectionItem, error) {
	u := dc.url

	u.Path = fmt.Sprintf("/items/%s/%d", item.GetCollectionName(), item.GetId())

	queryParams := u.Query()
	if item.GetCollectionFields() != "" {
		queryParams.Add("fields", item.GetCollectionFields())
	}

	u.RawQuery = queryParams.Encode()

	req, err := http.NewRequest("GET", u.String(), nil)
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

// GetPath get raw path with query parameters and returns raw response data
func (dc *Client) GetPath(path string, queryParams url.Values) (response []byte, err error) {
	u := dc.url
	u.Path = path

	if queryParams != nil {
		u.RawQuery = queryParams.Encode()
	}

	req, err := http.NewRequest("GET", u.String(), nil)
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

func (dc *Client) GetSingleton(item ISingletonItem) (ISingletonItem, error) {
	u := dc.url

	u.Path = fmt.Sprintf("/items/%s", item.GetCollectionName())

	queryParams := u.Query()
	queryParams.Add("fields", item.GetCollectionFields())
	u.RawQuery = queryParams.Encode()

	req, err := http.NewRequest("GET", u.String(), nil)
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

	singleton := item
	err = json.Unmarshal(body, &singleton)
	if err != nil {
		return nil, err
	}

	return singleton, nil
}

// UpdateItem updates an existing item in directus collection.
func (dc *Client) UpdateItem(item ICollectionItem) (ICollectionItem, error) {
	u := dc.url
	u.Path = fmt.Sprintf("/items/%s/%d", item.GetCollectionName(), item.GetId())
	// item.SetId(0)

	queryParams := u.Query()
	queryParams.Add("fields", item.GetCollectionFields())
	u.RawQuery = queryParams.Encode()

	dataBytes, err := SerializeItem(item) // json.Marshal(item)
	if err != nil {
		return nil, err
	}

	req, err := http.NewRequest("PATCH", u.String(), bytes.NewBuffer(dataBytes))
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

// UpdateSingleton updates an existing singleton in directus collection.
func (dc *Client) UpdateSingleton(item ISingletonItem) (ISingletonItem, error) {
	u := dc.url
	u.Path = fmt.Sprintf("/items/%s", item.GetCollectionName())

	queryParams := u.Query()
	queryParams.Add("fields", item.GetCollectionFields())
	u.RawQuery = queryParams.Encode()

	dataBytes, err := json.Marshal(item)
	if err != nil {
		return nil, err
	}

	req, err := http.NewRequest("PATCH", u.String(), bytes.NewBuffer(dataBytes))
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

// UpsertItem upserts an existing item in directus collection.
func (dc *Client) UpsertItem(item ICollectionItem) (ICollectionItem, error) {
	if item.GetId() == 0 {
		return dc.CreateItem(item)
	} else {
		return dc.UpdateItem(item)
	}
}

// Executes HTTP request with Directus API
func (dc *Client) sendRequest(request *http.Request, maxRetries int, retryCounter int) ([]byte, error) {
	var data []byte

	request.Header.Set("Authorization", fmt.Sprintf("Bearer %s", dc.accessToken))
	request.Header.Set("Content-Type", "application/json")

	resp, err := dc.httpClient.Do(request)
	if err != nil {
		fmt.Println("error sending request to directus API")
		fmt.Println(err.Error())
		return data, err
	}

	if resp.StatusCode == 500 {
		if retryCounter <= maxRetries {
			retryResponse, err := dc.sendRequest(request, maxRetries, retryCounter+1)
			if err != nil {
				return nil, err
			}
			return retryResponse, nil
		} else {
			return nil, fmt.Errorf("internal server error: directus api max attempts reached")
		}
		//return dc.sendRequest(request, maxRetries, retryCounter+1)
	} else if resp.StatusCode >= 300 {
		data, err := io.ReadAll(resp.Body)
		if err != nil {
			return nil, err
		}
		return nil, fmt.Errorf(string(data))
	}

	defer func(Body io.ReadCloser) {
		_ = Body.Close()
	}(resp.Body)

	return io.ReadAll(resp.Body)
}

func SerializeItem(item any) ([]byte, error) {
	type Alias any
	a := (*Alias)(&item)
	return json.Marshal(a)
}

func (dc *Client) GetCurrentUser(token string) (User, error) {
	var user User

	u := dc.url
	u.Path = "/users/me"

	req, err := http.NewRequest("GET", u.String(), nil)
	if err != nil {
		return user, err
	}
	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", token))

	resp, err := dc.httpClient.Do(req)
	if err != nil {
		fmt.Println("error sending request to directus API")
		fmt.Println(err.Error())
		return user, err
	}

	defer func(Body io.ReadCloser) {
		_ = Body.Close()
	}(resp.Body)

	data, err := io.ReadAll(resp.Body)
	if err != nil {
		return user, err
	}

	if resp.StatusCode >= 300 {
		return user, fmt.Errorf(string(data))
	}

	// Remove data wrapper
	data = data[8:]
	data = data[:len(data)-1]

	err = json.Unmarshal(data, &user)
	if err != nil {
		return User{}, err
	}

	return user, nil
}

func (dc *Client) GetRole(roleId string) (Role, error) {
	var role Role

	u := dc.url

	u.Path = fmt.Sprintf("/roles/%s", roleId)

	req, err := http.NewRequest("GET", u.String(), nil)
	if err != nil {
		return role, err
	}

	body, err := dc.sendRequest(req, 5, 0)
	if err != nil {
		return role, err
	}

	// Remove data wrapper
	body = body[8:]
	body = body[:len(body)-1]

	err = json.Unmarshal(body, &role)

	return role, err
}
