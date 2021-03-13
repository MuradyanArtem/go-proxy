package app

import "proxy/internal/domain/repository"

type App struct {
	Request repository.Request
}

func NewApp(r repository.Request) *App {
	return &App{Request: r}
}
