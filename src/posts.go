// EmikoTV ~ Written & maintained by Harvey S. Coombs
package main

import (
	"database/sql"
	"fmt"
	"log"
	"strings"

	_ "github.com/go-sql-driver/mysql"
)

type post struct {
	ID       int      `json:"id"`
	Likes    int      `json:"likes"`
	Views    int      `json:"views"`
	Title    string   `json:"title"`
	Created  string   `json:"created"`
	File     string   `json:"file"`
	Author   string   `json:"author"`
	AuthorID int      `json:"authorid"`
	Tags     []string `json:"tags"`
	Comments int 	  `json:"comments"`
}

type comment struct {
	ID        int    `json:"id"`
	Post      int    `json:"post"`
	Author    string `json:"author"`
	AuthorID  int    `json:"authorid"`
	Created   string `json:"created"`
	Content   string `json:"content"`
	PostTitle string `json:"postTitle"`
}

type tag struct {
	ID    int    `json:"id"`
	Label string `json:"label"`
}

func FetchFeaturedPosts(offset int) []post {
	records := []post{}

	db, err := sql.Open("mysql", db_str)

	if err != nil {
		log.Fatal(err)
	} else {
		detail, err := db.Query("SELECT posts.id, userbase.username AS author_name, posts.author AS author_id, posts.created, posts.title, posts.filename, (SELECT DISTINCT COUNT(*) FROM likes WHERE likes.post = posts.id) AS likes, posts.views, (SELECT COUNT(*) FROM comments WHERE comments.post=posts.id AND comments.deleted=0) AS comments FROM posts INNER join userbase ON posts.author=userbase.id WHERE posts.deleted=0 ORDER BY posts.likes DESC, posts.views DESC, posts.created DESC LIMIT 16 OFFSET " + fmt.Sprint(offset))

		var id int
		var author_name string
		var author_id int
		var created string
		var title string
		var filename string
		var likes int
		var views int
		var comments int

		if err != nil {
			log.Fatal(err)
		} else {
			for detail.Next() {
				if err := detail.Scan(&id, &author_name, &author_id, &created, &title, &filename, &likes, &views, &comments); err != nil {
					fmt.Printf("ERR! %s", err)
				} else {
					record := post{
						ID:       id,
						Author:   author_name,
						AuthorID: author_id,
						Created:  created,
						Title:    title,
						File:     filename,
						Likes:    likes,
						Views:    views,
						Comments: comments,
					}

					records = append(records, record)
				}
			}
		}

		defer detail.Close()
		defer db.Close()
	}

	return records
}

func FecthPostsByAuthor(author_id int) []post {
	records := []post{}

	db, err := sql.Open("mysql", db_str)

	if err != nil {
		log.Fatal(err)
	} else {
		detail, err := db.Query("SELECT posts.id, userbase.username AS author_name, posts.created, posts.title, posts.filename, (SELECT DISTINCT COUNT(*) FROM likes WHERE likes.post = posts.id) AS likes, (SELECT COUNT(*) FROM comments WHERE comments.post=posts.id AND comments.deleted=0) AS comments, posts.views FROM posts INNER join userbase ON posts.author=userbase.id WHERE posts.author=" + fmt.Sprint(author_id) + " AND posts.deleted=0 ORDER BY posts.likes DESC, posts.views DESC, posts.created DESC")

		var id int
		var author_name string
		var created string
		var title string
		var filename string
		var likes int
		var views int
		var comments int

		if err != nil {
			log.Fatal(err)
		} else {
			for detail.Next() {
				if err := detail.Scan(&id, &author_name, &created, &title, &filename, &likes, &comments, &views); err != nil {
					fmt.Printf("ERR! %s", err)
				} else {
					record := post{
						ID:      id,
						Author:  author_name,
						AuthorID: author_id,
						Created: created,
						Title:   title,
						File:    filename,
						Likes:   likes,
						Comments: comments,
						Views:   views,
					}

					records = append(records, record)
				}
			}
		}

		defer detail.Close()
		defer db.Close()
	}

	return records
}

