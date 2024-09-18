
{{ template "head.html.tpl" }}
<body>
	{{ template "toolbar-logs.html.tpl" . }}
	{{ $page := "logs" }}
	{{ template "sidebar.html.tpl" $page }}
	<div id="content" class="logs">
		<div class="view-options">
			<select>
				<option>type</option>
				<option>time</option>
			</select>
			<input type="search" placeholder="filter" />
		</div>
		{{ range .Logs }}
			<p class="log">
				<span class="created">{{ .Created }}:</span>
				<span class="message">{{ .Message }}</span>
			</p>
		{{ else }}
		{{ end }}
	</div>
</body>
