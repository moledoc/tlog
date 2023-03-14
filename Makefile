
test:
	! go test -v -cpu=4 -race -count=1 && go run results_compare/results_compare.go

record:
	! go test -v -cpu=4 -race -count=1 -record