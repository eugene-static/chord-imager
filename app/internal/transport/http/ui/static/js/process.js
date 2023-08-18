function collapse(el) {
    el.classList.add("hidden")
}

function expand(el) {
    el.classList.remove("hidden")
}

function hide(el) {
    el.classList.remove("visible")
}

function show(el) {
    el.classList.add("visible")
}

function construct(c) {
    let result = ""
    if (c.root !== undefined) {
        result += c.root
    }
    if (c.quality !== undefined) {
        if (c.quality === "sus2" || c.quality === "sus4") {
            if (c.extended !== undefined) {
                result += c.extended
            }
            result += c.quality
        } else {
            result += c.quality
            if (c.extended !== undefined) {
                result += c.extended
            }
        }
    } else {
        if (c.extended !== undefined) {
            result += c.extended
        }
    }
    if (c.altered !== undefined) {
        result += "(" + c.altered + ")"
    }
    if (c.omitted !== undefined) {
        result += c.omitted
    }
    return result
}

function construct_html(c) {
    let result = ""
    if (c.root !== undefined) {
        result += c.root
    }
    if (c.quality !== undefined) {
        if (c.quality === "sus2" || c.quality === "sus4") {
            if (c.extended !== undefined) {
                result += "<sup>" + c.extended + "</sup>"
            }
            result += c.quality
        } else {
            result += c.quality
            if (c.extended !== undefined) {
                result += "<sup>" + c.extended + "</sup>"
            }
        }
    } else {
        if (c.extended !== undefined) {
            result += "<sup>" + c.extended + "</sup>"
        }
    }
    if (c.altered !== undefined) {
        result += "<sup>(" + c.altered + ")</sup>"
    }
    if (c.omitted !== undefined) {
        result += "<sup>" + c.omitted + "</sup>"
    }
    return result
}

