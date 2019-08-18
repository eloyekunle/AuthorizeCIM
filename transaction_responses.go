package AuthorizeCIM

import "fmt"

func (transx TransactionResponse) TransactionID() string {
	return transx.Response.TransID
}

func (transx TransactionResponse) Message() string {
	return transx.Response.Errors[0].ErrorText
}

func (transx TransactionResponse) AVS() AVS {
	out := AVS{
		avsResultCode:  transx.Response.AvsResultCode,
		cvvResultCode:  transx.Response.CvvResultCode,
		cavvResultCode: transx.Response.CavvResultCode,
	}
	return out
}

type AVS struct {
	avsResultCode  string
	cvvResultCode  string
	cavvResultCode string
}

type TransxReponse interface {
	ErrorMessage()
	Approved()
}

// Error represents an Authorize.NET error.
type Error struct {
	ResultCode string
	Message    []Message
}

func (e Error) Error() string {
	return fmt.Sprintf("resultCode: %s, messageCode: %s, messageText: %s", e.ResultCode, e.Message[0].Code, e.Message[0].Text)
}

func (r MessagesResponse) ErrorMessage() string {
	return r.Messages.Message[0].Text
}

// Error derives an error from MessagesResponse.
func (r MessagesResponse) Error() error {
	return Error{
		ResultCode: r.Messages.ResultCode,
		Message:    r.Messages.Message,
	}
}

func (r TransactionResponse) Approved() bool {
	if r.Response.ResponseCode == "1" || r.Response.ResponseCode == "4" {
		return true
	}
	return false
}

func (r TransactionResponse) Held() bool {
	return r.Response.ResponseCode == "4"
}

func (r MessagesResponse) Ok() bool {
	return r.Messages.ResultCode == "Ok" && r.Messages.Message[0].Code == "I00001"
}
