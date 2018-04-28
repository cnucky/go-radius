package db

import (
	"database/sql"
	_ "github.com/go-sql-driver/mysql"
	"github.com/hel2o/go-radius/g"
	"log"
)

var RadiusDb, FireSystemDb *sql.DB

func InitDB() {
	radiusDb, err := sql.Open("mysql", g.Config().GoRadius.RadiusDb)
	if err != nil {
		log.Println(err)
		return
	}
	err = radiusDb.Ping()
	if err != nil {
		log.Println(err)
		return
	}
	RadiusDb = radiusDb
	fireSystemDb, err := sql.Open("mysql", g.Config().GoRadius.FireSystemDb)
	if err != nil {
		log.Println(err)
		return
	}
	err = fireSystemDb.Ping()
	if err != nil {
		log.Println(err)
		return
	}
	FireSystemDb = fireSystemDb
}

//对比radius用户名和密码是否一致
func CheckUserPassword(db *sql.DB, username, password string) bool {
	defer func() {
		if err := recover(); err != nil {
			log.Println(err)
		}
	}()
	var count int64
	q := "SELECT COUNT(*) FROM radcheck WHERE username = ? AND value = ?"
	err := db.QueryRow(q, username, password).Scan(&count)
	if err != nil {
		log.Println(err)
		return false
	}
	if count == 1 {
		return true
	}
	return false
}

//记录用户登录成功
func Login(db *sql.DB, userName, password, nasIPAddress, nasIdentifier, framedIPAddress, acctSessionId string) (int64, error) {
	defer func() {
		if err := recover(); err != nil {
			log.Println(err)
		}
	}()
	var w = "INSERT INTO `radpostauth` (username, pass, reply, authdate, nasipaddress, clientipaddress, nasidentifier, acctstarttime, acctsessionid) VALUES (?,?,'Access-Accept',now(),?,?,?,now(),?)"
	stmt, err := db.Prepare(w)
	defer stmt.Close()
	HandleErr(err)
	ret, err := stmt.Exec(userName, password, nasIPAddress, framedIPAddress, nasIdentifier, acctSessionId)
	if err != nil {
		log.Println(err)
		return 0, err
	}
	rowsAffected, err := ret.RowsAffected()
	HandleErr(err)
	return rowsAffected, err
}

//记录用户退出登录
func Logout(db *sql.DB, framedIPAddress, acctSessionId string) (int64, error) {
	defer func() {
		if err := recover(); err != nil {
			log.Println(err)
		}
	}()
	var w = "UPDATE `radpostauth` SET acctstoptime = now() WHERE acctsessionid = ? AND clientipaddress =?"
	stmt, err := db.Prepare(w)
	defer stmt.Close()
	HandleErr(err)
	ret, err := stmt.Exec(acctSessionId, framedIPAddress)
	if err != nil {
		log.Println(err)
		return 0, err
	}
	rowsAffected, err := ret.RowsAffected()
	HandleErr(err)
	return rowsAffected, err
}

//记录用户登录失败
func LoginFail(db *sql.DB, userName, password, nasIPAddress, nasIdentifier, framedIPAddress string) (int64, error) {
	defer func() {
		if err := recover(); err != nil {
			log.Println(err)
		}
	}()
	var w = "INSERT INTO `radpostauth` (username, pass, reply, authdate, nasipaddress, clientipaddress, nasidentifier) VALUES (?,?,'Access-Reject',now(),?,?,?)"
	stmt, err := db.Prepare(w)
	defer stmt.Close()
	HandleErr(err)
	ret, err := stmt.Exec(userName, password, nasIPAddress, framedIPAddress, nasIdentifier)
	if err != nil {
		log.Println(err)
		return 0, err
	}
	rowsAffected, err := ret.RowsAffected()
	HandleErr(err)
	return rowsAffected, err
}

func HandleErr(err error) {
	if err != nil {
		log.Println(err)
		return
	}
}
