import { Tools } from "./tools.js";

let registrationForm = document.querySelector(".outer-registration-area"),
    usernameField = registrationForm.querySelector("#username_in"),
    emailField = registrationForm.querySelector("#email_in"),
    locationField = registrationForm.querySelector("#location_in"),
    genderField = registrationForm.querySelector("#gender_in"),
    passwordField = registrationForm.querySelector("#password_in"),
    submitRegistrationBtn = registrationForm.querySelector("#submit_registration_btn");
    
Tools.populateMenuWithLocations(locationField);

usernameField.focus();

usernameField.addEventListener("input", validateUniqueField);
submitRegistrationBtn.addEventListener("click", submitAccountRegistration);

let usernamePattern = new RegExp("[^A-Za-z0-9_]");

function validateUniqueField(e) {
    submitRegistrationBtn.disabled = (e.target.value.length == 0);
    if (!e.target.value.length) return;

    let endpoint = "";
    let validSyntax = true;

    switch (e.target.id) {
        case "username_in":
            endpoint = "verify/username";
            if (usernamePattern.test(e.target.value)) {
                e.target.classList.add("invalid");
                validSyntax = false;
            } else e.target.classList.remove("invalid");

            submitRegistrationBtn.disabled = usernamePattern.test(e.target.value);
            break;
        case "email_in":
            endpoint = "verify/email";
            break;
    }

    if (!endpoint.length || !validSyntax) return;

    fetch(`/api/account/details/${endpoint}?input=${e.target.value}`).then(response => response.json()).then((result) => {
        result.unique ? e.target.classList.remove("invalid") : e.target.classList.add("invalid");
        submitRegistrationBtn.disabled = !result.unique;
    }).catch((ex) => {
        console.log(ex);
    });
}

function submitAccountRegistration() {
    let details = new URLSearchParams();
    details.set("username", usernameField.value);
    details.set("email", emailField.value);
    details.set("location", locationField.childNodes[locationField.selectedIndex].value);
    details.set("password", passwordField.value);
    details.set("gender", genderField.childNodes[genderField.selectedIndex].value);

    fetch("/register", {
        method: "POST",
        body: details
    }).then(response => response.json()).then((result) => {
        if (result.success) window.location.href = "/login";
    }).catch((ex) => {
        console.log(ex);
    });
}