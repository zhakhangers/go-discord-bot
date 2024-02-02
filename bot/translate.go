package bot

import (
	"encoding/json"
	"errors"
	"fmt"
	"html"
	"net/http"
	"net/url"
	"strings"
)

// googleTranslateAPI is the base URL for the Google Translate API.
const googleTranslateAPI = "https://www.googleapis.com/language/translate/v2"

// translate() translates text using Google Translate API.
// It returns the detected source language and the translation.
func translate(apiKey, target, text string) (string, string, error) {
	params := url.Values{}
	// Set translation parameters
	params.Set("target", target)
	params.Set("q", text)

	// Encode the API key for the URL
	encodedAPIKey := url.QueryEscape(apiKey)

	// Construct the full URL for the API request
	url := fmt.Sprintf("%s?key=%s", googleTranslateAPI, encodedAPIKey)

	// Create a request body with the encoded parameters
	body := strings.NewReader(params.Encode())

	// Perform an HTTP POST request to the Google Translate API
	resp, err := http.Post(url, "application/x-www-form-urlencoded", body)
	if err != nil {
		return "", "", err
	}
	defer resp.Body.Close()

	// Check the HTTP response status code
	if resp.StatusCode != http.StatusOK {
		var v struct {
			Error struct {
				Code    int
				Message string
			}
		}
		if err := json.NewDecoder(resp.Body).Decode(&v); err != nil {
			return "", "", err
		}
		return "", "", fmt.Errorf("%v: %v", v.Error.Code, v.Error.Message)
	}

	// Decode the JSON response from the API
	var v struct {
		Data struct {
			Translations []struct {
				TranslatedText         string
				DetectedSourceLanguage string
			}
		}
	}
	if err := json.NewDecoder(resp.Body).Decode(&v); err != nil {
		return "", "", err
	}

	// Check if there are translations in the response
	if len(v.Data.Translations) == 0 {
		return "", "", errors.New("no translation")
	}

	// Extract the detected source language and translated text
	source := v.Data.Translations[0].DetectedSourceLanguage
	translated := html.UnescapeString(v.Data.Translations[0].TranslatedText)
	return source, translated, nil
}
