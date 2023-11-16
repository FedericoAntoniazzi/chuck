package models

import "strings"

// Container describes common parameters of a container
type Container struct {
	Image  string
	Labels map[string]string
	Name   string
}

func (c Container) ImageName() string {
	imgSplit := strings.Split(c.Image, ":")
	return imgSplit[0]
}

func (c Container) ImageTag() string {
	imgSplit := strings.Split(c.Image, ":")

	if len(imgSplit) == 1 {
		return "latest"
	}

	return imgSplit[1]
}
