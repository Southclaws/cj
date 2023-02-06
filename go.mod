module github.com/Southclaws/cj

go 1.12

replace github.com/go-resty/resty => gopkg.in/resty.v1 v1.12.0

replace gopkg.in/russross/blackfriday.v2 => github.com/russross/blackfriday/v2 v2.0.1

require (
	github.com/Bios-Marcel/discordemojimap v1.0.1
	github.com/bwmarrin/discordgo v0.23.3-0.20210515023446-8dc42757bea5
	github.com/globalsign/mgo v0.0.0-20181015135952-eeefdecb41b8
	github.com/google/go-github/v28 v28.1.1
	github.com/google/go-querystring v1.1.0 // indirect
	github.com/joho/godotenv v1.5.1
	github.com/kelseyhightower/envconfig v1.4.0
	github.com/kr/pretty v0.1.0 // indirect
	github.com/kristoisberg/gonesyntees v0.0.0-20190527174556-0595a02f9399
	github.com/patrickmn/go-cache v2.1.0+incompatible
	github.com/pkg/errors v0.9.1
	github.com/robfig/cron v1.2.0
	github.com/stretchr/testify v1.7.1
	go.uber.org/multierr v1.7.0 // indirect
	go.uber.org/zap v1.17.0
	golang.org/x/crypto v0.0.0-20210513164829-c07d793c2f9a // indirect
	golang.org/x/net v0.0.0-20210510120150-4163338589ed // indirect
	golang.org/x/sys v0.0.0-20210514084401-e8d321eab015 // indirect
	gopkg.in/check.v1 v1.0.0-20180628173108-788fd7840127 // indirect
	gopkg.in/xmlpath.v2 v2.0.0-20150820204837-860cbeca3ebc
)
