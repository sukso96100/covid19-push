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
	Checking  int
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

func CreateStatMsg(prev StatData, current StatData) string {
	// See documentation on defining a message payload.
	var confirmedIncSig string
	var curedIncSig string
	var deathIncSig string
	var checkingIncSig string
	if current.Confirmed-prev.Confirmed > 0 {
		confirmedIncSig = fmt.Sprintf("+%d", current.Confirmed-prev.Confirmed)
	} else {
		confirmedIncSig = fmt.Sprintf("%d", current.Confirmed-prev.Confirmed)
	}
	if current.Cured-prev.Cured > 0 {
		curedIncSig = fmt.Sprintf("+%d", current.Cured-prev.Cured)
	} else {
		curedIncSig = fmt.Sprintf("%d", current.Cured-prev.Cured)
	}
	if current.Death-prev.Death > 0 {
		deathIncSig = fmt.Sprintf("+%d", current.Death-prev.Death)
	} else {
		deathIncSig = fmt.Sprintf("%d", current.Death-prev.Death)
	}
	if current.Checking-prev.Checking > 0 {
		checkingIncSig = fmt.Sprintf("+%d", current.Checking-prev.Checking)
	} else {
		checkingIncSig = fmt.Sprintf("%d", current.Checking-prev.Checking)
	}
	tmpl := "확진:%d명 (%s), 완치:%d(%s), 사망:%d(%s), 검사진행:%d(%s)"
	return fmt.Sprintf(tmpl,
		current.Confirmed, confirmedIncSig,
		current.Cured, curedIncSig,
		current.Death, deathIncSig,
		current.Checking, checkingIncSig,
	)
}
