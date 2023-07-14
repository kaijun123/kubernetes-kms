module github.com/kaijun123/kubernetes-kms/mock/plugins

go 1.20

replace k8s.io/kms => github.com/kaijun123/kubernetes-kms v0.0.0-20230713060447-d5a4726b8121

require (
	github.com/kaijun123/kubernetes-kms v0.0.0-20230713174642-ce449db4f752
	k8s.io/klog/v2 v2.100.1
)

require (
	github.com/go-logr/logr v1.2.4 // indirect
	github.com/gogo/protobuf v1.3.2 // indirect
	github.com/golang/protobuf v1.5.3 // indirect
	golang.org/x/net v0.9.0 // indirect
	golang.org/x/sys v0.8.0 // indirect
	golang.org/x/text v0.9.0 // indirect
	google.golang.org/genproto/googleapis/rpc v0.0.0-20230525234030-28d5490b6b19 // indirect
	google.golang.org/grpc v1.54.0 // indirect
	google.golang.org/protobuf v1.30.0 // indirect
)
