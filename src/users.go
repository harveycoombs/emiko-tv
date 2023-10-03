// EmikoTV ~ Written & maintained by Harvey S. Coombs
package main

import (
	"fmt"
	"log"
	"time"

	//"time"
	"database/sql"

	_ "github.com/go-sql-driver/mysql"
)

type user struct {
	ID        int    `json:"id"`
	UserName  string `json:"username"`
	Created   string `json:"created"`
	Location  string `json:"location"`
	Flare     int    `json:"flare"`
	Karma     int    `json:"karma"`
	IPAddress string `json:"ip_address"`
	Gender    int    `json:"gender"`
}

type account struct {
	ID        int    `json:"id"`
	UserName  string `json:"username"`
	Email     string `json:"email"`
	Phone     string `json:"phone"`
	Created   string `json:"created"`
	Location  string `json:"location"`
	Flare     int    `json:"flare"`
	IPAddress string `json:"ip_address"`
	Gender    int    `json:"gender"`
}

type credentials struct {
	ID    int    `json:"id"`
	Token string `json:"token"`
}

type loginhistory struct {
	AccountID   int    `json:"accountid"`
	IPAddress   string `json:"ip"`
	Platform    string `json:"platform"`
	DateAndTime string `json:"datetime"`
}

func AttemptLogin(username string, password string) credentials {
	record := credentials{}
	record.ID = 0

	db, err := sql.Open("mysql", db_str)

	if err != nil {
		log.Fatal(err)
	} else {
		creds, err := db.Query("SELECT token, id FROM userbase WHERE password=\"" + GenerateSHA256(password) + "\" AND (username=\"" + EscapeSQL(username) + "\" OR email=\"" + EscapeSQL(username) + "\") AND deleted=0")

		var token string
		var userid int

		if err != nil {
			log.Fatal(err)
		} else {
			for creds.Next() {
				if err := creds.Scan(&token, &userid); err != nil {
					fmt.Printf("ERR! %s", err)
				} else {
					record.ID = userid
					record.Token = token
				}
			}
		}

		defer creds.Close()
		defer db.Close()
	}

	return record
}

func VerifySession(token string, id int) user {
	verified_user := user{}

	if len(token) > 0 && id > 0 {
		db, err := sql.Open("mysql", db_str)

		if err != nil {
			log.Fatal(err)
		} else {
			session, err := db.Query("SELECT username, flare, last_ip FROM userbase WHERE token=\"" + EscapeSQL(token) + "\" AND id=" + fmt.Sprint(id) + " AND deleted=0")

			var username string
			var flare int
			var last_ip string

			if err != nil {
				log.Fatal(err)
			} else {
				for session.Next() {
					if err := session.Scan(&username, &flare, &last_ip); err != nil {
						fmt.Printf("ERR! %s", err)
					} else {
						verified_user.ID = id
						verified_user.UserName = username
						verified_user.Flare = flare
						verified_user.IPAddress = last_ip
					}
				}
			}

			defer session.Close()
			defer db.Close()
		}
	}

	return verified_user
}

func VerifyCredentials(token string, id int) bool {
	verified := false

	if len(token) > 0 && id > 0 {
		db, err := sql.Open("mysql", db_str)

		if err != nil {
			log.Fatal(err)
		} else {
			session, err := db.Query("SELECT COUNT(*) AS x FROM userbase WHERE token=\"" + EscapeSQL(token) + "\" AND id=" + fmt.Sprint(id) + " AND deleted=0")

			var x int

			if err != nil {
				log.Fatal(err)
			} else {
				for session.Next() {
					if err := session.Scan(&x); err != nil {
						fmt.Printf("ERR! %s", err)
					} else {
						verified = (x > 0)
					}
				}
			}

			defer session.Close()
			defer db.Close()
		}
	}

	return verified
}

func FetchUserFromUsername(username string) user {
	target_user := user{}

	if len(username) > 0 {
		db, err := sql.Open("mysql", db_str)

		if err != nil {
			log.Fatal(err)
		} else {
			session, err := db.Query("SELECT id, flare, created, location, karma, gender FROM userbase WHERE username=\"" + EscapeSQL(username) + "\" AND deleted=0")

			var id int
			var flare int
			var created string
			var location string
			var karma int
			var gender int

			if err != nil {
				log.Fatal(err)
			} else {
				for session.Next() {
					if err := session.Scan(&id, &flare, &created, &location, &karma, &gender); err != nil {
						fmt.Printf("ERR! %s", err)
					} else {
						target_user.ID = id
						target_user.UserName = username
						target_user.Flare = flare
						target_user.Created = created
						target_user.Location = location
						target_user.Karma = karma
						target_user.Gender = gender
					}
				}
			}

			defer session.Close()
			defer db.Close()
		}
	}

	return target_user
}

