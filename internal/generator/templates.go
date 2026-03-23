package generator

import (
	"embed"
	"io"
	"os"
)

//go:embed templates/*
var templatesFS embed.FS

func CopyTemplate(templateName, destPath string) error {
	src, err := templatesFS.Open("templates/" + templateName + ".md")
	if err != nil {
		return err
	}
	defer src.Close()

	dest, err := os.Create(destPath)
	if err != nil {
		return err
	}
	defer dest.Close()

	_, err = io.Copy(dest, src)
	return err
}

func GetTemplateContent(name string) (string, error) {
	content, err := templatesFS.ReadFile("templates/" + name + ".md")
	return string(content), err
}
