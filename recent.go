package main

import (
	"net/http"

	"github.com/labstack/echo/v4"
	"github.com/sukso96100/covid19-push/database"
)

func CurrentStat(c echo.Context) error {
	current := database.GetLastStat()
	return c.JSON(http.StatusOK, map[string]interface{}{
		"confirmed":     current.Confirmed,
		"confirmedIncr": current.ConfirmedIncr,
		"cured":         current.Cured,
		"curedIncr":     current.CuredIncr,
		"death":         current.Death,
		"deathIncr":     current.DeathIncr,
		"checking":      current.Checking,
		"patients":      current.Patients,
		"patientsIncr":  current.PatientsIncr,
		"resultNeg":     current.ResultNegative,
	})
}

func RecentNews(c echo.Context) error {
	recent := database.GetRecentNews()
	result := []map[string]string{}
	for _, item := range recent {
		result = append(result, map[string]string{
			"title": item.Title,
			"dept":  item.Department,
			"link":  item.Link,
		})
	}
	return c.JSON(http.StatusOK, result)
}
