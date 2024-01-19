package backend

import (
	"html/template"
	"log"
	"net/http"
)

type Application struct {
	SrvMux    *http.ServeMux
	Info_log  *log.Logger
	Error_log *log.Logger

	TemplateTextKz map[string]map[string]string
	TemplateTextRu map[string]map[string]string
	TemplateTextEn map[string]map[string]string

	TemplateCache map[string]*template.Template
	ReqInstances  map[string]*TemplateDataStruct
}

type TemplateDataStruct struct {
	Key          string
	Fname        string
	Lname        string
	Results      [6]string
	TemplateText map[string]map[string]string
	Links        map[string]string
	ResultsTxt   string
}

func (reqInstance *TemplateDataStruct) ConcStrings(s1, s2 string) string {
	return s1 + s2
}
