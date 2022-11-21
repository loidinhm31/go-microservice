module authentication

go 1.19

require (
	github.com/go-chi/chi/v5 v5.0.7
	github.com/go-chi/cors v1.2.1
	github.com/jackc/pgconn v1.13.0
	github.com/jackc/pgx/v5 v5.1.1
	github.com/loidinhm31/go-micro/common v0.0.0-00010101000000-000000000000
	golang.org/x/crypto v0.3.0
)

require (
	github.com/jackc/chunkreader/v2 v2.0.1 // indirect
	github.com/jackc/pgio v1.0.0 // indirect
	github.com/jackc/pgpassfile v1.0.0 // indirect
	github.com/jackc/pgproto3/v2 v2.3.1 // indirect
	github.com/jackc/pgservicefile v0.0.0-20200714003250-2b9c44734f2b // indirect
	golang.org/x/text v0.4.0 // indirect
)

replace github.com/loidinhm31/go-micro/common => ../common