func TotalPosts() int {
	result := 0
	db, err := sql.Open("mysql", db_str)

	if err != nil {
		log.Fatal(err)
	} else {
		detail, err := db.Query("SELECT COUNT(*) AS total FROM posts") // WHERE reply = 0 AND deleted = 0

		var total int

		if err != nil {
			log.Fatal(err)
		} else {
			for detail.Next() {
				if err := detail.Scan(&total); err != nil {
					fmt.Printf("ERR! %s", err)
				} else {
					result = total
				}
			}
		}

		defer detail.Close()
		defer db.Close()
	}

	return result
}

func FetchPost(post_id int) post {
	record := post{}

	db, err := sql.Open("mysql", db_str)

	if err != nil {
		log.Fatal(err)
	} else {
		detail, err := db.Query("SELECT posts.author AS author_id, userbase.username AS author_name, posts.created, posts.title, posts.filename, (SELECT COUNT(*) FROM (SELECT DISTINCT * FROM likes WHERE likes.post = " + fmt.Sprint(post_id) + ") AS l) AS likes, (SELECT COUNT(*) FROM (SELECT DISTINCT * FROM views WHERE views.post = " + fmt.Sprint(post_id) + ") AS v) AS views, (SELECT COUNT(*) FROM comments WHERE comments.post=posts.id AND comments.deleted=0) AS comments, posts.tags FROM posts INNER join userbase ON posts.author=userbase.id WHERE posts.id=" + fmt.Sprint(post_id) + " AND posts.deleted=0")

		var author_id int
		var author_name string
		var created string
		var title string
		var filename string
		var likes int
		var views int
		var comments int
		var tags string

		if err != nil {
			log.Fatal(err)
		} else {
			for detail.Next() {
				if err := detail.Scan(&author_id, &author_name, &created, &title, &filename, &likes, &views, &comments, &tags); err != nil {
					fmt.Printf("ERR! %s", err)
				} else {
					record.ID = post_id
					record.Author = author_name
					record.AuthorID = author_id
					record.Created = created
					record.Title = title
					record.File = filename
					record.Likes = likes
					record.Views = views
					record.Comments = comments
					record.Tags = ResolveTagLabels(strings.Split(tags, ","))
				}
			}
		}

		defer detail.Close()
		defer db.Close()
	}

	return record
}

func FetchComments(post_id int, offset int) []comment {
	comments := []comment{}

	db, err := sql.Open("mysql", db_str)

	if err != nil {
		log.Fatal(err)
	} else {
		detail, err := db.Query("SELECT comments.id, userbase.username AS author_name, comments.created, comments.content FROM comments INNER JOIN userbase ON userbase.id=comments.author WHERE comments.post=" + fmt.Sprint(post_id) + " AND comments.deleted=0 ORDER BY created DESC LIMIT 10 OFFSET " + fmt.Sprint(offset))

		var id int
		var author_name string
		var created string
		var content string

		if err != nil {
			log.Fatal(err)
		} else {
			for detail.Next() {
				if err := detail.Scan(&id, &author_name, &created, &content); err != nil {
					fmt.Printf("ERR! %s", err)
				} else {
					cmt := comment{
						ID:      id,
						Post:    post_id,
						Author:  author_name,
						Created: created,
						Content: content,
					}

					comments = append(comments, cmt)
				}
			}
		}

		defer detail.Close()
		defer db.Close()
	}

	return comments
}

