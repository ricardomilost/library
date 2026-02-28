package sqlite

import (
	"database/sql"
	"fmt"

	"github.com/mattn/go-sqlite3"
	//_ "github.com/mattn/go-sqlite3" // init sqlite3 driver

	"library/internal/storage"
)

/*
	Sqlite — это структура-обёртка над подключением к базе данных.

db *sql.DB — объект, через который мы выполняем SQL-запросы.
*/
type Storage struct {
	db *sql.DB
}

/*
	New — "конструктор" для нашего хранилища.

Он должен открыть sqlite-базу и вернуть готовый объект Sqlite.
*/
func New(storagePath string) (*Storage, error) {
	/* op — строка с названием операции.
	Её обычно добавляют в текст ошибок, чтобы проще понимать,
	в каком месте произошла проблема. */
	const op = "storage.sqlite.New"

	/* Открываем (или создаём, если файла ещё нет) SQLite-базу.
	"sqlite3" — это драйвер, который умеет работать с форматом SQLite.
	Второй параметр — путь к файлу базы данных на диске. */
	db, err := sql.Open("sqlite3", storagePath)
	if err != nil {
		/* Если подключиться к базе не удалось, возвращаем ошибку с контекстом:
		op показывает, в каком месте произошла проблема.
		Возвращаем nil, потому что готовый объект хранилища не создан. */
		return nil, fmt.Errorf("%s: %w", op, err)
	}
	/* Подготавливаем SQL-запрос(ы) для инициализации базы:
	   1) создаём таблицу url, если её ещё нет;
	   2) создаём индекс по alias, чтобы поиск по alias работал быстрее. */
	stmt, err := db.Prepare(`
CREATE TABLE IF NOT EXISTS url (
	id INTEGER PRIMARY KEY,
	alias TEXT NOT NULL UNIQUE,
	url TEXT NOT NULL);
CREATE INDEX IF NOT EXISTS idx_alias ON url(alias);
`)
	if err != nil {
		/* Если подготовка SQL не удалась, возвращаем ошибку с контекстом op,
		чтобы в логах было понятно, где именно произошёл сбой. */
		return nil, fmt.Errorf("%s: %w", op, err)
	}
	// Выполняем подготовленный SQL-запрос (например, создание таблицы).
	_, err = stmt.Exec()
	if err != nil {
		// Добавляем контекст операции и пробрасываем исходную ошибку выше.
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	// Инициализация прошла успешно: возвращаем объект хранилища с открытой БД.
	return &Storage{db: db}, nil
}

func (s *Storage) SaveURL(urlToSave string, alias string) (int64, error) {
	const op = "storage.sqlite.SaveURL"

	// Подготавливаем INSERT с плейсхолдерами, чтобы безопасно передать значения.
	stmt, err := s.db.Prepare("INSERT INTO url(url, alias) VALUES (?, ?)")
	if err != nil {
		// Ошибка на этапе подготовки SQL.
		return 0, fmt.Errorf("%s: %w", op, err)
	}

	// Выполняем INSERT и передаём реальные значения url и alias.
	res, err := stmt.Exec(urlToSave, alias)
	if err != nil {
		// Если нарушено уникальное ограничение (дубликат alias/url),
		// возвращаем понятную бизнес-ошибку ErrURLExists.
		if sqliteErr, ok := err.(sqlite3.Error); ok && sqliteErr.ExtendedCode == sqlite3.ErrConstraintUnique {
			return 0, fmt.Errorf("%s: %w", op, storage.ErrURLExists)
		}
		// Любую другую SQL-ошибку возвращаем как есть, с контекстом операции.
		return 0, fmt.Errorf("%s: %w", op, err)
	}

	// Получаем ID только что добавленной записи.
	id, err := res.LastInsertId()
	if err != nil {
		return 0, fmt.Errorf("%s: failed to get last insert id: %w", op, err)
	}

	// Успех: возвращаем ID новой записи.
	return id, nil
}
