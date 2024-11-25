package firehose

import (
	"encoding/json"
)

type RepoMessage struct {
	DID    string `json:"did"`
	TimeUS int64  `json:"time_us"`
	Kind   string `json:"kind"`
	Commit Commit `json:"commit"`
}

type Commit struct {
	Rev        string          `json:"rev"`
	Operation  string          `json:"operation"`
	Collection string          `json:"collection"`
	RKey       string          `json:"rkey"`
	Record     json.RawMessage `json:"record"`
	CID        string          `json:"cid"`
}

type Record struct {
	Type      string `json:"$type"`
	CreatedAt string `json:"createdAt"`
	Text      string `json:"text"`
	Reply     Reply  `json:"reply,omitempty"`
}

type CidAndUri struct {
	CID string `json:"cid"`
	URI string `json:"uri"`
}

type Reply struct {
	Parent CidAndUri `json:"parent"`
	Root   CidAndUri `json:"root"`
}

type Feature struct {
	Type string `json:"$type"`
	Tag  string `json:"tag"`
}

func (r *RepoMessage) URI() string {
	return "at://" + r.DID + "/app.bsky.feed.post/" + r.Commit.RKey;
}