package main

import (
	"github.com/labstack/echo/v4"
	"net/http"
)

type User struct {
	ID   string `json:"id"`
	Name string `json:"name"`
	Age  int    `json:"age"`
}

type UserPatch struct {
	ID   *string `json:"id"`
	Name *string `json:"name"`
	Age  *int    `json:"age"`
}

var users = make(map[string]*User)

func main() {
	e := echo.New()
	e.POST("/simple", func(context echo.Context) error {
		// リクエストを取得する
		user := new(User)
		_ = context.Bind(user)
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
		_ = context.Bind(user)
		users[id] = user
		return context.JSON(http.StatusOK, user)
	})
	e.PATCH("/simple/:id", func(context echo.Context) error {
		// リクエストを取得する
		id := context.Param("id")
		user := new(UserPatch)
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
