package wallet

import (
	"errors"

	"github.com/bahodurnazarov/findAccByID/pkg/types"
)

var ErrPhoneRegister = errors.New("Phone already registred")
var ErrAmountMustBePositive = errors.New("amount must be greater than zero")
var ErrAccountNotFount = errors.New("account not found")

type Service struct {
	nextAccountID int64
	accounts      []*types.Account
	payment       []*types.Payment
}

type Messenger interface {
	Send(message string) bool
	Receive() (message string, ok bool)
}

type Telegram struct {
}

func (t *Telegram) Send(message string) bool {
	return true
}

func (t *Telegram) Receive() (message string, ok bool) {
	return "", true
}

type Error string

func (e Error) Error() string {
	return string(e)
}

// func RegisterAccount(service *Service, phone types.Phone) {
// 	for _, account := range service.accounts {
// 		if account.Phone == phone {
// 			return
// 		}
// 	}
// 	service.nextAccountID++
// 	service.accounts = append(service.accounts, &types.Account{
// 		ID:      service.nextAccountID,
// 		Phone:   phone,
// 		Balance: 0,
// 	})
// }
func (s *Service) RegisterAccount(phone types.Phone) (*types.Account, error) {
	for _, account := range s.accounts {
		if account.Phone == phone {
			return nil, ErrPhoneRegister
		}
	}
	s.nextAccountID++
	account := &types.Account{
		ID:      s.nextAccountID,
		Phone:   phone,
		Balance: 0,
	}
	s.accounts = append(s.accounts, account)
	return account, nil
}

func (s *Service) Deposit(accountID int64, amount types.Money) error {
	if amount <= 0 {
		return ErrAmountMustBePositive
	}

	var account *types.Account
	for _, acc := range s.accounts {
		if acc.ID == accountID {
			account = acc
			break
		}
	}
	if account == nil {
		return ErrAccountNotFount
	}

	account.Balance += amount
	return nil
}
