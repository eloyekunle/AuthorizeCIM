package AuthorizeCIM

import (
	"encoding/json"
)

func (response SubscriptionResponse) Approved() bool {
	return response.Messages.ResultCode == "Ok"
}

func (response SubscriptionResponse) CustomerProfileId() string {
	return response.Profile.CustomerProfileID
}

func (response SubscriptionResponse) CustomerPaymentProfileId() string {
	return response.Profile.CustomerPaymentProfileID
}

func (response SubscriptionResponse) ErrorMessage() string {
	return response.Messages.Message[0].Text
}

func (sub Subscription) Charge() (*SubscriptionResponse, error) {
	response, err := SendSubscription(sub)
	return response, err
}

func (sub Subscription) Update() (*SubscriptionResponse, error) {
	response, err := UpdateSubscription(sub)
	return response, err
}

func (response SubscriptionResponse) Info() string {
	return response.Messages.Message[0].Text
}

type UpdateSubscriptionRequest struct {
	ARBUpdateSubscriptionRequest ARBUpdateSubscriptionRequest `json:"ARBUpdateSubscriptionRequest"`
}

type ARBUpdateSubscriptionRequest struct {
	MerchantAuthentication MerchantAuthentication `json:"merchantAuthentication"`
	RefID                  string                 `json:"refId,omitempty"`
	SubscriptionId         string                 `json:"subscriptionId,omitempty"`
	Subscription           Subscription           `json:"subscription,omitempty"`
}

type CreateSubscriptionRequest struct {
	ARBCreateSubscriptionRequest ARBCreateSubscriptionRequest `json:"ARBCreateSubscriptionRequest"`
}

type ARBCreateSubscriptionRequest struct {
	MerchantAuthentication MerchantAuthentication `json:"merchantAuthentication"`
	RefID                  string                 `json:"refId,omitempty"`
	Subscription           Subscription           `json:"subscription"`
}

type GetSubscriptionStatusRequest struct {
	ARBGetSubscriptionStatusRequest ARBGetSubscriptionRequest `json:"ARBGetSubscriptionStatusRequest"`
}

type GetSubscriptionCancelRequest struct {
	ARBCancelSubscriptionRequest ARBGetSubscriptionRequest `json:"ARBCancelSubscriptionRequest"`
}

type GetSubscriptionRequest struct {
	ARBGetSubscriptionRequest ARBGetSubscriptionRequest `json:"ARBGetSubscriptionRequest"`
}

type ARBGetSubscriptionRequest struct {
	MerchantAuthentication MerchantAuthentication `json:"merchantAuthentication"`
	RefID                  string                 `json:"refId"`
	SubscriptionID         string                 `json:"subscriptionId"`
}

type SetSubscription struct {
	Id string `json:"subscriptionId"`
}

type Subscription struct {
	Name            string            `json:"name,omitempty"`
	PaymentSchedule *PaymentSchedule  `json:"paymentSchedule,omitempty"`
	Amount          string            `json:"amount,omitempty"`
	TrialAmount     string            `json:"trialAmount,omitempty"`
	Payment         *Payment          `json:"payment,omitempty"`
	BillTo          *BillTo           `json:"billTo,omitempty"`
	SubscriptionId  string            `json:"subscriptionId,omitempty"`
	Profile         *CustomerProfiler `json:"profile,omitempty"`
}

type BillTo struct {
	FirstName   string `json:"firstName,omitempty"`
	LastName    string `json:"lastName,omitempty"`
	Address     string `json:"address,omitempty"`
	City        string `json:"city,omitempty"`
	State       string `json:"state,omitempty"`
	Zip         string `json:"zip,omitempty"`
	Country     string `json:"country,omitempty"`
	PhoneNumber string `json:"phoneNumber,omitempty"`
}

type PaymentSchedule struct {
	Interval         Interval `json:"interval,omitempty"`
	StartDate        string   `json:"startDate,omitempty"`
	TotalOccurrences string   `json:"totalOccurrences,omitempty"`
	TrialOccurrences string   `json:"trialOccurrences,omitempty"`
}

type Interval struct {
	Length string `json:"length,omitempty"`
	Unit   string `json:"unit,omitempty"`
}

