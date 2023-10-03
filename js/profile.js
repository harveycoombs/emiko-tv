import { Tools } from "./tools.js";
import { Posts } from "./posts.js";

Tools.updateHeader();

const profileUserId = parseInt(document.querySelector("#user_id").value);

let profileSectionTabsArea = document.querySelector(".profile-section-tabs"),
    profileCreationDateDisplay = document.querySelector("#profile_creation_date");

profileSectionTabsArea.addEventListener("click", (e) => {
    if (!e.target.matches(".section-tab")) return;

    let currentlyOpenTab = profileSectionTabsArea.querySelector(".section-tab.current"),
        currentSection = document.querySelector(`.profile-section[data-section="${currentlyOpenTab.dataset.section}"]`);

    currentlyOpenTab.classList.remove("current");
    currentSection.classList.add("hidden");
    
    e.target.classList.add("current");

    document.querySelector(`.profile-section[data-section="${e.target.dataset.section}"]`).classList.remove("hidden");
});

let profileCommentsArea = document.querySelector("#profile_comments_area");

fetchProfilePosts();
fetchProfileComments();

let observer = new IntersectionObserver((entries) => {
    for (let entry of entries) {
        if (entry.isIntersecting) fetchProfileComments(parseInt(profileCommentsArea.dataset.offset) + 1);
    }
}, { threshold: 0.5 });

function fetchProfilePosts() {
    let profilePostsArea = document.querySelector("#profile_videos_area");

    profilePostsArea.innerHTML = "";

    let details = new URLSearchParams();

    details.set("authorid", profileUserId);
    details.set("userid", sessionStorage.getItem("userid") ?? "0");
    details.set("token", sessionStorage.getItem("token") ?? "");

    fetch("/api/profile/posts", {
        method: "POST",
        body: details
    }).then(response => response.json()).then((data) => {
        data.self ? "" : "";

        for (let post of data.posts) {
            let profilePost = Posts.display(post, true);
            profilePostsArea.append(profilePost);

            let avatar = profilePost.querySelector(".post-avatar img");

            profilePost.addEventListener("click", (e) => {
                if (e.target.matches(".post-action")) e.preventDefault();
            });
        
            avatar.addEventListener("error", () => {
                avatar.src = "/content/avatars/default.jpg";
            });
        }

        if (!data.posts.length) {
            let postsPlaceholder = document.createElement("span");
            postsPlaceholder.classList = "subtle-txt";
            postsPlaceholder.innerText = "No videos";

            profilePostsArea.append(postsPlaceholder);
        }

        document.addEventListener("click", (e) => {
            if (e.target.matches(".post-action")) e.stopPropagation();
            if (e.target.matches(".delete-post-btn")) confirmPostDeletion(e.target.dataset.post);
        });
    }).catch((ex) => {
        console.log(ex);
    });
}

function fetchProfileComments(offset=0) {
    if (offset == 0) profileCommentsArea.innerHTML = "";

    Tools.appendLoader(profileCommentsArea);

    fetch(`/api/profile/comments?authorid=${profileUserId}&offset=${offset}`).then(response => response.json()).then((data) => {
        Tools.removeLoader(profileCommentsArea);

        for (let comment of data.comments) {
            let profileComment = document.createElement("div");
            profileComment.classList = "comment";
            profileComment.dataset.id = comment.id;
            profileComment.innerHTML = `<div class="comment-detail"><div class="comment-header"><a href="/posts/${comment.post}" class="link-button comment-post"><i class="fa-solid fa-video"></i> ${comment.postTitle}</a> &middot; <span>${Tools.formatDate(comment.created)}</span></div><div class="comment-content">${comment.content}</div></div>`;

            profileCommentsArea.append(profileComment);
        }

        if (!data.comments.length && offset == 0) {
            let commentsPlaceholder = document.createElement("span");
            commentsPlaceholder.classList = "subtle-txt";
            commentsPlaceholder.innerText = "No comments";

            profileCommentsArea.append(commentsPlaceholder);
        }

        profileCommentsArea.dataset.offset = offset;

        observer.observe(profileCommentsArea.querySelector(".comment:last-child"));
    }).catch((ex) => {
        console.log(ex);
    });
}

let profileAvatarContainer = document.querySelector(".profile-avatar"),
    profileAvatar = document.querySelector(".profile-avatar img");

profileAvatar.addEventListener("error", () => {
    profileAvatar.src = "/content/avatars/default.jpg";
});

profileCreationDateDisplay.innerHTML = Tools.formatDate(profileCreationDateDisplay.innerText);

checkProfileOwnership();

function checkProfileOwnership() {
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
        if (user.id != profileUserId) {
            profileAvatarContainer.addEventListener("click", () => {
                Tools.expandMedia(profileAvatar.src);
            });

            return;
        }

        profileAvatarContainer.classList.add("editable");

        let avatarUploader = document.createElement("input");
        avatarUploader.type = "file";
        avatarUploader.id = "avatar_uploader";
        avatarUploader.classList = "hidden";

        profileAvatarContainer.append(avatarUploader);
        profileAvatarContainer.addEventListener("click", () => {
            avatarUploader.click();
        });

        avatarUploader.addEventListener("change", () => {
            Tools.handleAvatarUpload(e, profileAvatar);
        });
    }).catch((ex) => {
        console.log(ex);
    });
}

function confirmPostDeletion(id) {
    let targetPost = document.querySelector(`.individual-post[data-id="${id}"]`);
    if (!targetPost) return;

    let popup = document.createElement("div");
    popup.classList = "popup post-deletion-confirmation";
    popup.innerHTML = `<div class="panel internal-popup"><strong>Delete Post</strong><div>Are you sure you would like to delete the following post: <span class="highlighted">${targetPost.querySelector(".post-title").textContent}</span>?</div><div class="notice"><i class="fa-solid fa-triangle-exclamation"></i> <span>This action is irreversible</span></div><div class="popup-options"><button class="button" id="confirm_post_deletion_btn" data-post="${id}"><i class="fa-solid fa-trash-can"></i> Delete</button><button class="button alt" id="close_popup_btn"><i class="fa-solid fa-xmark"></i> Cancel</button></div></div>`;

    if (!document.querySelector(".popup")) document.body.append(popup);

    let confirmPostDeletionBtn = popup.querySelector("#confirm_post_deletion_btn");
    confirmPostDeletionBtn.addEventListener("click", deletePost);
}

function deletePost(e) {
    let details = new URLSearchParams();
    details.set("postid", e.target.dataset.post);
    details.set("userid", sessionStorage.getItem("userid"));
    details.set("token", sessionStorage.getItem("token"));

    fetch("/api/posts/delete", {
        method: "POST",
        body: details
    }).then(response => response.json()).then((result) => {
        if (result.success) {
            if (document.querySelector(".popup")) document.querySelector(".popup").remove();
            fetchProfilePosts();
        }
    }).catch((ex) =>{
        console.log(ex);
    });
}