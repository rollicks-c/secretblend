package secretblend

type protocol string

var (
	protocolProviders map[protocol]Provider
	globalProviders   []Provider
)

func init() {
	protocolProviders = make(map[protocol]Provider)
	globalProviders = make([]Provider, 0)
}
