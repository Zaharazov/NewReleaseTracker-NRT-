package main

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"net/url"
	"os"
	"strings"
	"time"
	"unicode"

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
	// A
	"Adam Jensen",
	"Against The Current",
	"AJR",
	"Alec Benjamin",
	"All Time Low",
	"Arrested Youth",
	"Ava Max",
	"AViVA",
	// B
	"BOOKER",
	"Bring Me The Horizon",
	// C
	"Call Me Karizma",
	"Cjbeards",
	"Clean Bandit",
	"Confetti",
	// D
	"Dean Lewis",
	"Deepend",
	"DEIIN",
	"Demi the Daredevil",
	// E
	"Ed Sheeran",
	// F
	"Foo Fighters",
	// G
	"Grandson",
	// H
	"Halsey",
	// I
	"Ichika Nito",
	"Imagine Dragons",
	"ItaloBrothers",
	// J
	"Jake Daniels",
	"Justin Bieber",
	// K
	"Klaas",
	// L
	"Like Saturn",
	"Lucas Estrada",
	"LUNAX",
	// M
	"mgk",
	"Mike Candys",
	"Mike Williams",
	// N
	"New Medicine",
	// O
	"Oliver Tree",
	// P
	"Powfu",
	// Q
	// R
	"Rag'N'Bone Man",
	"Robert Grace",
	"Rudimental",
	// S
	"Silent Child",
	"Simple Plan",
	"Sub Urban",
	// T
	"The Killers",
	"The Hatters",
	"The Hives",
	"The Rasmus",
	"Tom Grennan",
	"Tom Morello",
	"tooboe",
	// U
	"Unlike Pluto",
	// V
	"VOILA",
	// W
	"Whethan",
	// X
	"X Ambassadors",
	// Y
	"Ya Rick",
	// Z
	"Zombie Americana",
}

func main() {
	client := &http.Client{}

	err := godotenv.Load()

	if err != nil {
		log.Fatal("Ошибка загрузки .env файла: ", err)
	}

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

		request, err := http.NewRequest("POST", "https://accounts.spotify.com/api/token", strings.NewReader(data.Encode()))
		if err != nil {
			log.Fatal("Ошибка создания запроса: ", err)
		}
		request.Header.Set("Content-Type", "application/x-www-form-urlencoded")

		response, err := client.Do(request)
		if err != nil {
			log.Fatal("Ошибка получения ответа: ", err)
		}
		body, err := io.ReadAll(response.Body)
		if err != nil {
			log.Fatal("Ошибка получения тела ответа: ", err)
		}

		var responseData map[string]interface{}
		err = json.Unmarshal(body, &responseData)
		if err != nil {
			log.Fatal("Ошибка преобразования тела запроса: ", err)
		}

		accessToken = responseData["access_token"].(string)
		refreshToken = responseData["refresh_token"].(string)

		fmt.Printf("\nAccess Token и Refresh Token установлены. \nПожалуйста, замените ваш ACCESS_TOKEN в env-файле на: %s \nА REFRESH_TOKEN на: %s \n\n", accessToken, refreshToken)
	}

	// Проверяем актуальность токена

	request, err := http.NewRequest("GET", "https://api.spotify.com/v1/me", nil)
	if err != nil {
		log.Fatal("Ошибка создания запроса: ", err)
	}
	request.Header.Set("Authorization", "Bearer "+accessToken)

	response, err := client.Do(request)
	if err != nil {
		log.Fatal("Ошибка получения ответа: ", err)
	}

	if response.StatusCode == 401 {
		data := url.Values{}
		data.Set("grant_type", "refresh_token")
		data.Set("redirect_uri", redirectURI)
		data.Set("client_id", clientID)
		data.Set("client_secret", clientSecret)
		data.Set("refresh_token", refreshToken)

		request, err := http.NewRequest("POST", "https://accounts.spotify.com/api/token", strings.NewReader(data.Encode()))
		if err != nil {
			log.Fatal("Ошибка создания запроса: ", err)
		}
		request.Header.Set("Content-Type", "application/x-www-form-urlencoded")

		response, err := client.Do(request)
		if err != nil {
			log.Fatal("Ошибка получения ответа: ", err)
		}
		body, err := io.ReadAll(response.Body)
		if err != nil {
			log.Fatal("Ошибка получения тела ответа: ", err)
		}

		var responseData map[string]interface{}
		err = json.Unmarshal(body, &responseData)
		if err != nil {
			log.Fatal("Ошибка преобразования тела запроса: ", err)
		}

		accessToken = responseData["access_token"].(string)
		fmt.Printf("Access Token обновлён. \nПожалуйста, замените ваш ACCESS_TOKEN в env-файле на: %s \n\n", accessToken)
	}

	// Анализ артистов

	atLeastOne := false

	for _, artist := range artists {

		// Получаем id артиста

		request, err := http.NewRequest("GET", "https://api.spotify.com/v1/search?q="+url.QueryEscape(artist)+"&type=artist", nil)
		if err != nil {
			log.Println("Ошибка при получении id артиста: ", err)
			continue
		}
		request.Header.Set("Authorization", "Bearer "+accessToken)

		response, err := client.Do(request)
		if err != nil {
			log.Println("Ошибка получения ответа: ", err)
			continue
		}
		body, err := io.ReadAll(response.Body)
		if err != nil {
			log.Println("Ошибка получения тела ответа: ", err)
			continue
		}

		var requestData SearchResponse
		err = json.Unmarshal(body, &requestData)
		if err != nil {
			log.Println("Ошибка преобразования тела запроса: ", err)
			continue
		}

		// Получаем последние релизы артиста

		if len(requestData.Artists.Items) == 0 {
			continue
		}
		request, err = http.NewRequest("GET", "https://api.spotify.com/v1/artists/"+requestData.Artists.Items[0].ID+"/albums", nil)
		if err != nil {
			log.Println("Ошибка создания запроса: ", err)
			continue
		}

		request.Header.Set("Authorization", "Bearer "+accessToken)

		response, err = client.Do(request)
		if err != nil {
			log.Println("Ошибка получения ответа: ", err)
			continue
		}
		body, err = io.ReadAll(response.Body)
		if err != nil {
			log.Println("Ошибка получения тела ответа: ", err)
			continue
		}

		var releaseResponse ReleaseResponse
		err = json.Unmarshal(body, &releaseResponse)
		if err != nil {
			log.Println("Ошибка преобразования тела запроса: ", err)
			continue
		}

		now := time.Now()
		oneWeekAgo := now.AddDate(0, 0, -7)

		for _, release := range releaseResponse.Items {
			releaseDate, err := time.Parse("2006-01-02", release.ReleaseDate)
			if err != nil {
				continue
			}
			if releaseDate.After(oneWeekAgo) && releaseDate.Before(now) {
				atLeastOne = true
				fmt.Printf("------------------------------------------------------------------------------------------------------\n")
				fmt.Printf("|%c| %s - %s (%s)\n|#| Дата выхода: %s \n", unicode.ToUpper(rune((requestData.Artists.Items[0].Name)[0])), requestData.Artists.Items[0].Name, release.Name, release.AlbumType, release.ReleaseDate)
			}
		}
	}
	if atLeastOne {
		fmt.Printf("------------------------------------------------------------------------------------------------------\n")
	}

	os.Exit(0)
}
