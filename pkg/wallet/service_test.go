package wallet

import (
	"reflect"
	"testing"

	"github.com/bahodurnazarov/findAccByID/pkg/types"
	"github.com/google/uuid"
)

func TestService_Reject_success(t *testing.T) {
	s := newTestService()
	_, payments, err := s.addAccount(defaultTestAccounts)
	if err != nil {
		t.Error(err)
		return
	}

	payment := payments[0]
	err = s.Reject(payment.ID)
	if err != nil {
		t.Errorf("Reject(): err = %v", err)
		return
	}

	savedPayment, err := s.FindPaymentByID(payment.ID)
	if err != nil {
		t.Errorf("Reject(): can't find payment, err = %v", err)
		return
	}
	if savedPayment.Status != types.PaymentStatusFail {
		t.Errorf("Reject(): wrong payment status = %v", savedPayment)
		return
	}

	savedAccount, err := s.FindAccountByID(payment.AccountID)
	if err != nil {
		t.Errorf("Reject(): can't find account, err = %v", err)
		return
	}
	if savedAccount.Balance != defaultTestAccounts.balance {
		t.Errorf("Reject(): wrong account balance = %v", savedAccount)
		return
	}
}

func TestService_FindPaymentByID_success(t *testing.T) {

	s := newTestService()
	_, payments, err := s.addAccount(defaultTestAccounts)
	if err != nil {
		t.Error(err)
		return
	}

	payment := payments[0]
	got, err := s.FindPaymentByID(payment.ID)
	if err != nil {
		t.Errorf("FindPaymentByID(): err = %v", err)
		return
	}

	if !reflect.DeepEqual(payment, got) {
		t.Errorf("FindPaymentByID():wrong payment returned = %v", err)
		return
	}
}

func TestService_FindPaymentByID_fail(t *testing.T) {
	s := newTestService()
	_, _, err := s.addAccount(defaultTestAccounts)
	if err != nil {
		t.Error(err)
		return
	}
	
	_, err = s.FindPaymentByID(uuid.New().String())
	if err == nil {
		t.Errorf("FindPaymentByID(): wrong error returned = %v", err)
		return
	}

	if err != ErrAccountNotFount {
		t.Errorf("FindPaymentByID(): wrong error returned = %v", err)
		return
	}
}
