package dto

// @Description Account creation payload
//
//	type AccountRequest struct {
//		// Account ID
//		// @example 123
//		AccountID int64 `json:"account_id"`
//		// Initial balance (string to allow decimal format)
//		// @example 100.23344
//		Balance string `json:"initial_balance"`
//	}
type UserRequest struct {
	Email    string `validate:"required,max=200,min=1" json:"email"`
	Passcode string `validate:"max=8,min=1" json:"passcode"`
}

// type InternalAccountRequest struct {
// 	AccountID int64
// 	Balance   decimal.Decimal
// }
