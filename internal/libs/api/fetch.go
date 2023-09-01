package api

import (
	"encoding/json"
	"errors"
	"io"
	"log"
	"net/http"
	"strings"
	"yyyoichi/Collo-API/pkg/apperror"
)

type FetchResult struct {
	SpeachJson *SpeachJson
	Err        error
	URL        string
}

type FetchError struct {
	error
}

func (fr *FetchResult) GetSpeachs() string {
	var result string
	for _, s := range fr.SpeachJson.SpeechRecord {
		result += s.Speech
	}
	return result
}

var (
	ErrUnExpected       = errors.New("unexpected")
	ErrCongestedServer  = errors.New("congested")
	ErrBadRequest       = errors.New("bad request")
	ErrParse            = errors.New("parse")
	ErrExcessCharacters = errors.New("excess characters")
	ErrNoData           = errors.New("no data")
)

type SpeachJson struct {
	NumberOfRecords    int    `json:"numberOfRecords"`
	NumberOfReturn     int    `json:"numberOfReturn"`
	StartRecord        int    `json:"startRecord"`
	NextRecordPosition int    `json:"nextRecordPosition"`
	Message            string `json:"message"`
	SpeechRecord       []struct {
		Speech     string      `json:"speach"`
		OtherField interface{} `json:"-"`
	} `json:"speechRecord"`
	OtherField interface{} `json:"-"`
}

func Fetch(url string) *FetchResult {
	var fetchResult *FetchResult
	fetchResult.URL = url
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		fetchResult.Err = FetchError{apperror.WrapError(err, err.Error())}
		return fetchResult
	}

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		fetchResult.Err = FetchError{apperror.WrapError(err, err.Error())}
		return fetchResult
	}

	defer resp.Body.Close()
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		fetchResult.Err = FetchError{apperror.WrapError(err, err.Error())}
		return fetchResult
	}

	var speachJson *SpeachJson
	if err := json.Unmarshal(body, &speachJson); err != nil {
		fetchResult.Err = FetchError{apperror.WrapError(err, err.Error())}
		return fetchResult
	}
	// regular
	if speachJson.Message == "" {
		fetchResult.SpeachJson = speachJson
		return fetchResult
	}
	// get error
	err = ErrUnExpected
	if strings.Contains(speachJson.Message, "19001") {
		err = ErrCongestedServer
	} else if strings.Contains(speachJson.Message, "19020") {
		err = ErrExcessCharacters
	} else if strings.Contains(speachJson.Message, "19011") {
		err = ErrBadRequest
	}
	log.Printf("Got Result Error Message: %s [%s]\n", err, speachJson.Message)
	fetchResult.Err = FetchError{apperror.WrapError(err, err.Error())}
	return fetchResult
}
