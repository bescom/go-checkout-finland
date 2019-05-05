package checkout

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"net/http/httputil"
	"strings"
	"time"
)

// constants
const (
	defaultAPIURL  = "https://api.checkout.fi"
	defaultTimeout = 10 // seconds
)

// Client represents an API client
type Client struct {
	merchantID string
	secretKey  string
	apiURL     string
	httpClient *http.Client
}

// New creates ands returns a new client configured with the specified merchant account
func New(merchantID, secretKey string) *Client {

	client := Client{}

	client.merchantID = merchantID
	client.secretKey = secretKey
	client.apiURL = defaultAPIURL

	client.httpClient = &http.Client{}
	client.httpClient.Timeout = time.Second * defaultTimeout

	return &client
}

// CreatePayment performs a POST request ...
func (c *Client) CreatePayment(p *Payment) (*PaymentResponse, error) {

	var res PaymentResponse

	payment, err := json.Marshal(p)

	r, err := c.executeRequest(http.MethodPost, "/payments", payment)

	if err != nil {
		return nil, errors.New("executing request failed")
	}

	err = json.Unmarshal(r, &res)

	if err != nil {
		return nil, errors.New("unmarshalling payment response failed")
	}

	return &res, err
}

// Performs the specified HTTP request and returns the response through handleResponse()
func (c *Client) executeRequest(method string, endpoint string, body []byte) ([]byte, error) {

	apiURL := fmt.Sprintf("%s%s", c.apiURL, endpoint)

	request, err := http.NewRequest(method, apiURL, bytes.NewBuffer(body))

	if err != nil {
		return nil, err
	}

	t := time.Now().Format(time.RFC3339)

	request.Header.Add("Accept", "application/json")
	request.Header.Add("Content-Type", "application/json; charset=utf-8")
	request.Header.Add("checkout-account", c.merchantID)
	request.Header.Add("checkout-algorithm", "sha256")
	request.Header.Add("checkout-method", method)
	request.Header.Add("checkout-nonce", "12345")
	request.Header.Add("checkout-timestamp", t)
	request.Header.Add("signature", getHeaders(body, t))

	dump, err := httputil.DumpRequestOut(request, true)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println("-- dump --")
	fmt.Printf("%q", dump)

	response, err := c.httpClient.Do(request)

	if err != nil {
		return nil, err
	}

	return handleResponse(response)
}

// Parses the response and returns either the response body or an error
func handleResponse(response *http.Response) ([]byte, error) {
	defer response.Body.Close()

	fmt.Println("-- response.StatusCode : ", response.StatusCode)

	// Return an error on unsuccessful requests

	if response.StatusCode < 200 || response.StatusCode > 201 {
		errorBody, _ := ioutil.ReadAll(response.Body)

		fmt.Println("-- error in handleResponse --")
		fmt.Println(response.StatusCode)
		fmt.Println(string(errorBody))

		return nil, errors.New("xxx")
	}

	responseBody, err := ioutil.ReadAll(response.Body)

	fmt.Println()
	fmt.Println(string(responseBody))
	fmt.Println()

	return responseBody, err
}

func getHeaders(content []byte, time string) string {
	var sb strings.Builder

	sb.WriteString("checkout-account:375917\n")
	sb.WriteString("checkout-algorithm:sha256\n")
	sb.WriteString("checkout-method:POST\n")
	sb.WriteString("checkout-nonce:12345\n")
	sb.WriteString("checkout-timestamp:" + time + "\n")
	sb.WriteString(string(content))

	h, _ := GenerateHMAC("SAIPPUAKAUPPIAS", sb.String())

	return h
}
