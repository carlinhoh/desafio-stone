package handlers

import (
	"net/http"
	"github.com/gorilla/mux"
	"strconv"
	"github.com/leoferlopes/desafio-stone/database"
	"github.com/leoferlopes/desafio-stone/model"
	"github.com/leoferlopes/desafio-stone/config"
	"encoding/json"
	"strings"
	"time"
	"fmt"
)

func checkErr(err error) bool{
	if err != nil {
		panic(err)
		return true
	}
	return false
}

func Index(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintln(w, "Bem Vindo! Leia a documentação em: https://github.com/leoferlopes/desafio-stone")
}

func authenticate(w http.ResponseWriter, r *http.Request) bool{
	token := r.FormValue("token")
	if token != config.Settings.Token{
		Err401(w, r)
		return false
	}
	return true
}

func GetInvoice(w http.ResponseWriter, r *http.Request) {
	if !authenticate(w, r){
		return
	}

	vars := mux.Vars(r)
	sId := vars["id"]

	if len(sId) > 0 {
		id, err := strconv.Atoi(sId)
		if err != nil{
			Err400(w, r)
			return
		}
		invoice, err := database.ReadById(id)
		if err != nil{
			Err500(w,r)
			return
		}
		if invoice == nil{
			Err404(w,r)
			return
		}

		w.Header().Set("Content-Type", "application/json; charset=UTF-8")
		if err := json.NewEncoder(w).Encode(invoice); checkErr(err) {
			Err500(w,r)
			return
		}
	} else{
		Err404(w, r)
		return
	}
}

func DeleteInvoice(w http.ResponseWriter, r *http.Request) {
	if !authenticate(w, r){
		return
	}

	vars := mux.Vars(r)
	sId := vars["id"]

	if len(sId) > 0 {
		id, err := strconv.ParseInt(sId, 10, 64)
		if err != nil{
			Err400(w, r)
			return
		}

		count, err := database.DeleteById(id)
		if err != nil{
			Err500(w,r)
			return
		}
		if count <= 0{
			Err404(w,r)
			return
		}

		resp := map[string]int64{
			"deleted": count,
		}

		w.Header().Set("Content-Type", "application/json; charset=UTF-8")
		w.WriteHeader(200) // ok

		if err := json.NewEncoder(w).Encode(resp); checkErr(err) {
			Err500(w,r)
			return
		}
	} else{
		Err404(w, r)
		return
	}
}

