package database

import (
	"database/sql"
	_ "github.com/go-sql-driver/mysql"
	"log"
	"github.com/go-sql-driver/mysql"
	"time"

	"github.com/leoferlopes/desafio-stone/model"
)

type OrderBy uint8

const (
	MONTH_ASC OrderBy = iota
	MONTH_DESC OrderBy = iota
	YEAR_ASC OrderBy = iota
	YEAR_DESC OrderBy = iota
	DOCUMENT_ASC OrderBy = iota
	DOCUMENT_DESC OrderBy = iota
)


type InvoiceDAO struct {

}

func checkErr(err error) bool{
	if err != nil {
		panic(err)
		return true
	}
	return false
}

func Read(month *int, year *int, document *string, orderBy []OrderBy, position int, size int) ([]model.Invoice, error){
	db, err := /*DbFactory.*/GetInstance()
	if checkErr(err) {
		return nil, err
	}

	var orderByClause string
	if(len(orderBy) > 0){
		orderByClause = " ORDER BY "
		for i := 0; i < len(orderBy); i++ {
			if i > 0{
				orderByClause += ", "
			}
			switch orderBy[i] {
			case MONTH_ASC:
				orderByClause += "ReferenceMonth ASC"
				break
			case MONTH_DESC:
				orderByClause += "ReferenceMonth DESC"
				break
			case YEAR_ASC:
				orderByClause += "ReferenceYear ASC"
				break
			case YEAR_DESC:
				orderByClause += "ReferenceYear DESC"
				break
			case DOCUMENT_ASC:
				orderByClause += "Document ASC"
				break
			case DOCUMENT_DESC:
				orderByClause += "Document DESC"
				break
			}
		}
	}else{
		orderByClause = ""
	}
	orderByClause += " LIMIT ?,?"
	log.Printf("LIMIT %d,%d", position, size)

	query := "SELECT Id, CreatedAt, ReferenceMonth, ReferenceYear, Document, Description, Amount, IsActive, DeactiveAt FROM mydb.Invoice WHERE IsActive <> 0"

	var rows *sql.Rows

	if month != nil {
		query += " AND ReferenceMonth = ?"
		if year != nil {
			query += " AND ReferenceYear = ?"
			if document != nil {
				query += " AND Document = ?"
				rows, err = db.Query(query + orderByClause, month, year, document, position, size)
			}else {
				rows, err = db.Query(query + orderByClause, month, year, position, size)
			}
		}else{
			if document != nil {
				query += " AND Document = ?"
				rows, err = db.Query(query + orderByClause, month, document, position, size)
			}else {
				rows, err = db.Query(query + orderByClause, month, position, size)
			}
		}
	}else{
		if year != nil {
			query += " AND ReferenceYear = ?"
			if document != nil {
				query += " AND Document = ?"
				rows, err = db.Query(query + orderByClause, year, document, position, size)
			}else {
				rows, err = db.Query(query + orderByClause, year, position, size)
			}
		}else{
			if document != nil {
				query += " AND Document = ?"
				rows, err = db.Query(query + orderByClause, document, position, size)
			}else {
				rows, err = db.Query(query + orderByClause, position, size)
			}
		}
	}

	if checkErr(err) {
		return nil, err
	}

	invoices := []model.Invoice{}
	for rows.Next() {
		invoice := new(model.Invoice)
		description := new(sql.NullString)
		amount := new(sql.NullFloat64)
		deactiveAt := new(mysql.NullTime)
		invoice.Id = new(int64)
		err = rows.Scan(invoice.Id,
			&(invoice.CreatedAt),
			&(invoice.ReferenceMonth),
			&(invoice.ReferenceYear),
			&(invoice.Document),
			description,
			amount,
			&(invoice.IsActive),
			deactiveAt)
		switch {
		case description.Valid:
			invoice.Description = &description.String
			break
		case amount.Valid:
			invoice.Amount = &amount.Float64
			break
		case deactiveAt.Valid:
			invoice.DeactiveAt = &deactiveAt.Time
			break
		}
		checkErr(err)

		invoices = append(invoices,*invoice)
	}
	return invoices, nil
}

