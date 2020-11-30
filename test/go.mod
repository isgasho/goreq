module github.com/aiscrm/goreq/test

go 1.15

replace (
	github.com/aiscrm/goreq => ../
	github.com/aiscrm/goreq/vo => ../vo
)

require (
	github.com/aiscrm/goreq v0.0.0-00010101000000-000000000000
	github.com/aiscrm/goreq/vo v0.0.0-20201129134153-907e71c1b7d6
)
