package files_sdk

import (
	"encoding/json"
	"time"

	lib "github.com/Files-com/files-sdk-go/v2/lib"
)

type Account struct {
	Name             string     `json:"name,omitempty" path:"name,omitempty" url:"name,omitempty"`
	Address          string     `json:"address,omitempty" path:"address,omitempty" url:"address,omitempty"`
	Address2         string     `json:"address_2,omitempty" path:"address_2,omitempty" url:"address_2,omitempty"`
	CardNumber       string     `json:"card_number,omitempty" path:"card_number,omitempty" url:"card_number,omitempty"`
	CardType         string     `json:"card_type,omitempty" path:"card_type,omitempty" url:"card_type,omitempty"`
	City             string     `json:"city,omitempty" path:"city,omitempty" url:"city,omitempty"`
	CompanyName      string     `json:"company_name,omitempty" path:"company_name,omitempty" url:"company_name,omitempty"`
	Country          string     `json:"country,omitempty" path:"country,omitempty" url:"country,omitempty"`
	CreatedAt        *time.Time `json:"created_at,omitempty" path:"created_at,omitempty" url:"created_at,omitempty"`
	Currency         string     `json:"currency,omitempty" path:"currency,omitempty" url:"currency,omitempty"`
	Email            string     `json:"email,omitempty" path:"email,omitempty" url:"email,omitempty"`
	PhoneNumber      string     `json:"phone_number,omitempty" path:"phone_number,omitempty" url:"phone_number,omitempty"`
	ProcessorType    string     `json:"processor_type,omitempty" path:"processor_type,omitempty" url:"processor_type,omitempty"`
	State            string     `json:"state,omitempty" path:"state,omitempty" url:"state,omitempty"`
	Zip              string     `json:"zip,omitempty" path:"zip,omitempty" url:"zip,omitempty"`
	BillingFrequency int64      `json:"billing_frequency,omitempty" path:"billing_frequency,omitempty" url:"billing_frequency,omitempty"`
	ExpirationYear   string     `json:"expiration_year,omitempty" path:"expiration_year,omitempty" url:"expiration_year,omitempty"`
	ExpirationMonth  string     `json:"expiration_month,omitempty" path:"expiration_month,omitempty" url:"expiration_month,omitempty"`
	StartYear        string     `json:"start_year,omitempty" path:"start_year,omitempty" url:"start_year,omitempty"`
	StartMonth       string     `json:"start_month,omitempty" path:"start_month,omitempty" url:"start_month,omitempty"`
	Cvv              string     `json:"cvv,omitempty" path:"cvv,omitempty" url:"cvv,omitempty"`
	PaypalToken      string     `json:"paypal_token,omitempty" path:"paypal_token,omitempty" url:"paypal_token,omitempty"`
	PaypalPayerId    string     `json:"paypal_payer_id,omitempty" path:"paypal_payer_id,omitempty" url:"paypal_payer_id,omitempty"`
	PlanId           int64      `json:"plan_id,omitempty" path:"plan_id,omitempty" url:"plan_id,omitempty"`
	SwitchToPlanId   int64      `json:"switch_to_plan_id,omitempty" path:"switch_to_plan_id,omitempty" url:"switch_to_plan_id,omitempty"`
	CreateAccount    *bool      `json:"create_account,omitempty" path:"create_account,omitempty" url:"create_account,omitempty"`
}

// Identifier no path or id

type AccountCollection []Account

