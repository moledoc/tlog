
test:
	! GORACE="exitcode=1" go test -v -cpu=1 -race -count=1 -shuffle=off ./... && go run results_compare/results_compare.go

record:
	! GORACE="exitcode=1" go test -v -cpu=1 -race -count=1 -shuffle=off ./... -record=true