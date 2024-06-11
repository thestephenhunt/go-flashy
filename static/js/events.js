
function openOptMenu() {
    document.getElementById("options-menu").classList.toggle('opened')
}

function openUserMenu() {
    const container = document.getElementById("login-container")
    container.classList.toggle("opened")
}

function openNavMenu() {
    const menu = document.getElementById("nav-menu")
    menu.classList.toggle("opened")
}

htmx.on("htmx:load", function(e) {
    if (!!document.getElementById("login-form")) {
        document.getElementsByName("username")[0].focus()
    }
})

window.onload = function() {
    let eventSource = new EventSource("/events")
    eventSource.onopen = (e) => {
        console.log("OPENED")
    }
    eventSource.onmessage = (e) => {
        console.log(e.data)
        let incData = JSON.parse(e.data)
        for (const key in incData) {
            sessionStorage.setItem(key, incData[key])
        }
    }
}