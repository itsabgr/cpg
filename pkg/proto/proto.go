package proto

//go:generate protoc --proto_path=../../proto --go_opt=Mcpg.proto=../proto --go-grpc_opt=Mcpg.proto=../proto --go_out=. --go-grpc_out=. cpg.proto
