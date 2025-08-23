# Proto generation

Dependencies:
- `protoc`
- Go plugins: `protoc-gen-go`, `protoc-gen-go-grpc`, `protoc-gen-grpc-gateway`
- Google API annotations under `third_party/googleapis`

Install tools:
```bash
go install google.golang.org/protobuf/cmd/protoc-gen-go@latest
go install google.golang.org/grpc/cmd/protoc-gen-go-grpc@latest
go install github.com/grpc-ecosystem/grpc-gateway/v2/protoc-gen-grpc-gateway@latest
```

Fetch googleapis (annotations.proto):
```bash
mkdir -p third_party && cd third_party
git clone https://github.com/googleapis/googleapis.git
cd ..
```

Generate:
```bash
make proto
``` 