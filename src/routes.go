package main

import (
	"fmt"
	"net/http"
	"os"
	"strings"
	"time"
	"io"
)

func IndexRoute(w http.ResponseWriter, r *http.Request) {
	resp := OpenFile(root + "/views/home.html")
	resp = strings.Replace(resp, "{{header}}", OpenFile(root+"/views/shared/header.html"), 1)
	resp = strings.Replace(resp, "{{footer}}", OpenFile(root+"/views/shared/footer.html"), 1)

	popular_tags_list := ""
	popular_tags := FetchPopularTags()

	for _, popular_tag := range popular_tags {
		popular_tags_list += "<div class=\"tag\">" + popular_tag.Label + "</div>"
	}

	resp = strings.Replace(resp, "{{popular_tags}}", popular_tags_list, 1)

	w.Header().Set("Access-Control-Allow-Origin", "*")
	fmt.Fprint(w, resp)
}

func LoginRoute(w http.ResponseWriter, r *http.Request) {
	resp := ""

	if r.Method == "POST" {
		r.ParseForm()
		username_or_email := r.FormValue("usernameoremail")
		password_hash := r.FormValue("password")

		creds := AttemptLogin(username_or_email, password_hash)
		resp = SerializeToJSON(creds)

		if creds.ID > 0 {
			UpdateLoginHistory(creds.ID, r.RemoteAddr, r.FormValue("platform"))
		}

		w.Header().Set("Content-Type", "application/json")
	} else {
		resp = OpenFile(root + "/views/login.html")
		resp = strings.Replace(resp, "{{footer}}", OpenFile(root+"/views/shared/footer.html"), 1)

		w.Header().Set("Content-Type", "text/html")
		w.Header().Set("Access-Control-Allow-Origin", "*")
	}

	fmt.Fprint(w, resp)
}

func IndividualPostRoute(w http.ResponseWriter, r *http.Request) {
	post_id := StringToInt(r.URL.Path[strings.Index(r.URL.Path, "/posts/")+7:])
	individual_post := FetchPost(post_id)

	resp := OpenFile(root + "/views/post.html")
	resp = strings.Replace(resp, "{{header}}", OpenFile(root+"/views/shared/header.html"), 1)
	resp = strings.Replace(resp, "{{footer}}", OpenFile(root+"/views/shared/footer.html"), 1)

	resp = strings.ReplaceAll(resp, "{{title}}", individual_post.Title)
	resp = strings.ReplaceAll(resp, "{{author}}", individual_post.Author)
	resp = strings.ReplaceAll(resp, "{{author_id}}", fmt.Sprint(individual_post.AuthorID))
	resp = strings.ReplaceAll(resp, "{{created}}", individual_post.Created)
	resp = strings.ReplaceAll(resp, "{{likes}}", fmt.Sprint(individual_post.Likes))
	resp = strings.ReplaceAll(resp, "{{post_id}}", fmt.Sprint(individual_post.ID))
	resp = strings.Replace(resp, "{{views}}", fmt.Sprint(individual_post.Views), 1)
	resp = strings.Replace(resp, "{{file}}", individual_post.File, 1)
	resp = strings.Replace(resp, "{{total_comments}}", fmt.Sprint(individual_post.Comments), 1)
	
	viewing_key := GenerateSHA256(r.RemoteAddr + r.UserAgent() + fmt.Sprint(post_id) + fmt.Sprint(time.Now().Unix()))
	resp = strings.Replace(resp, "{{viewer_key}}", viewing_key, 1)

	tags_list := ""

	for _, post_tag := range individual_post.Tags {
		tags_list += "<div class=\"tag\">" + post_tag + "</div>"
	}

	if len(individual_post.Tags) == 0 {
		tags_list = "<span class=\"subtle-txt\">No tags</div>"
	}

	resp = strings.Replace(resp, "{{tags}}", tags_list, 1)

	fmt.Fprint(w, resp)
}

func IndividualVideoRoute(w http.ResponseWriter, r *http.Request) {
	post_id := StringToInt(r.URL.Path[strings.Index(r.URL.Path, "/videos/")+8:])
	target_post := FetchPost(post_id)

	content_type := DetermineContentType(root + "/content/uploads/" + target_post.File)
	w.Header().Add("Content-Type", content_type)
	
	video, err := os.Open(root + "/content/uploads/" + target_post.File)
	if err != nil {
		w.WriteHeader(404)
	}

	_, err = io.Copy(w, video)
	if err != nil {
		w.WriteHeader(404)
	}
}

