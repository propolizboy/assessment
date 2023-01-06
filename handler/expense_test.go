//go:build integration

package handler

import (
	"bytes"
	"context"
	"database/sql"
	"encoding/json"
	"fmt"
	"github.com/labstack/echo"
	"github.com/labstack/echo/middleware"
	"github.com/propolizboy/assessment/expense"
	"github.com/stretchr/testify/assert"
	"io"
	"log"
	"net"
	"net/http"
	"strconv"
	"strings"
	"testing"
	"time"
)

type Response struct {
	*http.Response
	err error
}

const serverPort = 80

func TestCreateExpense(t *testing.T) {
	eh := setupServerTest()

	body := bytes.NewBufferString(`{
		"title": "strawberry smoothie",
		"amount": 79,
		"note": "night market promotion discount 10 bath", 
		"tags": ["food", "beverage"]
	}`)

	var e expense.Expenses

	res := request(http.MethodPost, uri("expenses"), body)
	err := res.Decode(&e)

	assert.Nil(t, err)
	assert.Equal(t, http.StatusCreated, res.StatusCode)
	assert.NotEqual(t, 0, e.ID)
	assert.Equal(t, "strawberry smoothie", e.Title)
	assert.Equal(t, float64(79), e.Amount)
	assert.Equal(t, "night market promotion discount 10 bath", e.Note)
	assert.Equal(t, []string{"food", "beverage"}, e.Tags)

	err = shutdownServerTest(eh)
	assert.NoError(t, err)
}

func TestGetEnpenseByID(t *testing.T) {
	eh := setupServerTest()

	e := seedExpense(t)
	var latest expense.Expenses
	res := request(http.MethodGet, uri("expenses", strconv.Itoa(e.ID)), nil)
	err := res.Decode(&latest)

	assert.Nil(t, err)
	assert.Equal(t, http.StatusOK, res.StatusCode)
	assert.Equal(t, e.ID, latest.ID)
	assert.Equal(t, e.Title, latest.Title)
	assert.Equal(t, e.Amount, latest.Amount)
	assert.Equal(t, e.Note, latest.Note)
	assert.Equal(t, e.Tags, latest.Tags)

	err = shutdownServerTest(eh)
	assert.NoError(t, err)
}

func TestUpdateEnpenseByID(t *testing.T) {
	eh := setupServerTest()

	e := seedExpense(t)
	body := bytes.NewBufferString(`{
		"title": "apple smoothie",
		"amount": 89,
		"note": "no discount",
		"tags": ["beverage"]
	}`)
	var latest expense.Expenses
	res := request(http.MethodPut, uri("expenses", strconv.Itoa(e.ID)), body)
	err := res.Decode(&latest)

	assert.Nil(t, err)
	assert.Equal(t, http.StatusOK, res.StatusCode)
	assert.Equal(t, e.ID, latest.ID)
	assert.Equal(t, "apple smoothie", latest.Title)
	assert.Equal(t, float64(89), latest.Amount)
	assert.Equal(t, "no discount", latest.Note)
	assert.Equal(t, []string{"beverage"}, latest.Tags)

	err = shutdownServerTest(eh)
	assert.NoError(t, err)
}

func TestGetAllUser(t *testing.T) {
	eh := setupServerTest()

	seedExpense(t)
	var es []expense.Expenses
	res := request(http.MethodGet, uri("expenses"), nil)
	err := res.Decode(&es)
	assert.Nil(t, err)
	assert.EqualValues(t, http.StatusOK, res.StatusCode)
	assert.Greater(t, len(es), 0)

	err = shutdownServerTest(eh)
	assert.NoError(t, err)
}

func seedExpense(t *testing.T) expense.Expenses {
	var c expense.Expenses
	body := bytes.NewBufferString(`{
		"title": "strawberry smoothie",
		"amount": 79,
		"note": "night market promotion discount 10 bath",
		"tags": ["food", "beverage"]
	}`)
	url := uri("expenses")
	err := request(http.MethodPost, url, body).Decode(&c)
	if err != nil {
		t.Fatal("can't create uomer:", err)
	}
	return c
}

func setupServerTest() *echo.Echo {
	eh := echo.New()
	eh.Use(middleware.KeyAuthWithConfig(middleware.KeyAuthConfig{
		AuthScheme: "November",
		Validator:  AuthMiddleware,
	}))
	db, err := sql.Open("postgres", "postgresql://root:root@db/go-example-db?sslmode=disable")
	if err != nil {
		log.Fatal(err)
	}
	h := NewHandler(db)
	h.SetupRoute(eh)
	// Setup server
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

func (r *Response) Decode(v interface{}) error {
	if r.err != nil {
		return r.err
	}

	return json.NewDecoder(r.Body).Decode(v)
}
