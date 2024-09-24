{{ define "settings-content" }}
<div id="content" class="settings">
	<h2>Libraries</h2>
	<div class="controls">
	<button id="new-library-button" type="button">New Library</button>
	</div>
	{{ template "libraries" .Libraries }}
	{{ template "new-library-modal" }}
</div>
{{ end }}

{{ define "new-library-modal" }}
<div id="modal">
<form id="new-library">
	<h2>New Library</h2>
	<p>Enter the name of the audiobook library</p>
	<input name="name" type="text" placeholder="my audiobooks" value="my audiobooks" />
	<p>Enter the source path of your audiobook files.</p>
	<input name="import_path" type="text" placeholder="/import" value="/import" />
	<p>Enter the path to store converted and indexed audiobooks. Audiolib automatically converts audiobooks to format m4b including matching metadata.</p>
	<input name="converted_path" type="text" placeholder="/converted" value="/converted" />
	<div class="controls">
	<span class="htmx-indicator">{{ template "loaders.grid.svg" }}</span>
	<button id="cancel-new-library" type="button">Cancel</button>
	<button type="submit" hx-post="/library" hx-target="#libraries" hx-indicator=".htmx-indicator">Submit</button>
	</div>
</form>
</div>
<script>

</script>
{{ end }}

{{ define "libraries" }}
<table id="libraries">
	<thead>
	<tr>
	<th style="width:20%">Library Name</th>
	<th style="width:30%">Import Path</th>
	<th style="width:30%">Converted Path</th>
	<th style="width:20%"></th>
	</tr>
	</thead>
	<tbody>
	{{ range . }}<tr class="library"><td>{{ .Name }}</td><td>{{ .ImportPath }}</td><td>{{ .ConvertedPath }}</td></tr>
	{{ else }}<tr><td>No library</td></tr>{{ end }}
	</tbody>
</table>
{{ end }}