package cryptex

import (
	"bufio"
	"encoding/base64"
	"os"
	"strings"

	"cryptex/cryptopasta"
)

func EncryptFile(src, dest string) error {
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
			return err
		}
		line = strings.TrimSuffix(line, "\n")
		line = encryptLine(line)
		bwd.WriteString(line)
		bwd.WriteByte('\n')
	}

	return nil
}

func encryptLine(in string) string {
	kvArr := strings.SplitAfter(in, "=")
	if len(kvArr) < 2 {
		return in
	}
	secret, err := cryptopasta.Encrypt([]byte(kvArr[1]), &key)
	if err != nil {
		panic(err)
	}

	kvArr[1] = base64.StdEncoding.EncodeToString(secret)
	return strings.Join(kvArr, " ")
}

func DecryptFile(src, dest string) error {
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
			return err
		}
		line = strings.TrimSuffix(line, "\n")
		line = decryptLine(line)
		bwd.WriteString(line)
		bwd.WriteByte('\n')
	}

	return nil
}

func decryptLine(in string) string {
	kvArr := strings.SplitAfter(in, "=")
	if len(kvArr) < 2 {
		return in
	}
	secret, _ := base64.StdEncoding.DecodeString(strings.TrimSpace(kvArr[1]))
	byteArr, _ := cryptopasta.Decrypt(secret, &key)

	kvArr[1] = string(byteArr)
	return strings.Join(kvArr[:2], "")
}
