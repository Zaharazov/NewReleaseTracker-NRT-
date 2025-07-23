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
	"strconv"
	"time"
	"unicode"
)

type ArtistSearchResponse struct {
	Data []struct {
		ID   int    `json:"id"`
		Name string `json:"name"`
	} `json:"data"`
	Total int `json:"total"`
}

type AlbumSearchResponse struct {
	Data []struct {
		ID          int    `json:"id"`
		Title       string `json:"title"`
		Image       string `json:"cover_xl"`
		ReleaseDate string `json:"release_date"`
		RecordType  string `json:"record_type"`
		Tracklist   string `json:"tracklist"`
	} `json:"data"`
	Total int `json:"total"`
}

type TracklistSearchResponse struct {
	Data []struct {
		ID      int    `json:"id"`
		Title   string `json:"title"`
		Preview string `json:"preview"`
	} `json:"data"`
	Total int `json:"total"`
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

	filePath := filepath.Join(imageFolderName, releaseName+".jpg")
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

func SavePreviewFromRelease(previewURL, releaseFolder, releaseName string, client *http.Client) error {
	request, err := http.NewRequest("GET", previewURL, nil)

	if err != nil {
		log.Println("Ошибка получения превью: ", err)
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
	releaseFolder = re.ReplaceAllString(releaseFolder, "")

	releaseFolderPath := filepath.Join(previewFolderName, releaseFolder)

	err = os.MkdirAll(releaseFolderPath, os.ModePerm)
	if err != nil {
		log.Fatal("Ошибка при создании папки: ", err)
	}

	filePath := filepath.Join(previewFolderName, releaseFolder, releaseName+".m4a")
	file, err := os.Create(filePath)
	if err != nil {
		log.Println("Ошибка создания файла с превью: ", err)
		return err
	}
	defer file.Close()

	_, err = io.Copy(file, response.Body)
	if err != nil {
		log.Println("Ошибка загрузки превью в файл: ", err)
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

const TGFormat = true
const ImageDownload = true
const PreviewDownload = true
const imageFolderName = "images"
const previewFolderName = "previews"

func main() {
	client := &http.Client{}

	// Создаем хранилище для изображений

	_, err := os.Stat(imageFolderName)

	if err == nil {
		err = os.RemoveAll(imageFolderName)
		if err != nil {
			log.Fatal("Ошибка удаления старых данных: ", err)
		}
	}

	err = os.MkdirAll(imageFolderName, os.ModePerm)
	if err != nil {
		log.Fatal("Ошибка при создании папки: ", err)
	}

	// Создаем хранилище для превью

	_, err = os.Stat(previewFolderName)

	if err == nil {
		err = os.RemoveAll(previewFolderName)
		if err != nil {
			log.Fatal("Ошибка удаления старых данных: ", err)
		}
	}

	err = os.MkdirAll(previewFolderName, os.ModePerm)
	if err != nil {
		log.Fatal("Ошибка при создании папки: ", err)
	}

	// Анализ артистов

	atLeastOne := false

	for _, artist := range artists {

		// Получаем id артиста

		request, err := http.NewRequest("GET", "https://api.deezer.com/search/artist?q="+url.QueryEscape(artist), nil)
		if err != nil {
			log.Println("Ошибка создания запроса: ", err)
			continue
		}

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

		var artistResponseData ArtistSearchResponse
		err = json.Unmarshal(body, &artistResponseData)
		if err != nil {
			log.Println("Ошибка преобразования тела запроса: ", err)
			continue
		}

		// Получаем последние релизы артиста

		if artistResponseData.Total == 0 {
			continue
		}
		request, err = http.NewRequest("GET", "https://api.deezer.com/artist/"+strconv.Itoa(artistResponseData.Data[0].ID)+"/albums?index=0&limit=100", nil)
		if err != nil {
			log.Println("Ошибка создания запроса: ", err)
			continue
		}

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

		var albumResponseData AlbumSearchResponse
		err = json.Unmarshal(body, &albumResponseData)
		if err != nil {
			log.Println("Ошибка преобразования тела запроса: ", err)
			continue
		}

		now := time.Now()
		oneWeekAgo := now.AddDate(0, 0, -7)

		for _, release := range albumResponseData.Data {
			releaseDate, err := time.Parse("2006-01-02", release.ReleaseDate)
			if err != nil {
				continue
			}
			if releaseDate.After(oneWeekAgo) && releaseDate.Before(now) {
				if TGFormat {
					fmt.Printf("• %s - %s\n", artistResponseData.Data[0].Name, release.Title)
				} else {
					atLeastOne = true
					fmt.Printf("------------------------------------------------------------------------------------------------------\n")
					fmt.Printf("|%c| %s - %s (%s)\n|#| Дата выхода: %s \n", unicode.ToUpper(rune((artistResponseData.Data[0].Name)[0])), artistResponseData.Data[0].Name, release.Title, release.RecordType, release.ReleaseDate)
				}

				// Получаем изображение релиза

				if ImageDownload {
					imageURL := release.Image
					err := SaveImageFromRelease(imageURL, release.Title, client)

					if err != nil {
						continue
					}
				}

				// Получаем превью релиза

				if PreviewDownload {
					request, err := http.NewRequest("GET", release.Tracklist+"?limit=100", nil)
					if err != nil {
						log.Println("Ошибка создания запроса: ", err)
						continue
					}

					response, err := client.Do(request)
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

					var tracklistResponseData TracklistSearchResponse
					err = json.Unmarshal(body, &tracklistResponseData)
					if err != nil {
						log.Println("Ошибка преобразования тела запроса: ", err)
						continue
					}

					for i := 0; i < tracklistResponseData.Total; i++ {
						track := tracklistResponseData.Data[i]
						err = SavePreviewFromRelease(track.Preview, release.Title, track.Title, client)
					}
				}
			}
		}
	}
	if atLeastOne {
		fmt.Printf("------------------------------------------------------------------------------------------------------\n")
	}

	os.Exit(0)
}
