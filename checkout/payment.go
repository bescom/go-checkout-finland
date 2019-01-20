package checkout

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
	UnitPrice     int         `json:"unitPrice"`     // m
	Units         int         `json:"units"`         // m
	VatPercentage int         `json:"vatPercentage"` // m
	ProductCode   string      `json:"productCode"`   // m
	DeliveryDate  string      `json:"deliveryDate"`  // m
	Description   string      `json:"description"`
	Category      string      `json:"category"`
	Stamp         string      `json:"stamp"`
	Reference     string      `json:"reference"`            // m for shop-in-shop
	Merchant      string      `json:"merchant,omitempty"`   // only for shop-in-shop
	Commission    *Commission `json:"commission,omitempty"` // only for shop-in-shop, pointer until omitempty works with structs
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

// PaymentResponse represents ...
type PaymentResponse struct {
	TransactionID string     `json:"transactionid"`
	Href          string     `json:"href"`
	Providers     []Provider `json:"providers"`
}

// Provider represents ...
type Provider struct {
	URL        string      `json:"url"`
	Icon       string      `json:"icon"`
	Svg        string      `json:"svg"`
	Name       string      `json:"name"`
	Group      string      `json:"group"`
	ID         string      `json:"id"`
	Parameters []Parameter `json:"parameters"`
}

// Parameter represents ...
type Parameter struct {
	Name  string `json:"name"`
	Value string `json:"value"`
}