type AccountCreateParams struct {
	Name             string `url:"name,omitempty" required:"false" json:"name,omitempty" path:"name"`
	CompanyName      string `url:"company_name,omitempty" required:"false" json:"company_name,omitempty" path:"company_name"`
	Address          string `url:"address,omitempty" required:"false" json:"address,omitempty" path:"address"`
	Address2         string `url:"address_2,omitempty" required:"false" json:"address_2,omitempty" path:"address_2"`
	City             string `url:"city,omitempty" required:"false" json:"city,omitempty" path:"city"`
	State            string `url:"state,omitempty" required:"false" json:"state,omitempty" path:"state"`
	Zip              string `url:"zip,omitempty" required:"false" json:"zip,omitempty" path:"zip"`
	Country          string `url:"country,omitempty" required:"false" json:"country,omitempty" path:"country"`
	Email            string `url:"email,omitempty" required:"false" json:"email,omitempty" path:"email"`
	PhoneNumber      string `url:"phone_number,omitempty" required:"false" json:"phone_number,omitempty" path:"phone_number"`
	CardNumber       string `url:"card_number,omitempty" required:"false" json:"card_number,omitempty" path:"card_number"`
	CardType         string `url:"card_type,omitempty" required:"false" json:"card_type,omitempty" path:"card_type"`
	ExpirationYear   string `url:"expiration_year,omitempty" required:"false" json:"expiration_year,omitempty" path:"expiration_year"`
	ExpirationMonth  string `url:"expiration_month,omitempty" required:"false" json:"expiration_month,omitempty" path:"expiration_month"`
	StartYear        string `url:"start_year,omitempty" required:"false" json:"start_year,omitempty" path:"start_year"`
	StartMonth       string `url:"start_month,omitempty" required:"false" json:"start_month,omitempty" path:"start_month"`
	Cvv              string `url:"cvv,omitempty" required:"false" json:"cvv,omitempty" path:"cvv"`
	PaypalToken      string `url:"paypal_token,omitempty" required:"false" json:"paypal_token,omitempty" path:"paypal_token"`
	PaypalPayerId    string `url:"paypal_payer_id,omitempty" required:"false" json:"paypal_payer_id,omitempty" path:"paypal_payer_id"`
	PlanId           int64  `url:"plan_id,omitempty" required:"false" json:"plan_id,omitempty" path:"plan_id"`
	BillingFrequency int64  `url:"billing_frequency,omitempty" required:"false" json:"billing_frequency,omitempty" path:"billing_frequency"`
	Currency         string `url:"currency,omitempty" required:"false" json:"currency,omitempty" path:"currency"`
	SwitchToPlanId   int64  `url:"switch_to_plan_id,omitempty" required:"false" json:"switch_to_plan_id,omitempty" path:"switch_to_plan_id"`
	CreateAccount    *bool  `url:"create_account,omitempty" required:"false" json:"create_account,omitempty" path:"create_account"`
}

type AccountUpdateParams struct {
	Name            string `url:"name,omitempty" required:"false" json:"name,omitempty" path:"name"`
	CompanyName     string `url:"company_name,omitempty" required:"false" json:"company_name,omitempty" path:"company_name"`
	Address         string `url:"address,omitempty" required:"false" json:"address,omitempty" path:"address"`
	Address2        string `url:"address_2,omitempty" required:"false" json:"address_2,omitempty" path:"address_2"`
	City            string `url:"city,omitempty" required:"false" json:"city,omitempty" path:"city"`
	State           string `url:"state,omitempty" required:"false" json:"state,omitempty" path:"state"`
	Zip             string `url:"zip,omitempty" required:"false" json:"zip,omitempty" path:"zip"`
	Country         string `url:"country,omitempty" required:"false" json:"country,omitempty" path:"country"`
	Email           string `url:"email,omitempty" required:"false" json:"email,omitempty" path:"email"`
	PhoneNumber     string `url:"phone_number,omitempty" required:"false" json:"phone_number,omitempty" path:"phone_number"`
	CardNumber      string `url:"card_number,omitempty" required:"false" json:"card_number,omitempty" path:"card_number"`
	CardType        string `url:"card_type,omitempty" required:"false" json:"card_type,omitempty" path:"card_type"`
	ExpirationYear  string `url:"expiration_year,omitempty" required:"false" json:"expiration_year,omitempty" path:"expiration_year"`
	ExpirationMonth string `url:"expiration_month,omitempty" required:"false" json:"expiration_month,omitempty" path:"expiration_month"`
	StartYear       string `url:"start_year,omitempty" required:"false" json:"start_year,omitempty" path:"start_year"`
	StartMonth      string `url:"start_month,omitempty" required:"false" json:"start_month,omitempty" path:"start_month"`
	Cvv             string `url:"cvv,omitempty" required:"false" json:"cvv,omitempty" path:"cvv"`
	PaypalToken     string `url:"paypal_token,omitempty" required:"false" json:"paypal_token,omitempty" path:"paypal_token"`
	PaypalPayerId   string `url:"paypal_payer_id,omitempty" required:"false" json:"paypal_payer_id,omitempty" path:"paypal_payer_id"`
}

func (a *Account) UnmarshalJSON(data []byte) error {
	type account Account
	var v account
	if err := json.Unmarshal(data, &v); err != nil {
		return lib.ErrorWithOriginalResponse{}.ProcessError(data, err, map[string]interface{}{})
	}

	*a = Account(v)
	return nil
}

func (a *AccountCollection) UnmarshalJSON(data []byte) error {
	type accounts AccountCollection
	var v accounts
	if err := json.Unmarshal(data, &v); err != nil {
		return lib.ErrorWithOriginalResponse{}.ProcessError(data, err, []map[string]interface{}{})
	}

	*a = AccountCollection(v)
	return nil
}

func (a *AccountCollection) ToSlice() *[]interface{} {
	ret := make([]interface{}, len(*a))
	for i, v := range *a {
		ret[i] = v
	}

	return &ret
}
