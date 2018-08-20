package layout

const outerDefault = `<html>
	<head>
		{{ template "head" . }}
	</head>
	<body>
		{{ template "content" . }}
	</body>
</html>`

const postDefault = `
{{ define "head" }}
<title>{{ .Title }}</title>
{{ end }}

{{ define "content" }}
<article>
	<header>
		{{ if .Cover }}<img src="{{ .Cover }}" /> {{ end }}
		<h1>{{ .Title }}</h1>
	</header>
	<div>{{ .Content }}</div>
	<footer>
		<ul>
			{{ range .Tags }}
				<li>
					<a href="{{ .URL }}">{{ .Title }}</a>
				</li>
			{{ end }}
		</ul>
	</footer>
</article>
{{ end }}`

const indexDefault = `
{{ define "head" }}
<title>{{ .Title }}</title>
{{ end }}

{{ define "content" }}
<main>
	<header><h1>{{ .Title }}</h1></header>
	{{ range .Posts }}
		<article>
			<header>
				<h1>
					<a href={{ .URL }}>{{ .Title }}</a>
				</h1>
			</header>
		</article>
	{{ end }}
</main>
{{ end }}`

const tagDefault = `
{{ define "head" }}
<title>{{ .Title }}</title>
{{ end }}

{{ define "content" }}
<main>
	<header><h1>{{ .Title }}</h1></header>
	{{ range .Posts }}
		<article>
			<header>
				<h1>
					<a href={{ .URL }}>{{ .Title }}</a>
				</h1>
			</header>
		</article>
	{{ end }}
</main>
{{ end }}`
