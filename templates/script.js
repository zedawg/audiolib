
document.addEventListener("DOMContentLoaded", () => {
	let page = '{{.Name}}';

	// toolbar active button based on page name
	document.querySelector(`a[name=${page}]`).classList.add("active");
	document.querySelector(`a[name=${page}] svg`).setAttribute("stroke-width", 2);

	// event listeners for settings page
	if (page == "settings") {
		const toggleModal = () => {
			document.getElementById('modal').classList.toggle('show');
		}
		const resetNewLibraryForm = () => {
			document.forms["new-library"].reset();
		}
		const focusNewLibraryNameInput = () => {
			document.querySelector("input[name='name'][type='text']").focus()
		}
		document.getElementById("new-library-button").addEventListener("click", e => {
			resetNewLibraryForm();
			toggleModal();
			focusNewLibraryNameInput();
		})
		document.getElementById("modal").addEventListener("click", e => {
			toggleModal();
		})
		document.getElementById("cancel-new-library").addEventListener("click", e => {
			toggleModal();
		})
		document.getElementById("new-library").addEventListener("click", e => {
			e.stopPropagation();
		})
		document.getElementById("new-library").addEventListener("closeModal", e => {
			toggleModal();
		})
	}
});
