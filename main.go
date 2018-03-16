package main

import (
	"fmt"
	"io/ioutil"
	"os"
	"os/exec"
	"strings"

	"github.com/toolateforteddy/cryptex/cryptopasta"
	exp "github.com/toolateforteddy/cryptex/export"
	"github.com/toolateforteddy/cryptex/process"
	kingpin "gopkg.in/alecthomas/kingpin.v2"
)

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

	keyEnvName = "CRYPTEX_KEY"
)

var (
	verbose    = kingpin.Flag("verbose", "Verbose mode.").Short('v').Bool()
	keyCmdFlag = kingpin.Flag("key", "the key!").Short('k').String()

	decrypt     = kingpin.Command("decrypt", "create the unprotected file.")
	decryptFile = decrypt.Flag("file", "encrypted file name").Short('f').
			Default(defaultSecretFileName).String()

	edit     = kingpin.Command("edit", "open vim to edit the encrypted file")
	editFile = edit.Flag("file", "encrypted file name").Short('f').
			Default(defaultSecretFileName).String()

	encrypt     = kingpin.Command("encrypt", "create the secret file.")
	encryptFile = encrypt.Flag("file", "unprotected file name").Short('f').
			Default(defaultDecryptedFileName).String()
	encryptOutput = encrypt.Flag("output", "output file").Short('o').
			Default(defaultSecretFileName).String()

	export     = kingpin.Command("export", "export the decrypted data from the protected file.")
	exportFile = export.Flag("file", "encrypted file name").Short('f').
			Default(defaultSecretFileName).String()
	exportSep = export.Flag("seperator", "string to join the env vars on").
			Default("\n").String()

	initialize     = kingpin.Command("init", "initialize the cryptex storage")
	initializeFile = initialize.Flag("file", "encrypted file name").Short('f').
			Default(defaultSecretFileName).String()

	keygen     = kingpin.Command("keygen", "create a new keyfile")
	keygenSave = keygen.Flag("save", "save key to disk").Short('s').Bool()

	saveKey     = kingpin.Command("save", "save the key to disk")
	saveKeyFile = saveKey.Flag("file", "file to save to").Short('f').
			Default(keyFile).String()
)

func main() {
	switch kingpin.Parse() {
	case "init":
		f, err := os.OpenFile(
			*initializeFile,
			os.O_RDWR|os.O_CREATE,
			os.ModeAppend|os.ModePerm)
		kingpin.FatalIfError(err, "Error creating secrets file [%s]", *initializeFile)
		f.Close()
	case "keygen":
		k := cryptopasta.NewEncryptionKey()
		pKey := process.NewKey(k)
		fmt.Println("New Key Generated:")
		fmt.Printf("%s\n", *pKey.Str())
		if *keygenSave {
			err := pKey.SaveToDisk(keyFile)
			kingpin.FatalIfError(err, "Error saving key to disk")
		}
	case "save":
		pKey := getKey()
		err := pKey.SaveToDisk(keyFile)
		kingpin.FatalIfError(err, "Error saving key to disk")
	case "export":
		msg, err := exportFunc()
		kingpin.FatalIfError(err, msg)
	case "encrypt":
		pKey := getKey()
		err := pKey.EncryptFile(*encryptFile, *encryptOutput)
		kingpin.FatalIfError(err, "Error encrypting file")
	case "decrypt":
		pKey := getKey()
		err := pKey.DecryptFile(*decryptFile, defaultDecryptedFileName)
		kingpin.FatalIfError(err, "Error decrypting file")
	case "edit":
		msg, err := editFunc()
		kingpin.FatalIfError(err, msg)
	}
}

func exportFunc() (string, error) {
	pKey := getKey()
	tmp, err := ioutil.TempFile("", "secret")
	if err != nil {
		return "Error making tmp file", err
	}
	tmpFileName := tmp.Name()
	tmp.Close()
	defer func() {
		os.Remove(tmpFileName)
	}()
	err = pKey.DecryptFile(*exportFile, tmpFileName)
	if err != nil {
		return "Error decrypting file", err
	}
	data, err := exp.LoadFile(tmpFileName)
	if err != nil {
		return "Error loading tmpfile", err
	}
	vars, err := exp.FormatForShellExport(data)
	if err != nil {
		return "Error creating shell exports", err
	}
	fmt.Printf("%v", strings.Join(vars, *exportSep))
	return "", nil
}

func editFunc() (string, error) {
	pKey := getKey()
	tmp, err := ioutil.TempFile("", "secret")
	if err != nil {
		return "Error making tmp file", err
	}
	tmpFileName := tmp.Name()
	tmp.Close()
	defer func() {
		os.Remove(tmpFileName)
	}()
	err = pKey.DecryptFile(*editFile, tmpFileName)
	if err != nil {
		return "Error decrypting file", err
	}

	vim := exec.Command("vim", tmpFileName)
	vim.Stdin = os.Stdin
	vim.Stdout = os.Stdout
	vim.Stderr = os.Stderr
	err = vim.Run()
	if err != nil {
		return "Error editing file", err
	}
	err = pKey.EncryptFile(tmpFileName, *editFile)
	if err != nil {
		return "Error encrypting file", err
	}
	return "", nil
}

func getKey() process.Key {
	if keyCmdFlag != nil && len(*keyCmdFlag) != 0 {
		return process.NewKeyFromPrintable(*keyCmdFlag)
	}
	if envKey := os.Getenv(keyEnvName); envKey != "" {
		return process.NewKeyFromPrintable(envKey)
	}
	k, err := process.NewFromDisk(keyFile)
	if err == nil {
		return k
	}
	fmt.Printf("Failed to load key from disk. %v\n", err)
	kingpin.Fatalf("No Key Found.")
	return process.Key{}
}