func FetchCommentsByAuthor(author_id int, offset int) []comment {
	comments := []comment{}

	db, err := sql.Open("mysql", db_str)

	if err != nil {
		log.Fatal(err)
	} else {
		detail, err := db.Query("SELECT comments.id, posts.title AS post_title, posts.id AS post_id, comments.created, comments.content FROM comments INNER JOIN posts ON posts.id=comments.post WHERE comments.author=" + fmt.Sprint(author_id) + " AND comments.deleted=0 ORDER BY created DESC LIMIT 10 OFFSET " + fmt.Sprint(offset))

		var id int
		var post_title string
		var post_id int
		var created string
		var content string

		if err != nil {
			log.Fatal(err)
		} else {
			for detail.Next() {
				if err := detail.Scan(&id, &post_title, &post_id, &created, &content); err != nil {
					fmt.Printf("ERR! %s", err)
				} else {
					cmt := comment{
						ID:        id,
						Post:      post_id,
						PostTitle: post_title,
						Created:   created,
						Content:   content,
					}

					comments = append(comments, cmt)
				}
			}
		}

		defer detail.Close()
		defer db.Close()
	}

	return comments
}

func AddCommentOnPost(post_id int, author_id int, text string) bool {
	result := 0
	db, err := sql.Open("mysql", db_str)

	if err != nil {
		log.Fatal(err)
	} else {
		like, err := db.Exec("INSERT INTO comments (post, author, created, content) VALUES(" + fmt.Sprint(post_id) + ", " + fmt.Sprint(author_id) + ", (SELECT NOW()), \"" + EscapeBoth(text) + "\")")

		if err != nil {
			log.Fatal(err)
		} else {
			affected, err := like.RowsAffected()

			if err == nil {
				result = int(affected)
			}
		}

		defer db.Close()
	}

	return (result > 0)
}

func RemoveCommentFromPost() bool {
	return true
}

func ResolveTagLabels(tag_ids []string) []string {
	tag_labels := []string{}

	if len(tag_ids) > 0 && len(tag_ids[0]) > 0 {
		db, err := sql.Open("mysql", db_str)

		if err != nil {
			log.Fatal(err)
		} else {
			detail, err := db.Query("SELECT label AS tag_label FROM tags WHERE id IN (" + strings.Join(tag_ids, ",") + ") ORDER BY label ASC")

			var tag_label string

			if err != nil {
				log.Fatal(err)
			} else {
				for detail.Next() {
					if err := detail.Scan(&tag_label); err != nil {
						fmt.Printf("ERR! %s", err)
					} else {
						tag_labels = append(tag_labels, tag_label)
					}
				}
			}

			defer detail.Close()
			defer db.Close()
		}
	}

	return tag_labels
}

func FetchPopularTags() []tag {
	popular_tags := []tag{}

	db, err := sql.Open("mysql", db_str)

	if err != nil {
		log.Fatal(err)
	} else {
		detail, err := db.Query("SELECT tags.id, tags.label, (SELECT COUNT(*) FROM posts WHERE posts.tags LIKE ('%' + tags.id + ',%') OR posts.tags LIKE ('%,' + tags.id + '%') OR posts.tags LIKE ('%' + tags.id) OR posts.tags LIKE (tags.id + '%')) AS total FROM tags ORDER BY total DESC")

		var id int
		var label string
		var total int

		if err != nil {
			log.Fatal(err)
		} else {
			for detail.Next() {
				if err := detail.Scan(&id, &label, &total); err != nil {
					fmt.Printf("ERR! %s", err)
				} else {
					popular_tag := tag{
						ID:    id,
						Label: label,
					}

					popular_tags = append(popular_tags, popular_tag)
				}
			}
		}

		defer detail.Close()
		defer db.Close()
	}

	return popular_tags
}

