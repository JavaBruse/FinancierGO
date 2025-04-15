package middleware

import (
	"bytes"
	"financierGo/internal/constants"
	"fmt"
	"io"
	"net/http"
	"runtime"
	"strings"
	"time"
)

type responseWriter struct {
	http.ResponseWriter
	statusCode int
	body       bytes.Buffer
}

func (rw *responseWriter) WriteHeader(code int) {
	rw.statusCode = code
	rw.ResponseWriter.WriteHeader(code)
}

func (rw *responseWriter) Write(b []byte) (int, error) {
	rw.body.Write(b)
	return rw.ResponseWriter.Write(b)
}

func getCallerInfo(skip int) string {
	var stack []string
	for i := 0; i < 4; i++ {
		pc, file, line, ok := runtime.Caller(skip + i)
		if !ok {
			break
		}

		fn := runtime.FuncForPC(pc)
		if fn == nil {
			break
		}

		// Получаем полное имя функции
		funcName := fn.Name()

		// Проверяем, является ли файл частью нашего проекта
		if strings.Contains(file, "FinancierGO") {
			// Получаем относительный путь от корня проекта
			parts := strings.Split(file, "FinancierGO/")
			if len(parts) > 1 {
				relPath := parts[1]
				stack = append(stack, fmt.Sprintf("%s:%s:%d", relPath, funcName, line))
			}
		} else {
			// Для внешних файлов показываем только имя функции
			stack = append(stack, fmt.Sprintf("ext:%s", funcName))
		}
	}

	if len(stack) == 0 {
		return "unknown"
	}

	return strings.Join(stack, " -> ")
}

func LoggingMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Создаем буфер для тела запроса
		var bodyBytes []byte
		if r.Body != nil {
			bodyBytes, _ = io.ReadAll(r.Body)
			r.Body = io.NopCloser(bytes.NewBuffer(bodyBytes))
		}

		// Создаем кастомный response writer
		rw := &responseWriter{
			ResponseWriter: w,
			statusCode:     http.StatusOK,
		}

		// Выполняем следующий обработчик
		next.ServeHTTP(rw, r)

		// Получаем текущее время
		now := time.Now().Format("2006-01-02 15:04:05")

		// Получаем ID пользователя из контекста
		userID := "anonymous"
		if id := r.Context().Value(constants.UserIDKey); id != nil {
			userID = fmt.Sprintf("%v", id)
		}

		// Убираем переносы строк из JSON
		bodyStr := strings.ReplaceAll(string(bodyBytes), "\n", "")
		bodyStr = strings.ReplaceAll(bodyStr, "\r", "")
		bodyStr = strings.ReplaceAll(bodyStr, "  ", " ")

		// Форматируем лог в зависимости от статуса ответа
		if rw.statusCode >= 400 {
			// Для ошибок показываем стек вызовов
			callerInfo := getCallerInfo(2)
			fmt.Printf("%s [%s] [user:id=%s] %s %s\n",
				now, callerInfo, userID, r.URL.Path, bodyStr)
		} else {
			// Для успешных запросов показываем только базовую информацию
			fmt.Printf("%s [user:id=%s] %s %s\n",
				now, userID, r.URL.Path, bodyStr)
		}
	})
}
