package core

type StorageType string

const (
	STORETYPE_UNKNOWN  StorageType = ""
	STORETYPE_POSTGRES StorageType = "POSTGRES"
	STORETYPE_MYSQL    StorageType = "MYSQL"
)
