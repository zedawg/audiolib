// websocket config
const WS = {
	instance: null,
	url: 'ws://localhost:8000/socket',
	handlers: {
		error: err => {
			console.error('WebSocket error:', err)
		},
		close: ({code}) => {
			switch (code) {
			case 1001: // closed by client
				break
			default:
				console.log('WebSocket connection closed:', code)
			}
		},
		open: e => {
			if (localStorage.getItem("session") == null) {

			} else {
				// window.socket.send(JSON.stringify({action: "list-entries"}))
			}
			WS.instance.send(`{"action": "list-audiobooks"}`)
		},
		message: e => {
		    const {action, state, data, string} = JSON.parse(e.data)
		    // check session
		    const {session} = state
		    if (session == null) {
		    	console.log('no session')
		    }
		    // route action
		    switch (action) {
		    case "list-audiobooks":
		    	try {
			        let audiobooks = JSON.parse(string)
			        MainContent.render(audiobooks)
		    	} catch (err) {
					console.error(`socket.onmessage error: ${err}`)
		    	}
		        break
		    default:
		    	console.log(`unhandled socket action ${action}`)
		        break
		    }
		}
	}
}

const MainContent = {
	render: audiobooks => {
		audiobooks.map(({files}, i) => {
		    audiobooks[i].files = JSON.parse(files)
		})
		document.querySelector("main#content").innerHTML = ''
		for (let b of audiobooks) {
			Audiobook.render(b)
		}
	}
}

const Audiobook = {
	render: ({id, name, files, props, image_id}) => {
		let img = document.createElement("img")
		img.src = image_id != null ? `/images/${image_id}` : ""
		img.alt = name.split(" ").map(word => word.charAt(0).toUpperCase() + word.substring(1)).join(" ")
		img.addEventListener("click", e => {
			Modal.render({id, name, files, props, image_id})
		})
		let div = document.createElement("div")
		div.dataset.id = id
		div.appendChild(img)
		document.querySelector("main#content").appendChild(div)
	}
}

const Modal = {
	handlers: {
		keydown: e => {
			if (e.key == "Escape") {
				Modal.remove()
			}
		},
		click: e => {
			const r = document.getElementById("modal").getBoundingClientRect()
			if (!(e.clientX >= r.left && e.clientX <= r.right && e.clientY >= r.top && e.clientY <= r.bottom)) {
			    Modal.remove()
			}
		}
	},
	render: ({id, name, files, props, image_id}) => {
		let img = document.createElement("img")
		img.src = image_id != null ? `/images/${image_id}` : ""
		img.alt = "No Image"

		let title = document.createElement("span")
		title.innerText = name.replaceAll(" - ", " ").replaceAll("/", "\n")

		let fileList = document.createElement("div")
		files.sort((a, b) => a.name.localeCompare(b.name, undefined, { numeric: true })).map(({id, name, ext, details}, i) => {
			if (['jpeg','jpg','png'].includes(ext)) {
				return
			}
			let fileEl = document.createElement("li")
			fileEl.dataset.id = id
			fileEl.dataset.name = name
			fileEl.dataset.ext = ext
			fileEl.innerText = name
			fileList.appendChild(fileEl)
		})

		let modal = document.querySelector("#modal")
		modal.innerHTML = ''
		// modal.appendChild(back)
		modal.appendChild(img)
		modal.appendChild(title)
		modal.appendChild(fileList)
		modal.classList.add("show")

		document.body.style.overflow = "hidden"
		document.body.addEventListener("click", Modal.handlers.click)
		document.body.addEventListener("keydown", Modal.handlers.keydown)
	},
	remove: () => {
		document.querySelector("#modal").classList.remove("show")
		document.body.style.overflow = "scroll"
		document.body.removeEventListener("click", Modal.handlers.click)
		document.body.removeEventListener("keydown", Modal.handlers.keydown)
	}
}

const Sidebar = {
	setSelected: name => {
		document.querySelectorAll("#sidebar svg").forEach(n => {
			n.getAttribute("name") == name ? n.classList.add("selected") : n.classList.remove("selected")
		})
	},
	handlers: {
		DOMContentLoaded: () => {
			document.querySelectorAll("#sidebar svg").forEach(svg => {
				let name = svg.getAttribute("name")
				svg.addEventListener("click", () => {
					Sidebar.setSelected(name)
					Modal.remove()
				})
			})
		}
	}
}

WS.instance = new WebSocket(WS.url);
WS.instance.addEventListener("error", WS.handlers.error)
WS.instance.addEventListener("close", WS.handlers.close)
WS.instance.addEventListener("open", WS.handlers.open)
WS.instance.addEventListener("message", WS.handlers.message)

document.addEventListener("DOMContentLoaded", Sidebar.handlers.DOMContentLoaded)
