package main

import (
	"flag"
	"os"
)

// неэкспорованна переменная flagRunAddr содержит адрес и порт для запуска сервера
var (
	flagRunAddr  string
	flagLogLevel string
)

// parseFlags обрабатывает аргументы командной строки
// и сохраняет их значения в соответствующих переменных

func parseFlags() {
	// Регистрируем переменную flagRunAddr
	// как аргумент -a со значением :8080 по умолчанию
	flag.StringVar(&flagRunAddr, "a", ":8080", "address and port to run server")
	flag.StringVar(&flagLogLevel, "l", "debug", "log level")
	// Парсим переданные серверу аргументы в зарегистрированные переменные
	flag.Parse()

	// для случаев, когда в переменной окружения RUN_ADDR присутствует непустое значение,
	// переопределим адрес запуска сервера,
	// даже если он был передан через аргумент командной строки
	if envRunAddr := os.Getenv("RUN_ADDR"); envRunAddr != "" {
		flagRunAddr = envRunAddr
	}

	if envLogLevel := os.Getenv("LOG_LEVEL"); envLogLevel != "" {
		flagLogLevel = envLogLevel
	}

}
