// пакеты исполняемых приложений должны называться main
package main

import (
	"encoding/json"
	"github.com/eampleev23/alice-skill.git/internal/logger"
	"github.com/eampleev23/alice-skill.git/internal/models"
	"go.uber.org/zap"
	"log"
	"net/http"
)

// функция main вызывается автоматически при запуске приложения
func main() {

	parseFlags()
	if err := run(); err != nil {
		log.Fatal(err)
	}
}

// функция run будет полезна при инициализации зависимостей сервера перед запуском
func run() error {

	if err := logger.Initialize(flagLogLevel); err != nil {
		return err
	}

	logger.Log.Info("Running server", zap.String("address", flagRunAddr))

	return http.ListenAndServe(flagRunAddr, logger.RequestLogger(webhook))
}

// функция webhook — обработчик HTTP-запроса
func webhook(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		// разрешаем только POST-запросы
		logger.Log.Debug("got request with bad method", zap.String("method", r.Method))
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}

	// десериализуем запрос в структуру модели
	logger.Log.Debug("decoding request")
	var req models.Request
	dec := json.NewDecoder(r.Body)
	if err := dec.Decode(&req); err != nil {
		logger.Log.Debug("cannot decode request JSON body", zap.Error(err))
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	// Проверяем что пришел запрос понятного типа
	if req.Request.Type != models.TypeSimpleUtterance {
		logger.Log.Debug("unsupported request type", zap.String("type", req.Request.Type))
		w.WriteHeader(http.StatusUnprocessableEntity)
		return
	}
	// Заполням модель ответа
	resp := models.Response{
		Response: models.ResponsePayload{
			Text: "Извините, я пока ничего не умею",
		},
		Version: "1.0",
	}
	w.Header().Set("Content-Type", "application/json")

	// Сериализуем ответ сервера
	enc := json.NewEncoder(w)
	if err := enc.Encode(resp); err != nil {
		logger.Log.Debug("error encoding response", zap.Error(err))
		return
	}
	logger.Log.Debug("Sending HTTP 200 response")
}
