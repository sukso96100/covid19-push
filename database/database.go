package database

import (
	"fmt"

	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/mysql"
)

// conn is a singleton object that store db connection
var DbConn *gorm.DB = nil

// InitDatabase is a function that initializes database connection with provided connection information
func InitDatabase(host string, database string,
	username string, password string,
	charset string) error {
	fmt.Println("Init db connection")
	connInfo := fmt.Sprintf("%s:%s@(%s)/%s?charset=%s&parseTime=True&loc=UTC",
		username, password, host, database, charset)
	db, err := gorm.Open("mysql", connInfo)
	if err != nil {
		return err
	}
	DbConn = db
	return nil
}

func MigrateDb() {
	DbConn.AutoMigrate(&StatData{}, NewsData{})
}

type StatData struct {
	gorm.Model
	Confirmed int
	Cured     int
	Death     int
}

type NewsData struct {
	gorm.Model
	PostId     int
	Title      string
	Department string
	Link       string
}

func GetLastStat() StatData {
	result := StatData{}
	DbConn.Last(&result)
	return result
}

func (d StatData) Create() {
	DbConn.Create((&d))
}

func GetLastNews() NewsData {
	result := NewsData{}
	DbConn.Last(&result)
	return result
}

func GetRecentNews() []NewsData {
	result := []NewsData{}
	DbConn.Limit(10).Find(&result)
	return result
}

func (d NewsData) Create() {
	DbConn.Create((&d))
}
