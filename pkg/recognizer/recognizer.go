package recognizer

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"mime/multipart"
	"net/http"
	"os"
	"path/filepath"
	"time"

	"github.com/openalpr/openalpr/src/bindings/go/openalpr"
)

type RecognizeResult struct {
	Version           int     `json:"version"`
	DataType          string  `json:"data_type"`
	EpochTime         int64   `json:"epoch_time"`
	ImgWidth          int     `json:"img_width"`
	ImgHeight         int     `json:"img_height"`
	ProcessingTimeMs  float64 `json:"processing_time_ms"`
	RegionsOfInterest []struct {
		X      int `json:"x"`
		Y      int `json:"y"`
		Width  int `json:"width"`
		Height int `json:"height"`
	} `json:"regions_of_interest"`
	Results []struct {
		Plate            string  `json:"plate"`
		Confidence       float64 `json:"confidence"`
		MatchesTemplate  int     `json:"matches_template"`
		PlateIndex       int     `json:"plate_index"`
		Region           string  `json:"region"`
		RegionConfidence int     `json:"region_confidence"`
		ProcessingTimeMs float64 `json:"processing_time_ms"`
		RequestedTopn    int     `json:"requested_topn"`
		Coordinates      []struct {
			X int `json:"x"`
			Y int `json:"y"`
		} `json:"coordinates"`
		Candidates []struct {
			Plate           string  `json:"plate"`
			Confidence      float64 `json:"confidence"`
			MatchesTemplate int     `json:"matches_template"`
		} `json:"candidates"`
	} `json:"results"`
}

type RecognizeRes struct {
	ProcessingTime float64 `json:"processing_time"`
	Results        []struct {
		Box struct {
			Xmin int `json:"xmin"`
			Ymin int `json:"ymin"`
			Xmax int `json:"xmax"`
			Ymax int `json:"ymax"`
		} `json:"box"`
		Plate  string `json:"plate"`
		Region struct {
			Code  string  `json:"code"`
			Score float64 `json:"score"`
		} `json:"region"`
		Score      float64 `json:"score"`
		Candidates []struct {
			Score float64 `json:"score"`
			Plate string  `json:"plate"`
		} `json:"candidates"`
		Dscore  float64 `json:"dscore"`
		Vehicle struct {
			Score float64 `json:"score"`
			Type  string  `json:"type"`
			Box   struct {
				Xmin int `json:"xmin"`
				Ymin int `json:"ymin"`
				Xmax int `json:"xmax"`
				Ymax int `json:"ymax"`
			} `json:"box"`
		} `json:"vehicle"`
	} `json:"results"`
	Filename  string      `json:"filename"`
	Version   int         `json:"version"`
	CameraID  interface{} `json:"camera_id"`
	Timestamp time.Time   `json:"timestamp"`
}

func DirectRecognize(photoLocation string) (res openalpr.AlprResults, err error) {
	alpr := openalpr.NewAlpr("id", "", "/usr/share/openalpr/runtime_data")
	defer alpr.Unload()

	if !alpr.IsLoaded() {
		fmt.Println("OpenAlpr failed to load!")
		err = errors.New("failed to load openalr")
		return
	}
	alpr.SetTopN(20)

	res, err = alpr.RecognizeByFilePath(photoLocation)
	if err != nil {
		return
	}

	return
}

func Recognize(photoLocation, token string) (res *RecognizeRes, err error) {
	client := &http.Client{}

	file, err := os.Open(photoLocation)
	if err != nil {
		return
	}
	defer file.Close()

	body := &bytes.Buffer{}
	writer := multipart.NewWriter(body)
	part, err := writer.CreateFormFile("upload", filepath.Base(photoLocation))
	if err != nil {
		return
	}
	_, err = io.Copy(part, file)

	_ = writer.WriteField("regions", "id")
	_ = writer.WriteField("config", `{"region":"strict"}`)
	err = writer.Close()
	if err != nil {
		return
	}

	req, err := http.NewRequest("POST", "https://api.platerecognizer.com/v1/plate-reader/", body)
	if err != nil {
		return
	}
	req.Header.Set("Content-Type", writer.FormDataContentType())
	req.Header.Set("Authorization", "Token "+token)

	resp, err := client.Do(req)
	if err != nil {
		return
	}
	defer resp.Body.Close()

	if err = json.NewDecoder(resp.Body).Decode(&res); err != nil {
		return
	}

	if resp.StatusCode != http.StatusCreated {
		err = fmt.Errorf("bad status: %s %+v", resp.Status, res)
		return
	}

	return
}
