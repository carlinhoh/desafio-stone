package main

import (
	"fmt"
	"log"
	"net/http"

	// importando o gorilla web toolkit (http://www.gorillatoolkit.org/)
	"github.com/gorilla/mux"

	"encoding/json"

	"strconv"
	"strings"
	"time"
	"os"
	"io/ioutil"
)

var config Config

func main() {
	file, e := ioutil.ReadFile("./config.json")
	if e != nil {
		fmt.Printf("File error: %v\n", e)
		os.Exit(1)
	}

	err := json.Unmarshal(file, &config)
	if err != nil{
		fmt.Printf("File sintax error: %v\n", e)
		os.Exit(1)
	}

	dbFactory.DBConfig.Params = make(map[string]string)
	dbFactory.DBConfig.Params["charset"] = "utf8"
	dbFactory.DBConfig.Params["parseTime"] = "true"
	dbFactory.DBConfig.Addr = config.MySqlConfig.Address
	dbFactory.DBConfig.User = config.MySqlConfig.User
	dbFactory.DBConfig.Passwd = config.MySqlConfig.Password
	dbFactory.DBConfig.DBName = config.MySqlConfig.Schema

	fmt.Printf("Connection String: %s\n", dbFactory.DBConfig.FormatDSN())

	// router do gorilla web toolkit
	router := mux.NewRouter().StrictSlash(true)
	router.HandleFunc("/", Index)
	router.HandleFunc("/invoices", InvoicesPath)
	router.HandleFunc("/invoices/{id}", InvoicePath)

	log.Fatal(http.ListenAndServe(":8080", router))
}

func checkErr(err error) bool{
	if err != nil {
		panic(err)
		return true
	}
	return false
}

func Index(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintln(w, "Welcome!")
}

// TODO: token de aplicacao para autenticar
func InvoicesPath(w http.ResponseWriter, r *http.Request) {
	token := r.FormValue("document")
	if token != config.Token{
		err401(w, r)
		return
	}

	switch r.Method {
	case "GET":
		GetInvoices(w,r)
	case "POST":
		PostInvoices(w,r)
		break
	case "DELETE":
		DeleteInvoices(w,r)
		break
	default:
		err405(w, r)
		break
	}
}

func InvoicePath(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case "GET":
		GetInvoice(w,r)
	case "DELETE":
		DeleteInvoice(w,r)
		break
	default:
		err405(w, r)
		break
	}
}

func GetInvoice(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	sId := vars["id"]

	if len(sId) > 0 {
		id, err := strconv.Atoi(sId)
		if err != nil{
			err400(w, r)
			return
		}
		invoice, err := readById(id)
		if err != nil{
			err500(w,r)
			return
		}
		if invoice == nil{
			err404(w,r)
			return
		}

		w.Header().Set("Content-Type", "application/json; charset=UTF-8")
		if err := json.NewEncoder(w).Encode(invoice); checkErr(err) {
			err500(w,r)
			return
		}
	} else{
		err404(w, r)
		return
	}
}

func DeleteInvoice(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	sId := vars["id"]

	if len(sId) > 0 {
		id, err := strconv.ParseInt(sId, 10, 64)
		if err != nil{
			err400(w, r)
			return
		}

		count, err := deleteById(id)
		if err != nil{
			err500(w,r)
			return
		}

		resp := map[string]int64{
			"deleted": count,
		}

		w.Header().Set("Content-Type", "application/json; charset=UTF-8")
		w.WriteHeader(200) // ok

		if err := json.NewEncoder(w).Encode(resp); checkErr(err) {
			err500(w,r)
			return
		}
	} else{
		err404(w, r)
		return
	}
}

func GetInvoices(w http.ResponseWriter, r *http.Request) {
	var i []Invoice

	var sMonth, sYear, sDocument, sOrderBy, sPage, sPageSize string
	sMonth = r.FormValue("month")
	sYear = r.FormValue("year")
	sDocument = r.FormValue("document")
	sOrderBy = r.FormValue("orderBy")
	sPage = r.FormValue("page")
	sPageSize = r.FormValue("pageSize")

	var month, year *int
	var document *string
	position := 0
	limit := 10

	if len(sMonth) > 0{
		m, err := strconv.Atoi(sMonth)
		if err != nil {
			err400(w, r)
			return
		}
		month = &m
	} else {
		month = nil
	}
	if len(sYear) > 0{
		y, err := strconv.Atoi(sYear)
		if err != nil {
			err400(w, r)
			return
		}
		year = &y
	} else {
		year = nil
	}
	if len(sDocument) > 0{
		document = &sDocument
	} else {
		document = nil
	}

	orderBy := []OrderBy{}
	if len(sOrderBy) > 0 {
		orderParams := strings.Split(sOrderBy, ",")
		for i := 0; i < len(orderParams); i++ {
			switch orderParams[i] {
			case "month:asc":
				orderBy = append(orderBy, MONTH_ASC)
				break
			case "month:desc":
				orderBy = append(orderBy, MONTH_DESC)
				break
			case "year:asc":
				orderBy = append(orderBy, YEAR_ASC)
				break
			case "year:desc":
				orderBy = append(orderBy, YEAR_DESC)
				break
			case "document:asc":
				orderBy = append(orderBy, DOCUMENT_ASC)
				break
			case "document:desc":
				orderBy = append(orderBy, DOCUMENT_DESC)
				break
			default:
				err400(w, r)
				return
			}
		}
	}
	if len(sPageSize) > 0{
		ps, err := strconv.Atoi(sPageSize)
		if err != nil || ps < 0{
			err400(w, r)
			return
		}
		limit = ps
	}
	if len(sPage) > 0{
		p, err := strconv.Atoi(sPage)
		if err != nil || p < 0{
			err400(w, r)
			return
		}
		position += p*limit
	}

	i, err := read(month, year, document, orderBy, position, limit)
	if err != nil{
		err500(w, r)
		return
	}
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")

	if err := json.NewEncoder(w).Encode(i); err != nil {
		panic(err)
	}
}

