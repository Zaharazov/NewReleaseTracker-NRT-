package main

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"net/url"
	"os"
	"path/filepath"
	"regexp"
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
		Images      []struct {
			URL string `json:"url"`
		} `json:"images"`
	} `json:"items"`
}

func SaveImageFromRelease(imageURL, releaseName string, client *http.Client) error {
	request, err := http.NewRequest("GET", imageURL, nil)

	if err != nil {
		log.Println("Ошибка получения изображения: ", err)
		return err
	}

	response, err := client.Do(request)

	if err != nil {
		log.Println("Ошибка получения ответа: ", err)
		return err
	}
	defer response.Body.Close()

	re := regexp.MustCompile(`[<>:"/\\|?*]`)
	releaseName = re.ReplaceAllString(releaseName, "")

	filePath := filepath.Join(folderName, releaseName+".jpg")
	file, err := os.Create(filePath)
	if err != nil {
		log.Println("Ошибка создания файла с изображением: ", err)
		return err
	}
	defer file.Close()

	_, err = io.Copy(file, response.Body)
	if err != nil {
		log.Println("Ошибка загрузки изображения в файл: ", err)
		return err
	}
	return nil
}

var artists = []string{
	// A
	"Adam Jensen",
	"Against The Current",
	"AJR",
	"Alan Walker",
	"Alec Benjamin",
	"Alle Farben",
	"All Good Things",
	"All Time Low",
	"Alok",
	"Alter.",
	"angelbaby",
	"Ariana Grande",
	"Arrested Youth",
	"Arrows in Actions",
	"Asking Alexandria",
	"AURORA",
	"AVAION",
	"Ava Max",
	"Avicii",
	"AViVA",
	// B
	"Badflower",
	"Barns Courtney",
	"Bastille",
	"BB Cooper",
	"bbno$",
	"Bear Ghost",
	"Bella Poarch",
	"BENNE",
	"Besomorph",
	"blackbear",
	"Blanks",
	"BOOKER",
	"BoyWithUke",
	"Brainstorm",
	"Bring Me The Horizon",
	"BRKN LOVE",
	"Bryce Fox",
	"Burn The Ballroom",
	"bulow",
	// C
	"Call Me Karizma",
	"Cheat Codes",
	"Christian French",
	"Cjbeards",
	"Clean Bandit",
	"COIN",
	"Coldplay",
	"Confetti",
	"Connor Kauffman",
	"СТИНТ",
	// D
	"Daya",
	"Dead Poet Society",
	"Dean Lewis",
	"Deepend",
	"Deep Thrills",
	"DEIIN",
	"Demi the Daredevil",
	"Denis First",
	"Des Rocs",
	"DNMO",
	"Don Diablo",
	"Dotan",
	"Dreamers",
	// E
	"Ed Sheeran",
	"Elderbrook",
	"Ellie Gouldling",
	"Ethan Bortnick",
	// F
	"Fall Out Boy",
	"Fitz and The Tantrums",
	"Foo Fighters",
	"Foreign Air",
	// G
	"Good Kid",
	"Grabbitz",
	"Grandson",
	"8 Graves",
	"Green Day",
	// H
	"Halsey",
	"half•alive",
	"Hollywood Undead",
	"Hozier",
	"Hurts",
	// I
	"Ichika Nito",
	"I DONT KNOW HOW BUT THEY FOUND ME",
	"Imagine Dragons",
	"I'm Geist",
	"In Her Own Words",
	"I Prevail",
	"Islandis",
	"ItaloBrothers",
	// J
	"Jagwar Twin",
	"Jake Daniels",
	"James Arthur",
	"Janieck",
	"JAOVA",
	"Jonas Brothers",
	"Joywave",
	"Justin Bieber",
	// K
	"KALEO",
	"Katy Perry",
	"Kesha",
	"K.Flay",
	"Khalid",
	"Klaas",
	"Koste",
	"KSHMR",
	"Kygo",
	// L
	"Layto",
	"Lady Gaga",
	"Lemaitre",
	"Letdown.",
	"Lewis Capaldi",
	"Like Saturn",
	"Linkin Park",
	"Little Big",
	"Lost Frequencies",
	"Lucas Estrada",
	"LUNAX",
	"Lyre le temps",
	// M
	"Mahmood",
	"Marnik",
	"Maroon 5",
	"Marshmello",
	"Max Leone",
	"Meltt",
	"mgk",
	"Michael Patrick Kelly",
	"Michael Shynes",
	"Mike Candys",
	"Mike Williams",
	"Milky Chance",
	"MISSIO",
	"Mosimann",
	"Mother Mother",
	"My Chemical Romance",
	"Myles Smith",
	"Mzlff",
	"Мукка",
	// N
	"Nessa Barrett",
	"New Hope Club",
	"New Medicine",
	"Nico & Chelsea",
	"Nico Collins",
	"Nico Santos",
	"Noan Kahan",
	"Noize MC",
	"NOTHING MORE",
	// O
	"Odd Chap",
	"Oliver Tree",
	"One Republic",
	// P
	"Paper Idol",
	"Peter McPoland",
	"Phem",
	"Pierce The Vell",
	"Placebo",
	"Portugal. The Man",
	"Powfu",
	"pyrokinesis",
	// Q
	// R
	"Rag'N'Bone Man",
	"Rare Americans",
	"Ricky Montgomery",
	"Robert Grace",
	"Roe Kapara",
	"Rompasso",
	"Royal Republic",
	"Royel Otis",
	"Rudimental",
	"RYYZN",
	// S
	"Saint Motel",
	"Sam Feldt",
	"Scissors Sisters",
	"Sea Girls",
	"Selena Gomez",
	"Self Deception",
	"Set It Off",
	"SHANGUY",
	"Shawn Mendes",
	"SIAMES",
	"Silent Child",
	"Simple Plan",
	"Skillet",
	"Sleeping Wolf",
	"Sofi Tukker",
	"Sting",
	"Sub Urban",
	// T
	"Teddy Swims",
	"The Band CAMINO",
	"The Black Keys",
	"The Blue Stones",
	"The Chalkeaters",
	"The Hatters",
	"The Hives",
	"The Kid LAROI",
	"The Killers",
	"The Living Tombstone",
	"The Maine",
	"The Midnight",
	"The People's Thieves",
	"The Rasmus",
	"The Score",
	"The Unlikely Candidates",
	"The Used",
	"The Weekend",
	"The Wrecks",
	"Three Days Grace",
	"Tokio Hotel",
	"Tokio Project",
	"Tom Gregory",
	"Tom Grennan",
	"Tom Morello",
	"Tommy Cash",
	"Tom Odell",
	"Tom Walker",
	"tooboe",
	"Trevor Daniel",
	"Trinix",
	"twenty one pilots",
	"Two Feet",
	"Три Дня Дождя",
	// U
	"Unlike Pluto",
	"updog",
	// V
	"Vanotek",
	"Vicetone",
	"VOILA",
	// W
	"WALK THE MOON",
	"We Are Scientists",
	"Weathers",
	"Whethan",
	"Will Jay",
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

const TGFormat = false
const folderName = "images"

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

	// Создаем хранилище для изображений

	_, err = os.Stat(folderName)

	if err == nil {
		err = os.RemoveAll(folderName)
		if err != nil {
			log.Fatal("Ошибка удаления старых данных: ", err)
		}
	}

	err = os.MkdirAll(folderName, os.ModePerm)
	if err != nil {
		log.Fatal("Ошибка при создании папки: ", err)
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
				if TGFormat {
					fmt.Printf("• %s - %s\n", requestData.Artists.Items[0].Name, release.Name)
				} else {
					atLeastOne = true
					fmt.Printf("------------------------------------------------------------------------------------------------------\n")
					fmt.Printf("|%c| %s - %s (%s)\n|#| Дата выхода: %s \n", unicode.ToUpper(rune((requestData.Artists.Items[0].Name)[0])), requestData.Artists.Items[0].Name, release.Name, release.AlbumType, release.ReleaseDate)
				}

				// Получаем изображение релиза

				imageURL := release.Images[0].URL
				err := SaveImageFromRelease(imageURL, release.Name, client)

				if err != nil {
					continue
				}
			}
		}
	}
	if atLeastOne {
		fmt.Printf("------------------------------------------------------------------------------------------------------\n")
	}

	os.Exit(0)
}
