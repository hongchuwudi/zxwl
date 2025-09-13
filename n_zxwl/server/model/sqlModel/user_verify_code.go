// Package sqlModel file: model/sqlModel/user_verify_code.go
package sqlModel

import "time"

type UserVerifyCode struct {
	ID        int       `json:"id" gorm:"primaryKey"`
	Email     string    `json:"email" gorm:"size:88"`
	Code      string    `json:"code" gorm:"size:20"`
	ExpiresAt time.Time `json:"expires_at"`
	CreatedAt time.Time `json:"created_at"`
}
