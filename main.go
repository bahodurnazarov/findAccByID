package main

import (
	"fmt"

	"github.com/bahodurnazarov/findAccByID/pkg/wallet"
)

func main() {
	svc := &wallet.Service{}
	account, err := svc.RegisterAccount("+9920000001")
	if err != nil {
		fmt.Println(err)
		return
	}

	err = svc.Deposit(account.ID, 10)
	if err != nil {
		switch err {
		case wallet.ErrAmountMustBePositive:
			fmt.Println(" Сумма должна быть положительной")
		case wallet.ErrAccountNotFount:
			fmt.Println(" Aккаунт пользователя не найден")
		}
		return
	}
	fmt.Println(account.Balance)
}
