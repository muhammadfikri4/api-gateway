package model

import "time"

type Service struct {
	ID          uint      `json:"id" gorm:"primaryKey"`
	Name        string    `json:"name" gorm:"uniqueIndex;not null"`
	Description string    `json:"description"`
	TargetURL   string    `json:"target_url" gorm:"not null"`
	RoutePrefix string    `json:"route_prefix" gorm:"uniqueIndex;not null"`
	IsActive    bool      `json:"is_active" gorm:"default:true"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}

type CreateServiceRequest struct {
	Name        string `json:"name" binding:"required"`
	Description string `json:"description"`
	TargetURL   string `json:"target_url" binding:"required,url"`
	RoutePrefix string `json:"route_prefix" binding:"required"`
}

type UpdateServiceRequest struct {
	Name        string `json:"name"`
	Description string `json:"description"`
	TargetURL   string `json:"target_url"`
	RoutePrefix string `json:"route_prefix"`
	IsActive    *bool  `json:"is_active"`
}
