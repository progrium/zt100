package zt100

import (
	_ "image/png"
)

type AppLibrary struct {
}

type BlockLibrary struct {
}

type Theme struct {
	Name string
}

type MenuItem struct {
	Title string
	Page  string
}

type Section struct {
	Block     *Block
	Overrides map[string]string
}

var Template = `<html>
<head>
  <meta charset="UTF-8" />
  <meta name="viewport" content="width=device-width, initial-scale=1.0" />
  <link href="/static/tailwind-2.0.1.css" rel="stylesheet">
  <script src="/vnd/mithril-2.0.4.min.js"></script>
  <script type="module" src="/lib/util/h.js"></script>
  <script src="//cdn.jsdelivr.net/npm/medium-editor@latest/dist/js/medium-editor.min.js"></script>
  <link rel="stylesheet" href="//cdn.jsdelivr.net/npm/medium-editor@latest/dist/css/medium-editor.min.css" type="text/css" media="screen" charset="utf-8">
  <link rel="stylesheet" href="//cdn.jsdelivr.net/npm/medium-editor@latest/dist/css/themes/default.min.css" type="text/css" media="screen" charset="utf-8">
  <script>var config = %s;</script>
  <style>
	:root {
		--color-primary: %s;
	}
  </style>
</head>
<body>
	<main></main>
	<script type="module">
		(async function() {
			h.render(document.querySelector("main"), h("main", {}, [
				%s
			]))
		}());
		new EventSource(location.href).onmessage = (e) => location.reload();
	</script>
</body>
</html>`
