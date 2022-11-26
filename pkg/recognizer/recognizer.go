package recognizer

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"os/exec"
	"path"

	"github.com/openalpr/openalpr/src/bindings/go/openalpr"
	"github.com/vincentwijaya/go-ocr/pkg/log"
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

func Recognize(photoLocation string) (res RecognizeResult, err error) {
	pwd, _ := os.Getwd()
	log.Info("docker", "run", "--rm", "-v", pwd+"/files/images/vehicle:/data:ro", "openalpr/openalpr", "-j", "-c", "sg", path.Base(photoLocation))
	cmd := exec.Command("docker", "run", "--rm", "--platform", "linux/amd64", "-v", pwd+"/files/images/vehicle:/data:ro", "openalpr/openalpr", "-j", "-c", "sg", path.Base(photoLocation))
	var stdout bytes.Buffer
	var stderr bytes.Buffer
	cmd.Stdout = &stdout
	cmd.Stderr = &stderr
	err = cmd.Run()
	if err != nil {
		err = errors.New(stderr.String())
		return
	}

	if err = json.Unmarshal(stdout.Bytes(), &res); err != nil {
		return
	}

	return
}

func DirectRecognize(photoLocation string) (res openalpr.AlprResults, err error) {
	alpr := openalpr.NewAlpr("sg", "", "/usr/share/openalpr/runtime_data")
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