func FetchAccount(id int, token string) account {
	target_account := account{}

	if id > 0 && len(token) > 0 {
		db, err := sql.Open("mysql", db_str)

		if err != nil {
			log.Fatal(err)
		} else {
			session, err := db.Query("SELECT username, email, phone, created, location, gender, flare FROM userbase WHERE token=\"" + EscapeSQL(token) + "\" AND id=" + fmt.Sprint(id) + " AND deleted=0")

			var flare int
			var created string
			var email string
			var phone string
			var username string
			var location string
			var gender int

			if err != nil {
				log.Fatal(err)
			} else {
				for session.Next() {
					if err := session.Scan(&username, &email, &phone, &created, &location, &gender, &flare); err != nil {
						fmt.Printf("ERR! %s", err)
					} else {
						target_account.ID = id
						target_account.UserName = username
						target_account.Email = email
						target_account.Phone = phone
						target_account.Flare = flare
						target_account.Created = created
						target_account.Location = location
						target_account.Gender = gender
					}
				}
			}

			defer session.Close()
			defer db.Close()
		}
	}

	return target_account
}

func FetchLoginHistory(id int) []loginhistory {
	login_history := []loginhistory{}

	if id > 0 {
		db, err := sql.Open("mysql", db_str)

		if err != nil {
			log.Fatal(err)
		} else {
			session, err := db.Query("SELECT account AS account_id, ip_address, platform, login_date_time FROM logins WHERE account=" + fmt.Sprint(id) + " ORDER BY login_date_time DESC LIMIT 5 OFFSET 0")

			var account_id int
			var ip_address string
			var platform string
			var login_date_time string

			if err != nil {
				log.Fatal(err)
			} else {
				for session.Next() {
					if err := session.Scan(&account_id, &ip_address, &platform, &login_date_time); err != nil {
						fmt.Printf("ERR! %s", err)
					} else {
						login := loginhistory{
							AccountID:   account_id,
							IPAddress:   ip_address,
							Platform:    platform,
							DateAndTime: login_date_time,
						}

						login_history = append(login_history, login)
					}
				}
			}

			defer session.Close()
			defer db.Close()
		}
	}

	return login_history
}

func UpdateLoginHistory(account_id int, ip_address string, platform string) {
	db, err := sql.Open("mysql", db_str)

	if err != nil {
		log.Fatal(err)
	} else {
		_, err := db.Exec("INSERT INTO logins (account, ip_address, platform, login_date_time) VALUES(" + fmt.Sprint(account_id) + ", \"" + EscapeSQL(ip_address) + "\", \"" + EscapeBoth(platform) + "\", (SELECT NOW()))")

		if err != nil {
			log.Fatal(err)
		}

		defer db.Close()
	}
}

func CreateAccount(username string, email string, password string, location string, gender int) int {
	new_account_id := 0

	db, err := sql.Open("mysql", db_str)

	if err != nil {
		log.Fatal(err)
	} else {
		new_token := GenerateSHA256(username + email + fmt.Sprint(time.Now().Unix()))
		creation, err := db.Exec("INSERT INTO userbase (username, email, token, password, location, gender, created) VALUES(\"" + EscapeBoth(username) + "\", \"" + EscapeBoth(email) + "\", \"" + new_token + "\", \"" + GenerateSHA256(password) + "\", \"" + EscapeBoth(location) + "\", " + fmt.Sprint(gender) + ", (SELECT NOW()))")

		if err != nil {
			log.Fatal(err)
		} else {
			lii, err := creation.LastInsertId()

			if err == nil {
				new_account_id = int(lii)
			}
		}

		defer db.Close()
	}

	return new_account_id
}

func CheckUsernameIsUnqiue(username string) bool {
	result := false
	db, err := sql.Open("mysql", db_str)

	if err != nil {
		log.Fatal(err)
	} else {
		check, err := db.Query("SELECT COUNT(*) AS total FROM userbase WHERE username=\"" + EscapeBoth(username) + "\" AND deleted=0")

		var total int

		if err != nil {
			log.Fatal(err)
		} else {
			for check.Next() {
				if err := check.Scan(&total); err != nil {
					fmt.Printf("ERR! %s", err)
				} else {
					result = (total == 0)
				}
			}
		}

		defer check.Close()
		defer db.Close()
	}

	return result
}

func CheckEmailIsUnqiue(email string) bool {
	result := false
	db, err := sql.Open("mysql", db_str)

	if err != nil {
		log.Fatal(err)
	} else {
		check, err := db.Query("SELECT COUNT(*) AS total FROM userbase WHERE email=\"" + EscapeBoth(email) + "\" AND deleted=0")

		var total int

		if err != nil {
			log.Fatal(err)
		} else {
			for check.Next() {
				if err := check.Scan(&total); err != nil {
					fmt.Printf("ERR! %s", err)
				} else {
					result = (total == 0)
				}
			}
		}

		defer check.Close()
		defer db.Close()
	}

	return result
}

