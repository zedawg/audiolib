
{{ template "head.html.tpl" }}
<body>
	{{ template "toolbar-user.html.tpl" . }}
	{{ $page := "user" }}
	{{ template "sidebar.html.tpl" $page }}
	<div id="content">

	</div>
</body>
