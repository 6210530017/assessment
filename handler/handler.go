package handler

import (
	"database/sql"
	"fmt"
	"net/http"
	"strings"

	"github.com/lib/pq"

	"github.com/labstack/echo/v4"
)

type handler struct {
	DB *sql.DB
}

func NewHandler(db *sql.DB) *handler {
	return &handler{db}
}

type Expense struct {
	ID     int      `json:"id"`
	Title  string   `json:"title"`
	Amount float64  `json:"amount"`
	Note   string   `json:"note"`
	Tags   pq.StringArray `json:"tags"`
}

type Err struct {
	Message string `json:"message"`
}

func (h *handler) CreateExpense(c echo.Context) error {
	data := Expense{}
	err := c.Bind(&data)
	if err != nil {
		return c.JSON(http.StatusBadRequest, Err{Message: err.Error()})
	}

	tag := fmt.Sprintf("{%v}", strings.Join(data.Tags, ", "))

	row := h.DB.QueryRow(`INSERT INTO expenses (title, amount, note, tags) values ($1, $2, $3, $4)  RETURNING id`, data.Title, data.Amount, data.Note, tag)
	err = row.Scan(&data.ID)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, Err{Message: err.Error()})
	}

	return c.JSON(http.StatusCreated, data)
}