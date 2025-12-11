package entity

type User struct {
	ID           uint64 `json:"id" gorm:"primaryKey;autoIncrement;not null;type:bigint"`
	UserName     string `json:"userName" gorm:"not null;type:varchar(191)"`
	Email        string `json:"email" gorm:"uniqueIndex;not null;type:varchar(191)"`
	Password     string `json:"password" gorm:"not null;type:varchar(191)"`
	Avatar       string `json:"avatar" gorm:"type:varchar(191);default:''"`
	CreatedAt    int64  `json:"createdAt" gorm:"autoCreateTime:milli"`
	UpdatedAt    int64  `json:"updatedAt" gorm:"autoUpdateTime:milli"`
	LastActiveAt int64  `json:"lastActiveAt" gorm:"autoUpdateTime:milli"`
}
