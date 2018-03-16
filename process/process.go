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
	k            *[32]byte
	printableKey *string
}

func NewKey(k *[32]byte) Key {
	return Key{k: k}.initialize()
}

func NewKeyFromPrintable(str string) Key {
	return Key{printableKey: &str}.initialize()
}

func (k *Key) key() *[32]byte {
	return k.k
}
func (k *Key) Str() *string {
	return k.printableKey
}

func (k Key) initialize() Key {
	if k.k != nil {
		b64Key := base64.StdEncoding.EncodeToString(k.k[:])
		k.printableKey = &b64Key
	} else if k.printableKey != nil {
		k.k, _ = keyFromString(*k.printableKey)
	}
	return k
}

func keyFromString(str string) (*[32]byte, error) {
	secret, err := base64.StdEncoding.DecodeString(str)
	if err != nil {
		return nil, errortrace.Wrap(err)
	}
	if len(secret) != 32 {
		return nil, errortrace.Errorf("incorrect data format")
	}
	key := [32]byte{}
	copy(key[:], secret)
	return &key, nil
}

func (k *Key) SaveToDisk(file string) error {
	dFile, _ := os.OpenFile(file,
		os.O_RDWR|os.O_CREATE|os.O_TRUNC,
		os.ModeAppend|os.ModePerm)
	defer dFile.Close()
	_, err := dFile.WriteString(*k.Str())
	return errortrace.Wrap(err)
}

func NewFromDisk(file string) (Key, error) {
	sFile, err := os.Open(file)
	defer sFile.Close()
	brs := bufio.NewReader(sFile)
	line, _, err := brs.ReadLine()
	if err != nil {
		return Key{}, errortrace.Wrap(err)
	}
	if len(line) == 0 {
		return Key{}, errortrace.Errorf("empty keyfile")
	}
	b64Key := string(line)
	k, err := keyFromString(b64Key)
	if err != nil {
		return Key{}, errortrace.Wrap(err)
	}
	return Key{k: k, printableKey: &b64Key}, nil
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
