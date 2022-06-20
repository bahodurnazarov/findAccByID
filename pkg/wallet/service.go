package wallet

import (
	"errors"
	"fmt"

	"github.com/bahodurnazarov/findAccByID/pkg/types"
	"github.com/google/uuid"
)

var ErrPhoneRegister = errors.New("phone already registred")
var ErrAmountMustBePositive = errors.New("amount must be greater than zero")
var ErrAccountNotFount = errors.New("account not found")
var ErrNotEnoughBalance = errors.New("not Enought Balance")

type Service struct {
	nextAccountID int64
	accounts      []*types.Account
	payments      []*types.Payment
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

func (s *Service) Pay(accountID int64, amount types.Money, category types.PaymentCategory) (*types.Payment, error) {
	if amount <= 0 {
		return nil, ErrAmountMustBePositive
	}

	var account *types.Account
	for _, acc := range s.accounts {
		if acc.ID == accountID {
			account = acc
			break
		}
	}
	if account == nil {
		return nil, ErrAccountNotFount
	}

	if account.Balance < amount {
		return nil, ErrNotEnoughBalance
	}

	account.Balance -= amount
	paymentID := uuid.New().String()
	payment := &types.Payment{
		ID:        paymentID,
		AccountID: accountID,
		Amount:    amount,
		Category:  category,
		Status:    types.PaymentStatusInProgress,
	}
	s.payments = append(s.payments, payment)
	return payment, nil
}

func (s *Service) FindAccountByID(accountID int64) (*types.Account, error) {
	for _, account := range s.accounts {
		if account.ID == accountID {
			return account, nil
		}
	}
	return nil, ErrAccountNotFount
}

func (s *Service) FindPaymentByID(paymentID string) (*types.Payment, error) {
	for _, payment := range s.payments {
		if payment.ID == paymentID {
			return payment, nil
		}
	}
	return nil, ErrAccountNotFount
}

func (s *Service) Reject(paymentID string) error {
	payment, err := s.FindPaymentByID(paymentID)
	if err != nil {
		return err
	}

	account, err := s.FindAccountByID(payment.AccountID)
	if err != nil {
		return err
	}

	payment.Status = types.PaymentStatusFail
	account.Balance += payment.Amount
	return nil
}

type testService struct {
	*Service
}

func newTestService() *testService {
	return &testService{Service: &Service{}}
}

// func (s *testService) addAccountWithBalance(phone types.Phone, balance types.Money) (*types.Account, error) {
// 	account, err := s.RegisterAccount(phone)
// 	if err != nil {
// 		return nil, fmt.Errorf("can't register account, error = %v", err)
// 	}

// 	err = s.Deposit(account.ID, balance)
// 	if err != nil {
// 		return nil, fmt.Errorf("can't deposit, error = %v", err)
// 	}
// 	return account, nil
// }

type testAccount struct {
	phone types.Phone
	balance types.Money
	payments []struct {
		amount types.Money
		category types.PaymentCategory
	}
}

var defaultTestAccounts = testAccount{
	phone: "+9920000001",
	balance: 10_000_00,
	payments: []struct {
		amount types.Money
		category types.PaymentCategory
	}{
		{amount: 1_000_00, category: "auto"},
	},
}

func (s *testService) addAccount(data testAccount) (*types.Account, []*types.Payment, error) {
	account, err := s.RegisterAccount(data.phone)
	if err != nil {
		return nil, nil, fmt.Errorf("can't register account, error = %v", err)
	}

	err = s.Deposit(account.ID, data.balance)
	if err != nil {
		return nil, nil, fmt.Errorf("can't deposit, error = %v", err)
	}

	payments := make([]*types.Payment, len(data.payments))
	for i, payment := range data.payments {
		payments[i], err = s.Pay(account.ID, payment.amount, payment.category)
		if err != nil {
			return nil, nil, fmt.Errorf("can't pay, error = %v", err)
		}
	}
	return account, payments, nil
}

