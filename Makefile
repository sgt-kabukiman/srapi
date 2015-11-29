default: build

build: fix
	go build -v .

test: fix
	go test -v

fix: *.go
	goimports -l -w .
	gofmt -l -w .

makegen:
	cd generate && make

gen: *.go
	generate/generate -type Category -plural Categories
	generate/generate -type Game -plural Games
	generate/generate -type Leaderboard -plural Leaderboards
	generate/generate -type Level -plural Levels
	generate/generate -type PersonalBest -plural PersonalBests
	generate/generate -type Platform -plural Platforms
	generate/generate -type Region -plural Regions
	generate/generate -type Run -plural Runs
	generate/generate -type Series -plural ManySeries
	generate/generate -type User -plural Users
	generate/generate -type Variable -plural Variables
