// EmikoTV ~ Written & maintained by Harvey Coombs
package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strings"
	"time"
)

func main() {
	fmt.Println("EmikoTV - http://emiko.tv/ - Written & maintained by Harvey Coombs, 2023")

	http.HandleFunc("/", IndexRoute)
	http.HandleFunc("/login", LoginRoute)
	http.HandleFunc("/posts/", IndividualPostRoute)
	http.HandleFunc("/videos/", IndividualVideoRoute)
	http.HandleFunc("/users/", IndividualUserRoute)
	http.HandleFunc("/settings", SettingsRoute)
	http.HandleFunc("/register", RegistrationRoute)

	http.HandleFunc("/leaderboard", LeaderboardRoute)

	http.HandleFunc("/api/", func(w http.ResponseWriter, r *http.Request) {
		endpoint := r.URL.Path[strings.Index(r.URL.Path, "/api/")+5:]

		switch endpoint {
			case "posts/featured":
				offset := (StringToInt(r.URL.Query().Get("offset"))  * 16)

				posts := FetchFeaturedPosts(offset)
				total := TotalPosts()
				content_size := DirSize(root + "/content/uploads/")

				fmt.Fprintf(w, "{ \"posts\": %s, \"total\": %d, \"content_size\": %d }", SerializeToJSON(posts), total, content_size)
				break
			case "users/featured":
				users := FetchFeaturedUsers()
				fmt.Fprint(w, SerializeToJSON(users))
				break
			case "session/verify":
				r.ParseForm()
				token := r.FormValue("token")
				userid := StringToInt(r.FormValue("userid"))

				user_details := VerifySession(token, userid)
				fmt.Fprint(w, SerializeToJSON(user_details))
				break
			case "search":
				search_term := r.URL.Query().Get("term")

				searched_posts := SerializeToJSON(SearchPosts(search_term))
				searched_users := SerializeToJSON(SearchUsers(search_term))
				searched_comments := SerializeToJSON(SearchComments(search_term))

				fmt.Fprintf(w, "{ \"posts\": %s, \"users\": %s, \"comments\": %s }", searched_posts, searched_users, searched_comments)
				break
			case "posts/comments":
				offset := (StringToInt(r.URL.Query().Get("offset"))  * 10)
				comments := FetchComments(StringToInt(r.URL.Query().Get("postid")), offset)

				fmt.Fprint(w, SerializeToJSON(comments))
				break
			case "posts/comments/add":
				r.ParseForm()

				author_id := StringToInt(r.FormValue("author"))

				if VerifyCredentials(r.FormValue("token"), author_id) {
					UpdateKarma(author_id, r.FormValue("token"), 1, false)

					successful_comment := AddCommentOnPost(StringToInt(r.FormValue("post")), author_id, r.FormValue("content"))
					fmt.Fprintf(w, "{ \"success\": %t }", successful_comment)
				}
				break
			case "posts/comments/remove":
				r.ParseForm()

				author_id := StringToInt(r.FormValue("author"))

				if VerifyCredentials(r.FormValue("token"), author_id) {
				}
				break
			case "posts/create":
				if r.Method == "POST" {
					r.ParseMultipartForm(5 * 1024 * 1024 * 1024)

					author_id := StringToInt(r.FormValue("author"))

					if VerifyCredentials(r.FormValue("token"), author_id) {
						upload, h, err := r.FormFile("file")

						uploaded_file_name := ""

						if err == nil {
							uploaded_file_name = UploadFile(upload, "/content/uploads/", fmt.Sprint(author_id)+"-"+fmt.Sprint(time.Now().UnixNano())+"-"+h.Filename)
						}

						tags := []int{}

						err2 := json.Unmarshal([]byte(r.FormValue("tags")), &tags)
						if err2 != nil {
							tags = []int{}
						}

						new_post_title := r.FormValue("title")

						if len(new_post_title) > 35 {
							new_post_title = new_post_title[:35]
						}

						new_post_id := CreateNewPost(author_id, new_post_title, uploaded_file_name, tags)

						if new_post_id > 0 {
							UpdateKarma(author_id, r.FormValue("token"), 2, false)

							successfully_trimmed := TrimVideo(root + "/content/uploads/" + uploaded_file_name, root + "/content/trimmed/" + uploaded_file_name, 10)
							
							if successfully_trimmed {
								if !CompressVideo(root + "/content/trimmed/" + uploaded_file_name, root + "/content/previews/preview-" + uploaded_file_name, "ceil(iw/2):ceil(ih/2)") {
									CreateFileFromBytes(root + "/content/previews/preview-" + uploaded_file_name, OpenFileBytes(root + "/content/trimmed/" + uploaded_file_name))
								}

								DeleteFile(root + "/content/trimmed/" + uploaded_file_name)
							}
							
							GenerateThumbnail(root + "/content/uploads/" + uploaded_file_name, root + "/content/thumbnails/thumb-" + fmt.Sprint(new_post_id) + ".jpg")
							http.Redirect(w, r, "/", 301)
						}
					}
				}
				break
			case "posts/delete":
				r.ParseForm()

				post_id := StringToInt(r.FormValue("postid"))
				user_id := StringToInt(r.FormValue("userid"))
				token := r.FormValue("token")

				successful_deletion := false

				if VerifyCredentials(token, user_id) {
					successful_deletion = DeletePost(post_id, user_id)
				}

				fmt.Fprintf(w, "{ \"success\": %t }", successful_deletion)
				break
			case "tags":
				tags := FetchAllTags()
				fmt.Fprint(w, SerializeToJSON(tags))
				break
			case "account/details":
				r.ParseForm()
				account_details := FetchAccount(StringToInt(r.FormValue("id")), r.FormValue("token"))
				fmt.Fprint(w, SerializeToJSON(account_details))
				break
			case "account/logins":
				r.ParseForm()

				account_id := StringToInt(r.FormValue("id"))

				if VerifyCredentials(r.FormValue("token"), account_id) {
					login_history := FetchLoginHistory(account_id)
					fmt.Fprint(w, SerializeToJSON(login_history))
				}
				break
			case "posts/likes/add":
				r.ParseForm()

				account_id := StringToInt(r.FormValue("userid"))
				post_id := StringToInt(r.FormValue("postid"))

				if VerifyCredentials(r.FormValue("token"), account_id) {
					successful_like := AddLikeToPost(account_id, post_id)
					fmt.Fprintf(w, "{ \"success\": %t }", successful_like)
				}
				break
			case "posts/likes/remove":
				r.ParseForm()

				account_id := StringToInt(r.FormValue("userid"))
				post_id := StringToInt(r.FormValue("postid"))

				if VerifyCredentials(r.FormValue("token"), account_id) {
					successful_removal := RemoveLikeFromPost(account_id, post_id)
					fmt.Fprintf(w, "{ \"success\": %t }", successful_removal)
				}
				break
			case "posts/likes/verify":
				r.ParseForm()

				account_id := StringToInt(r.FormValue("userid"))
				post_id := StringToInt(r.FormValue("postid"))

				if VerifyCredentials(r.FormValue("token"), account_id) {
					like_exists := CheckUserHasLikedPost(account_id, post_id)
					fmt.Fprintf(w, "{ \"liked\": %t }", like_exists)
				}
				break
			case "account/details/verify/username":
				username := r.URL.Query().Get("input")
				unique := CheckUsernameIsUnqiue(username)

				fmt.Fprintf(w, "{ \"unique\": %t }", unique)
				break
			case "account/details/verify/email":
				email := r.URL.Query().Get("input")
				unique := CheckEmailIsUnqiue(email)

				fmt.Fprintf(w, "{ \"unique\": %t }", unique)
				break
			case "account/details/update":
				r.ParseForm()

				account_id := StringToInt(r.FormValue("id"))
				token := (r.FormValue("token"))
				username := r.FormValue("username")
				email := r.FormValue("email")
				phone := r.FormValue("phone")
				location := r.FormValue("location")
				gender := StringToInt(r.FormValue("gender"))

				successful_update := UpdateAccountDetails(account_id, token, username, email, phone, location, gender)
				fmt.Fprintf(w, "{ \"success\": %t }", successful_update)
				break
			case "account/avatar/update":
				r.ParseMultipartForm(80 << 40)

				account_id := StringToInt(r.FormValue("id"))
				token := (r.FormValue("token"))

				successful_upload := false

				if VerifyCredentials(token, account_id) {
					upload, _, err := r.FormFile("file")

					uploaded_file_name := ""

					if err == nil {
						uploaded_file_name = UploadFile(upload, "/content/avatars/", fmt.Sprint(account_id) + ".png")
					} else {
						fmt.Println(err)
					}

					successful_upload = len(uploaded_file_name) > 0
				}

				fmt.Fprintf(w, "{ \"success\": %t }", successful_upload)
				break
			case "account/password/update":
				successful_update := false

				account_id := StringToInt(r.FormValue("id"))
				token := (r.FormValue("token"))

				old_password := r.FormValue("oldpassword")
				new_password := r.FormValue("newpassword")

				successful_update = ChangePassword(account_id, token, old_password, new_password)

				fmt.Fprintf(w, "{ \"success\": %t }", successful_update)
				break
			case "profile/posts":
				r.ParseForm()

				author_id := StringToInt(r.FormValue("authorid"))
				user_id := StringToInt(r.FormValue("userid"))
				token := r.FormValue("token")

				is_self := VerifyCredentials(token, user_id) && (author_id == user_id)
				posts := FecthPostsByAuthor(author_id)

				fmt.Fprintf(w, "{ \"posts\": %s, \"self\": %t }", SerializeToJSON(posts), is_self)
				break
			case "profile/comments":
				author_id := StringToInt(r.URL.Query().Get("authorid"))
				offset := (StringToInt(r.URL.Query().Get("offset"))  * 10)

				comments := FetchCommentsByAuthor(author_id, offset)

				fmt.Fprintf(w, "{ \"comments\": %s }", SerializeToJSON(comments))
				break
			case "video/views/add":
				r.ParseForm()

				post_id := StringToInt(r.FormValue("postid"))
				user_id := StringToInt(r.FormValue("userid"))
				token := r.FormValue("token")
				key := r.FormValue("key")
				ip_address := r.RemoteAddr

				if !ExistingPostView(post_id, ip_address, key, user_id, token) {
					RecordPostView(post_id, ip_address, key, user_id, token)
				}
				break
		}
	})

	ServeStatic("/css/")
	ServeStatic("/js/")
	ServeStatic("/assets/")
	ServeStatic("/assets/fonts/")
	ServeStatic("/assets/branding/")
	ServeStatic("/content/uploads/")
	ServeStatic("/content/thumbnails/")
	ServeStatic("/content/avatars/")
	ServeStatic("/content/previews/")
	ServeStatic("/content/trimmed/")

	log.Fatal(http.ListenAndServe(":6060", nil))
}

func ServeStatic(path string) {
	static := http.FileServer(http.Dir(root + path))
	http.Handle(path, http.StripPrefix(path, static))
}