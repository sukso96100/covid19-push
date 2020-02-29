package main

import (
	"net/http"

	"github.com/labstack/echo/v4"
	"github.com/sukso96100/covid19-push/database"
)

func CurrentStat(c echo.Context) error {
	current := database.GetLastStat()
	return c.JSON(http.StatusOK, map[string]int{
		"confirmed": current.Confirmed,
		"cured": current.Cured,
		"death": current.Death,
	})
}

func RecentNews(c echo.Context) error{
	recent := database.GetRecentNews()
	result := []map[string]string{}
	for _,item := range recent {
		result = append(result, map[string]string{
			"title":item.Title,
			"dept":item.Department,
			"link":item.Link,
		})
	}
	return c.JSON(http.StatusOK, result)
}