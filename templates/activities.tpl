<div id="content" class="activities">
	<h2>Activity</h2>
	<table id="libraries">
		<thead>
		<tr>
		<th style="width:65%">Value</th>
		<th style="width:10%">Type</th>
		<th style="width:25%">Created</th>
		</tr>
		</thead>
		<tbody>
		{{ range .Activities }}
		<tr>
			<td>{{ .Value }}</td>
			<td>{{ .Type }}</td>
			<td>{{ .Created }}</td>
		</tr>
		{{ end }}
		</tbody>
	</table>
</div>
