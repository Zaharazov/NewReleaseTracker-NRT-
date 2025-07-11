# NewReleaseTracker-NRT-

## Перед использованием скрипта

1. Перейдите на [Spotify Developer Dashboard](https://developer.spotify.com/dashboard), войдите в свою учетную запись.
2. На главной странице нажмите `Create app`.
3. Введите произвольные:
* `App name`
* `App description`
* `Redirect URIs` (рекомендую указывать `http://127.0.0.1:8080/callback`)
4. В пункте `Which API/SDKs are you planning to use?` отметьте `Web API`.
5. Сохраните dashboard.
6. После обновления страницы появятся `Client ID` и `Client secret`, вместе с `Redirect URIs` укажите их в env-файле:
```
REDIRECT_URI=http://127.0.0.1:8080/callback
CLIENT_ID=abcdefabcdefabcdefabc1234567890
CLIENT_SECRET=abcdefabcdefabcdefabc1234567890

ACCESS_TOKEN=
REFRESH_TOKEN=
```
7. Сохраните env-файл, запустите скрипт через `go run main.go`, он подскажет, как получить `ACCESS_TOKEN` и `REFRESH_TOKEN`.
8. Получите результат анализа.