package app

import (
	"context"
	"errors"
	"net/http"
	"os"
	"os/signal"
	"syscall"

	"github.com/Gerfey/shortener/internal/app/handler"
	"github.com/Gerfey/shortener/internal/app/middleware"
	"github.com/Gerfey/shortener/internal/app/service"
	"github.com/Gerfey/shortener/internal/app/settings"
	"github.com/Gerfey/shortener/internal/models"
	chi "github.com/go-chi/chi/v5"
	"github.com/sirupsen/logrus"
	"golang.org/x/crypto/acme/autocert"
)

// ShortenerApp основной класс приложения
type ShortenerApp struct {
	settings   *settings.Settings
	router     *chi.Mux
	handler    *handler.URLHandler
	server     *http.Server
	strategy   models.StorageStrategy
	repository models.Repository
}

// NewShortenerApp создает новое приложение
func NewShortenerApp(settings *settings.Settings, strategy models.StorageStrategy) (*ShortenerApp, error) {
	logrus.SetFormatter(&logrus.TextFormatter{
		FullTimestamp: true,
	})
	logrus.SetLevel(logrus.InfoLevel)

	repository, err := strategy.Initialize()
	if err != nil {
		return nil, err
	}

	shortenerService := service.NewShortenerService(repository)
	urlService := service.NewURLService(settings)
	urlHandler := handler.NewURLHandler(shortenerService, urlService, settings, repository)

	router := chi.NewRouter()
	router.Use(middleware.LoggingMiddleware)
	router.Use(middleware.GzipMiddleware)

	application := &ShortenerApp{
		settings:   settings,
		router:     router,
		handler:    urlHandler,
		strategy:   strategy,
		repository: repository,
		server: &http.Server{
			Addr:    settings.ServerAddress(),
			Handler: router,
		},
	}

	return application, nil
}

// configureRouter настраивает маршруты
func (a *ShortenerApp) configureRouter() {
	a.router.Route("/", func(r chi.Router) {
		r.Post("/", middleware.AuthMiddleware(a.handler.ShortenHandler))
		r.Post("/api/shorten", middleware.AuthMiddleware(a.handler.ShortenJSONHandler))
		r.Post("/api/shorten/batch", middleware.AuthMiddleware(a.handler.ShortenBatchHandler))
		r.Get("/api/user/urls", middleware.AuthMiddleware(a.handler.GetUserURLsHandler))
		r.Delete("/api/user/urls", middleware.AuthMiddleware(a.handler.DeleteUserURLsHandler))
		r.Get("/ping", a.handler.PingHandler)
		r.Get("/{id}", a.handler.RedirectURLHandler)
	})
}

// Run запускает приложение
func (a *ShortenerApp) Run() {
	logrus.Printf("Starting server: %v", a.settings.ServerAddress())

	a.configureRouter()

	stop := make(chan os.Signal, 1)
	signal.Notify(stop, os.Interrupt, syscall.SIGTERM, syscall.SIGQUIT)

	go func() {
		var err error

		if a.settings.Server.EnableHTTPS {
			logrus.Info("HTTPS enabled")

			manager := &autocert.Manager{
				Cache:  autocert.DirCache("certs-cache"),
				Prompt: autocert.AcceptTOS,
			}

			a.server.TLSConfig = manager.TLSConfig()
			err = a.server.ListenAndServeTLS("", "")
		} else {
			err = a.server.ListenAndServe()
		}

		if err != nil && !errors.Is(err, http.ErrServerClosed) {
			logrus.Fatal(err)
		}
	}()

	sig := <-stop
	logrus.Infof("Получен сигнал завершения: %v", sig)
	logrus.Info("Начинаем корректное завершение работы сервера...")

	ctx, cancel := context.WithTimeout(context.Background(), a.settings.ShutdownTimeout())
	defer cancel()

	logrus.Info("Ожидаем завершения всех обрабатываемых запросов...")
	if err := a.server.Shutdown(ctx); err != nil {
		logrus.Error("Сервер принудительно завершен:", err)
	} else {
		logrus.Info("Все запросы успешно обработаны")
	}

	logrus.Info("Сохраняем данные и закрываем хранилище...")
	if err := a.strategy.Close(); err != nil {
		logrus.Error("Ошибка при закрытии хранилища:", err)
	} else {
		logrus.Info("Хранилище успешно закрыто")
	}

	logrus.Info("Сервер успешно остановлен")
}
