import { Tools } from "./tools.js";

export class Posts {
    static setupPlayer(player) {
        if (!player) return;

        let playerControls = player.querySelector(".player-controls"),
            playerInterface = player.querySelector(".post-interface"),
            video = player.querySelector("video"),
            durationDisplay = player.querySelector("#duration_display"),
            elapsedDisplay = player.querySelector("#elapsed_display"),
            playBtn = player.querySelector(".play-btn"),
            expansionBtn = player.querySelector(".expand-btn"),
            progressBar = player.querySelector(".player-progress-bar"),
            elapsedBar = progressBar.querySelector(".elapsed-bar");
    
        video.addEventListener("loadedmetadata", function() {
            let duration = Tools.secondsToHMS(Math.abs(video.duration));
            durationDisplay.innerText = Tools.formatTime(duration.hours, duration.minutes, duration.seconds);
        });
        
        let watchedHalf = false;
    
        video.addEventListener("timeupdate", () => {
            let elapsed = Tools.secondsToHMS(Math.abs(video.currentTime));
            elapsedDisplay.innerText = Tools.formatTime(elapsed.hours, elapsed.minutes, elapsed.seconds);
    
            let progressionRatio = (video.currentTime / video.duration);
            elapsedBar.setAttribute("style", `width: ${progressBar.clientWidth * progressionRatio}px;`);
    
            if (video.currentTime == video.duration) playBtn.innerHTML = '<i class="fa-solid fa-rotate-left"></i>';

            let duration = Tools.secondsToHMS(Math.abs(video.duration));
            durationDisplay.innerText = Tools.formatTime(duration.hours, duration.minutes, duration.seconds);
        });

        player.addEventListener("click", (e) => {
            switch (true) {
                case (!e.target.matches(".player-controls, .player-controls *")):
                    pauseVideo(true);
                    return;
                case (e.target.matches(".volume-btn")):
                    appendVolumeSlider(e.target);
                    break;
            }
        });
    
        playBtn.addEventListener("click", pauseVideo);
    
        expansionBtn.addEventListener("click", () => {
            video.requestFullscreen();
        });
    
        let playerControlsTimeout = setTimeout(() => {
            if (!video.paused) playerInterface.classList.add("hidden");
        }, 3000);
    
        document.addEventListener("mousemove", (e) => {
            clearTimeout(playerControlsTimeout);
            playerInterface.classList.remove("hidden");
    
            playerControlsTimeout = setTimeout(() => {
                if (!video.paused && !e.target.matches(".player-interface, .player-interface *")) playerInterface.classList.add("hidden");
            }, 3000);
        });
    
        progressBar.addEventListener("click", (e) => {
            elapsedBar.setAttribute("style", `width: ${e.offsetX}px`);
            video.currentTime = video.duration / (progressBar.clientWidth / elapsedBar.clientWidth);
        });
    
        let currentlyDraggingBar = false;
    
        progressBar.addEventListener("mousedown", () => { currentlyDraggingBar = true; });
        progressBar.addEventListener("mouseup", () => { currentlyDraggingBar = false; });
    
        progressBar.addEventListener("mousemove", (e) => {
            if (currentlyDraggingBar) {
                elapsedBar.setAttribute("style", `width: ${e.offsetX}px`);
                video.currentTime = video.duration / (progressBar.clientWidth / elapsedBar.clientWidth);
            }
        });
    
        progressBar.addEventListener("mouseleave", () => {
            currentlyDraggingBar = false;
        });

        function pauseVideo() {
            if (video.paused) {
                video.play();
                playBtn.innerHTML = '<i class="fa-solid fa-pause"></i>';
    
                if (playerControlsTimeout) clearTimeout(playerControlsTimeout);
    
                player.classList.remove("pause-anim");
                player.classList.add("play-anim");

                setTimeout(() => {
                    player.classList.remove("play-anim");
                }, 500);
                
                playerControlsTimeout = setTimeout(() => {
                    if (!video.paused) playerControls.classList.add("hidden");
                }, 3000);
            } else {
                video.pause();
                playBtn.innerHTML = '<i class="fa-solid fa-play"></i>';

                player.classList.remove("play-anim");
                player.classList.add("pause-anim");

                setTimeout(() => {
                    player.classList.remove("pause-anim");
                }, 500);
            }
        }
    }

    static display(post, displayDateAsAge=false) {
        let postCreationDate = displayDateAsAge ? `${Tools.timeDifference(new Date(post.created))} ago` : Tools.formatDate(post.created);

        let postDisplay = document.createElement("a");
        postDisplay.classList = "featured-post";
        postDisplay.href = `/posts/${post.id}`;
        postDisplay.innerHTML = `<video src="http://emiko.tv/content/previews/preview-${post.file}" type="video/mp4" muted></video><div class="progress-bar"></div>`;

        postDisplay.addEventListener("mouseenter", () => {
            postDisplay.querySelector("video").play();
        });
        
        postDisplay.addEventListener("mouseleave", () => {
            let targetVideo = postDisplay.querySelector("video");
            targetVideo.pause();
            targetVideo.currentTime = 0;
        });

        return postDisplay;
    }

    static addLike(id) {
        let likePostBtn = document.querySelector(`.post-option[data-action="like"][data-target="${id}"]`);
        if (!likePostBtn) return;

        let details = new URLSearchParams();
        details.set("userid", sessionStorage.getItem("userid"));
        details.set("token", sessionStorage.getItem("token"));
        details.set("postid", id);

        fetch("/api/posts/likes/add", {
            method: "POST",
            body: details
        }).then(response => response.json()).then((result) => {
            if (!result.success) return;
    
            likePostBtn.innerHTML = '<i class="fa-solid fa-heart"></i>';
            likePostBtn.classList.add("liked");
    
            //postLikesCounter.innerText = `${(parseInt(postLikesCounter.innerText) + 1)}`;
        }).catch((ex) => {
            console.log(ex);
        });
    }
    
    static removeLike(id) {
        let likePostBtn = document.querySelector(`.post-option[data-action="like"][data-target="${id}"]`);
        if (!likePostBtn) return;

        let details = new URLSearchParams();
        details.set("userid", sessionStorage.getItem("userid"));
        details.set("token", sessionStorage.getItem("token"));
        details.set("postid", id);

        fetch("/api/posts/likes/remove", {
            method: "POST",
            body: details
        }).then(response => response.json()).then((result) => {
            if (!result.success) return;
    
            likePostBtn.innerHTML = '<i class="fa-regular fa-heart"></i>';
            likePostBtn.classList.remove("liked");
    
            //postLikesCounter.innerText = `${(parseInt(postLikesCounter.innerText) - 1)}`;
        }).catch((ex) => {
            console.log(ex);
        });
    }
}

function appendVolumeSlider(target) {
    let slider = document.createElement("div");
    slider.classList = "volume-slider";
    slider.innerHTML = '<input type="range" min="0" max="1" step="0.01" value="1" />';

    if (target && !target.querySelector(".volume-slider")) target.append(slider);

    let volumeControl = slider.querySelector('input[type="range"]');
    let video = document.querySelector(`.media-player[data-post="${target.dataset.id}"] video`);

    volumeControl.addEventListener("input", () => {
        if (video) video.volume = volumeControl.value;
    });
}