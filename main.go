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
		fmt.Println(err)
		return
	}
}
