function handleClickSettingsButton() {
	// document.querySelectorAll("#toolbar .fill").forEach(l => l.classList.remove("fill"));
	document.querySelector("#config-button svg").classList.toggle("fill");
	document.getElementById("config-pane").classList.toggle("active")
	document.getElementById("tasks-pane").classList.remove("active")
	document.querySelector("#tasks-button svg").classList.remove("fill");
}
function handleClickTasksButton() {
	// document.querySelectorAll("#toolbar .fill").forEach(l => l.classList.remove("fill"));
	document.querySelector("#tasks-button svg").classList.toggle("fill");
	document.getElementById("tasks-pane").classList.toggle("active")
}

function handleFocusInSearchInput() {
	document.getElementById("search-results-pane").classList.add("active");
}

function handleFocusOutSearchInput() {
	document.getElementById("search-results-pane").classList.remove("active");
}

document.addEventListener("DOMContentLoaded", () => {
	document.getElementById("search").addEventListener("focusin", handleFocusInSearchInput);
	document.getElementById("search").addEventListener("focusout", handleFocusOutSearchInput);
	document.getElementById("tasks-button").addEventListener("click", handleClickTasksButton);
	document.getElementById("config-button").addEventListener("click", handleClickSettingsButton);
});

