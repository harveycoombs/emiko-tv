package main

import (
	"time"
)

type ratelimitdata struct {
	ID int
	Token string
	InitialTimestamp int64
	RequestsCount int
}

var limited_credentials = []ratelimitdata{}

func AddCredentialsToTimeout(id int, token string) {
	if IndexOfRateLimitedCredentials(id, token) != -1 {
		return
	}

	new_creds := ratelimitdata{
		ID: id, 
		Token: token,
		InitialTimestamp: time.Now().UnixMilli(),
	}

	limited_credentials = append(limited_credentials, new_creds)
}

func RemoveCredentialsFromTimeout(id int, token string) {
	creds_index := IndexOfRateLimitedCredentials(id, token)

	if creds_index != -1 {
		limited_credentials = append(limited_credentials[:creds_index], limited_credentials[creds_index + 1:]...)
	}
}

func TimeoutHasExpired(id int, token string) bool {
	creds_index := IndexOfRateLimitedCredentials(id, token)

	if creds_index != -1 {
		creds := limited_credentials[creds_index]

		then := time.Unix(creds.InitialTimestamp, 0)
		now := time.Now()
		
		return then.Sub(now).Seconds() >= 60
	} else {
		return false
	}
}

func IndexOfRateLimitedCredentials(id int, token string) int {
	for x, creds := range limited_credentials {
		if creds.ID == id && creds.Token == token {
			return x
		}
	}

	return -1
}