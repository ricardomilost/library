package config

import (
	"log"
	"os"
	"time"

	"github.com/ilyakaznacheev/cleanenv"
)

/*
	Config — это "главный контейнер" всех настроек приложения.

Когда приложение стартует, оно читает yaml-файл и кладёт значения в эту структуру.
*/
type Config struct {
	/* Env — окружение, в котором запущено приложение.
	Обычно: local / dev / prod.
	yaml:"env" -> брать значение из поля "env" в yaml.
	env-default:"local" -> если "env" не указано, автоматически взять "local". */
	Env string `yaml:"env" env-default:"local"`

	/* StoragePath — путь к файлу/хранилищу данных (например, к .db файлу).
	yaml:"storage_path" -> связывает поле с ключом "storage_path" в yaml.
	env-required:"true" -> поле ОБЯЗАТЕЛЬНО, иначе загрузка конфига завершится ошибкой. */
	StoragePath string `yaml:"storage_path" env-required:"true"`

	/* HTTPServer — вложенные настройки HTTP-сервера
	(адрес, таймауты, логин/пароль и т.д.).
	Это "встроенная" структура: её поля доступны внутри Config.
	yaml:"http_server" -> данные берутся из блока "http_server" в yaml. */
	HTTPServer `yaml:"http_server"`
}

type HTTPServer struct {
	Address string `yaml:"address" env-default:"localhost:8080"`
	/*Адрес, где слушать сервер (host:port).
	yaml:"address" — берётся из поля address в yaml.
	env-default:"localhost:8080" — если не задано, подставится это значение.*/

	Timeout time.Duration `yaml:"timeout" env-default:"4s"`
	/*Таймаут обработки запроса (тип time.Duration, например 4s, 1m).
	yaml:"timeout" + дефолт 4s.*/

	IdleTimeout time.Duration `yaml:"idle_timeout" env-default:"60s"`
	/*IdleTimeout
	Сколько держать idle keep-alive соединение открытым.
	yaml:"idle_timeout" + дефолт 60s.*/

	User string `yaml:"user" env-required:"true"`
	/*User
	Логин для базовой авторизации (или просто обязательный параметр в твоей конфигурации).
	env-required:"true" — должен быть задан, иначе cleanenv вернёт ошибку загрузки конфига.*/

	Password string `yaml:"password" env-required:"true" env:"HTTP_SERVER_PASSWORD"`
	/*Password
	Пароль, тоже обязательный.
	env:"HTTP_SERVER_PASSWORD" — можно читать из env-переменной HTTP_SERVER_PASSWORD (удобно, чтобы не хранить пароль в yaml).
	env-required:"true" — обязателен.*/
}

func MustLoad() *Config {
	/* Читаем значение переменной окружения.
	ВАЖНО: os.Getenv ожидает ИМЯ переменной (например "CONFIG_PATH"),
	а не путь до папки/файла. */
	configPath := os.Getenv("CONFIG_PATH")

	/* Если переменная пустая, без конфига запускать приложение нельзя.
	log.Fatal печатает сообщение и завершает программу. */
	if configPath == "" {
		log.Fatal("CONFIG_PATH is not set")
	}

	/* Проверяем, существует ли файл конфига по указанному пути.
	Если файла нет — завершаем программу с понятной ошибкой. */
	if _, err := os.Stat(configPath); os.IsNotExist(err) {
		log.Fatalf("config file does not exist: %s", configPath)
	}

	/* Создаём переменную, куда будем загружать настройки. */
	var cfg Config

	/* Читаем yaml-конфиг из файла и заполняем структуру cfg.
	Если формат неверный или не хватает обязательных полей — завершаем программу. */
	if err := cleanenv.ReadConfig(configPath, &cfg); err != nil {
		log.Fatalf("cannot read config: %s", err)
	}

	/* Возвращаем указатель на готовую структуру с настройками. */
	return &cfg
}
