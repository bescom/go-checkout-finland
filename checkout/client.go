package checkout

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"strings"
	"time"
)

// constants
const (
	defaultAPIBaseURL = "https://api.checkout.fi"
	defaultTimeout    = 10 // seconds
)

// Payment represents ...
type Payment struct {
	Stamp            string      `json:"stamp"`     // m
	Reference        string      `json:"reference"` // m
	Amount           int         `json:"amount"`    // m
	Currency         string      `json:"currency"`  // m, EUR
	Language         string      `json:"language"`  // m, FI/SV/EN
	Items            []Item      `json:"items"`     // m
	Customer         Customer    `json:"customer"`  // m
	DeliveryAddress  Address     `json:"deliveryAddress"`
	InvoicingAddress Address     `json:"invoicingAddress"`
	RedirectUrls     CallbackURL `json:"redirectUrls"` // m
	CallbackUrls     CallbackURL `json:"callbackUrls"`
}

// Item represents ...
type Item struct {
	UnitPrice     int    `json:"unitPrice"`     // m
	Units         int    `json:"units"`         // m
	VatPercentage int    `json:"vatPercentage"` // m
	ProductCode   string `json:"productCode"`   // m
	DeliveryDate  string `json:"deliveryDate"`  // m
	Description   string `json:"description"`
	Category      string `json:"category"`
	Stamp         string `json:"stamp"`
	Reference     string `json:"reference"` // m for shop-in-shop
	//Merchant      string     `json:"merchant"` // only for shop-in-shop
	//Commission    Commission `json:"commission"` // only for shop-in-shop
}

// Customer represents ...
type Customer struct {
	Email     string `json:"email"` // m
	FirstName string `json:"firstName"`
	LastName  string `json:"lastName"`
	Phone     string `json:"phone"`
	VatID     string `json:"vatId"`
}

// Address represents ...
type Address struct {
	StreetAddress string `json:"streetAddress"` // m
	PostalCode    string `json:"postalCode"`    // m
	City          string `json:"city"`          // m
	County        string `json:"county"`
	Country       string `json:"country"` // m
}

// CallbackURL represents ...
type CallbackURL struct {
	Success string `json:"success"` // m
	Cancel  string `json:"cancel"`  // m
}

// Commission represents ...
type Commission struct {
	Merchant string `json:"merchant"` // m
	Amount   int    `json:"amount"`   // m
}

// Client represents an API client
type Client struct {
	merchantID string
	secretKey  string
	httpClient *http.Client

	apiBaseURL string
}

// New creates ands returns a new client configured with the specified merchant data
func New(merchantID, secretKey string) *Client {
	client := Client{}

	client.merchantID = merchantID
	client.secretKey = secretKey
	client.httpClient = &http.Client{}
	client.httpClient.Timeout = time.Second * defaultTimeout

	client.apiBaseURL = defaultAPIBaseURL

	return &client
}

// CreatePayment performs a POST request ...
func (c *Client) CreatePayment(p *Payment) ([]byte, error) {
	json, _ := json.Marshal(p)
	return c.executeRequest(http.MethodPost, "/payments", json)
}

// Performs the specified HTTP request and returns the response through handleResponse()
func (c *Client) executeRequest(method string, endpoint string, content []byte) ([]byte, error) {

	apiURL := fmt.Sprintf("%s%s", c.apiBaseURL, endpoint)

	fmt.Println("-- apiURL : ", apiURL)

	request, err := http.NewRequest(method, apiURL, bytes.NewBuffer(content))

	if err != nil {
		return nil, err
	}

	request.Header.Add("Accept", "application/json")
	request.Header.Add("Content-Type", "application/json; charset=utf-8")

	request.Header.Add("checkout-account", c.merchantID)
	request.Header.Add("checkout-algorithm", "sha256")
	request.Header.Add("checkout-method", "POST")
	request.Header.Add("checkout-nonce", "12345")
	request.Header.Add("checkout-timestamp", "2019-01-08T11:19:25.950Z")
	request.Header.Add("signature", getHeaders(content))

	/*
		dump, err := httputil.DumpRequestOut(request, true)
		if err != nil {
			log.Fatal(err)
		}

		fmt.Println("-- dump --")
		fmt.Printf("%q", dump)
	*/

	response, err := c.httpClient.Do(request)

	if err != nil {
		log.Fatal(err)

		return nil, err
	}

	return handleResponse(response)
}

// Parses the response and returns either the response body or an error
func handleResponse(response *http.Response) ([]byte, error) {
	defer response.Body.Close()

	// Return an error on unsuccessful requests
	if response.StatusCode < 200 || response.StatusCode > 299 {
		errorBody, _ := ioutil.ReadAll(response.Body)

		fmt.Println("-- error in handleResponse --")
		fmt.Println(response.StatusCode)
		fmt.Println(string(errorBody))

		//return nil, &Error{response.StatusCode, response.Status, errorBody}
	}

	responseBody, err := ioutil.ReadAll(response.Body)

	fmt.Println(string(responseBody))

	return responseBody, err
}

func getHeaders(content []byte) string {
	var sb strings.Builder

	sb.WriteString("checkout-account:375917\n")
	sb.WriteString("checkout-algorithm:sha256\n")
	sb.WriteString("checkout-method:POST\n")
	sb.WriteString("checkout-nonce:12345\n")
	sb.WriteString("checkout-timestamp:2019-01-08T11:19:25.950Z\n")
	sb.WriteString(string(content))

	h, _ := GenerateHMAC("SAIPPUAKAUPPIAS", sb.String())
	return h
}
