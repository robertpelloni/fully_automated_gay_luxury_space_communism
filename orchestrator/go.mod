module github.com/robertpelloni/hustle/orchestrator

go 1.24.3

require (
	github.com/mattn/go-sqlite3 v1.14.44
	github.com/robertpelloni/hustle/hustle/content v0.0.0-00010101000000-000000000000
	github.com/robertpelloni/hustle/hustle/curation v0.0.0-00010101000000-000000000000
	github.com/robertpelloni/hustle/hustle/research v0.0.0-00010101000000-000000000000
	github.com/robertpelloni/hustle/hustle/social v0.0.0-00010101000000-000000000000
	github.com/robertpelloni/hustle/hustle/trading v0.0.0-00010101000000-000000000000
	golang.org/x/time v0.10.0
)

require (
	github.com/PuerkitoBio/goquery v1.8.0 // indirect
	github.com/andybalholm/cascadia v1.3.1 // indirect
	github.com/dghubble/oauth1 v0.7.3 // indirect
	github.com/json-iterator/go v1.1.12 // indirect
	github.com/jung-kurt/gofpdf v1.16.2 // indirect
	github.com/mmcdole/gofeed v1.3.0 // indirect
	github.com/mmcdole/goxpp v1.1.1-0.20240225020742-a0c311522b23 // indirect
	github.com/modern-go/concurrent v0.0.0-20180306012644-bacd9c7ef1dd // indirect
	github.com/modern-go/reflect2 v1.0.2 // indirect
	golang.org/x/net v0.38.0 // indirect
	golang.org/x/oauth2 v0.27.0 // indirect
	golang.org/x/text v0.23.0 // indirect
)

replace github.com/robertpelloni/hustle/hustle/content => ../hustle/content

replace github.com/robertpelloni/hustle/hustle/curation => ../hustle/curation

replace github.com/robertpelloni/hustle/hustle/research => ../hustle/research

replace github.com/robertpelloni/hustle/hustle/social => ../hustle/social

replace github.com/robertpelloni/hustle/hustle/trading => ../hustle/trading
