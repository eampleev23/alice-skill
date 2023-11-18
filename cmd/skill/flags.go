package main

import "flag"

// неэкспорованна переменная flagRunAddr содержит адрес и порт для запуска сервера
var flagRunAddr string

// parseFlags обрабатывает аргументы командной строки
// и сохраняет их значения в соответствующих переменных

func parseFlags() {
	// Регистрируем переменную flagRunAddr
	// как аргумент -a со значением :8080 по умолчанию
	flag.StringVar(&flagRunAddr, "a", ":8080", "address and port to run server")
	// Парсим переданные серверу аргументы в зарегистрированные переменные
	flag.Parse()
}
