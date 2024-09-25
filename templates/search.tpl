{{define "search"}}
	<div id="search">
		<input id="search-input" type="search" name="search" placeholder="Search..." hx-post="/search" hx-trigger="input changed delay:500ms, search" hx-target="#search-results" hx-indicator=".htmx-indicator">
		<div class="htmx-indicator">{{template "loaders.grid.svg"}}</div>
	</div>
	<div id="search-results"></div>

	<script>
	let s = document.getElementById("search");
	let r = document.getElementById("search-results");
	s.addEventListener("focusin", () => r.classList.add("visible"));
	s.addEventListener("focusout", () => r.classList.remove("visible"));
	</script>

	<div style="flex-grow:1;"></div>
{{end}}

{{define "search-results"}}
	{{range .}}
	<div class="search-result" onclick="alert('TODO: go to [{{.Title}}]')">
		<img src="/public/nocover.jpg" class="image" />
		<div class="content">
			<div class="name">{{.Title}}</div>
			<div class="details">{{.Authors}}</div>
		</div>
	</div>
	{{else}}
		<div class="search-no-result">No match found</div>
	{{end}}
{{end}}