type SubscriptionResponse struct {
	SubscriptionID string `json:"subscriptionId"`
	Profile        struct {
		CustomerProfileID        string `json:"customerProfileId"`
		CustomerPaymentProfileID string `json:"customerPaymentProfileId"`
	} `json:"profile"`
	Messages struct {
		ResultCode string    `json:"resultCode"`
		Message    []Message `json:"message"`
	} `json:"messages"`
}

func SendSubscription(sub Subscription) (*SubscriptionResponse, error) {
	action := CreateSubscriptionRequest{
		ARBCreateSubscriptionRequest: ARBCreateSubscriptionRequest{
			MerchantAuthentication: GetAuthentication(),
			Subscription:           sub,
		},
	}
	jsoned, err := json.Marshal(action)
	if err != nil {
		return nil, err
	}
	response, err := SendRequest(jsoned)
	if err != nil {
		return nil, err
	}

	var dat SubscriptionResponse
	err = json.Unmarshal(response, &dat)
	if err != nil {
		return nil, err
	}
	return &dat, err
}

func UpdateSubscription(sub Subscription) (*SubscriptionResponse, error) {
	action := UpdateSubscriptionRequest{
		ARBUpdateSubscriptionRequest: ARBUpdateSubscriptionRequest{
			MerchantAuthentication: GetAuthentication(),
			SubscriptionId:         sub.SubscriptionId,
			Subscription: Subscription{
				Payment: &Payment{
					CreditCard: sub.Payment.CreditCard,
				},
			},
		},
	}
	jsoned, err := json.Marshal(action)
	if err != nil {
		return nil, err
	}
	response, err := SendRequest(jsoned)
	if err != nil {
		return nil, err
	}

	var dat SubscriptionResponse
	err = json.Unmarshal(response, &dat)
	if err != nil {
		return nil, err
	}
	return &dat, err
}

func (sub GetSubscriptionList) Count() int {
	return sub.TotalNumInResultSet
}

func (sub SetSubscription) Info() (*GetSubscriptionResponse, error) {
	action := GetSubscriptionRequest{
		ARBGetSubscriptionRequest: ARBGetSubscriptionRequest{
			MerchantAuthentication: GetAuthentication(),
			SubscriptionID:         sub.Id,
		},
	}
	jsoned, err := json.Marshal(action)
	if err != nil {
		return nil, err
	}
	response, err := SendRequest(jsoned)
	if err != nil {
		return nil, err
	}

	var dat GetSubscriptionResponse
	err = json.Unmarshal(response, &dat)
	if err != nil {
		return nil, err
	}
	return &dat, err
}

func (s SubscriptionStatus) Active() bool {
	return s.Status == "active"
}

func (sub SetSubscription) Status() (*SubscriptionStatus, error) {
	action := GetSubscriptionStatusRequest{
		ARBGetSubscriptionStatusRequest: ARBGetSubscriptionRequest{
			MerchantAuthentication: GetAuthentication(),
			SubscriptionID:         sub.Id,
		},
	}
	jsoned, err := json.Marshal(action)
	if err != nil {
		return nil, err
	}
	response, err := SendRequest(jsoned)
	if err != nil {
		return nil, err
	}

	var dat SubscriptionStatus
	err = json.Unmarshal(response, &dat)
	if err != nil {
		return nil, err
	}
	return &dat, err
}

func (sub SetSubscription) Cancel() (*SubscriptionCancel, error) {
	action := GetSubscriptionCancelRequest{
		ARBCancelSubscriptionRequest: ARBGetSubscriptionRequest{
			MerchantAuthentication: GetAuthentication(),
			SubscriptionID:         sub.Id,
		},
	}
	jsoned, err := json.Marshal(action)
	if err != nil {
		return nil, err
	}
	response, err := SendRequest(jsoned)
	if err != nil {
		return nil, err
	}

	var dat SubscriptionCancel
	err = json.Unmarshal(response, &dat)
	if err != nil {
		return nil, err
	}
	return &dat, err
}

