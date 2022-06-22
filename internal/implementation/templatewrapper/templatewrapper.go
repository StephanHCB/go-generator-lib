package templatewrapper

import (
	"github.com/Masterminds/sprig"
	"io"
	"text/template"
)

type TemplateWrapper struct {
	isRawFile       bool
	templateContent []byte
	templateName    string
	templatePath    string
	tmpl            *template.Template
}

// New allocates a new templateWrapper with the given name.
func New(isRawFile bool, templateContent []byte, templateName string, templatePath string) *TemplateWrapper {
	t := &TemplateWrapper{
		isRawFile:       isRawFile,
		templateContent: templateContent,
		templateName:    templateName,
		templatePath:    templatePath,
	}
	return t
}

func (i *TemplateWrapper) Write(wr io.Writer, name string, data interface{}) error {
	if i.isRawFile {
		_, err := wr.Write(i.templateContent)
		return err
	} else {
		return i.tmpl.ExecuteTemplate(wr, name, data)
	}
}

func (i *TemplateWrapper) Parse() (*TemplateWrapper, error) {
	if !i.isRawFile && i.tmpl == nil {
		tmpl, err := template.New(i.templateName).Funcs(sprig.TxtFuncMap()).Parse(string(i.templateContent))
		i.tmpl = tmpl
		return i, err
	}
	return i, nil
}
