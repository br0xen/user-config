// Package userConfig eases the use of config files in a user's home directory
package userConfig

import (
	"bytes"
	"errors"
	"io/ioutil"
	"strings"

	"github.com/BurntSushi/toml"
)

// GeneralConfig is the basic config structure
// All configs make with package userConfig will have this file
type GeneralConfig struct {
	Name        string            `toml:"-"`
	Path        string            `toml:"-"`
	ConfigFiles []string          `toml:"additional_config"`
	RawFiles    []string          `toml:"raw_files"`
	Values      map[string]string `toml:"general"`
}

// NewGeneralConfig generates a General Config struct
func NewGeneralConfig(name, path string) (*GeneralConfig, error) {
	gf := &GeneralConfig{Name: name, Path: path}
	gf.ConfigFiles = []string{}
	gf.RawFiles = []string{}
	gf.Values = make(map[string]string)

	if err := gf.Load(); err != nil {
		return gf, err
	}
	return gf, nil
}

// Load loads config files into the config
func (gf *GeneralConfig) Load() error {
	if strings.TrimSpace(gf.Name) == "" || strings.TrimSpace(gf.Path) == "" {
		return errors.New("Invalid ConfigFile Name: " + gf.Path + "/" + gf.Name)
	}

	// Config files end with .conf
	cfgPath := gf.Path + "/" + gf.Name + ".conf"
	tomlData, err := ioutil.ReadFile(cfgPath)
	if err != nil {
		return err
	}
	if _, err := toml.Decode(string(tomlData), &gf); err != nil {
		return err
	}
	return nil
}

// Save writes the config to file(s)
func (gf *GeneralConfig) Save() error {
	buf := new(bytes.Buffer)
	cfgPath := gf.Path + "/" + gf.Name + ".conf"
	if err := toml.NewEncoder(buf).Encode(gf); err != nil {
		return err
	}
	return ioutil.WriteFile(cfgPath, buf.Bytes(), 0644)
}

// Set sets a key/value pair in gf, if unable to save, revert to old value
// (and return the error)
func (gf *GeneralConfig) Set(k, v string) error {
	oldVal := gf.Values[k]
	gf.Values[k] = v
	if err := gf.Save(); err != nil {
		gf.Values[k] = oldVal
		return err
	}
	return nil
}

// Get gets a key/value pair from gf
func (gf *GeneralConfig) Get(k string) string {
	return gf.Values[k]
}