func CreateNewPost(author_id int, title string, filename string, tags []int) int {
	result := 0
	db, err := sql.Open("mysql", db_str)

	if err != nil {
		log.Fatal(err)
	} else {
		tag_ids := []string{}
		for _, tag := range tags {
			tag_ids = append(tag_ids, fmt.Sprint(tag))
		}

		creation, err := db.Exec("INSERT INTO posts (author, title, created, filename, tags) VALUES(" + fmt.Sprint(author_id) + ", \"" + EscapeBoth(title) + "\", (SELECT NOW()), \"" + EscapeBoth(filename) + "\", \"" + strings.Join(tag_ids, ",") + "\")")

		if err != nil {
			log.Fatal(err)
		} else {
			lii, err := creation.LastInsertId()

			if err == nil {
				result = int(lii)
			}
		}

		defer db.Close()
	}

	return result
}

func FetchAllTags() []tag {
	tags := []tag{}

	db, err := sql.Open("mysql", db_str)

	if err != nil {
		log.Fatal(err)
	} else {
		detail, err := db.Query("SELECT id, label FROM tags ORDER BY label ASC")

		var id int
		var label string

		if err != nil {
			log.Fatal(err)
		} else {
			for detail.Next() {
				if err := detail.Scan(&id, &label); err != nil {
					fmt.Printf("ERR! %s", err)
				} else {
					individual_tag := tag{
						ID:    id,
						Label: label,
					}

					tags = append(tags, individual_tag)
				}
			}
		}

		defer detail.Close()
		defer db.Close()
	}

	return tags
}

func AddLikeToPost(user_id int, post_id int) bool {
	result := 0
	db, err := sql.Open("mysql", db_str)

	if err != nil {
		log.Fatal(err)
	} else {
		like, err := db.Exec("INSERT INTO likes (post, liker) VALUES(" + fmt.Sprint(post_id) + ", " + fmt.Sprint(user_id) + ")")

		if err != nil {
			log.Fatal(err)
		} else {
			affected, err := like.RowsAffected()

			if err == nil {
				result = int(affected)
			}
		}

		defer db.Close()
	}

	return (result > 0)
}

func RemoveLikeFromPost(user_id int, post_id int) bool {
	result := 0
	db, err := sql.Open("mysql", db_str)

	if err != nil {
		log.Fatal(err)
	} else {
		like, err := db.Exec("DELETE FROM likes WHERE post=" + fmt.Sprint(post_id) + " AND liker=" + fmt.Sprint(user_id))

		if err != nil {
			log.Fatal(err)
		} else {
			affected, err := like.RowsAffected()

			if err == nil {
				result = int(affected)
			}
		}

		defer db.Close()
	}

	return (result > 0)
}

func CheckUserHasLikedPost(user_id int, post_id int) bool {
	result := false
	db, err := sql.Open("mysql", db_str)

	if err != nil {
		log.Fatal(err)
	} else {
		check, err := db.Query("SELECT COUNT(*) AS total FROM likes WHERE post=" + fmt.Sprint(post_id) + " AND liker=" + fmt.Sprint(user_id))

		var total int

		if err != nil {
			log.Fatal(err)
		} else {
			for check.Next() {
				if err := check.Scan(&total); err != nil {
					fmt.Printf("ERR! %s", err)
				} else {
					result = (total > 0)
				}
			}
		}

		defer check.Close()
		defer db.Close()
	}

	return result
}

func ExistingPostView(post_id int, ip_address string, key string, user_id int, token string) bool {
	result := false
	db, err := sql.Open("mysql", db_str)

	if err != nil {
		log.Fatal(err)
	} else {
		check, err := db.Query("SELECT COUNT(*) AS total FROM views WHERE post=" + fmt.Sprint(post_id) + " AND (ip_address=\"" + EscapeBoth(ip_address) + "\" OR (token=\"" + EscapeBoth(token) + "\" AND user=" + fmt.Sprint(user_id) + "))")

		var total int

		if err != nil {
			log.Fatal(err)
		} else {
			for check.Next() {
				if err := check.Scan(&total); err != nil {
					fmt.Printf("ERR! %s", err)
				} else {
					result = (total > 0)
				}
			}
		}

		defer check.Close()
		defer db.Close()
	}

	return result
}

