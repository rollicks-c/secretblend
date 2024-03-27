package secretblend

type protocol string

var (
	providerList map[protocol]Provider
)

func init() {
	providerList = make(map[protocol]Provider)
}
