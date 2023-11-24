package app_test

import (
	"context"
	"io"
	"net/http"
	"testing"
	"time"
	"webhook/internal/app"
	"webhook/internal/config"

	"github.com/stretchr/testify/assert"
	"go.uber.org/fx"
)

func Test_CreateApp(t *testing.T) {

	localServerPort := "8081"

	cfgLoader := func() (config.Configuration, error) {

		cfg := config.Configuration{
			URL: "http://localhost:" + localServerPort,
		}
		cfg.Requests.Amount = 1
		cfg.Requests.PerSecond = 1

		return cfg, nil
	}

	app := app.CreateApp(cfgLoader)

	// Test that app is not nil
	assert.NotNil(t, app)

	// Test that the app is of type *fx.App
	assert.IsType(t, &fx.App{}, app)
}

func Test_NoRequest(t *testing.T) {

	localServerPort := "8081"
	var (
		actualResponse []byte
	)
	expectedResponse := ""

	//
	// Starting test server (httptest server is not very usefull here)
	//
	go func() {
		httpMux := http.NewServeMux()
		httpMux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
			var bodyContent string
			r.Body.Read([]byte(bodyContent))
			w.WriteHeader(200)

			actualResponse, _ = io.ReadAll(r.Body)

		})
		_ = http.ListenAndServe(":"+localServerPort, httpMux)
	}()

	//
	// Configuring our "mock"-app
	//
	cfgLoader := func() (config.Configuration, error) {

		cfg := config.Configuration{
			URL: "http://localhost:" + localServerPort,
		}
		cfg.Requests.Amount = 0
		cfg.Requests.PerSecond = 1

		return cfg, nil
	}

	//
	// we need some workaround here
	//
	app := app.CreateApp(cfgLoader)
	go app.Run()

	//
	// we should wait for server will process our one request
	// then we will stop app manually to avoid command context deadline appears
	//
	time.Sleep(1 * time.Second)
	app.Stop(context.Background())

	//
	// If there was no request -> no response should be either
	//
	if string(actualResponse) != expectedResponse {
		t.Errorf("request body expected %s got %s", expectedResponse, string(actualResponse))
	}
}

func Test_OneRequest(t *testing.T) {

	localServerPort := "8082"
	expectedResponse := "{ 'iteration': 0 }"

	//
	// Starting test server (httptest server is not very usefull here)
	//
	go func() {
		httpMux := http.NewServeMux()
		httpMux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
			var bodyContent string
			r.Body.Read([]byte(bodyContent))
			w.WriteHeader(200)

			actualResponse, _ := io.ReadAll(r.Body)
			if string(actualResponse) != expectedResponse {
				t.Errorf("request body expected %s got %s", expectedResponse, string(actualResponse))
			}

		})
		_ = http.ListenAndServe(":"+localServerPort, httpMux)
	}()

	//
	// Configuring our "mock"-app
	//
	cfgLoader := func() (config.Configuration, error) {

		cfg := config.Configuration{
			URL: "http://localhost:" + localServerPort,
		}
		cfg.Requests.Amount = 1
		cfg.Requests.PerSecond = 1

		return cfg, nil
	}

	//
	// we need some workaround here
	//
	app := app.CreateApp(cfgLoader)
	go app.Run()

	//
	// we should wait for server will process our one request
	// then we will stop app manually to avoid command context deadline appears
	//
	time.Sleep(1 * time.Second)
	app.Stop(context.Background())
}

func Test_TwoRequests(t *testing.T) {

	localServerPort := "8083"
	responses := make([]string, 0, 2)

	//
	// Starting test server (httptest server is not very usefull here)
	//
	go func() {
		httpMux := http.NewServeMux()
		httpMux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
			var bodyContent string
			r.Body.Read([]byte(bodyContent))
			w.WriteHeader(200)

			body, _ := io.ReadAll(r.Body)
			responses = append(responses, string(body))

		})
		_ = http.ListenAndServe(":"+localServerPort, httpMux)
	}()

	//
	// Configuring our "mock"-app
	//
	cfgLoader := func() (config.Configuration, error) {

		cfg := config.Configuration{
			URL: "http://localhost:" + localServerPort,
		}
		cfg.Requests.Amount = 2
		cfg.Requests.PerSecond = 2

		return cfg, nil
	}

	//
	// we need some workaround here
	//
	app := app.CreateApp(cfgLoader)
	go app.Run()

	//
	// we should wait for server will process our one request
	// then we will stop app manually to avoid command context deadline appears
	//
	time.Sleep(2 * time.Second)
	app.Stop(context.Background())

	if len(responses) != 2 {
		t.Errorf("expected 2 responses got %d (%v)", len(responses), responses)
	}

	if responses[0] != "{ 'iteration': 0 }" {
		t.Errorf("expected first response '{ 'iteration': 0 }' got %s", responses[0])
	}

	if responses[1] != "{ 'iteration': 1 }" {
		t.Errorf("expected second response '{ 'iteration': 1 }' got %s", responses[1])
	}

}
