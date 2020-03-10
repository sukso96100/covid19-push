package database

import (
	"fmt"

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
	Cured          int
	Death          int
	Checking       int
	Patients       int
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

func CreateStatMsg(current StatData, incr map[string]string) string {
	tmpl := `환자 현황
	- 완치(격리해제): %d %s
	- 치료중(격리중): %d %s
	- 사망: %d %s
	- 합계(확진): %d %s
	
	검사 현황
	- 검사중: %d
	- 결과 음성: %d`
	return fmt.Sprintf(tmpl,
		current.Cured, incr["Cured"],
		current.Patients, incr["Patients"],
		current.Death, incr["Death"],
		current.Confirmed, incr["Confirmed"],
		current.Checking, current.ResultNegative,
	)
}
