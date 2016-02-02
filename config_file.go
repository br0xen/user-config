// Package userConfig eases the use of config files in a user's home directory
package userConfig

import (
	"bufio"
	"bytes"
	"errors"
	"fmt"
	"io/ioutil"
	"os"
	"strings"

	"github.com/BurntSushi/toml"
)

type generalConfig struct {
	ConfigFiles []string          `toml:"additional_config"`
	RawFiles    []string          `toml:"raw_files"`
	Values      map[string]string `toml:"general"`
}

// ConfigFile is a file that has key/value pairs in the config
type ConfigFile struct {
	Name string
	Path string
	data interface{}
}

// NewConfigFile generates a Config struct
func NewConfigFile(name, path string, dt interface{}) (*ConfigFile, error) {
	cf := &ConfigFile{Name: name, Path: path, data: dt}
	if err := cf.Load(); err != nil {
		return nil, err
	}
	return cf, nil
}

// Set sets a key/value pair in cf, if unable to save, revert to old value
// (and return the error)
func (cf *ConfigFile) Set(k, v string) error {
	oldVal := cf.values[k]
	cf.values[k] = v
	if err := cf.Save(); err != nil {
		cf.values[k] = oldVal
	}
	return err
}

// Get gets a key/value pair from cf
func (cf *ConfigFile) Get(k string) string {
	return cf.values[k]
}

// Load loads config files into the config
func (cf *ConfigFile) Load() error {
	if strings.TrimSpace(cf.Name) == "" && strings.TrimSpace(cf.Path) {
		return errors.New("Invalid ConfigFile Name: " + cf.Path + "/" + cf.Name)
	}

	// Config files end with .conf
	cfgPath := cf.Path + "/" + cf.Name + ".conf"
	tomlData, err := ioutil.ReadFile(cfgPath)
	if err != nil {
		return err
	}
	if _, err := toml.Decode(tomlData, &cf.data); err != nil {
		return err
	}
	return nil
}

// Save writes the config to file(s)
func (cf *ConfigFile) Save() error {
	if strings.TrimSpace(cf.Name) == "" && strings.TrimSpace(cf.Path) {
		return errors.New("Invalid ConfigFile Name: " + cf.Path + "/" + cf.Name + ".conf")
	}

	filePath := cf.path + "/" + cf.Name + ".conf"
	var err error
	buf := new(bytes.Buffer)
	enc := toml.NewEncoder(buf).Encode(data)
	err = ioutil.WriteFile(filePath, buf, 0644)
	if err != nil {
		return err
	}

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
}