package handlers

/*import (
	"net/http"
	//"os"
	"text/template"
)

type TextHandler struct{}

func (e TextHandler) ServeHTTP(rw http.ResponseWriter, r *http.Request) {
	nameParam := r.FormValue("name")
	sayParam := r.FormValue("say")

	if sayParam == "Nothing" {
		rw.WriteHeader(404)
	} else {
		rw.Write([]byte(sayParam))
		//rw.Write([]byte("hello!\n"))
	}

	t, err := template.New("person").Parse(personTemplate)
	if err != nil {
		panic(err)
	}

	var people = []Person{
		{"John", "Smith", 22},
		{"Alice", "Smith", 25},
		{"Bob", "Baker", 24},
	}

	err = t.Execute(rw, people)
	if err != nil {
		panic(err)
	}

}

type Person struct {
	FirstName string
	LastName  string
	Age       int
}

const personTemplate = `{{range .}}
{{.FirstName}} {{.LastName}} is {{.Age}} years old.
{{end}}`
*/
/*func main() {
	t, err := template.New("person").Parse(personTemplate)
	if err != nil {
		panic(err)
	}
	err = t.Execute(os.Stdout, Person{
		FirstName: "John",
		LastName:  "Smith",
		Age:       22,
	})
	if err != nil {
		panic(err)
	}
}
*/
