package repository

import (
	"fmt"
	"os"

	_ "github.com/go-sql-driver/mysql"
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/postgres" // indirect
	_ "github.com/jinzhu/gorm/dialects/sqlite"   //
	"synapso/config"
	"synapso/logger"
)

// Repository defines a repository for access the database.
type Repository struct {
	Db *gorm.DB
}

var rep *Repository

const (
	// SQLITE represents SQLite3
	SQLITE = "sqlite3"
	// POSTGRES represents PostgreSQL
	POSTGRES = "postgres"
	// MYSQL represents MySQL
	MYSQL = "mysql"
)

func getConnection(config *config.Config) string {
	if config.Database.Dialect == POSTGRES {
		return os.Getenv("DATABASE_URL")
	} else if config.Database.Dialect == MYSQL {
		return fmt.Sprintf("%s:%s(%s)/%s?charset=utf8&parseTime=True&loc=Local", config.Database.Username, config.Database.Password, config.Database.Host, config.Database.Dbname)
	}
	return config.Database.Host
}

// InitDB initialize a database connection.
func InitDB() {
	logger.GetEchoLogger().Info("Try database connection")
	conf := config.GetConfig()
	fmt.Println("My connection:" + getConnection(conf))
	db, err := gorm.Open(conf.Database.Dialect, getConnection(conf))
	if err != nil {
		logger.GetEchoLogger().Error("Failure database connection")
		logger.GetEchoLogger().Error(err)
	}
	logger.GetEchoLogger().Info(fmt.Sprintf("Success database connection, %s:%s", conf.Database.Host, conf.Database.Port))
	db.LogMode(true)
	db.SetLogger(logger.GetLogger())
	rep = &Repository{}
	rep.Db = db
}

// GetRepository returns the object of repository.
func GetRepository() *Repository {
	return rep
}

// GetDB returns the object of gorm.DB.
func GetDB() *gorm.DB {
	return rep.Db
}

// Find find records that match given conditions.
func (rep *Repository) Find(out interface{}, where ...interface{}) *gorm.DB {
	return rep.Db.Find(out, where...)
}

// Exec exec given SQL using by gorm.DB.
func (rep *Repository) Exec(sql string, values ...interface{}) *gorm.DB {
	return rep.Db.Exec(sql, values...)
}

// First returns first record that match given conditions, order by primary key.
func (rep *Repository) First(out interface{}, where ...interface{}) *gorm.DB {
	return rep.Db.First(out, where...)
}

// Raw returns the record that executed the given SQL using gorm.DB.
func (rep *Repository) Raw(sql string, values ...interface{}) *gorm.DB {
	return rep.Db.Raw(sql, values...)
}

// Create insert the value into database.
func (rep *Repository) Create(value interface{}) *gorm.DB {
	return rep.Db.Create(value)
}

// Save update value in database, if the value doesn't have primary key, will insert it.
func (rep *Repository) Save(value interface{}) *gorm.DB {
	return rep.Db.Save(value)
}

// Update update value in database
func (rep *Repository) Update(value interface{}) *gorm.DB {
	return rep.Db.Update(value)
}

// Delete delete value match given conditions.
func (rep *Repository) Delete(value interface{}) *gorm.DB {
	return rep.Db.Delete(value)
}

// Where returns a new relation.
func (rep *Repository) Where(query interface{}, args ...interface{}) *gorm.DB {
	return rep.Db.Where(query, args...)
}

// Preload preload associations with given conditions.
func (rep *Repository) Preload(column string, conditions ...interface{}) *gorm.DB {
	return rep.Db.Preload(column, conditions...)
}

// Scopes pass current database connection to arguments `func(*DB) *DB`, which could be used to add conditions dynamically
func (rep *Repository) Scopes(funcs ...func(*gorm.DB) *gorm.DB) *gorm.DB {
	return rep.Db.Scopes(funcs...)
}

// Transaction start a transaction as a block.
// If it is failed, will rollback and return error.
// If it is sccuessed, will commit.
// ref: https://github.com/jinzhu/gorm/blob/master/main.go#L533
func (rep *Repository) Transaction(fc func(tx *Repository) error) (err error) {
	panicked := true
	tx := rep.Db.Begin()
	defer func() {
		if panicked || err != nil {
			tx.Rollback()
		}
	}()

	txrep := &Repository{}
	txrep.Db = tx
	err = fc(txrep)

	if err == nil {
		err = tx.Commit().Error
	}

	panicked = false
	return
}
