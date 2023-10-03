let header = document.querySelector("header"),
    footer = document.querySelector("footer")

if (header) document.documentElement.style.setProperty("--header-height", `${header.clientHeight}px`);
if (footer) document.documentElement.style.setProperty("--footer-height", `${footer.clientHeight}px`);

document.addEventListener("keyup", (e) => {
    if (e.key != "Escape") return;

    let existingPopup = document.querySelector(".popup");
    if (existingPopup) existingPopup.remove();
});

document.addEventListener("click", (e) => {
    if (e.target.matches(".close-popup-btn, #close_popup_btn, .popup")) {
        let existingPopup = document.querySelector(".popup");
        if (existingPopup) existingPopup.remove();
    }

    if (!e.target.matches("#clipboard_fallback_text")) {
        let clipboardFallback = document.querySelector("#clipboard_fallback_text");
        if (clipboardFallback) clipboardFallback.remove();
    }
});