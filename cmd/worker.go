package main

import (
	"fmt"
	"net/http"
	"os"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

func main() {
	e := echo.New()
	// ログ出力ミドルウェア（任意）
	e.Use(middleware.Logger())
	e.Use(middleware.Recover())

	// POST /tasks エンドポイント
	e.POST("/tasks", func(c echo.Context) error {
		// リクエストヘッダをログに出力
		for name, values := range c.Request().Header {
			for _, v := range values {
				c.Logger().Info(fmt.Sprintf("%s: %s", name, v))
			}
		}

		// リクエストボディを取得してログ出力
		bodyBytes := make([]byte, c.Request().ContentLength)
		if _, err := c.Request().Body.Read(bodyBytes); err != nil && err.Error() != "EOF" {
			c.Logger().Error("Error reading body: %v", err)
			return c.String(http.StatusInternalServerError, "Failed to read body")
		}
		c.Logger().Info(fmt.Sprintf("Request Body: %s", string(bodyBytes)))

		// 単純に200を返す
		return c.String(http.StatusOK, "OK")
	})

	// ポート指定（Cloud Runは$PORT環境変数を使用）
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}
	// サーバー起動
	e.Logger.Fatal(e.Start(fmt.Sprintf(":%s", port)))
}
