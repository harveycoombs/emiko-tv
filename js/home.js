import { Tools } from "./tools.js";
import { Posts } from "./posts.js";

Tools.updateHeader();

let feed = document.querySelector("#feed");
let featuredUsersArea = document.querySelector("#users");

Tools.appendLoader(featuredUsersArea);

fetchFeed();

fetch("/api/users/featured").then(response => response.json()).then((data) => {
    Tools.removeLoader(featuredUsersArea);
    populateFeaturedUsersArea(data);
}).catch((ex) => {
    console.log(ex);
});

let heroSearchField = document.querySelector("#hero_search_field"),
    searchBtn = document.querySelector("#init_search_btn");

searchBtn.addEventListener("click", () => {
    if (heroSearchField.value.length) submitSearchQuery(heroSearchField.value, true);
});

heroSearchField.addEventListener("keyup", (e) => {
    if (e.key == "Enter" && heroSearchField.value.length) submitSearchQuery(heroSearchField.value, true);
});

let searchDelay;

let observer = new IntersectionObserver((entries) => {
    for (let entry of entries) {
        if (entry.isIntersecting) fetchFeed(parseInt(feed.dataset.offset) + 1);
    }
}, { threshold: 0.5 });


function fetchFeed(offset=0) {
    if (offset == 0) Tools.appendLoader(feed);

    fetch(`/api/posts/featured?offset=${offset}`).then(response => response.json()).then((data) => {
        if (offset == 0) Tools.removeLoader(feed);
        populateFeed(data, offset);
    }).catch((ex) => {
        console.log(ex);
    });
}

function populateFeed(data, offset=0) {
    if (offset == 0) feed.innerHTML = "";

    for (let post of data.posts) {
        let featuredPost = Posts.display(post, true);

        feed.append(featuredPost);
    }

    feed.dataset.offset = offset;
    if (feed.querySelector(".featured-post:last-child")) observer.observe(feed.querySelector(".featured-post:last-child"));
}

function populateFeaturedUsersArea(users) {
    featuredUsersArea.innerHTML = "";

    for (let user of users) {
        let featuredUser = document.createElement("a");
        featuredUser.classList = "featured-user";
        featuredUser.href = `/users/${user.username}`;

        featuredUser.innerHTML = `<img src="/content/avatars/${user.id}.png" alt="${user.username}" draggable="false" class="user-avatar" /><strong class="user-name">${user.username}</strong><div class="user-karma">${user.karma} Karma</div>`;

        featuredUsersArea.append(featuredUser);

        let featuredUserAvatar = featuredUser.querySelector(".user-avatar");

        featuredUserAvatar.addEventListener("error", () => {
            featuredUserAvatar.src = `/content/avatars/default.jpg`;
        });
    }
}

