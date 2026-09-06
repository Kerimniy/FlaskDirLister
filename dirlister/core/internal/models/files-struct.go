package models

import "time"

type File struct {
	ID      uint   `gorm:"primaryKey"`
	Dir     string `gorm:"not null;index:idx_dir;uniqueIndex:idx_dir_name"`
	Name    string `gorm:"not null;index:idx_name;uniqueIndex:idx_dir_name"`
	Size    int64
	ModTime time.Time
	IsDir   bool
}

type EntryInfo struct {
	Name     string    `json:"name"`
	Type     string    `json:"type"`
	Size     int64     `json:"size"`
	FullName string    `json:"fullName"`
	ModTime  time.Time `json:"modTime"`
}

type UploadData struct {
	NewName string `json:"newName"`
	Content string `json:"content"`
}
