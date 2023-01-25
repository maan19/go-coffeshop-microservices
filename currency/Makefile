.PHONY: protos

protos-nouse:
	protoc -I protos/ protos/currency.proto  --go-grpc_out=protos/currency


protosmessages:
	protoc --go_out=protos/currency --proto_path=protos/ protos/currency_msg.proto

protosservice:
	protoc --go-grpc_out=require_unimplemented_servers=false:protos/currency --proto_path=protos/ protos/currency.proto

protosone:
	protoc --go_out=protos/currency  --go-grpc_out=require_unimplemented_servers=false:protos/currency  protos/currency.proto --proto_path=protos/ protos/currency.proto

