package backend

import (
	"html/template"
	"log"
)

type Application struct {
	Info_log      *log.Logger
	Error_log     *log.Logger
	TemplateCache map[string]*template.Template
	TemplateData  *TemplateDataStruct

	TemplateTextKz map[string]map[string]string
	TemplateTextRu map[string]map[string]string
	TemplateTextEn map[string]map[string]string
}

type TemplateDataStruct struct {
	Fname        string
	Lname        string
	TemplateText map[string]map[string]string
}
