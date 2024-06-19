
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

htmx.on("htmx:beforeSwap", function(e) {
    console.log(e)
    if (e.detail.target.id === "flash-card") {
        // we only want the width to get bigger if it's more than 6 digits
        const newWidth = sessionStorage.getItem('Digits')
        console.log(newWidth)
        switch (newWidth) {
            case "10":
                document.documentElement.style.setProperty('--varFontSize', '2.25rem')
                break;
            case "9":
                document.documentElement.style.setProperty('--varFontSize', '2.8rem')
                break;
            case "8":
            case "7":
                document.documentElement.style.setProperty('--varFontSize', '3rem')
                break;
            case "6":
                document.documentElement.style.setProperty('--varFontSize', '4rem')
                break;
            case "5":
                document.documentElement.style.setProperty('--varFontSize', '4.5rem')
                break;
            case "4":
                document.documentElement.style.setProperty('--varFontSize', '5rem')
                break;
            default:
                break;
        }
    }
})

function UpdateDigits(e) {
    sessionStorage.setItem('Digits', e.detail.parameters['digits'])
    e.detail.parameters['operator'] = sessionStorage.getItem('Operator')
}

function AddEquationParams(e) {
    e.detail.parameters['operator'] = sessionStorage.getItem('Operator')
    e.detail.parameters['digits'] = sessionStorage.getItem('Digits')
}

window.onload = function() {
    let eventSource = new EventSource("/events")
    eventSource.onopen = (e) => {
        console.log("OPENED")
    }
    eventSource.onmessage = (e) => {
        let incData = JSON.parse(e.data)
        for (const key in incData) {
            sessionStorage.setItem(key, incData[key])
        }
    }
}