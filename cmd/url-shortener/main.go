package main

import (
	"fmt"
	"library/internal/config"
	"library/internal/storage/sqlite"
	"log/slog"
	"os"
)

const (
	/* envLocal — режим локальной разработки на твоём компьютере.
	Обычно здесь самые подробные логи и простые настройки для отладки. */
	envLocal = "local"

	/* envDev — режим для тестового/дев-окружения (не прод).
	Часто используют JSON-логи и настройки, близкие к боевым. */
	envDev = "dev"

	/* envProd — боевое окружение для реальных пользователей.
	Обычно меньше "шума" в логах и более строгие рабочие настройки. */
	envProd = "prod"
)

func main() {
	/* Загружаем конфигурацию приложения.
	MustLoad читает настройки (yaml/env) и завершает программу,
	если конфиг некорректный или отсутствует. */
	cfg := config.MustLoad()

	/* Временно печатаем весь конфиг в консоль.
	Полезно на этапе обучения/отладки, чтобы убедиться,
	что значения реально подхватились. */
	fmt.Println(cfg)

	/* Создаём логгер на основе окружения из конфига
	(local/dev/prod), чтобы включить нужный формат и уровень логов. */
	log := setupLogger(cfg.Env)

	/* Пишем информационный лог о старте приложения.
	slog.String("env", cfg.Env) добавляет структурированное поле env=... */
	log.Info("starting url-shortener", slog.String("env", cfg.Env))

	/* Пишем тестовое debug-сообщение.
	Оно появится в логах только если у логгера включён уровень Debug
	(например в local/dev, в зависимости от твоей настройки). */
	log.Debug("debug messages are enabled")

	/* Создаём хранилище (подключаемся к БД) по пути из конфига.
	Если база не существует, SQLite обычно создаст файл автоматически. */
	storage, err := sqlite.New(cfg.StoragePath)
	if err != nil {
		/* Если инициализировать хранилище не удалось,
		пишем ошибку и завершаем приложение, потому что без БД работать нельзя. */
		log.Error("failed to init storage")
		os.Exit(1)
	}
	
}

/* setupLogger создаёт и возвращает логгер.
Логгер — это "инструмент", через который программа пишет сообщения в консоль:
что происходит, где ошибка, в каком режиме запущено приложение. */

func setupLogger(env string) *slog.Logger {

	/* Создаём переменную, куда положим готовый логгер.
	Пока она пустая (nil), дальше заполним её в switch. */
	var log *slog.Logger

	/* Проверяем, в каком окружении запущено приложение
	(например local, dev, prod). */

	switch env {

	/* Если приложение запущено в local:
	делаем "читаемые" текстовые логи и включаем самый подробный уровень Debug.
	Это удобно, когда ты разрабатываешь и хочешь видеть максимум деталей. */

	case envLocal:
		log = slog.New(
			slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelDebug}),
		)

		/* Если запущено в dev:
		   используем JSON-логи (их удобно парсить инструментами)
		   и тоже оставляем уровень Debug, чтобы видеть подробности. */
	case envDev:
		log = slog.New(
			slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelDebug}),
		)

		/* Если запущено в prod (боевой режим):
		   тоже JSON-логи, но уровень Info.
		   Info = меньше "шума", только важные рабочие сообщения. */
	case envProd:
		log = slog.New(
			slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelInfo}),
		)
	}
	return log
}
