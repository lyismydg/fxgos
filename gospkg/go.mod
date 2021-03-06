module github.com/fidelfly/gostool

go 1.12

require (
	github.com/fidelfly/gox v0.0.0-20191013050143-1cc15e873ba6
	github.com/go-sql-driver/mysql v1.4.1
	github.com/go-xorm/xorm v0.7.9
	github.com/golang/protobuf v1.3.2
	github.com/kr/pretty v0.1.0 // indirect
	github.com/sirupsen/logrus v1.4.2
	google.golang.org/genproto v0.0.0-20190819201941-24fa4b261c55 // indirect
	google.golang.org/grpc v1.19.0
	gopkg.in/alexcesaro/quotedprintable.v3 v3.0.0-20150716171945-2caba252f4dc // indirect
	gopkg.in/gomail.v2 v2.0.0-20160411212932-81ebce5c23df
	xorm.io/core v0.7.2
)

replace github.com/fidelfly/gox => ../../gox
