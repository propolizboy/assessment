//go:build unit

package handler

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net"
	"net/http"
	"strings"
	"testing"
	"time"

	"github.com/labstack/echo"
	"github.com/labstack/echo/middleware"
	"github.com/stretchr/testify/assert"
)

type Response struct {
	*http.Response
	err error
}

const serverPort = 80

func TestGetHealths(t *testing.T) {
	eh := setupServerTest()
	var result string
	res := request(http.MethodGet, uri("healths"), nil)
	err := res.decode(&result)
	assert.Equal(t, http.StatusOK, res.StatusCode)
	assert.Equal(t, "Hello Expenses", result)
	err = shutdownServerTest(eh)
	assert.NoError(t, err)
}

func setupServerTest() *echo.Echo {
	eh := echo.New()
	eh.Use(middleware.KeyAuthWithConfig(middleware.KeyAuthConfig{
		AuthScheme: "November",
		Validator:  AuthMiddleware,
	}))

	h := NewHandler(nil)
	h.SetupRoute(eh)
	go func(ech *echo.Echo) {

		eh.Start(fmt.Sprintf(":%d", serverPort))
	}(eh)
	for {
		conn, err := net.DialTimeout("tcp", fmt.Sprintf("localhost:%d", serverPort), 30*time.Second)
		if err != nil {
			log.Println(err)
		}
		if conn != nil {
			conn.Close()
			break
		}
	}

	return eh
}

func shutdownServerTest(eh *echo.Echo) error {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	err := eh.Shutdown(ctx)
	return err
}

func uri(paths ...string) string {
	host := "http://localhost:80"
	if paths == nil {
		return host
	}

	url := append([]string{host}, paths...)
	return strings.Join(url, "/")
}

func request(method, url string, body io.Reader) *Response {
	req, _ := http.NewRequest(method, url, body)
	req.Header.Add("Authorization", "November 10, 2009")
	req.Header.Add("Content-Type", "application/json")
	client := http.Client{}
	res, err := client.Do(req)
	return &Response{res, err}
}

func (r *Response) decode(v interface{}) error {
	if r.err != nil {
		return r.err
	}

	return json.NewDecoder(r.Body).Decode(v)
}
