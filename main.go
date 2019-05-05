package main

import (
	"encoding/json"
	"fmt"
	"math/rand"
	"strconv"
	"time"

	"github.com/bescom/go-checkout-finland/checkout"
)

func main() {

	//h, _ := checkout.GenerateHMAC("secret", "data")
	//fmt.Println("hmac: " + h)

	//headers, _ := checkout.GenerateHMAC("SAIPPUAKAUPPIAS", checkout.GetHeaders())
	//fmt.Println("headers: " + headers)

	c := checkout.New("375917", "SAIPPUAKAUPPIAS")

	/*
		type Payment struct {
		Stamp            string      `json:"stamp"`     // m
		Reference        string      `json:"reference"` // m
		Amount           int         `json:"amount"`    // m
		Currency         string      `json:"currency"`  // m, alpha3
		Language         string      `json:"language"`  // m, alpha2
		Items            []item      `json:"items"`     // m
		Customer         customer    `json:"customer"`  // m
		DeliveryAddress  address     `json:"deliveryAddress"`
		InvoicingAddress address     `json:"invoicingAddress"`
		RedirectUrls     callbackUrl `json:"redirectUrls"` // m
		CallbackUrls     callbackUrl `json:"callbackUrls"`
	}*/

	/*
		p := new(checkout.Payment)
		p.Stamp = "xxx"
		p.Reference = "ref"
		p.Amount = 100
		p.Currency = "EUR"
	*/

	rand.Seed(time.Now().UTC().UnixNano())
	s := strconv.Itoa(rand.Int())
	fmt.Println(s)

	p := checkout.Payment{Stamp: s, Reference: "ref", Amount: 100, Currency: "EUR", Language: "FI"}

	i := checkout.Item{Description: "tuotteen kuvaus", UnitPrice: 101, Units: 1, VatPercentage: 24, ProductCode: "pcode", DeliveryDate: "2019-01-15"}
	p.AddItem(i)

	p.Customer = checkout.Customer{FirstName: "Teppo", LastName: "Testaaja", Email: "test@test.com"}
	p.InvoicingAddress = checkout.Address{StreetAddress: "Testikatu 1 A 2", PostalCode: "12345", City: "Testilä", Country: "Suomi"}
	p.InvoicingAddress = checkout.Address{StreetAddress: "Testikatu 3 A 4", PostalCode: "12345", City: "Testilä", Country: "Suomi"}
	p.RedirectUrls = checkout.CallbackURL{Success: "https://tsadasdest987.com/success", Cancel: "https://tsadasdest987.com/cancel"}

	res1B, _ := json.Marshal(p)
	fmt.Println(string(res1B))

	_, err := c.CreatePayment(&p)

	if err != nil {
		panic(err)

		//fmt.Println(pr.TransactionID)
		//fmt.Println(pr.Href)

	}

}
