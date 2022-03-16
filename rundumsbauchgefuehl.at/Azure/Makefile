backend: main.go
	gofmt -w main.go
	go build

publish:
	func azure functionapp publish rundumsbauchgefuehl

logs:
	func azure functionapp logstream rundumsbauchgefuehl
