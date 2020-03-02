package database

import (
	"fmt"

	_ "github.com/GoogleCloudPlatform/cloudsql-proxy/proxy/dialers/mysql"
	"github.com/go-sql-driver/mysql"
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/mysql"
)

// conn is a singleton object that store db connection
var DbConn *gorm.DB = nil

// InitDatabase is a function that initializes database connection with provided connection information
func InitDatabase(protocol string, host string, database string,
	username string, password string,
	charset string) error {
	fmt.Println("Init db connection")
	config := mysql.NewConfig()
	config.User = username
	config.Passwd = password
	config.Net = protocol
	config.Addr = host
	config.DBName = database
	config.Params = map[string]string{
		"charset":   charset,
		"parseTime": "True",
	}
	fmt.Println(config.FormatDSN())
	db, err := gorm.Open("mysql", config.FormatDSN())
	if err != nil {
		return err
	}
	DbConn = db
	return nil
}

func MigrateDb() {
	DbConn.AutoMigrate(&StatData{}, &NewsData{})
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
	DbConn.Order("id desc").Limit(10).Find(&result)
	return result
}

func (d NewsData) Create() {
	DbConn.Create((&d))
}
