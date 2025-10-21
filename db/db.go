package db

import(
	"fmt"
	"github.com/jmoiron/sqlx"
	_ "github.com/go-sql-driver/mysql" // Blank import to register the driver
)

var DB *sqlx.DB

type Config struct{
	Username	 string
	Password	 string
	Server		 string
	DatabaseName string
}

func InitDB(config Config) error  {
	constr := fmt.Sprintf("%s:%s@tcp(%s)/%s?parseTime=true&loc=Asia%%2FBangkok&charset=utf8mb4",config.Username,config.Password,config.Server,config.DatabaseName)

	_db, err := sqlx.Open("mysql",constr)

	if err != nil {
		return err
		
	}

	DB = _db
	return nil
	
}