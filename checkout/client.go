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
func (c *Client) CreatePayment(payment *Payment) (*PaymentResponse, error) {
	p, _ := json.Marshal(payment)

	r, err := c.executeRequest(http.MethodPost, "/payments", p)

	var pr PaymentResponse
	err1 := json.Unmarshal(r, &pr)

	if err1 != nil {
		fmt.Println("-- err1 : ", err1)
	}

	return &pr, err
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

	fmt.Println("-- response.StatusCode : ", response.StatusCode)

	// Return an error on unsuccessful requests
	if response.StatusCode < 200 || response.StatusCode > 299 {
		errorBody, _ := ioutil.ReadAll(response.Body)

		fmt.Println("-- error in handleResponse --")
		fmt.Println(response.StatusCode)
		fmt.Println(string(errorBody))

		//return nil, &Error{response.StatusCode, response.Status, errorBody}
	}

	responseBody, err := ioutil.ReadAll(response.Body)

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
