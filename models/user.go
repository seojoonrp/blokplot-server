// models/user.go

package models

type Profile struct {
	AvatarUrl string `bson:"avatarUrl" json:"avatarUrl"`
	BannerUrl string `bson:"bannerUrl" json:"bannerUrl"`
}

type Stats struct {
	Wins int `bson:"wins" json:"wins"`
	LiarWins int `bson:"liar-wins" json:"liarWins"`
	Trophies int `bson:"trophies" json:"trophies"`
	Coins int `bson:"coins" json:"coins"`
}

type User struct {
	Username string `bson:"username" json:"username"`
	Profile Profile `bson:"profile" json:"profile"`
	Stats Stats `bson:"stats" json:"stats"`
	Friends []string `bson:"friends" json:"friends"`
}

type LoginRequest struct {
	Username string `json:"username"`
}