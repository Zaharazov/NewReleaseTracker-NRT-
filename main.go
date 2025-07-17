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
	"Airways",
	"AJR",
	"Alan Walker",
	"Alec Benjamin",
	"All Time Low",
	"Arrested Youth",
	"Asking Alexandria",
	"AURORA",
	"Ava Max",
	"AViVA",
	// B
	"Badflower",
	"Barns Courtney",
	"bbno$",
	"blackbear",
	"BOOKER",
	"Brainstorm",
	"Bring Me The Horizon",
	"Bryce Fox",
	"Burn The Ballroom",
	"bulow",
	// C
	"Call Me Karizma",
	"Christian French",
	"Cjbeards",
	"Clean Bandit",
	"COIN",
	"Confetti",
	// D
	"Dean Lewis",
	"Deepend",
	"DEIIN",
	"Demi the Daredevil",
	"Des Rocs",
	"DNMO",
	"Don Diablo",
	"Dreamers",
	// E
	"Ed Sheeran",
	"Elderbrook",
	"Ethan Bortnick",
	// F
	"Fall Out Boy",
	"Foo Fighters",
	// G
	"Grandson",
	"8 Graves",
	// H
	"Halsey",
	"half•alive",
	"Hozier",
	"Hurts",
	// I
	"Ichika Nito",
	"I DONT KNOW HOW BUT THEY FOUND ME",
	"Imagine Dragons",
	"I Prevail",
	"Islandis",
	"ItaloBrothers",
	// J
	"Jagwar Twin",
	"Jake Daniels",
	"Janieck",
	"JAOVA",
	"Justin Bieber",
	// K
	"K.Flay",
	"Klaas",
	"Koste",
	"KSHMR",
	// L
	"Lemaitre",
	"Like Saturn",
	"Lucas Estrada",
	"LUNAX",
	"Lyre le temps",
	// M
	"Mahmood",
	"Marnik",
	"Max Leone",
	"Meltt",
	"mgk",
	"Michael Patrick Kelly",
	"Mike Candys",
	"Mike Williams",
	"Mosimann",
	"My Chemical Romance",
	// N
	"New Hope Club",
	"New Medicine",
	"Nico Collins",
	"Noan Kahan",
	"Noize MC",
	// O
	"Oliver Tree",
	"One Republic",
	// P
	"Peter McPoland",
	"Phem",
	"Pierce The Vell",
	"Placebo",
	"Powfu",
	// Q
	// R
	"Rag'N'Bone Man",
	"Rare Americans",
	"Ricky Montgomery",
	"Robert Grace",
	"Rompasso",
	"Rudimental",
	// S
	"Saint Motel",
	"Sea Girls",
	"Set It Off",
	"SHANGUY",
	"Shawn Mendes",
	"SIAMES",
	"Silent Child",
	"Simple Plan",
	"Sleeping Wolf",
	"Sofi Tukker",
	"Sub Urban",
	// T
	"The Blue Stones",
	"The Hatters",
	"The Hives",
	"The Killers",
	"The Maine",
	"The People's Thieves",
	"The Rasmus",
	"The Score",
	"The Used",
	"The Wrecks",
	"Tom Gregory",
	"Tom Grennan",
	"Tom Morello",
	"tooboe",
	"Trevor Daniel",
	"twenty one pilots",
	"Two Feet",
	// U
	"Unlike Pluto",
	// V
	"VOILA",
	// W
	"WALK THE MOON",
	"Whethan",
	// X
	"X Ambassadors",
	// Y
	"Ya Rick",
	"Years & Years",
	"Yoshi Flower",
	"YUNGBLUD",
	// Z
	"Zero 9:36",
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

		defer response.Body.Close()

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

	defer response.Body.Close()

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

		defer response.Body.Close()

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

		defer response.Body.Close()

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

		defer response.Body.Close()

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
