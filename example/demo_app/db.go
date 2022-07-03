package main

import (
	"database/sql"
	"errors"
	"log"
	"os"

	_ "github.com/mattn/go-sqlite3"
)

type DB struct {
	store *sql.DB
}

func NewDB() *DB {
	d := &DB{}

	d.init()
	return d
}

func (d *DB) init() {
	// creates the db file if it doesn't exist
	if _, err := os.Stat("demo-app.db"); errors.Is(err, os.ErrNotExist) {
		file, err := os.Create("demo-app.db")
		if err != nil {
			log.Fatal(err.Error())
		}
		file.Close()
	}

	sqliteDatabase, _ := sql.Open("sqlite3", "./demo-app.db")

	d.store = sqliteDatabase
}

func (d *DB) Close() {
	d.store.Close()
}

func (d *DB) setupEmployeeTable() error {
	query := `CREATE TABLE employee (
        "id" INTEGER NOT NULL PRIMARY KEY AUTOINCREMENT,
        "first_name" TEXT,
        "last_name" TEXT,
        "company_email" TEXT DEFAULT "" ,
        "personal_email" TEXT,
        "phone" TEXT,
        "date_of_birth" TEXT,
        "start_date" TEXT DEFAULT "",
        "term_date" TEXT DEFAULT "",
        "department" TEXT DEFAULT "",
        "manager" INTEGER DEFAULT -1,
        "salary" INTEGER DEFAULT 0,
        "medical" TEXT DEFAULT "",
        "profile_picture" BLOB,
        "dl_picture" BLOB,
        "ssn" TEXT DEFAULT "",
        "status" TEXT DEFAULT 1,
        FOREIGN KEY(manager) REFERENCES employee(id));`

	log.Println("[DB.setupEmployeeTable] Creating employee table...")

	statement, err := d.store.Prepare(query)
	if err != nil {
		log.Println("[DB.setupEmployeeTable] prepare table failed")
		return err
	}
	_, err = statement.Exec()
	if err != nil {
		log.Println("[DB.setupEmployeeTable] create table failed")
		return err
	}

	log.Println("[DB.setupEmployeeTable] table created.")

	query1 := `CREATE UNIQUE INDEX idx_employee_pemail ON employee(personal_email)`
	statement, err = d.store.Prepare(query1)
	if err != nil {
		log.Println("[DB.setupEmployeeTable] prepare idx1 failed")
		return err
	}
	_, err = statement.Exec()
	if err != nil {
		log.Println("[DB.setupEmployeeTable] create idx1 failed")
		return err
	}

	return nil
}

func (d *DB) addEmployee(emp employee) error {
	query := `INSERT INTO employee ( first_name, last_name, personal_email, phone,
        date_of_birth, ssn, status) VALUES (?, ?, ?, ?, ?, ?, ?)`

	log.Println("[DB.addEmployee] Preparing to insert to employee")

	statement, err := d.store.Prepare(query)
	if err != nil {
		log.Println("[DB.addEmployee] prepare failed")
		return err
	}
	_, err = statement.Exec(emp.FirstName, emp.LastName, emp.PersonalEmail,
		emp.Phone, emp.DateOfBirth, emp.SSN, emp.Status)
	if err != nil {
		log.Println("[DB.addEmployee] insert failed")
		return err
	}

	log.Println("[DB.addEmployee] employee added")
	return nil
}

func (d *DB) lookupEmployee(email string) (*employee, error) {
	query := `SELECT first_name, last_name, company_email, personal_email,
        date_of_birth, start_date, term_date, department, manager,
        salary, medical, ssn, status, id
        FROM employee
        WHERE personal_email=? OR company_email=?`

	log.Println("[DB.lookupEmployee] looking for employee ", email)

	row := d.store.QueryRow(query, email, email)

	var emp employee

	err := row.Scan(&emp.FirstName, &emp.LastName, &emp.CompanyEmail, &emp.PersonalEmail,
		&emp.DateOfBirth, &emp.StartDate, &emp.TermDate, &emp.Department,
		&emp.ManagerId, &emp.Salary, &emp.Medical, &emp.SSN, &emp.Status, &emp.ID)

	if err != nil {
		log.Println("[DB.lookupEmployee] failed")
		return nil, err
	}

	return &emp, nil
}

func (d *DB) updateEmployee(emp employee) error {
	query := `UPDATE employee
        SET company_email = ?,
            start_date = ?,
            term_date = ?,
            department = ?,
            manager = ?,
            salary = ?,
            status = ?
        WHERE id = ?`

	log.Println("[DB.updateEmployee] Updating employee ", emp.PersonalEmail)

	statement, err := d.store.Prepare(query)
	if err != nil {
		log.Println("[DB.updateEmployee] prepare failed")
		return err
	}

	_, err = statement.Exec(emp.CompanyEmail, emp.StartDate, emp.TermDate, emp.Department,
		emp.ManagerId, emp.Salary, emp.Status, emp.ID)
	if err != nil {
		log.Println("[DB.updateEmployee] update failed")
		return err
	}

	log.Println("[DB.updateEmployee] Record updated")
	return nil
}
