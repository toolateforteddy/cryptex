package main

import (
	"fmt"

	"github.com/BurntSushi/toml"
)

// var cryptKey = [32]byte("XVlBzgbaiCMRAjWwhTHctcuAxhxKQFDa")

var key [32]byte

const strKey = "AD_-7nN8yRCmppoYyqmpfOItyxRKVQbf"

func init() {
	for i := range strKey {
		key[i] = strKey[i]
	}
}

func LoadFile(filename string) (interface{}, error) {
	var config interface{}
	if _, err := toml.DecodeFile(filename, &config); err != nil {
		fmt.Println(err)
		return nil, err
	}
	return config, nil
}

func main() {

	EncryptFile("secrets.toml", "tmp.toml")
	DecryptFile("tmp.toml", "decrypted.toml")
}
