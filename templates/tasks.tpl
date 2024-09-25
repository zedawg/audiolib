<div id="content" class="tasks">
	<h2>Tasks</h2>
	<table id="tasks">
		<thead>
		<tr>
		<th style="width:65%">Value</th>
		<th style="width:10%">Type</th>
		<th style="width:25%">Created</th>
		</tr>
		</thead>
		<tbody>
		{{ range . }}
		<tr>
			<td>{{ .Value }}</td>
			<td>{{ .Type }}</td>
			<td>{{ .Created }}</td>
		</tr>
		{{ end }}
		</tbody>
	</table>
</div>
