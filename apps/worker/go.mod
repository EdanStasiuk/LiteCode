module github.com/EdanStasiuk/LiteCode/apps/worker

go 1.24.5

replace github.com/EdanStasiuk/LiteCode/pkg => ../../pkg

replace github.com/EdanStasiuk/LiteCode/apps/backend/server => ../backend/server

require (
	github.com/EdanStasiuk/LiteCode/apps/backend/server v0.0.0-20250909030020-ec7e12baa137
	github.com/EdanStasiuk/LiteCode/pkg v0.0.0-00010101000000-000000000000
	github.com/segmentio/kafka-go v0.4.49
)

require (
	github.com/gocql/gocql v1.7.0 // indirect
	github.com/golang/snappy v0.0.3 // indirect
	github.com/hailocab/go-hostpool v0.0.0-20160125115350-e80d13ce29ed // indirect
	github.com/jinzhu/inflection v1.0.0 // indirect
	github.com/jinzhu/now v1.1.5 // indirect
	github.com/klauspost/compress v1.15.9 // indirect
	github.com/pierrec/lz4/v4 v4.1.15 // indirect
	golang.org/x/text v0.29.0 // indirect
	gopkg.in/inf.v0 v0.9.1 // indirect
	gorm.io/gorm v1.31.0 // indirect
)
