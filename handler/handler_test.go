package handler

import (
	"net/http"
	"net/http/httptest"
	"regexp"
	"strings"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/labstack/echo/v4"
	"github.com/stretchr/testify/assert"
)

var	testBody = `{"id":1,"title":"Board game","amount":60,"note":"Play board game with friends","tags":["Play","Social"]}`

func TestCreateExpense(t *testing.T) {
	e := echo.New()
	req := httptest.NewRequest(http.MethodPost, "/expenses", strings.NewReader(testBody))
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	rec := httptest.NewRecorder()

	expenseMockRows := sqlmock.NewRows([]string{"id"}).
		AddRow("1")

	db , mock, err := sqlmock.New()
	mock.ExpectQuery(regexp.QuoteMeta(`INSERT INTO expenses (title, amount, note, tags) values ($1, $2, $3, $4)  RETURNING id`)).
		WithArgs("Board game", 60.0, "Play board game with friends", `{Play, Social}`).
		WillReturnRows(expenseMockRows)

	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}

	h := handler{db}
	c := e.NewContext(req, rec)

	if assert.NoError(t, h.CreateExpense(c)) {
		assert.Equal(t, http.StatusCreated, rec.Code)
		assert.Equal(t, testBody, strings.TrimSpace(rec.Body.String()))
	}
}

func TestGetExpense(t *testing.T) {
	e := echo.New()
	req := httptest.NewRequest(http.MethodGet, "/expenses", strings.NewReader(""))
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	rec := httptest.NewRecorder()

	expenseMockRows := sqlmock.NewRows([]string{"id", "title", "amount", "note", "tags"}).
		AddRow("1", "Board game", 60, "Play board game with friends", `{Play,Social}`)

	db, mock, err := sqlmock.New()
	mock.ExpectPrepare(regexp.QuoteMeta(`SELECT id, title, amount, note, tags FROM expenses WHERE id = $1`)).
		ExpectQuery().
		WithArgs("1").
		WillReturnRows(expenseMockRows)
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	h := handler{db}
	c := e.NewContext(req, rec)
	c.SetPath("/expenses/:id")
	c.SetParamNames("id")
	c.SetParamValues("1")

	if assert.NoError(t, h.GetExpense(c)) {
		assert.Equal(t, http.StatusOK, rec.Code)
		assert.Equal(t, testBody, strings.TrimSpace(rec.Body.String()))
	}
}

func TestUpdateExpense(t *testing.T) {
	e := echo.New()
	req := httptest.NewRequest(http.MethodPut, "/expenses", strings.NewReader(testBody))
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	rec := httptest.NewRecorder()

	expenseMockRows := sqlmock.NewRows([]string{"id", "title", "amount", "note", "tags"}).
		AddRow("1", "Board game", 60, "Play board game with friends", `{Play,Social}`)

	db, mock, err := sqlmock.New()
	mock.ExpectPrepare(regexp.QuoteMeta(`UPDATE expenses SET title=$2, amount=$3, note=$4, tags=$5 WHERE id=$1 RETURNING id, title, amount, note, tags;`)).
		ExpectQuery().
		WithArgs("1", "Board game", 60.0, "Play board game with friends", `{Play, Social}`).
		WillReturnRows(expenseMockRows)

	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	h := handler{db}
	c := e.NewContext(req, rec)
	c.SetPath("/expenses/:id")
	c.SetParamNames("id")
	c.SetParamValues("1")

	if assert.NoError(t, h.UpdateExpense(c)) {
		assert.Equal(t, http.StatusOK, rec.Code)
		assert.Equal(t, testBody, strings.TrimSpace(rec.Body.String()))
	}
}


func TestGetExpenses(t *testing.T) {
	e := echo.New()
	req := httptest.NewRequest(http.MethodGet, "/expenses", strings.NewReader(""))
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	rec := httptest.NewRecorder()

	expenseMockRows := sqlmock.NewRows([]string{"id", "title", "amount", "note", "tags"}).
		AddRow("1", "Board game", 60, "Play board game with friends", `{Play,Social}`).
		AddRow("2", "apple smoothie", 89, "no discount", `{beverage}`)

	db, mock, err := sqlmock.New()
	mock.ExpectPrepare(regexp.QuoteMeta(`SELECT * FROM expenses`)).
		ExpectQuery().
		WillReturnRows(expenseMockRows)
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	expected := `[{"id":1,"title":"Board game","amount":60,"note":"Play board game with friends","tags":["Play","Social"]},{"id":2,"title":"apple smoothie","amount":89,"note":"no discount","tags":["beverage"]}]`
	h := handler{db}
	c := e.NewContext(req, rec)

	if assert.NoError(t, h.GetExpenses(c)) {
		assert.Equal(t, http.StatusOK, rec.Code)
		assert.Equal(t, expected, strings.TrimSpace(rec.Body.String()))
	}
}