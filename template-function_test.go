package learning

import (
	"fmt"
	"html/template"
	"io"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

type MyPage struct {
	Name string
}

func (myPage MyPage) SayHello(name string) string {
	return "Hello " + name + ", my name is " + myPage.Name
}

func TemplateFunction(writer http.ResponseWriter, request *http.Request) {
	t := template.Must(template.New("function").Parse(`{{.SayHello "Flo"}}`))
	t.ExecuteTemplate(writer, "function", MyPage{
		Name: "Billy",
	})
}

func TestTemplateFunction(t *testing.T) {
	request := httptest.NewRequest("GET", "localhost:8080", nil)
	recorder := httptest.NewRecorder()

	TemplateFunction(recorder, request)

	body, _ := io.ReadAll(recorder.Result().Body)
	fmt.Println(string(body))
}

func TemplateFunctionGlobal(writer http.ResponseWriter, request *http.Request) {
	t := template.Must(template.New("function").Parse(`{{len .Name}}`))
	t.ExecuteTemplate(writer, "function", MyPage{
		Name: "Billy",
	})
}

func TestTemplateFunctionGlobal(t *testing.T) {
	request := httptest.NewRequest("GET", "localhost:8080", nil)
	recorder := httptest.NewRecorder()

	TemplateFunctionGlobal(recorder, request)

	body, _ := io.ReadAll(recorder.Result().Body)
	fmt.Println(string(body))
}

func TemplateCreateGlobal(writer http.ResponseWriter, request *http.Request) {
	t := template.New("function")
	t = t.Funcs(map[string]interface{}{
		"upper": func(v string) string {
			return strings.ToUpper(v)
		},
	})

	t = template.Must(t.Parse("{{upper .Name}}"))

	t.ExecuteTemplate(writer, "function", MyPage{
		Name: "Evanbill Antonio Kore",
	})
}

func TestCreateGlobal(t *testing.T) {
	request := httptest.NewRequest("GET", "localhost:8080", nil)
	recorder := httptest.NewRecorder()

	TemplateCreateGlobal(recorder, request)

	body, _ := io.ReadAll(recorder.Result().Body)
	fmt.Println(string(body))
}

func TemplateFunctionPipelines(writer http.ResponseWriter, request *http.Request) {
	t := template.New("function")
	t = t.Funcs(map[string]interface{}{
		"sayHello": func(name string) string {
			return "Hello " + name
		},
		"upper": func(v string) string {
			return strings.ToUpper(v)
		},
	})

	t = template.Must(t.Parse("{{sayHello .Name | upper}}"))

	t.ExecuteTemplate(writer, "function", MyPage{
		Name: "Florence Fedora Agustina",
	})
}

func TestFunctionPipelines(t *testing.T) {
	request := httptest.NewRequest("GET", "localhost:8080", nil)
	recorder := httptest.NewRecorder()

	TemplateFunctionPipelines(recorder, request)

	body, _ := io.ReadAll(recorder.Result().Body)
	fmt.Println(string(body))
}