func GetInvoices(w http.ResponseWriter, r *http.Request) {
	if !authenticate(w, r){
		return
	}

	var i []model.Invoice

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
			Err400(w, r)
			return
		}
		month = &m
	} else {
		month = nil
	}
	if len(sYear) > 0{
		y, err := strconv.Atoi(sYear)
		if err != nil {
			Err400(w, r)
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

	orderBy := []database.OrderBy{}
	if len(sOrderBy) > 0 {
		orderParams := strings.Split(sOrderBy, ",")
		for i := 0; i < len(orderParams); i++ {
			switch orderParams[i] {
			case "month:asc":
				orderBy = append(orderBy, database.MONTH_ASC)
				break
			case "month:desc":
				orderBy = append(orderBy, database.MONTH_DESC)
				break
			case "year:asc":
				orderBy = append(orderBy, database.YEAR_ASC)
				break
			case "year:desc":
				orderBy = append(orderBy, database.YEAR_DESC)
				break
			case "document:asc":
				orderBy = append(orderBy, database.DOCUMENT_ASC)
				break
			case "document:desc":
				orderBy = append(orderBy, database.DOCUMENT_DESC)
				break
			default:
				Err400(w, r)
				return
			}
		}
	}
	if len(sPageSize) > 0{
		ps, err := strconv.Atoi(sPageSize)
		if err != nil || ps < 0{
			Err400(w, r)
			return
		}
		limit = ps
	}
	if len(sPage) > 0{
		p, err := strconv.Atoi(sPage)
		if err != nil || p < 0{
			Err400(w, r)
			return
		}
		position += p*limit
	}

	i, err := database.Read(month, year, document, orderBy, position, limit)
	if err != nil{
		Err500(w, r)
		return
	}
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")

	if err := json.NewEncoder(w).Encode(i); err != nil {
		panic(err)
	}
}

func PostInvoices(w http.ResponseWriter, r *http.Request) {
	if !authenticate(w, r){
		return
	}

	var sReferenceMonth, sReferenceYear, sDocument, sDescription, sAmount string
	sReferenceMonth = r.FormValue("month")
	sReferenceYear = r.FormValue("year")
	sDocument = r.FormValue("document")
	sDescription = r.FormValue("description")
	sAmount = r.FormValue("amount")
	if len(sReferenceMonth) > 0 && len(sReferenceYear) > 0 && len(sDocument) > 0 {
		i := new(model.Invoice)
		i.Document = sDocument
		var err error
		i.ReferenceMonth, err = strconv.Atoi(sReferenceMonth)
		if err != nil {
			Err400(w, r)
			return
		}
		i.ReferenceYear, err = strconv.Atoi(sReferenceYear)
		if err != nil {
			Err400(w, r)
			return
		}
		if len(sDescription) > 0{
			i.Description = &sDescription
		}
		if len(sAmount) > 0{
			amount, err := strconv.ParseFloat(sAmount, 64)
			if err != nil {
				Err400(w, r)
				return
			} else {
				i.Amount = &amount
			}
		}
		i.CreatedAt = time.Now().UTC()
		i.IsActive = 1

		err = database.Create(i)
		if err != nil{
			Err500(w, r)
			return
		}

		w.Header().Set("Content-Type", "application/json; charset=UTF-8")
		w.WriteHeader(201) // criado

		if err := json.NewEncoder(w).Encode(i); err != nil {
			panic(err)
		}
	} else {
		Err400(w, r)
		return
	}
}

func DeleteInvoices(w http.ResponseWriter, r *http.Request) {
	if !authenticate(w, r){
		return
	}

	var sMonth, sYear, sDocument string
	sMonth = r.FormValue("month")
	sYear = r.FormValue("year")
	sDocument = r.FormValue("document")

	var month, year *int
	var document *string

	if len(sMonth) > 0{
		m, err := strconv.Atoi(sMonth)
		if err != nil {
			Err400(w, r)
			return
		}
		month = &m
	} else {
		month = nil
	}
	if len(sYear) > 0{
		y, err := strconv.Atoi(sYear)
		if err != nil {
			Err400(w, r)
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

	count, err := database.Delete(month, year, document)

	if err != nil{
		Err500(w,r)
		return
	}
	resp := map[string]int64{
		"deleted": count,
	}

	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	w.WriteHeader(200) // ok

	if err := json.NewEncoder(w).Encode(resp); checkErr(err) {
		Err500(w,r)
		return
	}
}

func Err400(w http.ResponseWriter, r *http.Request) {
	resp := map[string]int{
		"error": 400,
	}

	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	w.WriteHeader(400)
	json.NewEncoder(w).Encode(resp)
}

func Err401(w http.ResponseWriter, r *http.Request) {
	resp := map[string]int{
		"error": 401,
	}

	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	w.WriteHeader(401)
	json.NewEncoder(w).Encode(resp)
}

func Err404(w http.ResponseWriter, r *http.Request) {
	resp := map[string]int{
		"error": 404,
	}

	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	w.WriteHeader(404)
	json.NewEncoder(w).Encode(resp)
}

func Err405(w http.ResponseWriter, r *http.Request) {
	resp := map[string]int{
		"error": 405,
	}

	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	w.WriteHeader(405)
	json.NewEncoder(w).Encode(resp)
}

func Err500(w http.ResponseWriter, r *http.Request) {
	resp := map[string]int{
		"error": 500,
	}

	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	w.WriteHeader(500)
	json.NewEncoder(w).Encode(resp)
}