package types

// Image describes a container image.
// A container image is composed by Registry/Namespace/Name:Tag (e.g docker.io/library/nginx:1.25)
type Image struct {
	Raw       string // Unparsed image reference
	Registry  string // Registry's URL
	Namespace string // Image namespace (In case of Docker Hub may be library, or the username)
	Name      string // Name of the image (e.g nginx, redis)
	Tag       string // Tag assigned to the image (e.g. latest, 1.15, v2.1.0)
}
