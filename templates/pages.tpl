{{ define "head" }}
<head>
<meta name="viewport" content="width=device-width, initial-scale=1.0">
<link rel="icon" type="image/svg+xml" href="/static/favicon.svg">
<!-- <link rel="icon" type="image/x-icon" href="/static/favicon.ico"> -->
<link rel="stylesheet" href="/static/style.css">
<script src="/static/htmx.min.js"></script>
<script>{{ template "script.js" . }}</script>
</head>
{{ end }}

{{ define "pages.home" }}
{{ template "head" . }}
<body>
{{ template "toolbar.tpl" . }}
{{ template "home.tpl" . }}
</body>
{{ end }}

{{ define "pages.activities"}}
{{ template "head" . }}
<body>
{{ template "toolbar.tpl" . }}
{{ template "activities.tpl" . }}
</body>
{{ end }}

{{ define "pages.settings"}}
{{ template "head" . }}
<body>
{{ template "toolbar.tpl" . }}
{{ template "settings-content" . }}
</body>
{{ end }}

{{ define "pages.user"}}
{{ template "head" . }}
<body>
{{ template "toolbar.tpl" . }}
{{ template "user.tpl" . }}
</body>
{{ end }}

{{ define "pages.new_library" }}
{{ template "head" . }}
<body>
{{ template "toolbar.tpl" . }}
{{ template "new_library.tpl" . }}
</body>
{{ end }}