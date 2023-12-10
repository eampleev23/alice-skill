package main

import (
	"encoding/json"
	"fmt"
	"github.com/eampleev23/alice-skill.git/internal/logger"
	"github.com/eampleev23/alice-skill.git/internal/models"
	"github.com/eampleev23/alice-skill.git/internal/store"
	"go.uber.org/zap"
	"net/http"
	"time"
)

/*
Лучше всего при проектировании приложения с внешними зависимостями инкапсулировать их внутрь объекта,
реализующего бизнес-логику приложения. Для этого вынесем webhook из файла cmd/skill/main.go в файл cmd/skill/app.go
и превратим его в метод структуры приложения:
*/

type app struct {
	store store.Store
}

func newApp(s *store.Store) *app {
	return &app{store: *s}
}

func (a *app) webhook(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	if r.Method != http.MethodPost {
		logger.Log.Debug("got request with bad method", zap.String("method", r.Method))
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}

	// Десериализуем запрос в структуру модели
	logger.Log.Debug("decoding request")
	var req models.Request
	dec := json.NewDecoder(r.Body)
	if err := dec.Decode(&req); err != nil {
		logger.Log.Debug("cannot decode request JSON body", zap.Error(err))
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	// проверяем что пришел запрос понтного типа
	if req.Request.Type != models.TypeSimpleUtterance {
		logger.Log.Debug("unsupported request type", zap.String("type", req.Request.Type))
		w.WriteHeader(http.StatusUnprocessableEntity)
		return
	}

	// Получаем список сообщений для текущего пользователя
	messages, err := a.store.ListMessages(ctx, req.Session.User.UserID)
	if err != nil {
		logger.Log.Debug("cannot load messages for user", zap.Error(err))
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	// Формируем текст с количеством сообщений
	text := "Для вас нет новых сообщений"
	if len(messages) > 0 {
		text = fmt.Sprintf("Для вас %d новых сообщений.", len(messages))
	}

	// Первый запрос новой сессии
	if req.Session.New {
		// Обработаем поле timezone запроса
		tz, err := time.LoadLocation(req.Timezone)
		if err != nil {
			logger.Log.Debug("cannot parse timezone")
			w.WriteHeader(http.StatusBadRequest)
			return
		}
		// Получаем текущее время в часовом поясе пользователя
		now := time.Now().In(tz)
		hour, minute, _ := now.Clock()
		// Формируем новый текст приветствия
		text = fmt.Sprintf("Точное время %d часов, %d минут. %s", hour, minute, text)
	}

	// Заполняем модель ответа
	resp := models.Response{
		Response: models.ResponsePayload{
			Text: text,
		},
		Version: "1.0",
	}
	w.Header().Set("Content-Type", "application/json")

	// сериализуем ответ сервера
	enc := json.NewEncoder(w)
	if err := enc.Encode(resp); err != nil {
		logger.Log.Debug("Error encoding response", zap.Error(err))
		return
	}
	logger.Log.Debug("sending HTTP 200 response")
}
