{{ define "search-control" }}
<div id="search">
	<input id="search-input" type="search" name="search" placeholder="Search..." hx-post="/search" hx-trigger="input changed delay:500ms, search" hx-target="#search-results" hx-indicator=".htmx-indicator">
	<div class="htmx-indicator">{{ template "loaders.grid.svg" }}</div>
</div>
<div id="search-results"></div>

<script>
let s = document.getElementById("search");
let r = document.getElementById("search-results");
s.addEventListener("focusin", () => r.classList.add("visible"));
s.addEventListener("focusout", () => r.classList.remove("visible"));
</script>

<div style="flex-grow:1;"></div>
{{ end }}

{{ define "search-results" }}
{{ range .SearchResults }}
<div class="search-result" onclick="alert('TODO: go to [{{.Name}}]')">
	<img src="{{ .Image }}" class="image" />
	<div class="content">
		<div class="name">{{ .Name }}</div>
		<div class="details">{{ .Details }} ({{ .Type }})</div>
	</div>
</div>
{{ else }}{{ if gt (len .Query) 2 }}<div class="search-no-result">No match found for "{{ .Query }}"</div>{{ end }}{{ end }}
{{ end }}
