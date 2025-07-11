package main

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"os"
	"strings"
	"time"

	"github.com/joho/godotenv"
)

type SearchResponse struct {
	Artists struct {
		Items []struct {
			ID   string `json:"id"`
			Name string `json:"name"`
		} `json:"items"`
	} `json:"artists"`
}

type ReleaseResponse struct {
	Items []struct {
		AlbumType   string `json:"album_type"`
		ID          string `json:"id"`
		Name        string `json:"name"`
		ReleaseDate string `json:"release_date"`
	} `json:"items"`
}

var artists = []string{
	"Imagine Dragons",
	"Call Me Karizma",
	"Arrested Youth",
	"DEIIN",
	"Ed Sheeran",
	"Ichika Nito",
	"BOOKER",
	"VOILA",
	"Adam Jensen",
	"Alec Benjamin",
	"Grandson",
}

func main() {

	analysisResults := make(map[string][]string)
	client := &http.Client{}

	_ = godotenv.Load()

	redirectURI := os.Getenv("REDIRECT_URI")
	clientID := os.Getenv("CLIENT_ID")
	clientSecret := os.Getenv("CLIENT_SECRET")

	accessToken := os.Getenv("ACCESS_TOKEN")
	refreshToken := os.Getenv("REFRESH_TOKEN")

	// Для тех, кто README не читал

	if redirectURI == "" || clientID == "" || clientSecret == "" {
		fmt.Printf("Необходимо заполнить обязательные поля в env-файле:\nredirectURI\nclientID\nclientSecret\n")
		fmt.Printf("Если не знаете, откуда их взять, загляните в README.md\n")
		os.Exit(1)
	}

	// Для новых пользователей

	if accessToken == "" || refreshToken == "" {

		fmt.Println("Для получения кода авторизации перейдите по ссылке ниже, скопируйте код из преобразованной ссылки и введите его в консоль.")
		fmt.Println("https://accounts.spotify.com/authorize?client_id=" + clientID + "&response_type=code&redirect_uri=" + url.QueryEscape(redirectURI) + "&scope=user-read-email&state=test123")

		var code string
		fmt.Printf("\nВведите ваш code: \n")
		fmt.Scanln(&code)

		data := url.Values{}
		data.Set("grant_type", "authorization_code")
		data.Set("redirect_uri", redirectURI)
		data.Set("client_id", clientID)
		data.Set("client_secret", clientSecret)
		data.Set("code", code)

		request, _ := http.NewRequest("POST", "https://accounts.spotify.com/api/token", strings.NewReader(data.Encode()))
		request.Header.Set("Content-Type", "application/x-www-form-urlencoded")

		response, _ := client.Do(request)
		body, _ := io.ReadAll(response.Body)

		var responseData map[string]interface{}
		_ = json.Unmarshal(body, &responseData)

		accessToken = responseData["access_token"].(string)
		refreshToken = responseData["refresh_token"].(string)

		fmt.Printf("\nAccess Token и Refresh Token установлены. \nПожалуйста, замените ваш ACCESS_TOKEN в env-файле на: %s \nА REFRESH_TOKEN на: %s \n\n", accessToken, refreshToken)
	}

	// Проверяем актуальность токена

	request, _ := http.NewRequest("GET", "https://api.spotify.com/v1/me", nil)
	request.Header.Set("Authorization", "Bearer "+accessToken)

	response, _ := client.Do(request)

	if response.StatusCode == 401 {
		data := url.Values{}
		data.Set("grant_type", "refresh_token")
		data.Set("redirect_uri", redirectURI)
		data.Set("client_id", clientID)
		data.Set("client_secret", clientSecret)
		data.Set("refresh_token", refreshToken)

		request, _ := http.NewRequest("POST", "https://accounts.spotify.com/api/token", strings.NewReader(data.Encode()))
		request.Header.Set("Content-Type", "application/x-www-form-urlencoded")

		response, _ := client.Do(request)
		body, _ := io.ReadAll(response.Body)

		var responseData map[string]interface{}
		_ = json.Unmarshal(body, &responseData)

		accessToken = responseData["access_token"].(string)
		fmt.Printf("Access Token обновлён. \nПожалуйста, замените ваш ACCESS_TOKEN в env-файле на: %s \n\n", accessToken)
	}

	// Анализ артистов

	for _, artist := range artists {

		// Получаем id артиста

		request, _ := http.NewRequest("GET", "https://api.spotify.com/v1/search?q="+url.QueryEscape(artist)+"&type=artist", nil)
		request.Header.Set("Authorization", "Bearer "+accessToken)

		response, _ := client.Do(request)
		body, _ := io.ReadAll(response.Body)

		var requestData SearchResponse
		_ = json.Unmarshal(body, &requestData)

		fmt.Printf("Анализ артиста: %s \n", requestData.Artists.Items[0].Name)

		// Получаем последние релизы артиста

		request, _ = http.NewRequest("GET", "https://api.spotify.com/v1/artists/"+requestData.Artists.Items[0].ID+"/albums", nil)

		request.Header.Set("Authorization", "Bearer "+accessToken)

		response, _ = client.Do(request)
		body, _ = io.ReadAll(response.Body)

		var releaseResponse ReleaseResponse
		_ = json.Unmarshal(body, &releaseResponse)

		now := time.Now()
		oneWeekAgo := now.AddDate(0, 0, -7)

		for _, release := range releaseResponse.Items {
			releaseDate, _ := time.Parse("2006-01-02", release.ReleaseDate)
			if releaseDate.After(oneWeekAgo) && releaseDate.Before(now) {
				fmt.Printf("Найден новый релиз: %s (%s). Дата выхода: %s \n", release.Name, release.AlbumType, release.ReleaseDate)
				analysisResults[requestData.Artists.Items[0].Name] = append(analysisResults[requestData.Artists.Items[0].Name], release.Name+" ("+release.AlbumType+")")
			}
		}
		fmt.Printf("\n")
	}

	fmt.Println("Новые релизы найдены у следующих артистов:")
	for artist, _ := range analysisResults {
		fmt.Println(artist)
	}
	os.Exit(0)
}
