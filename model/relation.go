package model

type Follow struct {
	Id         int64 `json:"id" form:"id"`
	UserId     int64 `json:"user_id" form:"user_id"`
	ToUserId   int64 `json:"to_user_id" form:"to_user_id"`
	ActionType int64 `json:"action_type" form:"action_type"`
}

type Follow_sql struct {
	Id       int64 `json:"id"`
	UserId   int64 `json:"user_id"`
	TargetId int64 `json:"target_id"`
}

type Followed_sql struct {
	Id       int64 `json:"id"`
	UserId   int64 `json:"user_id"`
	TargetId int64 `json:"target_id"`
}
