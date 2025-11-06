package oauth2

type Provider string

const (
	Google   Provider = "google"
	GitHub   Provider = "github"
	Facebook Provider = "facebook"
	Twitter  Provider = "twitter"
	LinkedIn Provider = "linkedin"
)

func (p Provider) String() string {
	return string(p)
}
