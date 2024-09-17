<html>
	<head>
		<meta name="viewport" content="width=device-width, initial-scale=1.0">
		<link rel="stylesheet" href="/static/style.css">
		<script src="/static/htmx.min.js"></script>
	</head>
	<body>
		{{ template "toolbar.html.tpl" . }}
		{{ template "sidebar.html.tpl" . }}
		{{ template "content.html.tpl" . }}
	</body>
	<script src="/static/script.js"></script>
</html>