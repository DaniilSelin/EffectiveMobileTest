/*
Пакет для взаимодействие с внешними API

Не стал оборачивать в супертип, и использовать интерейс на стороне сервисов,
Напрямую использую функции
*/

package enrichment

import (
	"encoding/json"
	"fmt"
	"time"
	"net/http"
)

func GetAge(name string) (int, error) {
	client := &http.Client{
		Timeout: 5 * time.Second,
	}

	url := fmt.Sprintf("https://api.agify.io/?name=%s", name)
	resp, err := client.Get(url)
	if err != nil {
		return 0, err
	}
	defer resp.Body.Close()

	var result struct {
		Age int `json:"age"`
	}

	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return 0, err
	}
	return result.Age, nil
}

func GetGender(name string) (string, error) {
	client := &http.Client{
		Timeout: 5 * time.Second, 
	}

	url := fmt.Sprintf("https://api.agify.io/?name=%s", name)
	resp, err := client.Get(url)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	var result struct {
		Gender string `json:"gender"`
	}

	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return "", err
	}

	return result.Gender, nil
}

func GetNationality(name string) (string, error) {
	client := &http.Client{
		Timeout: 5 * time.Second, 
	}

	url := fmt.Sprintf("https://api.agify.io/?name=%s", name)
	resp, err := client.Get(url)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	var result struct {
		Country []struct {
			CountryID string `json:"country_id"`
		} `json:"country"`
	}

	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return "", err
	}

	if len(result.Country) > 0 {
		return result.Country[0].CountryID, nil
	}

	return "unknown", nil
}
