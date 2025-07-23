# NewReleaseTracker-NRT-

## Как пользоваться скриптом?

1. Склонируйте репозиторий на свой компьютер:
```bash
git clone https://github.com/Zaharazov/NewReleaseTracker-NRT-.git
```
2. Укажите в переменной `artists` имена артистов, релизы которых вы хотите отслеживать.
3. Настройте параметры в коде:
```go
// Формат вывода: обычный (false) или упрощённый (true)
const TGFormat = false

// Скачивать обложки альбомов
const ImageDownload = true

// Скачивать музыкальные превью (30 сек.)
const PreviewDownload = true

// Папка для сохранения обложек
const imageFolderName = "images"

// Папка для сохранения превью
const previewFolderName = "previews"
```
4. Запустите скрипт через `go run main.go`.
5. Получите результат анализа.
