package templates

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/G-Villarinho/level-up-api/internal"
)

type TemplateService interface {
	RenderTemplate(templateName string, params map[string]string) (string, error)
}

type templateService struct {
	di   *internal.Di
	path string
}

func NewTemplateService(di *internal.Di) (TemplateService, error) {
	dir, err := os.Getwd()
	if err != nil {
		return nil, err
	}

	return &templateService{
		di:   di,
		path: filepath.Join(dir, "templates"),
	}, nil
}

func (t *templateService) RenderTemplate(templateName string, params map[string]string) (string, error) {
	content, err := os.ReadFile(filepath.Join(t.path, templateName+".html"))
	if err != nil {
		return "", errors.New("read email template: " + err.Error())
	}

	template := string(content)
	for key, value := range params {
		placeholder := fmt.Sprintf("#%s#", key)
		template = strings.ReplaceAll(template, placeholder, value)
	}

	return template, nil
}
