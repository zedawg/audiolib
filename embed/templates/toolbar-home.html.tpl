<div id="toolbar">

	<input type="search" name="search"
		placeholder="Search..."
		hx-post="/search"
		hx-trigger="input changed delay:500ms, search"
		hx-target="#search-results"
		hx-indicator=".htmx-indicator">

</div>