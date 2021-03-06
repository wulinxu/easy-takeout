package definition

import (
	"errors"
	"fmt"
	"strings"
	"time"

	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/mysql"
	_ "github.com/jinzhu/gorm/dialects/postgres"
	_ "github.com/jinzhu/gorm/dialects/sqlite"
)

func Connect(dbtype, dbaddress string, args map[string]interface{}) (*gorm.DB, error) {
	source := ""
	switch dbtype {
	case "mysql":
		source = dbaddress
		//combination user and password
		if user, ok := args["user"]; ok {
			preStr := user.(string)
			if pass, ok := args["password"]; ok {
				preStr += ":" + pass.(string)
				delete(args, "password")
			}
			source = fmt.Sprintf("%s@%s", preStr, dbaddress)
			delete(args, "user")
		}
		if dbname, ok := args["dbname"]; ok {
			source += fmt.Sprintf("/%s", dbname)
			delete(args, "dbname")
		}
		//combination args
		if len(args) > 0 {
			query := []string{}
			for k, v := range args {
				query = append(query, fmt.Sprintf("%s=%s", k, v))
			}
			source += "?" + strings.Join(query, "&")
		}
	case "sqlite3":
		source = dbaddress
	case "postgres":
		source = "host=" + dbaddress
		//combination args
		for k, v := range args {
			source += fmt.Sprintf(" %s=%s", k, v)
		}

	default:
		return nil, errors.New("database type error")
	}
	db, err := gorm.Open(dbtype, source)
	if err != nil {
		return nil, err
	}
	return db, nil
}

func BetweenCreateTime(start, end time.Time) func(db *gorm.DB) *gorm.DB {
	return func(db *gorm.DB) *gorm.DB {
		return db.Where("created_at BETWEEN ? AND ?", start, end)
	}
}
