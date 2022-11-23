package utils

import (
	"bytes"
	"encoding/base64"
	"errors"
	"fmt"
	"image/png"
	"os"
)

func Base64ToPNG(filename string, data string) error {
	unbased, err := base64.StdEncoding.DecodeString(data)
	if err != nil {
		panic("Cannot decode b64")
	}

	r := bytes.NewReader(unbased)
	im, err := png.Decode(r)
	if err != nil {
		return errors.New("Bad PNG")
	}

	f, err := os.OpenFile(fmt.Sprintf("%s.png", filename), os.O_WRONLY|os.O_CREATE, 0777)
	if err != nil {
		return errors.New("Can't open file")
	}

	png.Encode(f, im)
	return nil
}

func RemoveLocalFile(fileName string) error {
	e := os.Remove(fileName)
	return e
}
