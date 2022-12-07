package recognizer

import (
	"errors"
	"fmt"

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