func SubscriptionList(search string) (*GetSubscriptionList, error) {
	action := GetSubscriptionListRequest{
		ARBGetSubscriptionListRequest: ARBGetSubscriptionListRequest{
			MerchantAuthentication: GetAuthentication(),
			SearchType:             search,
			Sorting: Sorting{
				OrderBy:         "id",
				OrderDescending: "false",
			},
			Paging: Paging{
				Limit:  "1000",
				Offset: "1",
			},
		},
	}
	jsoned, err := json.Marshal(action)
	if err != nil {
		return nil, err
	}
	response, err := SendRequest(jsoned)
	if err != nil {
		return nil, err
	}

	var dat GetSubscriptionList
	err = json.Unmarshal(response, &dat)
	if err != nil {
		return nil, err
	}
	return &dat, err
}

func (sub GetSubscriptionResponse) ErrorMessage() string {
	return sub.Messages.Message[0].Text
}

type GetSubscriptionResponse struct {
	Subscription struct {
		Name            string `json:"name"`
		PaymentSchedule struct {
			Interval struct {
				Length int    `json:"length"`
				Unit   string `json:"unit"`
			} `json:"interval"`
			StartDate        string `json:"startDate"`
			TotalOccurrences int    `json:"totalOccurrences"`
			TrialOccurrences int    `json:"trialOccurrences"`
		} `json:"paymentSchedule"`
		Amount      float64 `json:"amount"`
		TrialAmount float64 `json:"trialAmount"`
		Status      string  `json:"status"`
		Profile     struct {
			PaymentProfile struct {
				CustomerPaymentProfileID string `json:"customerPaymentProfileId"`
				Payment                  struct {
					CreditCard struct {
						CardNumber     string `json:"cardNumber"`
						ExpirationDate string `json:"expirationDate"`
					} `json:"creditCard"`
				} `json:"payment"`
				CustomerType string `json:"customerType"`
				BillTo       struct {
					FirstName string `json:"firstName"`
					LastName  string `json:"lastName"`
				} `json:"billTo"`
			} `json:"paymentProfile"`
			CustomerProfileID string `json:"customerProfileId"`
			Description       string `json:"description"`
		} `json:"profile"`
	} `json:"subscription"`
	RefID string `json:"refId"`
	MessagesResponse
}

type SubscriptionStatus struct {
	Note            string `json:"note"`
	Status          string `json:"status"`
	StatusSpecified bool   `json:"statusSpecified"`
	RefID           string `json:"refId"`
	MessagesResponse
}

type SubscriptionCancel struct {
	RefID string `json:"refId"`
	MessagesResponse
}

type GetSubscriptionListRequest struct {
	ARBGetSubscriptionListRequest ARBGetSubscriptionListRequest `json:"ARBGetSubscriptionListRequest"`
}

type ARBGetSubscriptionListRequest struct {
	MerchantAuthentication MerchantAuthentication `json:"merchantAuthentication"`
	SearchType             string                 `json:"searchType"`
	Sorting                Sorting                `json:"sorting"`
	Paging                 Paging                 `json:"paging"`
}

type Sorting struct {
	OrderBy         string `json:"orderBy"`
	OrderDescending string `json:"orderDescending"`
}

type Paging struct {
	Limit  string `json:"limit"`
	Offset string `json:"offset"`
}

type GetSubscriptionList struct {
	TotalNumInResultSet int `json:"totalNumInResultSet"`
	SubscriptionDetails []struct {
		ID                        int     `json:"id"`
		Name                      string  `json:"name"`
		Status                    string  `json:"status"`
		CreateTimeStampUTC        string  `json:"createTimeStampUTC"`
		FirstName                 string  `json:"firstName"`
		LastName                  string  `json:"lastName"`
		TotalOccurrences          int     `json:"totalOccurrences"`
		PastOccurrences           int     `json:"pastOccurrences"`
		PaymentMethod             string  `json:"paymentMethod"`
		AccountNumber             string  `json:"accountNumber"`
		Invoice                   string  `json:"invoice"`
		Amount                    float64 `json:"amount"`
		CurrencyCode              string  `json:"currencyCode"`
		CustomerProfileID         int     `json:"customerProfileId"`
		CustomerPaymentProfileID  int     `json:"customerPaymentProfileId,omitempty"`
		CustomerShippingProfileID int     `json:"customerShippingProfileId,omitempty"`
	} `json:"subscriptionDetails"`
	MessagesResponse
}
