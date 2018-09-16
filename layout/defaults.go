package layout

const outerDefault = `<!DOCTYPE html>
<html>
	<meta charset="utf-8"> 
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

const rssDefault = `<?xml version="1.0"?>
<rss version="2.0">
	<channel>
    	<title>{{ .Title }}</title>
    	<link>{{ .URL }}</link>
    	<description>{{ .Description }}</description>
		<language>{{ .Lang }}</language>
		{{ $firstItem := index .Posts 0 }}
    	<pubDate>{{ formatDateRFC $firstItem.Created }}</pubDate>
    	<lastBuildDate>{{ formatDateRFC $firstItem.Created }}</lastBuildDate>
    	<docs>http://blogs.law.harvard.edu/tech/rss</docs>
    	<generator>GoBlogging</generator>
    	<managingEditor>{{ .Author }}</managingEditor>
    	<webMaster>{{ .Author }}</webMaster>
		{{ range .Posts }}
		<item>
			<title>{{ .Title }}</title>
			<link>{{ .URL }}</link>
			<description>{{ .Teaser }}</description>
			<pubDate>{{ formatDateRFC .Created }}</pubDate>
      		<guid>{{ .URL }}</guid>
		</item>
		{{ end }}
	</channel>
</rss>
`
