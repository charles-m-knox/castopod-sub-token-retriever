module git.cmcode.dev/cmcode/castopod-sub-token-retriever/examples/server

go 1.23.0

require (
	git.cmcode.dev/cmcode/castopod-sub-token-retriever v0.0.0-20240821193910-7945c5412320
	git.cmcode.dev/cmcode/go-castopod v0.0.3
	github.com/go-sql-driver/mysql v1.8.1
)

require filippo.io/edwards25519 v1.1.0 // indirect

replace git.cmcode.dev/cmcode/castopod-sub-token-retriever => ../../
