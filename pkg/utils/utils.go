package utils

import (
	"bytes"
	"encoding/base64"
	"errors"
	"fmt"
	"image/png"
	"io"
	"net/http"
	"os"
)

func Base64ToPNG(filename string, data string) error {
	fmt.Println(data)
	unbased, err := base64.StdEncoding.DecodeString(data)
	if err != nil {
		panic("Cannot decode b64")
	}

	r := bytes.NewReader(unbased)
	im, err := png.Decode(r)
	if err != nil {
		return errors.New("bad PNG")
	}

	f, err := os.OpenFile(fmt.Sprintf("%s.png", filename), os.O_WRONLY|os.O_CREATE, 0777)
	if err != nil {
		return errors.New("can't open file")
	}

	png.Encode(f, im)
	return nil
}

func RemoveLocalFile(fileName string) error {
	e := os.Remove(fileName)
	return e
}

func DownloadFile(URL, fileName string) error {
	response, err := http.Get(URL)
	if err != nil {
		return err
	}
	defer response.Body.Close()

	if response.StatusCode != 200 {
		return errors.New("received non 200 response code")
	}

	file, err := os.Create(fileName)
	if err != nil {
		return err
	}
	defer file.Close()

	_, err = io.Copy(file, response.Body)
	if err != nil {
		return err
	}

	return nil
}
