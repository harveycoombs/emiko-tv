import { locations } from "./locations.js";

export class Tools {
    static timeDifference(date) {
        date.setHours(date.getHours() + 8);

        let now = new Date();
        let diff = (now - date);
        
        let years = Math.floor(diff / 31536000000);
        let months = Math.floor(diff / 2592000000);
        let weeks = Math.floor(diff / 604800000);
        let days = Math.floor(diff / 86400000);
        let hours = Math.floor(diff / (1000 * 60 * 60));
        let minutes = Math.floor((diff % (1000 * 60 * 60)) / (1000 * 60));
        let seconds = Math.floor((diff % (1000 * 60)) / 1000);
    
        switch (true) {
            case (years > 0):
                return `${years}y`;
            case (months > 0):
                return `${months}M`;
            case (weeks > 0):
                return `${weeks}w`;
            case (days > 0):
                return `${days}d`;
            case (hours > 0):
                return `${hours}h`;
            case (minutes > 0):
                return `${minutes}m`;
            case (seconds > 0):
                return `${seconds}s`;
            default: 
                return "";
        }
    }

    static formatBytes(bytes) {
        switch (true) {
            case bytes < 1024:
                return `${bytes} B`;
            case bytes < 1024 * 1024:
                return `${(bytes / 1024).toFixed(2)} kB`;
            case bytes < 1024 * 1024 * 1024:
                return `${(bytes / (1024 * 1024)).toFixed(2)} MB`;
            default:
                return `${(bytes / (1024 * 1024 * 1024)).toFixed(2)} GB`;
        }
    }
    
    static formatDate(raw) {
        let dt = new Date(raw);
        dt.setHours(dt.getHours() + 8);
        dt.setMinutes(dt.getMinutes() + dt.getTimezoneOffset());

        let months = ["Jan", "Feb", "Mar", "Apr", "May", "Jun", "Jul", "Aug", "Sep", "Oct", "Nov", "Dec"];

        return `${months[dt.getMonth()]} ${dt.getDate()} ${dt.getFullYear()}, ${this.toAnalogTime(dt)}`;
    }

    static toAnalogTime(time) {
        let hour = time.getHours();
        let minute = (time.getMinutes() == 0) ? "00" : (time.getMinutes() < 10) ? `0${time.getMinutes()}` : time.getMinutes().toString();

        let meridiem = ((hour) >= 12) ? "PM" : "AM";

        if (meridiem == "PM") {
            hour += 12;

            if (hour >= 24) {
                hour -= 24;
            }
        }   
        
        return `${(hour == 0) ? 12 : hour}:${minute}${meridiem}`;
    }

    static updateHeader(redirectIfInvalidSession=false) {
        let sessionDetails = new URLSearchParams();
        sessionDetails.append("userid", sessionStorage.getItem("userid"));
        sessionDetails.append("token", sessionStorage.getItem("token"));

        fetch("/api/session/verify", {
            method: "POST",
            header: {
                "Content-Type": "application/x-www-form-urlencoded"
            },
            body: sessionDetails
        }).then(response => response.json()).then((user) => {
            if (!user.id) {
                if (redirectIfInvalidSession) {
                    window.location.href = "/" 
                } else return;
            }

            let headerOptionsArea = document.querySelector(".header-options");
            if (!headerOptionsArea) return;

            let headerLoginBtn = headerOptionsArea.querySelector("#header_login_btn"),
                headerRegisterBtn = headerOptionsArea.querySelector("#header_register_btn");

            if (headerLoginBtn) headerLoginBtn.remove();
            if (headerRegisterBtn) headerRegisterBtn.remove();

            let headerUserArea = document.createElement("div");
            headerUserArea.classList = "header-user-area inline-items";

            headerUserArea.innerHTML = `<div class="header-options inline-items"><a class="header-option" id="account_settings_btn" href="/settings"><i class="fa-solid fa-gear"></i></a><button class="button" id="new_post_btn"><i class="fa-solid fa-arrow-up-from-bracket"></i> Upload</button><a href="/users/${user.username}" class="avatar header-user"><img src="/content/avatars/${user.id}.png?t=${(new Date()).getTime()}" alt="${user.username}" draggable="false" /></a><div id="logout_btn">Log out</div></div>`;

            headerOptionsArea.append(headerUserArea);

            headerUserArea.addEventListener("click", (e) => {
                switch (e.target.id) {
                    case "new_post_btn":
                        revealPostCreationForm();
                        break;
                    case "logout_btn":
                        sessionStorage.clear();
                        window.location.href = "/login";
                        break;
                }
            });

            let headerUserAvatar = headerUserArea.querySelector(".header-user img");
            headerUserAvatar.addEventListener("error", () => {
                headerUserAvatar.src = "/content/avatars/default.jpg";
            });
        }).catch((ex) => {
            console.log(ex);
        });
    }

    static secondsToHMS(s) {
        return {
            hours: Math.floor(s / 3600),
            minutes: Math.floor((s % 3600) / 60),
            seconds: Math.floor(s % 60)
        };
    }

