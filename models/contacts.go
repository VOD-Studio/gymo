package models

type Contact struct {
	ID     uint `gorm:"primaryKey" json:"id,omitempty"`
	Uid    uint
	Target uint `gorm:"not null"`
}
