<div id="toolbar">
	<a href="/" name="audiobooks" class="button">{{ template "icons.music.svg" }}</a>
	{{ template "search" . }}
	<a href="/tasks" name="tasks" class="button">{{ template "icons.history.svg" }}</a>
	<a href="/settings" name="settings" class="button">{{ template "icons.settings.svg" }}</a>
	<a href="/user" name="user" class="button circle">{{ template "icons.user.svg" }}</a>
</div>
<script>
	const page = document.querySelector("body").getAttribute("name")
	document.querySelector(`a[name=${page}]`).classList.add("active");
	document.querySelector(`a[name=${page}] svg`).setAttribute("stroke-width", 2);
</script>