{{ define "head" }}
<head>
	<meta name="viewport" content="width=device-width, initial-scale=1.0">
	<link rel="icon" type="image/svg+xml" href="/public/favicon.svg">
	<link rel="stylesheet" href="/assets/style.css">
	<script src="/assets/htmx.min.js"></script>
	<script src="/assets/script.js"></script>
</head>
{{ end }}

{{ define "pages.audiobooks" }}
	{{ template "head" . }}
	<body name="audiobooks">
		{{ template "toolbar.tpl" . }}
		{{ template "audiobooks.tpl" . }}
	</body>
{{ end }}

{{ define "pages.tasks"}}
	{{ template "head" . }}
	<body name="tasks">
		{{ template "toolbar.tpl" . }}
		{{ template "tasks.tpl" . }}
	</body>
{{ end }}

{{ define "pages.settings"}}
	{{ template "head" . }}
	<body name="settings">
		{{ template "toolbar.tpl" . }}
		{{ template "settings-content" . }}
	</body>
{{ end }}

{{ define "pages.user"}}
	{{ template "head" . }}
	<body name="user">
		{{ template "toolbar.tpl" . }}
		{{ template "user.tpl" . }}
	</body>
{{ end }}

{{ define "pages.new_library" }}
	{{ template "head" . }}
	<body name="new-library">
		{{ template "toolbar.tpl" . }}
		{{ template "new_library.tpl" . }}
	</body>
{{ end }}