SET CGO_ENABLED=0
SET GOOS=linux
SET GOARCH=riscv64
go build  -v -o ./httpPush github.com/chuccp/smtp2http
