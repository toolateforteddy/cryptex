package process

import (
	"bufio"
	"encoding/base64"
	"fmt"
	"io"
	"os"
	"strings"

	"github.com/toolateforteddy/cryptex/cryptopasta"
	"github.com/toolateforteddy/errortrace"
)

type Key struct {
	k *[32]byte
}

func NewKey(k *[32]byte) Key {
	return Key{k}
}

func (k *Key) key() *[32]byte {
	return k.k
}

func (k *Key) SaveToDisk(file string) error {
	dFile, _ := os.OpenFile(file,
		os.O_RDWR|os.O_CREATE|os.O_TRUNC,
		os.ModeAppend|os.ModePerm)
	defer dFile.Close()
	_, err := dFile.Write(k.key()[:])
	return errortrace.Wrap(err)
}

func NewFromDisk(file string) (Key, error) {
	sFile, err := os.Open(file)
	defer sFile.Close()

	data := []byte{}
	numBytes, err := sFile.Read(data)
	if err != nil {
		return Key{}, errortrace.Wrap(err)
	}
	if numBytes == 0 {
		return Key{}, errortrace.Errorf("empty keyfile")
	}
	if numBytes < 32 {
		return Key{}, errortrace.Errorf("insufficient bytes in file [%d]", numBytes)
	}
	key := [32]byte{}
	for i := range key {
		key[i] = data[i]
	}
	return Key{&key}, nil
}

func (k Key) EncryptFile(src, dest string) error {
	sFile, _ := os.Open(src)
	defer sFile.Close()
	dFile, _ := os.OpenFile(dest,
		os.O_RDWR|os.O_CREATE|os.O_TRUNC,
		os.ModeAppend|os.ModePerm)
	defer dFile.Close()

	brs := bufio.NewReader(sFile)
	bwd := bufio.NewWriter(dFile)
	defer bwd.Flush()
	for {
		line, err := brs.ReadString('\n')
		if err != nil {
			if err == io.EOF {
				break
			}
			return errortrace.Wrap(err)
		}
		line = strings.TrimSuffix(line, "\n")
		line = k.encryptLine(line)
		_, err = bwd.WriteString(line)
		if err != nil {
			return errortrace.Wrap(err)
		}
		err = bwd.WriteByte('\n')
		if err != nil {
			return errortrace.Wrap(err)
		}
	}

	return nil
}

func (k *Key) encryptLine(in string) string {
	kvArr := strings.SplitAfterN(in, "=", 2)
	if len(kvArr) < 2 {
		return in
	}
	secret, err := cryptopasta.Encrypt([]byte(kvArr[1]), k.key())
	if err != nil {
		panic(err)
	}

	kvArr[1] = fmt.Sprintf("%q", base64.StdEncoding.EncodeToString(secret))
	return strings.Join(kvArr, " ")
}

func (k *Key) DecryptFile(src, dest string) error {
	sFile, _ := os.Open(src)
	defer sFile.Close()
	dFile, _ := os.OpenFile(dest,
		os.O_RDWR|os.O_CREATE|os.O_TRUNC,
		os.ModeAppend|os.ModePerm)
	defer dFile.Close()

	brs := bufio.NewReader(sFile)
	bwd := bufio.NewWriter(dFile)
	defer bwd.Flush()
	for {
		line, err := brs.ReadString('\n')
		if err != nil {
			if err == io.EOF {
				break
			}
			return errortrace.Wrap(err)
		}
		line = strings.TrimSuffix(line, "\n")
		line = k.decryptLine(line)
		_, err = bwd.WriteString(line)
		if err != nil {
			return errortrace.Wrap(err)
		}
		err = bwd.WriteByte('\n')
		if err != nil {
			return errortrace.Wrap(err)
		}
	}
	return nil
}

func (k *Key) decryptLine(in string) string {
	kvArr := strings.SplitAfterN(in, "=", 2)
	if len(kvArr) < 2 {
		return in
	}
	trimmedLine := strings.Trim(strings.TrimSpace(kvArr[1]), "\"")
	secret, _ := base64.StdEncoding.DecodeString(trimmedLine)
	byteArr, _ := cryptopasta.Decrypt(secret, k.key())

	kvArr[1] = string(byteArr)
	return strings.Join(kvArr[:2], "")
}
