package api

import (
	"fmt"
	"os"

	// QR CODES GENERATOR

	qrcode "github.com/yeqown/go-qrcode/v2"
	standard "github.com/yeqown/go-qrcode/writer/standard"
)

const (
	minimumNodes            = 3
	tofu_secrets_dir        = "./.tofu_secrets"
	unseal_keys_secrets_dir = tofu_secrets_dir + "/.unseal_keys"
	root_token_file         = tofu_secrets_dir + "/.root_token"
	qrcodes_prefix          = ""
	// qrcodes_prefix          = "https://kairos/qrcode/"
)

func generateQRCode(unsealKey_B64 string, key_name string) (string, error) {
	var imgFilePath string = fmt.Sprintf("%v/%v.jpeg", unseal_keys_secrets_dir, key_name)
	qrc, err := qrcode.New(fmt.Sprintf(qrcodes_prefix+"%v", unsealKey_B64))
	// qrc, err := qrcode.New(unsealKey_B64)
	if err != nil {
		fmt.Printf("could not generate QRCode: %v", err)
		return imgFilePath, err
	}
	// os.MkdirAll("/tmp/", FileMode)
	os.MkdirAll(unseal_keys_secrets_dir+"/", os.ModePerm)
	w, errCreatingFile := standard.New(imgFilePath)
	if errCreatingFile != nil {
		fmt.Printf("standard.New failed: %v", errCreatingFile)
		return imgFilePath, errCreatingFile
	}
	// save file
	if errSavingImgToFile := qrc.Save(w); errSavingImgToFile != nil {
		fmt.Printf("could not save image: %v", errSavingImgToFile)
		return imgFilePath, errSavingImgToFile
	}

	return imgFilePath, nil
}
