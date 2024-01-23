package main

import (
	"html/template"
	"log/slog"
	"net/http"
)

const root string = `
<h2>Aetos</h2>
<ul>
<li><a href="/{{ . }}">metrics</a></li>
<li><a href="/healthz">healthz</a></li>
<li><a href="/varz">varz</a></li>
</ul>`

var tmpl *template.Template = template.Must(template.New("root").Parse(root))

func index(w http.ResponseWriter, _ *http.Request) {
	slog.Info("root")
	w.Header().Set("Content-Type", "text/html; charset=UTF-8")
	if err := tmpl.Execute(w, *path); err != nil {
		slog.Info("root", "error", err)
	}
}
