package models

type ImageUpdate struct {
	Name        string
	CurrentTag  string
	UpdatedTags []string
}
