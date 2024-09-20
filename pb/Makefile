all: fxwatcher fxgateway swagger

fxwatcher:
	protoc -I . \
   	 	--go_out ./ --go_opt paths=source_relative \
    		--go-grpc_out ./ --go-grpc_opt paths=source_relative \
    		fxwatcher.proto

fxgateway:
	protoc -I . \
    		--go_out ./ --go_opt paths=source_relative \
    		--go-grpc_out ./ --go-grpc_opt paths=source_relative \
    		fxgateway.proto
	
	protoc -I . --grpc-gateway_out ./ \
    		--grpc-gateway_opt logtostderr=true \
    		--grpc-gateway_opt paths=source_relative \
    		--grpc-gateway_opt generate_unbound_methods=true \
 		fxgateway.proto

swagger:
	protoc -I . --go_out . --go-grpc_out . --grpc-gateway_out . --openapiv2_out . --openapiv2_opt use_go_templates=true fxgateway.proto
	go generate .
