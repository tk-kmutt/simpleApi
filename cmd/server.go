package main

import (
	"github.com/go-playground/validator/v10"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"net/http"
)

type User struct {
	ID   string `json:"id"`
	Name string `json:"name" validate:"required"`
	Age  int    `json:"age" validate:"max=150"`
}

type UserPatch struct {
	ID   *string `json:"id"`
	Name *string `json:"name"`
	Age  *int    `json:"age"`
}

var users = make(map[string]*User)

type CustomValidator struct {
	validator *validator.Validate
}

func NewValidator() *CustomValidator {
	return &CustomValidator{
		validator: validator.New(),
	}
}

func (cv *CustomValidator) Validate(i interface{}) error {
	if err := cv.validator.Struct(i); err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}
	return nil
}

func main() {
	e := echo.New()
	e.Use(middleware.Logger())
	e.Validator = NewValidator()
	e.POST("/simple", func(context echo.Context) error {
		// リクエストを取得する
		user := new(User)
		_ = context.Bind(user)
		// バリデーション
		if err := context.Validate(user); err != nil {
			return context.String(http.StatusBadRequest, err.Error())
		}
		users[user.ID] = user
		return context.JSON(http.StatusCreated, user)
	})
	e.GET("/simple/:id", func(context echo.Context) error {
		// リクエストを取得する
		id := context.Param("id")
		return context.JSON(http.StatusOK, users[id])
	})
	e.PUT("/simple/:id", func(context echo.Context) error {
		// リクエストを取得する
		id := context.Param("id")
		user := new(User)
		// バリデーション
		if err := context.Validate(user); err != nil {
			return context.JSON(http.StatusBadRequest, err)
		}
		_ = context.Bind(user)
		users[id] = user
		return context.JSON(http.StatusOK, user)
	})
	e.PATCH("/simple/:id", func(context echo.Context) error {
		// リクエストを取得する
		id := context.Param("id")
		user := new(UserPatch)
		// バリデーション
		if err := context.Validate(user); err != nil {
			return context.JSON(http.StatusBadRequest, err)
		}
		_ = context.Bind(user)
		base := users[id]
		if user.Age != nil {
			base.Age = *user.Age
		}
		if user.Name != nil {
			base.Name = *user.Name
		}
		users[id] = base
		return context.JSON(http.StatusOK, base)
	})
	e.DELETE("/simple/:id", func(context echo.Context) error {
		// リクエストを取得する
		id := context.Param("id")
		users[id] = nil
		return context.String(http.StatusNoContent, "")
	})
	e.Logger.Fatal(e.Start(":1232"))
}
