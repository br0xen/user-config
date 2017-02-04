// Package userConfig eases the use of config files in a user's home directory
package userConfig

import (
	"errors"
	"os"
	"strings"
	"time"

	"github.com/br0xen/user-config/ext/go-xdg"
)

// Config is a stuct for managing the config
type Config struct {
	name          string
	generalConfig *GeneralConfig
}

// NewConfig generates a Config struct
func NewConfig(name string) (*Config, error) {
	c := &Config{name: name}
	if err := c.Load(); err != nil {
		return c, err
	}
	return c, nil
}

// GetKeyList at the config level returns all keys in the <c.name>.conf file
func (c *Config) GetKeyList() []string {
	return c.generalConfig.GetKeyList()
}

// Set at the config level sets a value in the <c.name>.conf file
func (c *Config) Set(k, v string) error {
	return c.generalConfig.Set(k, v)
}

// SetBytes at the config level sets a value in the <c.name>.conf file
func (c *Config) SetBytes(k string, v []byte) error {
	return c.generalConfig.SetBytes(k, v)
}

// SetInt saves an integer (as a string) in the <c.name>.conf file
func (c *Config) SetInt(k string, v int) error {
	return c.generalConfig.SetInt(k, v)
}

// SetDateTime saves a time.Time (as a string) in the <c.name>.conf file
func (c *Config) SetDateTime(k string, v time.Time) error {
	return c.generalConfig.SetDateTime(k, v)
}

// Get at the config level retrieves a value from the <c.name>.conf file
func (c *Config) Get(k string) string {
	return c.generalConfig.Get(k)
}

// GetBytes at the config level retrieves a value from the <c.name>.conf file
// and returns it as a byte slice
func (c *Config) GetBytes(k string) []byte {
	return c.generalConfig.GetBytes(k)
}

// GetInt at the config level retrieves a value from the <c.name>.conf file
// and returns it as an integer (or an error if conversion fails)
func (c *Config) GetInt(k string) (int, error) {
	return c.generalConfig.GetInt(k)
}

// GetDateTime at the config level retrieves a value from the <c.name>.conf file
func (c *Config) GetDateTime(k string) (time.Time, error) {
	return c.generalConfig.GetDateTime(k)
}

// DeleteKey at the config level removes a key from the <c.name>.conf file
func (c *Config) DeleteKey(k string) {
	c.generalConfig.DeleteKey(k)
}

// GetConfigPath just returns the config path
func (c *Config) GetConfigPath() string {
	return c.generalConfig.Path
}

// Load loads config files into the config
func (c *Config) Load() error {
	var err error
	if strings.TrimSpace(c.name) == "" {
		return errors.New("Invalid Config Name: " + c.name)
	}

	var cfgPath string
	cfgPath = xdg.Config.Dirs()[0]
	if cfgPath != "" {
		cfgPath = cfgPath + string(os.PathSeparator) + c.name
		if err = c.verifyOrCreateDirectory(cfgPath); err != nil {
			return err
		}
		// We always have a <name>.conf file
		//cfgPath = cfgPath + "/" + c.name + ".conf"
	}
	// Load general config
	if c.generalConfig, err = NewGeneralConfig(c.name, cfgPath); err != nil {
		return err
	}

	return nil
}

// Save writes the config to file(s)
func (c *Config) Save() error {
	if c.generalConfig == nil {
		return errors.New("Bad setup.")
	}
	return c.generalConfig.Save()
}

// verifyOrCreateDirectory is a helper function for building an
// individual directory
func (c *Config) verifyOrCreateDirectory(path string) error {
	var tstDir *os.File
	var tstDirInfo os.FileInfo
	var err error
	if tstDir, err = os.Open(path); err != nil {
		if err = os.Mkdir(path, 0755); err != nil {
			return err
		}
		if tstDir, err = os.Open(path); err != nil {
			return err
		}
	}
	if tstDirInfo, err = tstDir.Stat(); err != nil {
		return err
	}
	if !tstDirInfo.IsDir() {
		return errors.New(path + " exists and is not a directory")
	}
	// We were able to open the path and it was a directory
	return nil
}
