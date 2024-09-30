{{define "html"}}
<!DOCTYPE html>
<html>
	<head>
		<meta name="viewport" content="width=device-width, initial-scale=1.0">
		<link rel="icon" type="image/svg+xml" href="/public/favicon.svg">
		<link rel="stylesheet" href="/assets/style.css">
		<script src="/assets/htmx.min.js"></script>
		<script src="/assets/script.js"></script>
	</head>
	<body>
		<div id="toolbar">
			<div id="search">
				<input id="search-input" type="search" name="search" placeholder="Search..." hx-post="/search" hx-trigger="input changed delay:500ms, search" hx-target="#search-results-pane" hx-indicator=".htmx-indicator">
				<div class="htmx-indicator">{{template "loaders.grid"}}</div>
			</div>

			<div style="flex-grow:1;"></div>
			<button id="tasks-button" type="button" hx-get="/tasks" hx-target="#tasks-pane">{{template "icons.bell"}}</button>
			<button id="config-button" type="button" hx-get="/config" hx-target="#config-pane">{{template "icons.settings"}}</button>
		</div>
		<div id="search-results-pane"><!-- --></div>
		<div id="tasks-pane"><!-- --></div>
		<div id="config-pane"><!-- --></div>
		<main id="content"></main>
	</body>
</html>
{{end}}

{{define "search-results"}}
{{range .}}
	<div class="search-result" onclick="alert('TODO: go to [{{.Title}}]')">
		<img src="/public/nocover.jpg" class="image" />
		<div class="content">
			<div class="name">{{.Title}}</div>
			<div class="details">{{.Author}}</div>
		</div>
	</div>
{{else}}
	<div class="search-no-result">No match found</div>
{{end}}
{{end}}

{{define "tasks"}}
<h3>Tasks</h3>
{{range .}}
<div class="task">
	<div class="name">{{.Name}}</div>
	<div class="status">{{.Status}}</div>
</div>
{{end}}
{{end}}

{{define "config"}}
<h3>Users</h3>
<h3>Sources</h3>
<h3>Database</h3>
<h3>App</h3>
<p>Port: <input type="text" defaultValue="{{.Port}}" /></p>
{{end}}