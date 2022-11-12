package router

import (
	"github.com/go-echarts/go-echarts/v2/opts"
	botGithub "github.com/tinode/chat/server/extra/bots/github"
	"github.com/tinode/chat/server/extra/vendors"
	"github.com/tinode/chat/server/extra/vendors/dropbox"
	"github.com/tinode/chat/server/extra/vendors/github"
	"github.com/tinode/chat/server/extra/vendors/pocket"
	"math/rand"
)

func newProvider(category string) vendors.OAuthProvider {
	var provider vendors.OAuthProvider

	switch category {
	case pocket.ID:
		p := pocket.NewPocket("", "", "", "")
		provider = p
	case github.ID:
		provider = github.NewGithub(botGithub.Config.ID, botGithub.Config.Secret, "", "")
	case dropbox.ID:
		provider = dropbox.NewDropbox("", "", "", "")
	default:
		return nil
	}

	return provider
}

// generate random data for bar chart
func generateBarItems() []opts.BarData {
	items := make([]opts.BarData, 0)
	for i := 0; i < 7; i++ {
		items = append(items, opts.BarData{Value: rand.Intn(300)})
	}
	return items
}
