# GRPC-W2-NGC

<b>Generate proto pb.go files</b>
```bash
protoc --go_out=. --go_opt=paths=source_relative --go-grpc_out=. --go-grpc_opt=paths=source_relative ./common/model/*.proto
```