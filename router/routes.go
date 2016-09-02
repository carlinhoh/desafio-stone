package router

import (
	"net/http"
	"github.com/leoferlopes/desafio-stone/handlers"
)

type Route struct {
	Name		string
	Method		string
	Pattern		string
	HandlerFunc	http.HandlerFunc
}

type Routes []Route

var routes = Routes{
	Route{
		"Index",
		"GET",
		"/",
		handlers.Index,
	},
	Route{
		"GetInvoices",
		"GET",
		"/invoices",
		handlers.GetInvoices,
	},
	Route{
		"PostInvoices",
		"POST",
		"/invoices",
		handlers.PostInvoices,
	},
	Route{
		"DeleteInvoices",
		"DELETE",
		"/invoices",
		handlers.DeleteInvoices,
	},
	Route{
		"405",
		"",
		"/invoices",
		handlers.Err405,
	},
	Route{
		"GetInvoice",
		"GET",
		"/invoices/{id}",
		handlers.GetInvoice,
	},
	Route{
		"DeleteInvoice",
		"DELETE",
		"/invoices/{id}",
		handlers.DeleteInvoice,
	},
	Route{
		"405",
		"",
		"/invoices/{id}",
		handlers.Err405,
	},
	Route{
		"404",
		"",
		"/{path}",
		handlers.Err404,
	},
}