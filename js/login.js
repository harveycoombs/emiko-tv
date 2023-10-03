import { Tools } from "./tools.js";

let usernameOrEmailField = document.querySelector("#username_or_email_in"),
    passwordField = document.querySelector("#password_in"),
    submitLoginAttemptBtn = document.querySelector("#submit_login_attempt_btn");

usernameOrEmailField.focus();
    
submitLoginAttemptBtn.addEventListener("click", attemptLogin);

usernameOrEmailField.addEventListener("keyup", (e) => {
    if (e.key == "Enter") attemptLogin();
});

passwordField.addEventListener("keyup", (e) => {
    if (e.key == "Enter") attemptLogin();
});

function attemptLogin() {
    let details = new URLSearchParams();
    details.append("usernameoremail", usernameOrEmailField.value);
    details.append("password", passwordField.value);
    details.append("platform", (navigator.userAgentData ?? navigator).platform);

    fetch("/login", {
        method: "POST",
        header: {
            "Content-Type": "application/x-www-form-urlencoded"
        },
        body: details
    }).then(response => response.json()).then((credentials) => {
        if (!credentials || !credentials.id || !credentials.token.length) {
            if (document.querySelector(".notice.negative")) document.querySelector(".notice.negative").remove();

            let invalidDetailsNotice = Tools.createNotice({
                content: "Invalid username/email or password",
                sentiment: "negative"
            });

            submitLoginAttemptBtn.before(invalidDetailsNotice);
            return;
        }

        sessionStorage.setItem("token", credentials.token);
        sessionStorage.setItem("userid", credentials.id);

        window.location.href = "/";
    }).catch((ex) => {
        console.log(ex);
    });
}