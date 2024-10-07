function handleClickTasksButton() {
	let btn = document.querySelector("#tasks-button");

	btn.classList.toggle("active");

	if (btn.classList.contains("active")) {
		document.querySelector("#tasks-pane").classList.add("active")
	} else {
		document.querySelector("#tasks-pane").classList.remove("active")
	}
}

function handleFocusInSearchInput() {
	document.querySelector("#search-results-pane").classList.add("active");
}

function handleFocusOutSearchInput() {
	document.querySelector("#search-results-pane").classList.remove("active");
}

document.addEventListener("DOMContentLoaded", () => {
	document.querySelector("#search").addEventListener("focusin", handleFocusInSearchInput);
	document.querySelector("#search").addEventListener("focusout", handleFocusOutSearchInput);
	document.querySelector("#tasks-button").addEventListener("click", handleClickTasksButton);
});

