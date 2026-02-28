package sl

import (
	"log/slog"
)

/*
Err — маленькая вспомогательная функция для логирования ошибок в slog.
Она принимает обычную ошибку `error` и превращает её в `slog.Attr`
с ключом "error", чтобы в логах поле ошибки было единообразным.
*/
func Err(err error) slog.Attr {
	return slog.Attr{
		/* Название поля в логе.
		В результате в записи лога будет что-то вроде: error="...текст..." */
		Key: "error",

		/* Значение поля.
		err.Error() берёт текст ошибки,
		slog.StringValue(...) упаковывает его в формат, понятный slog. */
		Value: slog.StringValue(err.Error()),
	}
}
