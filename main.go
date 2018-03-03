package main

import (
	"fmt"
	"io/ioutil"
	"os"
	"os/exec"
	"strings"

	exp "github.com/toolateforteddy/cryptex/export"
	"github.com/toolateforteddy/cryptex/process"
	kingpin "gopkg.in/alecthomas/kingpin.v2"
)

// var cryptKey = [32]byte("XVlBzgbaiCMRAjWwhTHctcuAxhxKQFDa")

var key [32]byte

const strKey = "AD_-7nN8yRCmppoYyqmpfOItyxRKVQbf"

func init() {
	for i := range strKey {
		key[i] = strKey[i]
	}
}

const (
	defaultSecretFileName    = ".cryptex/secrets.toml"
	defaultDecryptedFileName = ".cryptex/decrypted.toml"

	keyFile = ".cryptex/stored_key"
)

var (
	verbose = kingpin.Flag("verbose", "Verbose mode.").Short('v').Bool()

	decrypt     = kingpin.Command("decrypt", "create the unprotected file.")
	decryptFile = decrypt.Flag("file", "encrypted file name").Short('f').
			Default(defaultSecretFileName).String()

	encrypt     = kingpin.Command("encrypt", "create the secret file.")
	encryptFile = encrypt.Flag("file", "unprotected file name").Short('f').
			Default(defaultDecryptedFileName).String()
	encryptOutput = encrypt.Flag("output", "output file").Short('o').
			Default(defaultSecretFileName).String()

	export     = kingpin.Command("export", "export the decrypted data from the protected file.")
	exportFile = export.Flag("file", "encrypted file name").Short('f').
			Default(defaultSecretFileName).String()

	saveKey     = kingpin.Command("save", "save the key to disk")
	saveKeyFile = saveKey.Flag("file", "file to save to").Short('f').
			Default(keyFile).String()

	edit     = kingpin.Command("edit", "open vim to edit the encrypted file")
	editFile = edit.Flag("file", "encrypted file name").Short('f').
			Default(defaultSecretFileName).String()
)

func main() {
	switch kingpin.Parse() {
	case "save":
		pKey := process.NewKey(&key)
		err := pKey.SaveToDisk(keyFile)
		if err != nil {
			fmt.Println(err)
		}
	case "export":
		pKey := process.NewKey(&key)
		tmp, err := ioutil.TempFile("", "secret")
		kingpin.FatalIfError(err, "Error making tmp file")
		tmpFileName := tmp.Name()
		tmp.Close()
		defer func() {
			os.Remove(tmpFileName)
		}()
		err = pKey.DecryptFile(*exportFile, tmpFileName)
		kingpin.FatalIfError(err, "Error decrypting file")
		data, err := exp.LoadFile(tmpFileName)
		kingpin.FatalIfError(err, "Error loading tmpfile")
		vars, err := exp.FormatForShellExport(data)
		kingpin.FatalIfError(err, "Error creating shell exports")
		fmt.Printf("%v\n", strings.Join(vars, "\n"))
	case "encrypt":
		pKey := process.NewKey(&key)
		err := pKey.EncryptFile(*encryptFile, *encryptOutput)
		kingpin.FatalIfError(err, "Error encrypting file")
	case "decrypt":
		pKey := process.NewKey(&key)
		err := pKey.DecryptFile(*decryptFile, defaultDecryptedFileName)
		kingpin.FatalIfError(err, "Error decrypting file")
	case "edit":
		pKey := process.NewKey(&key)
		tmp, err := ioutil.TempFile("", "secret")
		kingpin.FatalIfError(err, "Error making tmp file")
		tmpFileName := tmp.Name()
		tmp.Close()
		defer func() {
			os.Remove(tmpFileName)
		}()
		err = pKey.DecryptFile(*editFile, tmpFileName)
		kingpin.FatalIfError(err, "Error decrypting file")
		fmt.Println(tmpFileName)
		vim := exec.Command("vim", tmpFileName)
		vim.Stdin = os.Stdin
		vim.Stdout = os.Stdout
		vim.Stderr = os.Stderr
		err = vim.Run()
		kingpin.FatalIfError(err, "Error editing file")
		err = pKey.EncryptFile(tmpFileName, *editFile)
		kingpin.FatalIfError(err, "Error encrypting file")
	}
	// pKey.EncryptFile("secrets.toml", "tmp.toml")
	// pKey.DecryptFile("tmp.toml", "decrypted.toml")
}