func UpdateAccountDetails(id int, token string, username string, email string, phone string, location string, gender int) bool {
	affected_rows := 0
	db, err := sql.Open("mysql", db_str)

	if err != nil {
		log.Fatal(err)
	} else {
		update, err := db.Exec("UPDATE userbase SET username=\"" + EscapeBoth(username) + "\", email=\"" + EscapeBoth(email) + "\", phone=\"" + EscapeBoth(phone) + "\", location=\"" + EscapeBoth(location) + "\", gender=" + fmt.Sprint(gender) + " WHERE id=" + fmt.Sprint(id) + " AND token=\"" + EscapeSQL(token) + "\" AND deleted=0")

		if err != nil {
			log.Fatal(err)
		} else {
			result, err := update.RowsAffected()

			if err == nil {
				affected_rows = int(result)
			}
		}

		defer db.Close()
	}

	return (affected_rows > 0)
}

func UpdateKarma(id int, token string, amount int, self bool) {
	db, err := sql.Open("mysql", db_str)

	if err != nil {
		log.Fatal(err)
	} else {
		_, err := db.Exec("UPDATE userbase SET karma=(karma + " + fmt.Sprint(amount) + ") WHERE id=" + fmt.Sprint(id) + " AND token=\"" + EscapeSQL(token) + "\" AND deleted=0")

		if err != nil {
			log.Fatal(err)
		}

		defer db.Close()
	}
}

func SearchUsers(term string) []user {
	records := []user{}

	if len(term) > 0 {
		db, err := sql.Open("mysql", db_str)

		if err != nil {
			log.Fatal(err)
		} else {
			session, err := db.Query("SELECT id, username, flare, created, location, karma, gender FROM userbase WHERE username LIKE \"%" + EscapeSQL(term) + "%\" AND deleted=0")

			var id int
			var username string
			var flare int
			var created string
			var location string
			var karma int
			var gender int

			if err != nil {
				log.Fatal(err)
			} else {
				for session.Next() {
					if err := session.Scan(&id, &username, &flare, &created, &location, &karma, &gender); err != nil {
						fmt.Printf("ERR! %s", err)
					} else {
						record := user{
							ID:       id,
							UserName: username,
							Flare:    flare,
							Created:  created,
							Location: location,
							Karma:    karma,
							Gender:   gender,
						}

						records = append(records, record)
					}
				}
			}

			defer session.Close()
			defer db.Close()
		}
	}

	return records
}

func FetchKarmaLeaderboard() []user {
	records := []user{}

	db, err := sql.Open("mysql", db_str)

	if err != nil {
		log.Fatal(err)
	} else {
		session, err := db.Query("SELECT id, username, karma, gender, location FROM userbase WHERE deleted=0 ORDER BY karma DESC")

		var id int
		var username string
		var location string
		var karma int
		var gender int

		if err != nil {
			log.Fatal(err)
		} else {
			for session.Next() {
				if err := session.Scan(&id, &username, &karma, &gender, &location); err != nil {
					fmt.Printf("ERR! %s", err)
				} else {
					record := user{
						ID:       id,
						UserName: username,
						Karma:    karma,
						Gender:   gender,
						Location: location,
					}

					records = append(records, record)
				}
			}
		}

		defer session.Close()
		defer db.Close()
	}

	return records
}

func FetchFeaturedUsers() []user {
	records := []user{}

	db, err := sql.Open("mysql", db_str)

	if err != nil {
		log.Fatal(err)
	} else {
		session, err := db.Query("SELECT id, username, karma FROM userbase WHERE deleted=0 ORDER BY karma DESC LIMIT 5 OFFSET 0")

		var id int
		var username string
		var karma int

		if err != nil {
			log.Fatal(err)
		} else {
			for session.Next() {
				if err := session.Scan(&id, &username, &karma); err != nil {
					fmt.Printf("ERR! %s", err)
				} else {
					record := user{
						ID:       id,
						UserName: username,
						Karma:    karma,
					}

					records = append(records, record)
				}
			}
		}

		defer session.Close()
		defer db.Close()
	}

	return records
}

func ChangePassword(id int, token string, old_password string, new_password string) bool {
	affected_rows := 0

	new_password_hash := GenerateSHA256(new_password)
	old_password_hash := GenerateSHA256(old_password)

	db, err := sql.Open("mysql", db_str)

	if err != nil {
		log.Fatal(err)
	} else {
		update, err := db.Exec("UPDATE userbase SET password=\"" + new_password_hash + "\" WHERE id=" + fmt.Sprint(id) + " AND token=\"" + EscapeSQL(token) + "\" AND password=\"" + old_password_hash + "\" AND deleted=0")

		if err != nil {
			log.Fatal(err)
		} else {
			result, err := update.RowsAffected()

			if err == nil {
				affected_rows = int(result)
			}
		}

		defer db.Close()
	}

	return (affected_rows > 0)
}