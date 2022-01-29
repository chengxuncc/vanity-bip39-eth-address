package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"os"
	"runtime"
	"strings"

	hdwallet "github.com/miguelmota/go-ethereum-hdwallet"
)

var prefix string
var caseSensitive bool

func init() {
	flag.StringVar(&prefix, "prefix", "cc", "vanity address prefix, default value is cc")
	flag.BoolVar(&caseSensitive, "case-sensitive", false, "case sensitive, default value is false")
}

var derivationPath = hdwallet.MustParseDerivationPath("m/44'/60'/0'/0/0")

func main() {
	flag.Parse()

	prefix = "0x" + prefix

	fmt.Println("Computing vanity address with prefix:", prefix)
	for i := 0; i < runtime.NumCPU()-1; i++ {
		go func() {
			for !genOneAndCheck() {
			}
			os.Exit(0)
		}()
	}
	for !genOneAndCheck() {

	}
}

func genOneAndCheck() bool {
	mnemonic, err := hdwallet.NewMnemonic(256)
	if err != nil {
		panic(err)
	}
	wallet, err := hdwallet.NewFromMnemonic(mnemonic)
	if err != nil {
		panic(err)
	}
	account, err := wallet.Derive(derivationPath, false)
	if err != nil {
		panic(err)
	}
	addr := account.Address.Hex()
	if !caseSensitive {
		addr = strings.ToLower(addr)
	}
	if strings.HasPrefix(addr, prefix) {
		jsonln(struct {
			Mnemonic string `json:"mnemonic"`
			Address  string `json:"address"`
		}{
			Mnemonic: mnemonic,
			Address:  addr,
		})
		return true
	}
	return false
}

func jsonln(v interface{}) {
	j, err := json.Marshal(v)
	if err != nil {
		panic(err)
	}
	fmt.Println(string(j))
}
