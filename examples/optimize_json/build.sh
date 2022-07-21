# you can use jsoniter faster than native json
go build -tags "jsoniter" -i -v -o output/bin/${RUN_NAME}

# you can use go_json faster than jsoniter
go build -tags "go_json" -i -v -o output/bin/${RUN_NAME}