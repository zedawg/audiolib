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
				console.log('WebSocket connection closed:', e)
			}
		},
		open: e => {
			if (localStorage.getItem("session") == null) {

			} else {
				// window.socket.send(JSON.stringify({action: "list-dir"}))
			}
			WS.instance.send(`{"action": "list-dir"}`)
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
		    case "list-dir":
		    	try {
			        let dirs = JSON.parse(string)
			        Main.render(dirs)
			        // Modal.render(window.dirs[22])
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

const Main = {
	render: dirs => {
		dirs.map((dir, i) => {
		    dirs[i].entries = JSON.parse(dirs[i].entries)
		})
		document.querySelector("main#content").innerHTML = ''
		for (let dir of dirs) {
			Dir.render(dir)
		}
		window.dirs = dirs
	}
}

const Dir = {
	render: ({id, name, entries, props}) => {
		let img = document.createElement("img")
		img.src = entries.reduce((acc, {id, ext}) => ['jpg', 'jpeg', 'png'].includes(ext) ? `/images/${id}` : acc, null)
		img.alt = "NO IMAGE"
		img.addEventListener("click", e => {
			document.querySelector(`[data-id='${id}']`).classList.toggle("selected")
		})

		let span = document.createElement("span")
		span.innerText = name.replaceAll(" - ", " ").replaceAll("/", "\n")
		span.addEventListener("click", e => {
			Modal.render({id, name, entries, props})
		})

		let div = document.createElement("div")
		div.dataset.id = id
		div.appendChild(img)
		div.appendChild(span)

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
	render: ({id, name, entries, props}) => {
		let back = document.createElement("a")
		back.addEventListener("click", Modal.remove)
		back.innerHTML = `<svg xmlns="http://www.w3.org/2000/svg" width="24" height="24" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round" class="icon icon-tabler icons-tabler-outline icon-tabler-chevron-left"><path stroke="none" d="M0 0h24v24H0z" fill="none"/><path d="M15 6l-6 6l6 6" /></svg>`

		let img = document.createElement("img")
		img.src = entries.reduce((acc, {id, ext}) => ['jpg', 'jpeg', 'png'].includes(ext) ? `/images/${id}` : acc, null)
		img.alt = "No Image"

		let title = document.createElement("span")
		title.innerText = name.replaceAll(" - ", " ").replaceAll("/", "\n")

		let fileList = document.createElement("div")
		entries.sort((a, b) => a.name.localeCompare(b.name, undefined, { numeric: true })).map(({id, name, ext, details}, i) => {
			if (['jpeg','jpg','png'].includes(ext)) {
				return
			}
			let file = document.createElement("li")
			file.dataset.id = id
			file.dataset.name = name
			file.dataset.ext = ext
			file.innerText = name
			fileList.appendChild(file)
		})

		let modal = document.querySelector("#modal")
		modal.innerHTML = ''
		modal.appendChild(back)
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

// if (event.data instanceof Blob) {
//     const imageURL = URL.createObjectURL(event.data);
//     const img = document.createElement('img');
//     img.src = imageURL;
//     img.alt = 'Requested Image';
//     document.body.appendChild(img);

//     // Revoke the object URL after the image loads to free up memory
//     img.onload = function() {
//         URL.revokeObjectURL(imageURL);
//     };
// }

WS.instance = new WebSocket(WS.url);
WS.instance.addEventListener("error", WS.handlers.error)
WS.instance.addEventListener("close", WS.handlers.close)
WS.instance.addEventListener("open", WS.handlers.open)
WS.instance.addEventListener("message", WS.handlers.message)

document.addEventListener("DOMContentLoaded", Sidebar.handlers.DOMContentLoaded)
