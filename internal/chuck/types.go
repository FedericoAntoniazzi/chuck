package chuck

type Image struct {
	Registry string
	Name     string
	Tag      string
}

type Container struct {
	Name         string
	Image        Image
	ImageUpdates []string
}