func ReadById(id int) (*model.Invoice, error){
	db, err := /*DbFactory.*/GetInstance()
	if checkErr(err) {
		return nil, err
	}

	query := "SELECT Id, CreatedAt, ReferenceMonth, ReferenceYear, Document, Description, Amount, IsActive, DeactiveAt FROM mydb.Invoice WHERE IsActive <> 0 AND Id = ?"

	invoice := new(model.Invoice)
	description := new(sql.NullString)
	amount := new(sql.NullFloat64)
	deactiveAt := new(mysql.NullTime)
	invoice.Id = new(int64)
	err = db.QueryRow(query, id).Scan(
		invoice.Id,
		&(invoice.CreatedAt),
		&(invoice.ReferenceMonth),
		&(invoice.ReferenceYear),
		&(invoice.Document),
		description,
		amount,
		&(invoice.IsActive),
		deactiveAt)
	switch {
	case err == sql.ErrNoRows:
		return nil, nil
	case err != nil:
		checkErr(err)
		return nil, err
	case description.Valid:
		invoice.Description = &description.String
		break
	case amount.Valid:
		invoice.Amount = &amount.Float64
		break
	case deactiveAt.Valid:
		invoice.DeactiveAt = &deactiveAt.Time
		break
	}
	return invoice, nil
}

func Delete(month *int, year *int, document *string) (int64, error){
	db, err := /*DbFactory.*/GetInstance()
	if checkErr(err) {
		return 0, err
	}

	query := "UPDATE `mydb`.`Invoice` SET `IsActive` = 0, `DeactiveAt` = ? WHERE IsActive <> 0"

	var res sql.Result

	if month != nil {
		query += " AND ReferenceMonth = ?"
		if year != nil {
			query += " AND ReferenceYear = ?"
			if document != nil {
				query += " AND Document = ?"
				res, err = db.Exec(query,time.Now(), month, year, document)
			}else {
				res, err = db.Exec(query,time.Now(), month, year)
			}
		}else{
			if document != nil {
				query += " AND Document = ?"
				res, err = db.Exec(query,time.Now(), month, document)
			}else {
				res, err = db.Exec(query,time.Now(), month)
			}
		}
	}else{
		if year != nil {
			query += " AND ReferenceYear = ?"
			if document != nil {
				query += " AND Document = ?"
				res, err = db.Exec(query,time.Now(), year, document)
			}else {
				res, err = db.Exec(query,time.Now(), year)
			}
		}else{
			if document != nil {
				query += " AND Document = ?"
				res, err = db.Exec(query,time.Now(), document)
			}else {
				res, err = db.Exec(query,time.Now())
			}
		}
	}

	if checkErr(err) {
		return 0, err
	}

	rowCnt, err := res.RowsAffected()
	if checkErr(err) {
		return rowCnt, err
	}
	log.Printf("DELETE affected = %d\n", rowCnt)
	return rowCnt, err
}

func DeleteById(id int64) (int64, error){
	db, err := /*DbFactory.*/GetInstance()
	if checkErr(err) {
		return 0, err
	}

	stmt, err := db.Prepare("UPDATE `mydb`.`Invoice` SET `IsActive` = 0, `DeactiveAt` = ? WHERE `Id` = ? AND IsActive <> 0;")
	if checkErr(err) {
		return 0, err
	}
	res, err := stmt.Exec(time.Now(), id)
	if checkErr(err) {
		return 0, err
	}
	rowCnt, err := res.RowsAffected()
	if checkErr(err) {
		return rowCnt, err
	}
	log.Printf("DELETE ID = %d, affected = %d\n", id, rowCnt)
	return rowCnt, err
}

//INSERT INTO `mydb`.`Invoice` (`Id`, `CreatedAt`, `ReferenceMonth`, `ReferenceYear`, `Document`, `Description`, `Amount`, `IsActive`) VALUES (NULL, '2016-08-28 12:00', '8', '2016', 'documento', 'descrição', '22.00', '1');
func Create(invoice *model.Invoice) error{
	db, err := /*DbFactory.*/GetInstance()
	if checkErr(err) {
		return err
	}

	stmt, err := db.Prepare("INSERT INTO `mydb`.`Invoice` (`Id`, `CreatedAt`, `ReferenceMonth`, `ReferenceYear`, `Document`, `Description`, `Amount`, `IsActive`) VALUES (?, ?, ?, ?, ?, ?, ?, ?);")
	if err != nil {
		log.Fatal(err)
	}
	res, err := stmt.Exec(invoice.Id, invoice.CreatedAt, invoice.ReferenceMonth, invoice.ReferenceYear, invoice.Document, invoice.Description, invoice.Amount, invoice.IsActive)
	if err != nil {
		log.Fatal(err)
	}
	lastId, err := res.LastInsertId()
	if err != nil {
		log.Fatal(err)
	} else{
		invoice.Id = &lastId
	}
	rowCnt, err := res.RowsAffected()
	if err != nil {
		log.Fatal(err)
	}
	log.Printf("ID = %d, affected = %d\n", lastId, rowCnt)
	return nil
}