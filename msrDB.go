package main

import (
	"database/sql"
	"log"
	"os"

	_ "github.com/mattn/go-sqlite3"
)

var DB *sql.DB
var dbOpen = false

/*
*******************
admin/msr
guest/msrguest
********************
*/
const (
	sqlCreate     = `CREATE DATABASE IF NOT EXISTS msr`
	sqlPrepareM   = `CREATE TABLE IF NOT EXISTS members (Id INTEGER NOT NULL PRIMARY KEY, Rank INTEGER, Name TEXT NOT NULL, Since Text, Phone TEXT, Address TEXT, City TEXT, State TEXT, Zip NUMERIC, Email TEXT, BirthDay TEXT, Notes TEXT, Status INTEGER)`
	sqlPrepareU   = `CREATE TABLE IF NOT EXISTS users (name text, pwd text, utype int)`
	sqlCleanup    = `DELETE FROM users where name='admin' and pwd='38b33779833838a98c2a241ce465fb07'`
	sqlCheckAdmin = `INSERT INTO users(name, pwd, utype) SELECT 'admin', 'b25f6e74c90cf2294e14997d593c6e0e', 9 WHERE NOT EXISTS(SELECT 1 FROM users WHERE name='admin' and pwd='b25f6e74c90cf2294e14997d593c6e0e')`
	sqlCheckGuest = `INSERT INTO users(name, pwd, utype) SELECT 'guest', '43265a3c6387a820b59e14069089be46', 1 WHERE NOT EXISTS(SELECT 1 FROM users WHERE name='guest')`
	sqlMembers    = `SELECT Id, Name, Rank, Since, Phone, Address, City, State, Zip, Email, BirthDay, Status, Notes FROM members`
	sqlMembersS   = `SELECT Name, Email, Phone FROM members`
	sqlGetMember  = `SELECT Id, Name, Rank, Since, Phone, Address, City, State, Zip, Email, BirthDay, Status, Notes  FROM members WHERE id = ?`
	sqlMemberSave = `UPDATE members set Name=?, Rank=?, Since=?, Phone=?, Address=?, City=?, State=?, Zip=?, Email=?, BirthDay=?, Status=?, Notes=? where id=?`
	sqlMemberAdd  = `INSERT into members (Name, Rank, Since, Phone, Address, City, State, Zip, Email, BirthDay, Status, Notes) values (?,?,?,?,?,?,?,?,?,?,?,?) RETURNING id`
	sqlLogin      = `select utype from users where name=$1 and pwd=$2;`
)

func checkErr(err error) {
	if err != nil {
		panic(err)
	}
}

// OpenDB : Open the Grooming database
func OpenDB() *sql.DB {
	log.Println("Create dir " + DBDir)
	err := os.Mkdir(DBDir, 0777)
	if err != nil {
		log.Printf("Error %s creating %s", err.Error(), DBDir)
	}
	log.Println("Open database sqlite3, " + DBConnection)
	db, err := sql.Open("sqlite3", DBConnection)
	checkErr(err)
	err = db.Ping()
	checkErr(err)
	log.Println("Create member table")
	_, err = db.Exec(sqlPrepareM)
	checkErr(err)
	log.Println("Create user table")
	_, err = db.Exec(sqlPrepareU)
	checkErr(err)

	log.Println("DB Cleanup")
	_, err = db.Exec(sqlCleanup)
	checkErr(err)
	log.Println("Insert admin user") // admin/msr@dmin
	_, err = db.Exec(sqlCheckAdmin)
	checkErr(err)
	log.Println("Insert guest user") // guest/msrguest
	_, err = db.Exec(sqlCheckGuest)
	checkErr(err)

	log.Println("DB is now Open")
	dbOpen = true
	return db
}

func GetMember(mid int) Member {
	var r Member
	if !dbOpen {
		DB = OpenDB()
	}
	row := DB.QueryRow(sqlGetMember, mid)
	err := row.Scan(&r.ID, &r.Name, &r.Rank, &r.Since, &r.Phone, &r.Address, &r.City, &r.State, &r.Zip, &r.Email, &r.BirthDay, &r.Status, &r.Notes)
	checkErr(err)
	return r
}

func GetMembers() MemberList {
	var lMembers MemberList
	if !dbOpen {
		DB = OpenDB()
	}
	//ins, _ := DB.Prepare("INSERT  into members (Name, Email, NickName, Prospect, Notes, Address, Zip, City) VALUES (?, ?, ?, ?, ?, ?, ?, ?)")
	//ins.Exec("Halldor", "aatlason@gmail.com", "Axel", 0, "", "201 Countryside Ln", 38040, "Halls")
	log.Println("DB Query=" + sqlMembers)
	rows, err := DB.Query(sqlMembers)
	checkErr(err)
	defer rows.Close()
	log.Println("populate the data")
	for rows.Next() {
		var r Member
		err = rows.Scan(&r.ID, &r.Name, &r.Rank, &r.Since, &r.Phone, &r.Address, &r.City, &r.State, &r.Zip, &r.Email, &r.BirthDay, &r.Status, &r.Notes)
		checkErr(err)
		lMembers = append(lMembers, r)
	}
	return lMembers
}

// UpdMember -- Update or add a single member
func UpdMember(mb Member) int {
	id := mb.ID
	if !dbOpen {
		DB = OpenDB()
	}
	if id > 0 {
		log.Println("DB Query=" + sqlMemberSave)
		log.Printf("%s,%d,%s,%s,%s,%s,%s,%d,%s,%s,%d,%s,%d", mb.Name, mb.Rank, mb.Since, mb.Phone, mb.Address, mb.City, mb.State, mb.Zip, mb.Email, mb.BirthDay, mb.Status, mb.Notes, mb.ID)
		_, err := DB.Exec(sqlMemberSave, mb.Name, mb.Rank, mb.Since, mb.Phone, mb.Address, mb.City, mb.State, mb.Zip, mb.Email, mb.BirthDay, mb.Status, mb.Notes, mb.ID)
		checkErr(err)
	} else {
		err := DB.QueryRow(sqlMemberAdd, mb.Name, mb.Rank, mb.Since, mb.Phone, mb.Address, mb.City, mb.State, mb.Zip, mb.Email, mb.BirthDay, mb.Status, mb.Notes).Scan(&id)
		checkErr(err)
	}
	return id
}

func GetMembersS() MemberSList {
	var lMembers MemberSList
	if !dbOpen {
		DB = OpenDB()
	}
	log.Println("DB Query=" + sqlMembersS)
	rows, err := DB.Query(sqlMembersS)
	checkErr(err)
	defer rows.Close()
	log.Println("populate the data")
	for rows.Next() {
		var r MemberS
		err = rows.Scan(&r.Name, &r.Email, &r.Phone)
		checkErr(err)
		lMembers = append(lMembers, r)
	}
	return lMembers
}

// LoginUser Check if the user can be authenticated
func LoginUser(uname string, pwd string) int {
	var utype int
	if !dbOpen {
		DB = OpenDB()
	}
	log.Printf("Login, uname=%s, pwd=%s, Qry=%s", uname, pwd, sqlLogin)
	row := DB.QueryRow(sqlLogin, uname, pwd)
	switch err := row.Scan(&utype); err {
	case sql.ErrNoRows:
		utype = 0
		log.Println("Login, Not found in DB")
	case nil:
		// We got all we need
	default:
		panic(err)
	}
	return utype
}
