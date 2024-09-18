
{{ template "head.html.tpl" }}
<body>
	{{ template "toolbar-home.html.tpl" . }}
	{{ $page := "home" }}
	{{ template "sidebar.html.tpl" $page }}
	<div id="content">

	</div>
</body>
