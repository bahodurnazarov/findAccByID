package main

import (
	"github.com/bahodurnazarov/findAccByID/pkg/wallet"
)

func main() {
	svc := &wallet.Service{}
	wallet.RegisterAccount(svc, "+992000000001")
}
