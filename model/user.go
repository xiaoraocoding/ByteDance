package model

type User struct {
	Id             int    `json:"id,omitempty"`
	Username       string `json:"name,omitempty"`
	Password       string
	Follow_Count   int  `json:"follow_count"`
	Follower_Count int  `json:"follower_count"`
	IsFollow       bool `json:"is_follow,omitempty"`
}
