package expense

import (
	"database/sql"
	"log"
	"net/http"

	"github.com/labstack/echo"
	"github.com/lib/pq"
)

var (
	errMsgInsertFailed = "Can't insert Data "
)

type Expenses struct {
	ID     int      `json:"id"`
	Title  string   `json:"title"`
	Amount float64  `json:"amount"`
	Note   string   `json:"note"`
	Tags   []string `json:"tags"`
}

func Create(c echo.Context, db *sql.DB) error {
	var e Expenses
	err := c.Bind(&e)

	if err != nil {
		return c.JSON(http.StatusBadRequest, err.Error())
	}

	row := db.QueryRow(`INSERT INTO expenses(title ,amount ,note ,tags) values($1, $2, $3, $4) RETURNING id`, e.Title, e.Amount, e.Note, pq.Array(e.Tags))

	err = row.Scan(&e.ID)
	if err != nil {
		log.Panic(errMsgInsertFailed, err)
		return c.JSON(http.StatusBadRequest, errMsgInsertFailed)
	}

	return c.JSON(http.StatusCreated, e)
}