func RecordPostView(post_id int, ip_address string, key string, user_id int, token string) {
	db, err := sql.Open("mysql", db_str)

	if err != nil {
		log.Fatal(err)
	} else {
		like, err := db.Exec("INSERT INTO views (post, ip_address, token, user, viewer_key, watched) VALUES (" + fmt.Sprint(post_id) + ", \"" + EscapeBoth(ip_address) + "\", \"" + EscapeBoth(token) + "\", " + fmt.Sprint(user_id) + ", \"" + EscapeBoth(key) + "\", (SELECT NOW()))")

		if err != nil {
			log.Fatal(err)
		} else {
			_, err := like.RowsAffected()

			if err != nil {
				log.Fatal(err)
			}
		}

		defer db.Close()
	}
}

func DeletePost(post_id int, author_id int) bool {
	result := false

	db, err := sql.Open("mysql", db_str)

	if err != nil {
		log.Fatal(err)
	} else {
		deletion, err := db.Exec("UPDATE posts SET deleted=1 WHERE id=" + fmt.Sprint(post_id) + " AND author=" + fmt.Sprint(author_id))

		if err != nil {
			log.Fatal(err)
		} else {
			_, err := deletion.RowsAffected()
			result = (err == nil)
		}

		defer db.Close()
	}

	return result
}

func SearchPosts(term string) []post {
	records := []post{}

	db, err := sql.Open("mysql", db_str)

	if err != nil {
		log.Fatal(err)
	} else {
		detail, err := db.Query("SELECT posts.id, userbase.username AS author_name, posts.created, posts.title, posts.filename, (SELECT DISTINCT COUNT(*) FROM likes WHERE likes.post = posts.id) AS likes, posts.views FROM posts INNER join userbase ON posts.author=userbase.id WHERE posts.deleted=0 AND posts.title LIKE \"%" + EscapeSQL(term) + "%\" ORDER BY posts.likes DESC, posts.views DESC, posts.created DESC")

		var id int
		var author_name string
		var created string
		var title string
		var filename string
		var likes int
		var views int

		if err != nil {
			log.Fatal(err)
		} else {
			for detail.Next() {
				if err := detail.Scan(&id, &author_name, &created, &title, &filename, &likes, &views); err != nil {
					fmt.Printf("ERR! %s", err)
				} else {
					record := post{
						ID:      id,
						Author:  author_name,
						Created: created,
						Title:   title,
						File:    filename,
						Likes:   likes,
						Views:   views,
					}

					records = append(records, record)
				}
			}
		}

		defer detail.Close()
		defer db.Close()
	}

	return records
}

func SearchComments(term string) []comment {
	comments := []comment{}

	db, err := sql.Open("mysql", db_str)

	if err != nil {
		log.Fatal(err)
	} else {
		detail, err := db.Query("SELECT comments.id, comments.post AS post_id, (SELECT title FROM posts WHERE posts.id = comments.post) AS post_title, userbase.username AS author_name, comments.created, comments.content FROM comments INNER JOIN userbase ON userbase.id=comments.author WHERE comments.content LIKE \"%" + EscapeSQL(term) + "%\" AND comments.deleted=0 ORDER BY created DESC")

		var id int
		var post_id int
		var post_title string
		var author_name string
		var created string
		var content string

		if err != nil {
			log.Fatal(err)
		} else {
			for detail.Next() {
				if err := detail.Scan(&id, &post_id, &post_title, &author_name, &created, &content); err != nil {
					fmt.Printf("ERR! %s", err)
				} else {
					cmt := comment{
						ID:        id,
						Post:      post_id,
						PostTitle: post_title,
						Author:    author_name,
						Created:   created,
						Content:   content,
					}

					comments = append(comments, cmt)
				}
			}
		}

		defer detail.Close()
		defer db.Close()
	}

	return comments
}

func PostMediaPlayer(post_data post) string {
	return ""
}
