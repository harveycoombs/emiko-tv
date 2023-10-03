import { Tools } from "./tools.js";

Tools.updateHeader(true);

document.addEventListener("click", (e) => {
    if (!e.target.matches(".settings-tab")) return;

    let targetTab = document.querySelector(`.settings-section[data-name="${e.target.dataset.tab}"]`);
    if (!targetTab) return;

    let existingTab = document.querySelector(".settings-section:not(.hidden)"),
        existingTabBtn = document.querySelector(".settings-tab.current");

    if (existingTab) existingTab.classList.add("hidden");
    if (existingTabBtn) existingTabBtn.classList.remove("current");

    e.target.classList.add("current");
    targetTab.classList.remove("hidden");
});

let accountSettingsArea = document.querySelector(".settings-sections");

let accountUsernameField = accountSettingsArea.querySelector("#account_username"),
    accountEmailField = accountSettingsArea.querySelector("#account_email"),
    accountPhoneField = accountSettingsArea.querySelector("#account_phone"),
    accountLocationField = accountSettingsArea.querySelector("#account_location"),
    accountGenderField = accountSettingsArea.querySelector("#account_gender"),
    accountAvatar = document.querySelector("#account_avatar"),
    accountAvatarUploader = document.querySelector("#avatar_uploader"),
    accountCreationDate = accountSettingsArea.querySelector("#account_creation_date"),
    updateDetailsBtn = accountSettingsArea.querySelector("#update_details_btn"),
    updatePasswordBtn = accountSettingsArea.querySelector("#change_password_btn");

Tools.populateMenuWithLocations(accountLocationField);

accountSettingsArea.addEventListener("input", () => {
    let existingNotice = document.querySelector(".notice");

    if (existingNotice) {
        existingNotice.remove();
        if (verifyChanges() == 0) return;
    }

    let notice = Tools.createNotice({
        content: '<i class="fa-solid fa-triangle-exclamation"></i> You have unsaved changes',
        sentiment: "neutral"
    });

    let settingsButtonsArea = document.querySelector(".account-settings-buttons");
    if (settingsButtonsArea) settingsButtonsArea.before(notice);
});

let credentials = new URLSearchParams();
credentials.set("token", sessionStorage.getItem("token"));
credentials.set("id", sessionStorage.getItem("userid"));

fetch("/api/account/details", {
    method: "POST",
    body: credentials
}).then(response => response.json()).then((details) => {
    accountUsernameField.value = details.username;
    accountEmailField.value = details.email;
    accountPhoneField.value = details.phone;
    accountLocationField.selectedIndex = accountLocationField.querySelector(`option[value="${details.location}"]`).index;
    accountGenderField.selectedIndex = accountGenderField.querySelector(`option[value="${details.gender}"]`).index;

    accountUsernameField.dataset.original = accountUsernameField.value;
    accountEmailField.dataset.original = accountEmailField.value;
    accountPhoneField.dataset.original = accountPhoneField.value;
    accountLocationField.dataset.original = accountLocationField.childNodes[accountLocationField.selectedIndex].value;
    accountGenderField.dataset.original = accountGenderField.childNodes[accountGenderField.selectedIndex].value;

    accountAvatar.src = `/content/avatars/${details.id}.png`;
    accountAvatar.alt = details.username;
    
    accountCreationDate.textContent = Tools.formatDate(details.created);

    accountAvatar.parentNode.classList.add("editable");

    updateDetailsBtn.addEventListener("click", updateAccountDetails);
    updatePasswordBtn.addEventListener("click", updateAccountPassword);

    accountAvatar.addEventListener("click", () => {
        accountAvatarUploader.click();
    });
    
    accountAvatarUploader.addEventListener("change", (e) => {
        Tools.handleAvatarUpload(e, accountAvatar);
    });
}).catch((ex) => {
    console.log(ex);
});

fetch("/api/account/logins", {
    method: "POST",
    body: credentials
}).then(response => response.json()).then((logins) => {
    let loginHistoryArea = document.querySelector("#login_history_area");

    for (let login of logins) {
        let loginHistoryRecord = document.createElement("div");
        loginHistoryRecord.classList = "individual-login-record inline-items";
    
        let platform = determinePlatform(login.platform);
        loginHistoryRecord.innerHTML = `${platform.icon}<div class="login-record-detail"><strong class="title">${platform.name}</strong><div>${Tools.timeDifference(new Date(login.datetime))} &middot; ${login.ip}</div></div>`;

        loginHistoryArea.append(loginHistoryRecord);
    }
}).catch((ex) => {
    console.log(ex);
});

accountAvatar.addEventListener("error", () => {
    accountAvatar.src = "/content/avatars/default.jpg";
});

function determinePlatform(name) {
    switch (name) {
        case "Win32":
            return { icon: '<i class="fa-brands fa-windows"></i>', name: "Windows" };
        case "macOS":
            return { icon: '<i class="fa-brands fa-apple"></i>', name: "MacOS" };
        case "Linux":
            return { icon: '<i class="fa-brands fa-linux"></i>', name: "Linux" };
        default:
            return { icon: '<i class="fa-solid fa-computer"></i>', name: "Unknown" };
    }
}

function updateAccountDetails() {
    let accountDetails = new URLSearchParams();

    accountDetails.set("token", sessionStorage.getItem("token"));
    accountDetails.set("id", sessionStorage.getItem("userid"));
    accountDetails.set("username", accountUsernameField.value);
    accountDetails.set("email", accountEmailField.value);
    accountDetails.set("phone", accountPhoneField.value);
    accountDetails.set("location", accountLocationField.value);
    accountDetails.set("gender", accountGenderField.value);

    fetch("/api/account/details/update", {
        method: "POST",
        body: accountDetails
    }).then(response => response.json()).then((result) => {
        if (result.success) window.location.reload();
    }).catch((ex) => {
        console.log(ex);
    });
}

function verifyChanges() {
    let changes = 0;
    let fields = document.querySelectorAll(".field[data-original]");

    for (let field of fields) {
        switch (field.nodeName) {
            case "INPUT":
                if (field.value != field.dataset.original) changes++;
                break;
            case "SELECT":
                if (field.childNodes[field.selectedIndex].value != field.dataset.original) changes++;
                break;
        }
    }
    
    return changes;
}

function updateAccountPassword() {
    let currentPasswordField = document.querySelector("#current_password"),
        newPasswordField = document.querySelector("#new_password"),
        newPasswordConfirmationField = document.querySelector("#confirm_new_password");

    if (newPasswordField.value != newPasswordConfirmationField.value) return;
    
    let accountPasswordDetails = new URLSearchParams();

    accountPasswordDetails.set("token", sessionStorage.getItem("token"));
    accountPasswordDetails.set("id", sessionStorage.getItem("userid"));
    accountPasswordDetails.set("oldpassword", currentPasswordField.value);
    accountPasswordDetails.set("newpassword", newPasswordField.value);

    fetch("/api/account/password/update", {
        method: "POST",
        body: accountPasswordDetails
    }).then(response => response.json()).then((result) => {
        if (result.success) window.location.reload();
    }).catch((ex) => {
        console.log(ex);
    });
}

/*function validateUniqueField(e) {
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
}*/