module github.com/vitaliy/file-storage/server

go 1.22.2

require (
	github.com/google/uuid v1.6.0
	github.com/vitaliy/file-storage/common v0.0.0-00010101000000-000000000000
)

replace github.com/vitaliy/file-storage/common => ../common
