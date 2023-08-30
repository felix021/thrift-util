module tests

go 1.20

replace github.com/felix021/thrift-util => ./..

require (
	github.com/apache/thrift v0.0.0-00010101000000-000000000000
	github.com/felix021/thrift-util v0.0.0-00010101000000-000000000000
)

replace github.com/apache/thrift => github.com/apache/thrift v0.13.0
