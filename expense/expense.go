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

type Err struct {
	Message string `json:"message"`
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

func GetById(c echo.Context, db *sql.DB) error {
	id := c.Param("id")
	stmt, err := db.Prepare("SELECT id,title ,amount ,note ,tags FROM expenses where id=$1")
	if err != nil {
		return c.JSON(http.StatusInternalServerError, Err{Message: "can't prepare query expenses statment:" + err.Error()})
	}
	row := stmt.QueryRow(id)
	var e Expenses
	err = row.Scan(&e.ID, &e.Title, &e.Amount, &e.Note, pq.Array(&e.Tags))
	switch err {
	case sql.ErrNoRows:
		return c.JSON(http.StatusNotFound, Err{Message: "expense not found"})
	case nil:
		return c.JSON(http.StatusOK, e)
	default:
		return c.JSON(http.StatusInternalServerError, Err{Message: "can't scan expense:" + err.Error()})
	}
}

func UpdateByID(c echo.Context, db *sql.DB) error {
	id := c.Param("id")
	var e Expenses
	err := c.Bind(&e)
	stmt, err := db.Prepare("Update expenses set title=$2 ,amount=$3 ,note=$4 ,tags=$5 where id=$1 RETURNING id,title ,amount ,note ,tags ")

	if err != nil {
		return c.JSON(http.StatusInternalServerError, Err{Message: "can't prepare query expenses statment:" + err.Error()})
	}
	row := stmt.QueryRow(id, e.Title, e.Amount, e.Note, pq.Array(e.Tags))
	err = row.Scan(&e.ID, &e.Title, &e.Amount, &e.Note, pq.Array(&e.Tags))
	switch err {
	case sql.ErrNoRows:
		return c.JSON(http.StatusNotFound, Err{Message: "expense not found"})
	case nil:
		return c.JSON(http.StatusOK, e)
	default:
		return c.JSON(http.StatusInternalServerError, Err{Message: "can't scan expenses:" + err.Error()})
	}
}

func GetAll(c echo.Context, db *sql.DB) error {
	stmt, err := db.Prepare("SELECT id,title ,amount ,note ,tags FROM expenses")
	if err != nil {
		return c.JSON(http.StatusInternalServerError, Err{Message: "can't prepare query expenses statment:" + err.Error()})
	}

	rows, err := stmt.Query()

	if err != nil {
		return c.JSON(http.StatusInternalServerError, Err{Message: "can't query all expenses" + err.Error()})
	}

	expenses := []Expenses{}
	for rows.Next() {
		var e Expenses
		err = rows.Scan(&e.ID, &e.Title, &e.Amount, &e.Note, pq.Array(&e.Tags))
		if err != nil {
			return c.JSON(http.StatusInternalServerError, Err{Message: "can't scan user:" + err.Error()})
		}
		expenses = append(expenses, e)
	}

	return c.JSON(http.StatusOK, expenses)

}