function submitSearchQuery(query, initial=true) {
    clearTimeout(searchDelay);

    let outerSearchArea = document.querySelector("#outer_search_area");

    if (!outerSearchArea) {
        outerSearchArea = document.createElement("div");
        outerSearchArea.classList = "popup";
        outerSearchArea.id = "outer_search_area";

        document.body.append(outerSearchArea);

        outerSearchArea.addEventListener("click", (e) => {
            if (!e.target.matches(".search-tab")) return;

            let targetTab = outerSearchArea.querySelector(`.search-section[data-name="${e.target.dataset.tab}"]`);
            if (!targetTab) return;

            let existingTab = outerSearchArea.querySelector(".search-section:not(.hidden)"),
                existingTabBtn = outerSearchArea.querySelector(".search-tab.current");

            if (existingTab) existingTab.classList.add("hidden");
            if (existingTabBtn) existingTabBtn.classList.remove("current");

            e.target.classList.add("current");
            targetTab.classList.remove("hidden");
        });
    }

    let popupSearchField = outerSearchArea.querySelector("#popup_search_field"),
        searchResultsArea = outerSearchArea.querySelector("#search_results");

    if (initial) {
        outerSearchArea.innerHTML = '<div class="panel search-results-area no-select"><strong class="panel-title sb-flexbox"><span>Search</span><div class="close-popup-btn"><i class="fa-solid fa-xmark"></i></div></strong><div class="search-field-area"><input type="text" class="field" id="popup_search_field" placeholder="Search" /></div><div class="tab-menu search-section-tabs"><div class="tab search-tab" data-tab="posts">Posts <span id="post_results_counter" class="tab-counter"></span></div><div class="tab search-tab" data-tab="users">Users <span id="user_results_counter" class="tab-counter"></span></div><div class="tab search-tab" data-tab="comments">Comments <span id="comment_results_counter" class="tab-counter"></span></div></div><div id="search_results"></div></div>';

        popupSearchField = outerSearchArea.querySelector("#popup_search_field");
        searchResultsArea = outerSearchArea.querySelector("#search_results");
    
        heroSearchField.blur();
    
        popupSearchField.focus();
        popupSearchField.value = query;
    
        popupSearchField.addEventListener("keyup", (e) => {
            if (e.key == "Enter" && popupSearchField.value.length) submitSearchQuery(popupSearchField.value, false);
        });
    }

    let postsTabBtn = outerSearchArea.querySelector('.search-tab[data-tab="posts"]'),
        usersTabBtn  = outerSearchArea.querySelector('.search-tab[data-tab="users"]'),
        commentsTabBtn = outerSearchArea.querySelector('.search-tab[data-tab="comments"]');

    Tools.appendLoader(searchResultsArea);

    fetch(`/api/search?term=${encodeURI(query)}`).then(response => response.json()).then((results) => {
        Tools.removeLoader(searchResultsArea);

        searchResultsArea.innerHTML = "";

        let postsCounter = outerSearchArea.querySelector("#post_results_counter"), 
            usersCounter = outerSearchArea.querySelector("#user_results_counter"),
            commentsCounter = outerSearchArea.querySelector("#comment_results_counter");

        postsCounter.innerText = results.posts.length;
        usersCounter.innerText = results.users.length;
        commentsCounter.innerText = results.comments.length;

        let postsSection, usersSection, commentsSection;

        postsSection = document.createElement("div");
        postsSection.id = "post_results";
        postsSection.classList = "search-section hidden";
        postsSection.dataset.name = "posts";

        usersSection = document.createElement("div");
        usersSection.id = "user_results";
        usersSection.classList = "search-section hidden";
        usersSection.dataset.name = "users";

        commentsSection = document.createElement("div");
        commentsSection.id = "comment_results";
        commentsSection.classList = "search-section hidden";
        commentsSection.dataset.name = "comments";

        searchResultsArea.append(postsSection, usersSection, commentsSection);
        
        for (let post of results.posts) {
            let result = document.createElement("div");
            result.classList = "inline-items search-result post";

            result.innerHTML = `<img class="result-thumbnail" src="/content/thumbnails/thumb-${post.id}.jpg" alt="${post.title}" draggable="false" /><div class="result-details"><strong class="result-title">${post.title}</strong><div class="result-subtitle"><a href="/users/${post.author}">${post.author}</a> &middot; ${post.views} Views</div></div>`;

            if (postsSection) postsSection.append(result);
        }

        for (let user of results.users) {
            let result = document.createElement("div");
            result.classList = "inline-items search-result user";

            result.innerHTML = `<div class="avatar"><img src="/content/avatars/${user.id}.png" alt="${user.username}" draggable="false" /></div><div class="user-details"><a href="/users/${user.username}" class="link-button">@${user.username}</a><div>${Tools.determineGender(user.gender)} &middot; ${user.location}</div></div>`;

            if (usersSection) usersSection.append(result);

        }

        for (let comment of results.comments) {
            let result = document.createElement("div");
            result.classList = "inline-items search-result comment"; 

            result.innerHTML = `<div class="comment-detail"><div class="comment-header"><a href="/users/${comment.author}" class="link-button comment-author">@${comment.author}</a> &middot; <a href="/posts/${comment.post}" class="link-button comment-post"><i class="fa-solid fa-video"></i> ${comment.postTitle}</a> &middot; <span>${Tools.formatDate(comment.created)}</span></div><div class="comment-content">${comment.content}</div></div>`;

            if (commentsSection) commentsSection.append(result);
        }
        
        if (!results.posts.length) postsSection.innerHTML = '<span style="color: var(--alt-text-color); display: block; padding-top: 8px;">No results</span>';
        if (!results.users.length) usersSection.innerHTML = '<span style="color: var(--alt-text-color); display: block; padding-top: 8px;">No results</span>';
        if (!results.comments.length) commentsSection.innerHTML = '<span style="color: var(--alt-text-color); display: block; padding-top: 8px;">No results</span>';

        if (usersSection) {
            let avatars = usersSection.querySelectorAll(".search-result.user .avatar img");

            for (let avatar of avatars) {
                avatar.addEventListener("error", () => {
                    avatar.src = "/content/avatars/default.jpg";
                });
            }
        }

        let currentTabBtn = outerSearchArea.querySelector(".search-tab.current");
        if (currentTabBtn) currentTabBtn.classList.remove("current");

        if (results.posts.length) {
            postsTabBtn.classList.add("current");
            postsSection.classList.remove("hidden");
        } else if (results.users.length) {
            usersTabBtn.classList.add("current");
            usersSection.classList.remove("hidden");
        } else if (results.comments.length) {
            commentsTabBtn.classList.add("current");
            commentsSection.classList.remove("hidden");
        }
    }).catch((ex) => {
        console.log(ex);
    });
}