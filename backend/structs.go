package backend

import (
	"html/template"
	"log"
)

type Application struct {
	Info_log      *log.Logger
	Error_log     *log.Logger
	templateCache map[string]*template.Template
}

type TemplateData struct {
	test string
}
