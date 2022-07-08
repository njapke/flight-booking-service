module github.com/christophwitzko/flight-booking-service

go 1.18

require (
	github.com/brianvoe/gofakeit/v6 v6.17.0
	github.com/dgraph-io/badger/v3 v3.2103.2
	github.com/go-chi/chi/v5 v5.0.7
	github.com/google/uuid v1.3.0
	github.com/stretchr/testify v1.8.0
	go.uber.org/zap v1.21.0
)

require (
	github.com/cespare/xxhash v1.1.0 // indirect
	github.com/cespare/xxhash/v2 v2.1.1 // indirect
	github.com/davecgh/go-spew v1.1.1 // indirect
	github.com/dgraph-io/ristretto v0.1.0 // indirect
	github.com/dustin/go-humanize v1.0.0 // indirect
	github.com/gogo/protobuf v1.3.2 // indirect
	github.com/golang/glog v0.0.0-20160126235308-23def4e6c14b // indirect
	github.com/golang/groupcache v0.0.0-20190702054246-869f871628b6 // indirect
	github.com/golang/protobuf v1.3.1 // indirect
	github.com/golang/snappy v0.0.3 // indirect
	github.com/google/flatbuffers v1.12.1 // indirect
	github.com/klauspost/compress v1.12.3 // indirect
	github.com/pkg/errors v0.9.1 // indirect
	github.com/pmezard/go-difflib v1.0.0 // indirect
	go.opencensus.io v0.22.5 // indirect
	go.uber.org/atomic v1.9.0 // indirect
	go.uber.org/multierr v1.8.0 // indirect
	golang.org/x/net v0.0.0-20210405180319-a5a99cb37ef4 // indirect
	golang.org/x/sys v0.0.0-20210510120138-977fb7262007 // indirect
	gopkg.in/yaml.v3 v3.0.1 // indirect
)

replace github.com/go-chi/chi/v5 => ./chi