func PostInvoices(w http.ResponseWriter, r *http.Request) {
	var sReferenceMonth, sReferenceYear, sDocument, sDescription, sAmount string
	sReferenceMonth = r.FormValue("month")
	sReferenceYear = r.FormValue("year")
	sDocument = r.FormValue("document")
	sDescription = r.FormValue("description")
	sAmount = r.FormValue("amount")
	if len(sReferenceMonth) > 0 && len(sReferenceYear) > 0 && len(sDocument) > 0 {
		i := new(Invoice)
		i.Document = sDocument
		var err error
		i.ReferenceMonth, err = strconv.Atoi(sReferenceMonth)
		if err != nil {
			err400(w, r)
			return
		}
		i.ReferenceYear, err = strconv.Atoi(sReferenceYear)
		if err != nil {
			err400(w, r)
			return
		}
		if len(sDescription) > 0{
			i.Description = &sDescription
		}
		if len(sAmount) > 0{
			amount, err := strconv.ParseFloat(sAmount, 64)
			if err != nil {
				err400(w, r)
				return
			} else {
				i.Amount = &amount
			}
		}
		i.CreatedAt = time.Now().UTC()
		i.IsActive = 1

		err = create(i)
		if err != nil{
			err500(w, r)
			return
		}

		w.Header().Set("Content-Type", "application/json; charset=UTF-8")
		w.WriteHeader(201) // criado

		if err := json.NewEncoder(w).Encode(i); err != nil {
			panic(err)
		}
	} else {
		err400(w, r)
		return
	}
}

func DeleteInvoices(w http.ResponseWriter, r *http.Request) {

	var sMonth, sYear, sDocument string
	sMonth = r.FormValue("month")
	sYear = r.FormValue("year")
	sDocument = r.FormValue("document")

	var month, year *int
	var document *string

	if len(sMonth) > 0{
		m, err := strconv.Atoi(sMonth)
		if err != nil {
			err400(w, r)
			return
		}
		month = &m
	} else {
		month = nil
	}
	if len(sYear) > 0{
		y, err := strconv.Atoi(sYear)
		if err != nil {
			err400(w, r)
			return
		}
		year = &y
	} else {
		year = nil
	}
	if len(sDocument) > 0{
		document = &sDocument
	} else {
		document = nil
	}

	count, err := delete(month, year, document)

	if err != nil{
		err500(w,r)
		return
	}
	resp := map[string]int64{
		"deleted": count,
	}

	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	w.WriteHeader(200) // ok

	if err := json.NewEncoder(w).Encode(resp); checkErr(err) {
		err500(w,r)
		return
	}
}

func err400(w http.ResponseWriter, r *http.Request) {
	resp := map[string]int{
		"error": 400,
	}

	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	w.WriteHeader(400)
	json.NewEncoder(w).Encode(resp)
}

func err401(w http.ResponseWriter, r *http.Request) {
	resp := map[string]int{
		"error": 400,
	}

	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	w.WriteHeader(400)
	json.NewEncoder(w).Encode(resp)
}

func err404(w http.ResponseWriter, r *http.Request) {
	resp := map[string]int{
		"error": 404,
	}

	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	w.WriteHeader(404)
	json.NewEncoder(w).Encode(resp)
}

func err405(w http.ResponseWriter, r *http.Request) {
	resp := map[string]int{
		"error": 405,
	}

	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	w.WriteHeader(405)
	json.NewEncoder(w).Encode(resp)
}

func err500(w http.ResponseWriter, r *http.Request) {
	resp := map[string]int{
		"error": 500,
	}

	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	w.WriteHeader(500)
	json.NewEncoder(w).Encode(resp)
}
/*
func err501(w http.ResponseWriter, r *http.Request) {
	resp := map[string]int{
		"error": 501,
	}

	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	w.WriteHeader(501)
	json.NewEncoder(w).Encode(resp)
}
*/