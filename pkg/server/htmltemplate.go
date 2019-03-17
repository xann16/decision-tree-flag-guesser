package server

import (
	"html/template"
	"io"
)

var tpl = template.Must(template.New("flaggerPage").Parse(`
<html>
<head>
<link rel="stylesheet" type="text/css" href="/static/style.css">
<title>{{.Title}}</title>
</head>
<body>
<div class="grid-container flag-grid">
{{range .Flags}}
<div class="grid-item {{.ItemClass}} tooltip">
{{if .IsButton}} <a href="/walk/{{.Tag}}/1"> {{end}}
  <img src="/static/img/{{.ThumbDir}}/{{.Tag}}.png">
{{if .IsButton}} </a> {{end}}
<span class="tooltiptext"> {{.Name}} </span>
</div>
{{end}}
</div>
<div class="grid-container label-grid">
<div class="grid-item active-chosen-flag">
<img src="/static/img/flags/{{.Choice}}.png" height="150px">
</div>
<div class="grid-item text-item">
<p class="main-question"> {{.Question}} </p>
{{if .IsWalk}}
<p class="main-links">
<a href=/walk/{{.Choice}}/{{.YesNode}}> Tak </a> &nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;
<a href=/walk/{{.Choice}}/{{.NoNode}}> Nie </a>
</p>
{{end}}
{{if .IsFinal}}
<p class="main-links">
<a href=/start> Jeszcze raz </a>
</p>
{{end}}
</div>
</div>
{{if .HasNotes}}
<div class="grid-container notes-grid">
{{range .Notes}}
<div class="grid-item">
<p class="note"> {{.}} </p>
</div>
{{end}}
</div>
{{end}}
</body>
</html>
`))

func WriteHtml( out io.Writer, state *FlaggerState ) error {
	return tpl.Execute( out, state )
}
