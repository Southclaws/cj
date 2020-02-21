module github.com/Southclaws/cj

go 1.12

replace github.com/go-resty/resty => gopkg.in/resty.v1 v1.12.0

replace gopkg.in/russross/blackfriday.v2 => github.com/russross/blackfriday/v2 v2.0.1

require (
	github.com/PuerkitoBio/goquery v1.5.0
	github.com/Southclaws/go-cloudflare-scraper v0.0.0-20171030160446-76acfe58205d
	github.com/Southclaws/samp-servers-api v0.0.0-20190501054307-50d4ce94e27b
	github.com/bwmarrin/discordgo v0.20.2
	github.com/globalsign/mgo v0.0.0-20181015135952-eeefdecb41b8
	github.com/google/go-github/v28 v28.1.1
	github.com/joho/godotenv v1.3.0
	github.com/kelseyhightower/envconfig v1.4.0
	github.com/kristoisberg/gonesyntees v0.0.0-20190301122441-7d230b161c5b
	github.com/mb-14/gomarkov v0.0.0-20190125094512-044dd0dcb5e7
	github.com/microcosm-cc/bluemonday v1.0.2
	github.com/patrickmn/go-cache v2.1.0+incompatible
	github.com/pkg/errors v0.9.1
	github.com/robertkrimen/otto v0.0.0-20180617131154-15f95af6e78d // indirect
	github.com/robfig/cron v0.0.0-20180505203441-b41be1df6967
	github.com/shurcooL/sanitized_anchor_name v1.0.0 // indirect
	github.com/stretchr/testify v1.4.0
	github.com/texttheater/golang-levenshtein v0.0.0-20180516184445-d188e65d659e
	go.uber.org/zap v1.14.0
	gopkg.in/resty.v1 v1.12.0
	gopkg.in/russross/blackfriday.v2 v2.0.0-00010101000000-000000000000
	gopkg.in/sourcemap.v1 v1.0.5 // indirect
	gopkg.in/src-d/go-git.v4 v4.13.1
	gopkg.in/xmlpath.v2 v2.0.0-20150820204837-860cbeca3ebc
)