func IndividualUserRoute(w http.ResponseWriter, r *http.Request) {
	username := r.URL.Path[strings.Index(r.URL.Path, "/users/")+7:]
	individual_user := FetchUserFromUsername(username)

	resp := OpenFile(root + "/views/profile.html")
	resp = strings.Replace(resp, "{{header}}", OpenFile(root+"/views/shared/header.html"), 1)
	resp = strings.Replace(resp, "{{footer}}", OpenFile(root+"/views/shared/footer.html"), 1)

	resp = strings.ReplaceAll(resp, "{{username}}", individual_user.UserName)
	resp = strings.Replace(resp, "{{user_id}}", fmt.Sprint(individual_user.ID), 1)
	resp = strings.Replace(resp, "{{karma}}", fmt.Sprint(individual_user.Karma), 1)
	resp = strings.Replace(resp, "{{created}}", individual_user.Created, 1)
	resp = strings.Replace(resp, "{{location}}", individual_user.Location, 1)

	user_flare := ""

	if individual_user.Flare == 3 && individual_user.ID == 1 {
		user_flare = "<span class=\"user-flare no-select\" title=\"Developer of Emiko.TV\"><i class=\"fa-solid fa-screwdriver-wrench\"></i></span>"
	}

	resp = strings.Replace(resp, "{{flare}}", user_flare, 1)

	avatar_url := "/content/avatars/" + fmt.Sprint(individual_user.ID) + ".png"

	if !FileExists(root + avatar_url) {
		avatar_url = "/content/avatars/default.jpg"
	} else {
		avatar_url += "?t=" + fmt.Sprint(time.Now().Unix())
	}

	resp = strings.ReplaceAll(resp, "{{avatar_url}}", avatar_url)

	gender_field := ""

	switch individual_user.Gender {
		case 1:
			gender_field = "<i class=\"fa-solid fa-mars\"></i> <span>Male</span>"
			break
		case 2:
			gender_field = "<i class=\"fa-solid fa-venus\"></i> <span>Female</span>"
			break
		default:
			gender_field = "<i class=\"fa-solid fa-transgender\"></i> <span>Unknown</span>"
			break
	}

	resp = strings.Replace(resp, "{{gender}}", gender_field, 1)

	fmt.Fprint(w, resp)
}

func SettingsRoute(w http.ResponseWriter, r *http.Request) {
	resp := OpenFile(root + "/views/settings.html")
	resp = strings.Replace(resp, "{{header}}", OpenFile(root+"/views/shared/header.html"), 1)
	resp = strings.Replace(resp, "{{footer}}", OpenFile(root+"/views/shared/footer.html"), 1)

	fmt.Fprint(w, resp)
}

func RegistrationRoute(w http.ResponseWriter, r *http.Request) {
	if r.Method == "POST" {
		r.ParseForm()

		username := r.FormValue("username")
		email := r.FormValue("email")
		location := r.FormValue("location")
		gender := StringToInt(r.FormValue("gender"))
		password := r.FormValue("password")

		if !CheckEmailIsUnqiue(email) || !CheckUsernameIsUnqiue(username) {
			fmt.Fprint(w, "{ \"success\": false }")
		} else {
			creation_result := (CreateAccount(username, email, password, location, gender) > 0)
			fmt.Fprintf(w, "{ \"success\": %t }", creation_result)
		}
	} else {
		resp := OpenFile(root + "/views/register.html")
		resp = strings.Replace(resp, "{{footer}}", OpenFile(root+"/views/shared/footer.html"), 1)

		fmt.Fprint(w, resp)
	}
}

func LeaderboardRoute(w http.ResponseWriter, r *http.Request) {
	resp := OpenFile(root + "/views/leaderboard.html")
	resp = strings.Replace(resp, "{{header}}", OpenFile(root+"/views/shared/header.html"), 1)
	resp = strings.Replace(resp, "{{footer}}", OpenFile(root+"/views/shared/footer.html"), 1)

	rows := ""
	karma_leaderboard := FetchKarmaLeaderboard()

	for x, individual_user := range karma_leaderboard {
		pos_label := "#" + fmt.Sprint(x + 1)
		classes := "sb-flexbox leaderboard-user no-select"

		if x == 0 {
			pos_label = "<i class=\"fa-solid fa-crown\"></i>"
			classes += " first"
		}

		gender := ""

		switch individual_user.Gender {
			case 1:
				gender = "<i class=\"fa-solid fa-mars\"></i> Male"
				break
			case 2:
				gender = "<i class=\"fa-solid fa-venus\"></i> Female"
				break
			default:
				gender = "<i class=\"fa-solid fa-transgender\"></i> Unknown"
				break
		}

		rows += "<div class=\"" + classes + "\"><div class=\"inline-items\"><div class=\"avatar leaderboard-avatar\"><img src=\"/content/avatars/" + fmt.Sprint(individual_user.ID) + ".png\" alt=\"" +  individual_user.UserName + "\" /></div><div><a href=\"\" class=\"link-button\">@" + individual_user.UserName + "</a><div class=\"leaderboard-user-details\">" + gender + " &middot; " + individual_user.Location + "</div></div></div> <div><div class=\"position-label\">" + pos_label + "</div><div class=\"karma\">" + fmt.Sprint(individual_user.Karma) + "</div></div></div>"
	}

	resp = strings.Replace(resp, "{{leaderboard}}", rows, 1)
	
	fmt.Fprint(w, resp)
}