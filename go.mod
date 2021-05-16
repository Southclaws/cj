module github.com/Southclaws/cj

go 1.12

replace github.com/go-resty/resty => gopkg.in/resty.v1 v1.12.0

replace gopkg.in/russross/blackfriday.v2 => github.com/russross/blackfriday/v2 v2.0.1

require (
	github.com/Bios-Marcel/discordemojimap v1.0.1
	github.com/bwmarrin/discordgo v0.23.1
	github.com/globalsign/mgo v0.0.0-20181015135952-eeefdecb41b8
	github.com/google/go-github/v28 v28.1.1
	github.com/gorilla/websocket v1.4.2 // indirect
	github.com/joho/godotenv v1.3.0
	github.com/kelseyhightower/envconfig v1.4.0
	github.com/kristoisberg/gonesyntees v0.0.0-20190301122441-7d230b161c5b
	github.com/patrickmn/go-cache v2.1.0+incompatible
	github.com/pkg/errors v0.9.1
	github.com/robfig/cron v0.0.0-20180505203441-b41be1df6967
	github.com/stretchr/testify v1.7.0
	github.com/texttheater/golang-levenshtein v0.0.0-20180516184445-d188e65d659e
	go.uber.org/zap v1.15.0
	golang.org/x/crypto v0.0.0-20201002170205-7f63de1d35b0 // indirect
	golang.org/x/net v0.0.0-20200202094626-16171245cfb2 // indirect
	golang.org/x/sys v0.0.0-20201008064518-c1f3e3309c71 // indirect
	gopkg.in/xmlpath.v2 v2.0.0-20150820204837-860cbeca3ebc
)