    static formatTime(hours, minutes, seconds) {
        let hourLabel = hours > 0 ? `${hours < 10 ? "0": ""}${hours}:` : "";
        return `${hourLabel}${minutes < 10 ? "0": ""}${minutes}:${seconds < 10 ? "0": ""}${seconds}`;
    }

    static appendLoader(target=document.body) {
        let loader = document.createElement("div");
        loader.classList = "loader";
        loader.innerHTML = '<i class="fa-solid fa-circle-notch"></i>';

        target.append(loader);
    }

    static removeLoader(target=document.body) {
        if (target.querySelector(".loader")) target.querySelector(".loader").remove();
    }

    static expandMedia(url) {
        let expansionContainer = document.createElement("div");
        expansionContainer.classList = "popup";
        expansionContainer.innerHTML = `<div class="expanded-media"><video src="${url}"></video><div id="close_expanded_media_btn"><i class="fa-solid fa-circle-xmark"></i></div></div>`;
        
        document.body.append(expansionContainer);

        expansionContainer.addEventListener("click", (e) => {
            if (e.target.matches("#close_expanded_media_btn")) expansionContainer.remove();
        });
    }

    static copyTextToClipboard(text, fallbackTarget=null) {
        if (!window.isSecureContext) {
            Tools.clipboardFallbackField(text, fallbackTarget);
            return;
        }

        let item = new ClipboardItem({ "text/plain": new Blob([text], { type: "text/plain" }) });

        navigator.clipboard.write([item]).then(() => {}).catch((ex) => {
            console.log(ex);
        });
    }

    static clipboardFallbackField(text, target) {
        if (!target || !text || !text.length) return;

        let field = document.createElement("input");
        field.type = "text";
        field.classList = "field";
        field.id = "clipboard_fallback_text";
        field.value = text;
        field.readOnly = true;

        setTimeout(() => {
            target.append(field);
            field.select();
        }, 100);
    }

    static determineGender(n) {
        switch (n) {
            case 1:
                return "Male";
            case 2:
                return "Female";
            default:
                return "Unknown";
        }
    }

    static appendNotice(data) {
        let notice = document.createElement("div");
        notice.classList = `notice ${data.sentiment}`;
        notice.innerHTML = data.content;
            
        data.target.append(notice);
    }

    static createNotice(data) {
        let notice = document.createElement("div");
        notice.classList = `notice ${data.sentiment}`;
        notice.innerHTML = data.content;
            
        return notice;
    }

    static createPopup(data) {
        let popup = document.createElement("div");
        popup.classList = "popup";
        popup.innerHTML = `<div class="panel internal-popup"><div class="popup-header sb-flexbox"><strong class="popup-title">${data.title}</strong><div class="close-popup-btn"><i class="fa-solid fa-xmark"></i></div></div>${data.content}</div>`;
            
        data.target.append(popup);
    }

    static populateMenuWithLocations(target) {
        if (!target) return;

        for (let location of locations) {
            let opt = document.createElement("option");
            opt.value = location;
            opt.innerText = location;

            target.append(opt);
        }
    }

    static handleAvatarUpload(e, target) {
        let upload = e.target.files[0];
        let data = new FormData();
    
        data.append("file", upload);
        data.append("id", sessionStorage.getItem("userid"));
        data.append("token", sessionStorage.getItem("token"));
    
        fetch("/api/account/avatar/update", {
            method: "POST",
            body: data
        }).then(response => response.json()).then((result) => {
            if (result.success) {
                let url = `/content/avatars/${profileUserId}.png?t=${(new Date()).getTime()}`;
    
                target.src = url;
                
                let headerUserAvatar = document.querySelector(".header-user img");
                if (headerUserAvatar) headerUserAvatar.src = url;
            }
        }).catch((ex) => {
            console.log(ex);
        });
    }
}

