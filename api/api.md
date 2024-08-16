## Installazione
Installa i plugin del compilatore di protocol buffer:
```bash
go install google.golang.org/protobuf/cmd/protoc-gen-go@latest
go install google.golang.org/grpc/cmd/protoc-gen-go-grpc@latest
```

Compila i file `.proto`:
```bash
protoc *.proto --go_out=./ --go-grpc_out=./
```



