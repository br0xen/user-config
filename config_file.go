// Package userConfig eases the use of config files in a user's home directory
package userConfig

// ConfigFile is the interface for all config files
type ConfigFile interface {
	SetName(string)
	GetName() string
	SetPath(string)
	GetPath() string
	Load() error
	Save() error
	Set(string, string)
	Get(string) string
	GetFullPath() string
}