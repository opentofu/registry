package verification

import (
	"bytes"
	"fmt"
	"io"
	"text/template"
)

var templates *template.Template

func initTemplates() {
	templates = template.New("").Funcs(template.FuncMap{
		"status": func(s Status) string {
			switch s {
			case StatusSuccess:
				return "Success"
			case StatusFailure:
				return "Failure"
			case StatusNotRun:
				return "Not Run"
			case StatusSkipped:
				return "Skipped"
			default:
				return ""
			}
		},
		"call": func(tmpl string, args ...interface{}) string {
			t, err := templates.New(tmpl).Parse("{{template \"" + tmpl + "\" .}}")
			if err != nil {
				fmt.Printf("Error parsing template %s: %s", tmpl, err)
				return ""
			}
			var output bytes.Buffer
			err = t.Execute(&output, args[0])
			if err != nil {
				fmt.Printf("Error executing template %s: %s", tmpl, err)
				return ""
			}
			return output.String()
		},
	})

	_, err := templates.Parse(markdownTemplate)
	if err != nil {
		fmt.Printf("Error parsing templates: %s", err)
	}
}

func (r *Result) RenderMarkdown() string {
	initTemplates() // Initialize templates
	var output bytes.Buffer
	for _, step := range r.Steps {
		renderTemplate(&output, "step", step)
		for _, subStep := range step.SubSteps {
			renderTemplate(&output, "subStep", subStep)
		}
		output.WriteString("\n")
	}
	return output.String()
}

func renderTemplate(wr io.Writer, tmplName string, data interface{}) {
	err := templates.ExecuteTemplate(wr, tmplName, data)
	if err != nil {
		fmt.Printf("Error executing template: %s", err)
	}
}
