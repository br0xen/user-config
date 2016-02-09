// Package userConfig eases the use of config files in a user's home directory
package userConfig

import (
	"errors"
	"os"
	"strings"

	"gogs.bullercodeworks.com/brian/user-config/ext/go-xdg"
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

// Set at the config level sets a value in the <c.name>.conf file
func (c *Config) Set(k, v string) error {
	return c.generalConfig.Set(k, v)
}

// Get at the config level retrieves a value from the <c.name>.conf file
func (c *Config) Get(k string) string {
	return c.generalConfig.Get(k)
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
		cfgPath = cfgPath + "/" + c.name
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
	/*
		var cfgPath string
		var configLines []string
		//configLines = append(configLines, "server="+client.ServerAddr)
		//configLines = append(configLines, "key="+client.ServerKey)
		cfgPath = os.Getenv("HOME")
		if cfgPath != "" {
			cfgPath = cfgPath + "/.config"
			if err := c.verifyOrCreateDirectory(cfgPath); err != nil {
				return err
			}
			cfgPath = cfgPath + "/" + c.name
		}
		if cfgPath != "" {
			file, err := os.Create(cfgPath)
			if err != nil {
				// Couldn't load config even though one was specified
				return err
			}
			defer file.Close()

			w := bufio.NewWriter(file)
			for _, line := range configLines {
				fmt.Fprintln(w, line)
			}
			if err = w.Flush(); err != nil {
				return err
			}
		}
		return nil
	*/
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
