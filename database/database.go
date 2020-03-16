package database

import (
	"fmt"
	"strings"

	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/sqlite"
)

// conn is a singleton object that store db connection
var DbConn *gorm.DB = nil

// InitDatabase is a function that initializes database connection with provided connection information
func InitDatabase(dsn string) error {
	fmt.Println("Init db connection")
	// config := mysql.NewConfig()
	// config.User = username
	// config.Passwd = password
	// config.Net = protocol
	// config.Addr = host
	// config.DBName = database
	// config.Params = map[string]string{
	// 	"charset":   charset,
	// 	"parseTime": "True",
	// }
	db, err := gorm.Open("sqlite3", dsn)
	if err != nil {
		fmt.Println(err)
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
	Confirmed      int
	ConfirmedIncr  string
	Cured          int
	CuredIncr      string
	Death          int
	DeathIncr      string
	Checking       int
	Patients       int
	PatientsIncr   string
	ResultNegative int
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

func CreateStatMsg(current StatData) string {
	var builder strings.Builder
	builder.WriteString("환자현황 - 격리해제: %d %s, 격리중: %d %s,\n")
	builder.WriteString("사망: %d %s, 확진합계: %d %s\n")
	builder.WriteString("검사현황 - 검사중: %d, 결과음성: %d")

	return fmt.Sprintf(builder.String(),
		current.Cured, current.CuredIncr,
		current.Patients, current.PatientsIncr,
		current.Death, current.DeathIncr,
		current.Confirmed, current.ConfirmedIncr,
		current.Checking, current.ResultNegative,
	)
}