function revealPostCreationForm() {
    let outerCreationArea = document.createElement("div");
    outerCreationArea.classList = "popup outer-post-creation-area";
    outerCreationArea.innerHTML = `<form class="panel internal-popup-area" id="post_creation_form" action="/api/posts/create" method="POST" enctype="multipart/form-data"><div class="post-creation-header sb-flexbox"><strong>New post</strong><div class="inline-items"><button class="button" id="submit_post_btn"><i class="fa-solid fa-feather-pointed"></i> Publish</button><button class="button alt" id="cancel_post_creation_btn"><i class="fa-solid fa-xmark"></i> Cancel</button></div></div><div class="post-creation-area even-flexbox"><div class="outer-upload-area"><div class="post-content-uploader" id="post_content_uploader"><span><i class="fa-solid fa-photo-film"></i>Drag files here or <u>browse</u></span><input type="file" id="post_file" class="hidden" name="file" accepts="video/mp4, video/avi, video/mkv, video/mov, video/wmv, video/webm, video/3gpp, video/mpeg, video/rm, video/rmvb, video/vob, video/divx, video/ogg, video/x-m4v, video/x-ms-asf, video/mpg, video/x-m2ts, video/mts, video/dv, application/x-shockwave-flash, video/x-flv" /></div></div><div class="post-detail-fields"><input type="text" class="field" placeholder="Title" id="post_title" name="title" maxlength="35" /><div id="post_tags_list" class="field">Tags: </div><select id="available_tags" class="field"><option selected disabled>Tag</option></select><input type="hidden" id="post_tags" name="tags" value="[]" /><input type="hidden" name="author" value="${sessionStorage.getItem("userid")}" /><input type="hidden" name="token" value="${sessionStorage.getItem("token")}" /></div></div></form>`;

    if (!document.querySelector(".outer-post-creation-area")) document.body.append(outerCreationArea);

    let uploader = document.querySelector("#post_file");  

    outerCreationArea.addEventListener("click", (e) => {
        switch (e.target.id) {
            case "post_content_uploader":
                if (e.target.classList.contains("with-content")) {
                    uploader.files = null;

                    e.target.classList.remove("with-content");
                    if (e.target.querySelector(".preview-video")) e.target.querySelector(".preview-video").remove();
                    e.target.removeAttribute("style");

                    if (!e.target.querySelector("span")) {
                        let placeholder = document.createElement("span");
                        placeholder.innerHTML = '<i class="fa-solid fa-photo-film"></i>Drag files here or <u>browse</u>';

                        e.target.append(placeholder);
                    }
                } else uploader.click();
                break;
            case "submit_post_btn":
                createPost(uploader);
                break;
            case "cancel_post_creation_btn":
                outerCreationArea.remove();
                break;
        }
    });

    uploader.addEventListener("change", (e) => {
        if (!e.target.files.length) return;

        let upload = e.target.files[0];
        let stream = new FileReader();

        let postTitleField = document.querySelector("#post_title");
        postTitleField.value = upload.name;

        switch (upload.name.substring(upload.name.lastIndexOf(".") + 1)) {
            case "mp4":
            case "avi":
            case "mkv":
            case "mov":
            case "wmv":
            case "webm":
            case "3gp":
            case "mpeg":
            case "rm":
            case "rmvb":
            case "vob":
            case "divx":
            case "ogv":
            case "m4v":
            case "asf":
            case "mpg":
            case "m2ts":
            case "mts":
            case "dv":
            case "swf":
            case "f4v":
            case "flv":
                let videoPlaceholder = document.createElement("video");
                videoPlaceholder.classList = "preview-video";
                
                stream.addEventListener("load", () => {
                    videoPlaceholder.src = stream.result;
                    uploader.parentNode.querySelector("span").replaceWith(videoPlaceholder);
                }, false);
                break;
            default: 
                stream.addEventListener("load", () => {
                    uploader.parentNode.setAttribute("style", `background-image: url('${stream.result}');`);
                    uploader.parentNode.querySelector("span").remove();
                }, false);
                break;
        }

        stream.readAsDataURL(e.target.files[0]);
        uploader.parentNode.classList.add("with-content");
    });

    populateAvailableTagsList();

    let postTagSelector = document.querySelector("#available_tags"),
        postTagsField = document.querySelector("#post_tags"),
        postTagsList = document.querySelector("#post_tags_list");

    postTagSelector.addEventListener("change", () => {
        let choice = postTagSelector.childNodes[postTagSelector.selectedIndex];

        let tags = JSON.parse(postTagsField.value);

        if (tags.indexOf(choice.value) != -1) return;

        tags.push(parseInt(choice.value));

        postTagsField.value = JSON.stringify(tags);

        let newTag = document.createElement("div");
        newTag.classList = "tag";
        newTag.dataset.id = choice.value;
        newTag.textContent = choice.textContent;

        postTagsList.append(newTag);

        postTagSelector.selectedIndex = 0;
        choice.classList.add("hidden");
    });

    postTagsList.addEventListener("click", (e) => {
        if (!e.target.matches(".tag")) return;
        postTagsField.value = JSON.stringify(JSON.parse(postTagsField.value).filter((id) => id != e.target.dataset.id));

        let tagInMenu = postTagSelector.querySelector(`option[value="${e.target.dataset.id}"]`);
        if (tagInMenu) tagInMenu.classList.remove("hidden");

        e.target.remove();
    });
}

function createPost(uploader) {
    let postCreationForm = document.querySelector("#post_creation_form");

    let title = postCreationForm.querySelector("#post_title").value,
        tags = postCreationForm.querySelector("#post_tags").value;

    if (!title.length || tags.length || tags != "[]" || uploader.files.length) return;
    postCreationForm.submit();
}

function populateAvailableTagsList() {
    let list = document.querySelector("#available_tags");

    fetch("/api/tags").then(response => response.json()).then((tags) => {
        for (let tag of tags) {
            let listTag = document.createElement("option");
            listTag.value = tag.id;
            listTag.textContent = tag.label;

            list.append(listTag);
        }
    }).catch((ex) => {
        console.log(ex);
    });
}