document.addEventListener("DOMContentLoaded", () => {
    window.onload = () => {
        document.querySelector('main').style.opacity = '100%'
        hide(document.querySelector('.loader'))
    }
    let strings = Array(6).fill('X')

    let info = {
        name: '',
        pattern: 'XXXXXX',
        fret: 0,
        capo: false
    }
    let requestData = {
        event: '',
        info: info,
    }

    let name = {
        root: "",
        quality: "",
        extended: "",
        altered: "",
        omitted: "",
    }

    let names = {
        base_name: name,
        variations: [name]
    }
    let tab = {
        tab: ''
    }
    let image = {
        url: ''
    }


    const refresh = document.querySelector(".refresh")
    const fretboard = document.querySelector("#fretboard_svg")
    const numbers = document.querySelectorAll(".fret_number")
    const range = document.querySelector("#meter")
    const rangeButtons = document.querySelectorAll(".meter_btn")
    const left_btn = document.querySelector("#left")
    const right_btn = document.querySelector("#right")
    const capo = document.querySelector(".capo_input")
    const notes = document.querySelectorAll('.fret_input')
    const title_inputs = document.querySelectorAll('input[name="title"]')
    const title_labels = document.querySelectorAll(".title_label")
    const vars = document.querySelectorAll(".var")
    const header = document.querySelector(".header_chord")
    const tab_button = document.querySelector("#build_tab")
    const img_button = document.querySelector("#build_png")
    const content = document.querySelector(".content")
    const content_cont = document.querySelector(".content_container")
    const copy_to_clipboard = document.querySelector(".copy_to_clipboard")
    const download = document.querySelector(".download")


    function name_request() {
        if (info.pattern !== "XXXXXX") {
            requestData.event = 'GetName'
            requestData.info = info
            let request = JSON.stringify(requestData)
            console.log(request)
            fetch('/', {
                method: 'POST',
                headers: {
                    'Content-Type': 'application/json;charset=utf-8'
                },
                body: request
            })
                .then(response => response.json())
                .then((data) => {
                    names = data
                    console.log(names)
                    if (names.base_name !== undefined) {
                        const name = construct(names.base_name)
                        const name_html = construct_html(names.base_name)
                        title_labels[0].innerHTML = name_html
                        header.innerHTML = name_html
                        info.name = name
                        title_inputs[0].value = name
                        title_inputs[0].checked = true
                        tab_button.disabled = false
                        img_button.disabled = false
                        show(title_labels[0])
                    } else {
                        info.name = ""
                        tab_button.disabled = true
                    }
                    if (names.variations !== undefined) {
                        for (let i = 0; i < names.variations.length; i++) {
                            vars[i].innerHTML = construct_html(names.variations[i])
                            title_inputs[i + 1].value = construct(names.variations[i])
                            show(vars[i])
                        }
                        for (let i = names.variations.length; i < vars.length; i++) {
                            hide(vars[i])
                        }

                    } else {
                        for (let i = 0; i < vars.length; i++) {
                            hide(vars[i])
                        }
                    }

                })
        } else {
            title_labels.forEach((lbl) => {
                hide(lbl)
            })
            img_button.disabled = true
            tab_button.disabled = true
            content.replaceChildren()
            collapse(content_cont)

        }
    }

    function tab_request() {
        if (content.childElementCount !== 0) {
            content.removeChild(content.lastElementChild)
        }
        const tab_container = document.createElement("code")
        if (info.pattern !== "XXXXXX") {
            requestData.event = 'GetTab'
            requestData.info = info
            let request = JSON.stringify(requestData)
            console.log(request)
            fetch('/', {
                method: 'POST',
                headers: {
                    'Content-Type': 'application/json;charset=utf-8'
                },
                body: request
            })
                .then(response => response.json())
                .then((data) => {
                    tab = data
                    tab_container.innerHTML = tab.tab.replaceAll("\n", "<br>")
                    content.appendChild(tab_container)
                    expand(content_cont)
                })
        } else {
            collapse(content_cont)
            tab_container.remove()
            tab_container.innerHTML = ""
        }
    }

    function png_request() {
        if (content.childElementCount !== 0) {
            content.replaceChildren()
        }
        const img_container = document.createElement("img")
        img_container.setAttribute("alt", info.name)
        img_container.classList.add("result")
        if (info.pattern !== "XXXXXX") {
            requestData.event = 'GetPNG'
            requestData.info = info
            let request = JSON.stringify(requestData)
            console.log(request)
            fetch('/', {
                method: 'POST',
                headers: {
                    'Content-Type': 'application/json;charset=utf-8'
                },
                body: request
            })
                .then(response => response.json())
                .then((data) => {
                    image = data
                    console.log("image url:", image.url)
                    img_container.setAttribute("src", image.url)
                    download.href = image.url
                    content.appendChild(img_container)
                    expand(content_cont)
                })
        } else {
            collapse(content_cont)
            img_container.remove()
        }
    }

    function change_fret(rv) {
        info.fret = rv
        if (rv === 0) {
            capo.checked = false
            info.capo = false
            hide(left_btn)
        } else {
            show(left_btn)
        }
        if (rv === 18) {
            hide(right_btn)
        } else {
            show(right_btn)
        }
        fretboard.style.left = -rv * 50 + 'px'
        for (let i = 0; i < numbers.length; i++) {
            numbers[i].innerHTML = `${rv + i + 1}`
        }
    }

    for (let i = 0; i < title_inputs.length; i++) {
        title_inputs[i].onchange = function () {
            if (title_inputs[i].checked) {
                header.innerHTML = title_labels[i].innerHTML
                info.name = title_inputs[i].value
            }
        }
    }
    capo.onclick = function () {
        if (range.value === '0') {
            capo.checked = false
        } else {
            info.capo = !!capo.checked;
            name_request()
        }
    }
    rangeButtons.forEach((btn) => {
        btn.onclick = () => {
            let rv = parseInt(range.value)
            const bv = parseInt(btn.value)
            range.value = `${rv + bv}`
            rv = parseInt(range.value)
            change_fret(rv)
            if (rv !== 0 && rv !== 18) {
                name_request()
            }
            return false
        }
    })

    notes.forEach((note) => {
        note.onclick = function () {
            if (note.checked) {
                let string_notes = document.querySelectorAll(`.fret_input[name=${note.name}]`)
                string_notes.forEach((n) => {
                    if (n.id !== note.id) {
                        n.checked = false
                        let finger = document.querySelector("." + n.id)
                        hide(finger)
                    }
                })
                let muted = document.querySelectorAll('.m_' + note.name)
                muted.forEach((m) => {
                    hide(m)
                })
                let finger = document.querySelector("." + note.id)
                show(finger)
                strings[note.name.replace("str", "") - 1] = note.value
            } else {
                strings[note.name.replace("str", "") - 1] = 'X'
                let finger = document.querySelector("." + note.id)
                let muted = document.querySelectorAll('.m_' + note.name)
                hide(finger)
                muted.forEach((m) => {
                    show(m)
                })
            }
            info.pattern = strings.join("")
            name_request()
        }
    })
    tab_button.onclick = () => {
        collapse(download)
        expand(copy_to_clipboard)
        tab_request()
    }
    img_button.onclick = () => {
        collapse(copy_to_clipboard)
        expand(download)
        png_request()
    }
    copy_to_clipboard.onclick = () => {
        const tabCont = document.querySelector("code")
        if (tabCont) {
            let data = tabCont.innerHTML
            data = data.replaceAll("<br>", "\n")
            data = data.replaceAll("&nbsp;", " ")
            navigator.clipboard.writeText(data).then(() => {
                console.log('Content copied to clipboard')
            })
        }
        return false
    }
    refresh.onclick = () => {
        info.pattern = 'XXXXXX'
        info.name = ''
        info.capo = false
        range.value = 0
        change_fret(0)
        name_request()
        strings.fill('X')
        header.innerHTML = ""
        document.querySelectorAll(".finger").forEach((el) => {
            hide(el)
        })
        document.querySelectorAll(".open").forEach((el) => {
            hide(el)
        })
        document.querySelectorAll(".muted").forEach((el) => {
            show(el)
        })
        notes.forEach((el) => {
            el.checked = false
        })
    }
})