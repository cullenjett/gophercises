package cyoa

import (
	"encoding/json"
	"io"
	"net/http"
	"text/template"
)

var defaultHandlerTmpl = `
<!DOCTYPE html>
<html>
<head>
  <meta charset="utf-8">
  <meta name="viewport" content="width=device-width, initial-scale=1">
  <title>Choose your own adventure</title>
</head>
<body>
  <h1>{{.Title}}</h1>

  {{range .Paragraphs}}
  <p>{{.}}</p>
  {{end}}

  <ul>
    {{range .Options}}
    <li>
      <a href="/{{.Chapter}}">{{.Text}}</a>
    </li>
    {{end}}
  </ul>
</body>
</html>`

func NewHandler(s Story) http.Handler {
	return handler{
		s,
	}
}

type handler struct {
	s Story
}

func (h handler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	tpl := template.Must(template.New("").Parse(defaultHandlerTmpl))

	err := tpl.Execute(w, h.s["intro"])
	if err != nil {
		panic(err)
	}

}

func JsonStory(r io.Reader) (Story, error) {
	d := json.NewDecoder(r)

	var story Story
	if err := d.Decode(&story); err != nil {
		return nil, err
	}

	return story, nil
}

type Story map[string]Chapter

type Chapter struct {
	Title      string    `json:"title"`
	Paragraphs []string  `json:"story"`
	Options    []Options `json:"options"`
}

type Options struct {
	Text    string `json:"text"`
	Chapter string `json:"arc"`
}
