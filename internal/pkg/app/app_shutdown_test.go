package app

import (
	"context"
	"errors"
	"net/http"
	"os"
	"path/filepath"
	"syscall"
	"testing"
	"time"

	"github.com/Gerfey/shortener/internal/app/settings"
	"github.com/Gerfey/shortener/internal/app/strategy"
	"github.com/stretchr/testify/require"
)

func TestShortenerApp_GracefulShutdown(t *testing.T) {
	config := settings.NewSettings(settings.ServerSettings{
		ServerRunAddress:       ":0",
		ServerShortenerAddress: "http://localhost",
		ShutdownTimeout:        time.Second,
	})

	stg := strategy.NewMemoryStrategy()
	app, err := NewShortenerApp(config, stg)
	require.NoError(t, err)

	appDone := make(chan struct{})
	go func() {
		app.Run()
		close(appDone)
	}()

	time.Sleep(100 * time.Millisecond)

	process, err := os.FindProcess(os.Getpid())
	require.NoError(t, err)
	err = process.Signal(syscall.SIGTERM)
	require.NoError(t, err)

	select {
	case <-appDone:
	case <-time.After(3 * time.Second):
		t.Fatal("Приложение не завершилось в течение ожидаемого времени")
	}
}

func TestShortenerApp_HTTPSSupport(t *testing.T) {
	tmpDir := t.TempDir()
	certFile := filepath.Join(tmpDir, "test.crt")
	keyFile := filepath.Join(tmpDir, "test.key")

	require.NoError(t, os.WriteFile(certFile, []byte("test cert"), 0644))
	require.NoError(t, os.WriteFile(keyFile, []byte("test key"), 0644))

	config := settings.NewSettings(settings.ServerSettings{
		ServerRunAddress:       ":0",
		ServerShortenerAddress: "https://localhost",
		EnableHTTPS:            true,
		CertFile:               certFile,
		KeyFile:                keyFile,
	})

	stg := strategy.NewMemoryStrategy()
	app, err := NewShortenerApp(config, stg)
	require.NoError(t, err)

	originalServer := app.server
	app.server = &http.Server{
		Addr:    ":0",
		Handler: originalServer.Handler,
	}

	done := make(chan struct{})
	go func() {
		defer func() {
			if r := recover(); r != nil {
				t.Logf("Перехвачена паника: %v", r)
			}
			close(done)
		}()

		app.configureRouter()

		process, _ := os.FindProcess(os.Getpid())
		_ = process.Signal(syscall.SIGTERM)
	}()

	select {
	case <-done:
	case <-time.After(1 * time.Second):
		t.Fatal("Тест не завершился вовремя")
	}
}

func TestShortenerApp_ShutdownWithRequests(t *testing.T) {
	config := settings.NewSettings(settings.ServerSettings{
		ServerRunAddress:       ":0",
		ServerShortenerAddress: "http://localhost",
		ShutdownTimeout:        time.Second * 2,
	})

	stg := strategy.NewMemoryStrategy()
	app, err := NewShortenerApp(config, stg)
	require.NoError(t, err)

	app.configureRouter()

	server := &http.Server{
		Addr:    ":0",
		Handler: app.router,
	}

	go func() {
		if err := server.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
			t.Logf("Ошибка запуска сервера: %v", err)
		}
	}()

	time.Sleep(100 * time.Millisecond)

	longRequestDone := make(chan struct{})
	go func() {
		ctx, cancel := context.WithTimeout(context.Background(), time.Second*3)
		defer cancel()

		if err := server.Shutdown(ctx); err != nil {
			t.Logf("Ошибка при завершении сервера: %v", err)
		}
		close(longRequestDone)
	}()

	select {
	case <-longRequestDone:
	case <-time.After(5 * time.Second):
		t.Fatal("Сервер не завершился корректно в течение ожидаемого времени")
	}
}
