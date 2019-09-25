// Copyright (c) 2018-2019 The Decred developers
// Use of this source code is governed by an ISC
// license that can be found in the LICENSE file.

package commstats

import (
	"context"
	"github.com/raedahgroup/dcrextdata/app/config"
	"net/http"
	"time"
)

type CommStat struct {
	Date               time.Time         `json:"date"`
	RedditStats        map[string]Reddit `json:"reddit_stats"`
	TwitterFollowers   int               `json:"twitter_followers"`
	YoutubeSubscribers int               `json:"youtube_subscribers"`
	GithubStars        int               `json:"github_stars"`
	GithubFolks        int               `json:"github_folks"`
}

type RedditResponse struct {
	Kind string `json:"kind"`
	Data Reddit `json:"data"`
}

type Reddit struct {
	Date time.Time `json:"date"`
	Subscribers    int `json:"subscribers"`
	AccountsActive int `json:"active_user_count"`
}

type Github struct {
	Date time.Time `json:"date"`
	Star int `json:"star"`
	Folks int `json:"folks"`
}

type Youtube struct {
	Date time.Time `json:"date"`
	Subscribers int `json:"subscribers"`
}

type Twitter struct {
	Date time.Time `json:"date"`
	Followers int `json:"followers"`
}

type DataStore interface {
	StoreCommStat(context.Context, CommStat) error
	LastCommStatEntry() (time time.Time)

}

type Collector struct {
	client    http.Client
	period    time.Duration
	dataStore DataStore
	options   *config.CommunityStatOptions
}
