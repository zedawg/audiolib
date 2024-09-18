
{{ template "head.html.tpl" }}
<body>
	{{ template "toolbar-tasks.html.tpl" . }}
	{{ $page := "tasks" }}
	{{ template "sidebar.html.tpl" $page }}
	<div id="content">
		<table>
			<thead>
				<tr>
					<th style="width:20%">Name</th>
					<th style="width:20%">Status</th>
					<th style="width:10%">Priority</th>
					<th style="width:25%">Params</th>
					<th style="width:25%">Result</th>
				</tr>
			</thead>
			<tbody>
				{{ range .Tasks }}
					<tr>
						<td>{{ .Name }}</td>
						<td>{{ .Status }}</td>
						<td>{{ .Priority }}</td>
						<td>{{ .Params }}</td>
						<td>{{ .Result }}</td>
					</tr>
				{{ else }}
				{{ end }}
			</tbody>
		</table>

	</div>
</body>
