import { Tools } from "./tools.js";
import { Posts } from "./posts.js";

const postId = document.querySelector("#post_id").value;
const viewerKey = document.querySelector("#viewer_key").value;

let details = new URLSearchParams();
    details.set("userid", sessionStorage.getItem("userid"));
    details.set("token", sessionStorage.getItem("token"));
    details.set("postid", postId);

let postCreationDateDisplay = document.querySelector("#post_creation_date");
if (postCreationDateDisplay) postCreationDateDisplay.innerHTML = Tools.formatDate(postCreationDateDisplay.textContent);

Tools.updateHeader();

let player = document.querySelector(".media-player");
Posts.setupPlayer(player);

document.addEventListener("load", () => {
    document.querySelector(".media-player video").play();
});

/*let shareBtn = document.querySelector("#share_post_btn");
shareBtn.addEventListener("click", () => {
    Tools.copyTextToClipboard(`http://emiko.tv/videos/${postId}`, shareBtn);
});*/

/*populatePostComments();

let postCommentsArea = document.querySelector("#post_comments");

let observer = new IntersectionObserver((entries) => {
    for (let entry of entries) {
        if (entry.isIntersecting) populatePostComments(false, parseInt(postCommentsArea.dataset.offset) + 1);
    }
}, { threshold: 0.5 });

let likePostBtn = document.querySelector("#like_post_btn"),
    reportPostBtn = document.querySelector("#report_post_btn"),
    postLikesCounter = document.querySelector("#post_likes_counter");

fetch("/api/posts/likes/verify", {
    method: "POST",
    body: details
}).then(response => response.json()).then((result) => {
    if (!result.liked) return;

    likePostBtn.innerHTML = '<i class="fa-solid fa-heart"></i>';
    likePostBtn.classList.add("liked");
}).catch((ex) => {
    console.log(ex);
});

likePostBtn.addEventListener("click", () => {
    likePostBtn.classList.contains("liked") ? removeLikeFromPost() : addLikeToPost();
});

function addLikeToPost() {
    fetch("/api/posts/likes/add", {
        method: "POST",
        body: details
    }).then(response => response.json()).then((result) => {
        if (!result.success) return;

        likePostBtn.innerHTML = '<i class="fa-solid fa-heart"></i>';
        likePostBtn.classList.add("liked");

        postLikesCounter.innerText = `${(parseInt(postLikesCounter.innerText) + 1)}`;
    }).catch((ex) => {
        console.log(ex);
    });
}

function removeLikeFromPost() {
    fetch("/api/posts/likes/remove", {
        method: "POST",
        body: details
    }).then(response => response.json()).then((result) => {
        if (!result.success) return;

        likePostBtn.innerHTML = '<i class="fa-regular fa-heart"></i>';
        likePostBtn.classList.remove("liked");

        postLikesCounter.innerText = `${(parseInt(postLikesCounter.innerText) - 1)}`;
    }).catch((ex) => {
        console.log(ex);
    });
}

let addCommentBtn = document.querySelector("#add_comment_btn");
if (addCommentBtn) addCommentBtn.addEventListener("click", addCommentOnPost);

function addCommentOnPost() {
    if (addCommentBtn.disabled) return;

    let commentContentField = document.querySelector("#comment_content");
    if (!commentContentField || !commentContentField.value.length) return;

    let args = new URLSearchParams();
    args.set("token", sessionStorage.getItem("token"));
    args.set("author", sessionStorage.getItem("userid"));
    args.set("post", postId);
    args.set("content", commentContentField.value);

    fetch("/api/posts/comments/add", {
        method: "POST",
        body: args
    }).then(response => response.json()).then((result) => {
        if (!result.success) return;

        populatePostComments(true);
    }).catch((ex) => {
        console.log(ex);
    });
}

function populatePostComments(delayButtonReactivation=false, offset=0) {
    let postCommentsArea = document.querySelector("#post_comments");

    if (offset == 0) postCommentsArea.innerHTML = "";
    Tools.appendLoader(postCommentsArea);

    fetch(`/api/posts/comments?postid=${postId}&offset=${offset}`).then(response => response.json()).then((comments) => {    
        Tools.removeLoader(postCommentsArea);
        addCommentBtn.disabled = true;

        for (let comment of comments) {
            let postComment = document.createElement("div");
            postComment.classList = "comment";
            postComment.dataset.commentid = comment.id;
            postComment.innerHTML = `<div class="comment-header"><a href="/users/${comment.author}" class="link-button">@${comment.author}</a> &middot; <span class="comment-timestamp">${Tools.formatDate(comment.created)}</span></div></div></div><div class="comment-content">${comment.content}</div>`;
    
            postCommentsArea.append(postComment);
        }

        if (!comments.length && offset == 0) {
            let commentsPlaceholder = document.createElement("span");
            commentsPlaceholder.classList = "subtle-txt";
            commentsPlaceholder.innerText = "No comments";

            postCommentsArea.append(commentsPlaceholder);
        }

        postCommentsArea.dataset.offset = offset;

        if (postCommentsArea.querySelector(".comment:last-child")) observer.observe(postCommentsArea.querySelector(".comment:last-child"));

        if (delayButtonReactivation) {
            setTimeout(() => {
                addCommentBtn.disabled = false;
            }, 4000);
        } else addCommentBtn.disabled = false;
    }).catch((ex) => {
        console.log(ex);
    });
}*/