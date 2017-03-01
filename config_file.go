// Package userConfig eases the use of config files in a user's home directory
package userConfig

import (
	"bytes"
	"encoding/json"
	"errors"
	"io/ioutil"
	"strconv"
	"strings"
	"time"

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

// GetKeyList returns a list of all keys in the config file
func (gf *GeneralConfig) GetKeyList() []string {
	var ret []string
	for k, _ := range gf.Values {
		ret = append(ret, k)
	}
	return ret
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

// SetBytes at the config level sets a value in the <c.name>.conf file
func (gf *GeneralConfig) SetBytes(k string, v []byte) error {
	return gf.Set(k, string(v))
}

// SetInt sets an integer value (as a string) in the config file
func (gf *GeneralConfig) SetInt(k string, v int) error {
	return gf.Set(k, strconv.Itoa(v))
}

// SetDateTime sets a DateTime value (as a string) in the config file
func (gf *GeneralConfig) SetDateTime(k string, v time.Time) error {
	return gf.Set(k, v.Format(time.RFC3339))
}

// SetArray sets a string slice value (as a string) in the config file
func (gf *GeneralConfig) SetArray(k string, v []string) error {
	b, e := json.Marshal(v)
	if e != nil {
		return e
	}
	return gf.SetBytes(k, b)
}

// Get gets a key/value pair from gf
func (gf *GeneralConfig) Get(k string) string {
	return gf.Values[k]
}

// GetInt gets a key/value pair from gf and return it as an integer
// An error if it can't be converted
func (gf *GeneralConfig) GetInt(k string) (int, error) {
	return strconv.Atoi(gf.Get(k))
}

// GetDateTime gets a key/value pair from gf and returns it as a time.Time
// An error if it can't be converted
func (gf *GeneralConfig) GetDateTime(k string) (time.Time, error) {
	return time.Parse(time.RFC3339, gf.Get(k))
}

// GetBytes gets a key/value pair from gf and returns it as a byte slice
// Or an error if it fails for whatever reason
func (gf *GeneralConfig) GetBytes(k string) []byte {
	return []byte(gf.Get(k))
}

func (gf *GeneralConfig) GetArray(k string) ([]string, error) {
	var ret []string
	err := json.Unmarshal(gf.GetBytes(k), &ret)
	return ret, err
}

// DeleteKey removes a key from the file
func (gf *GeneralConfig) DeleteKey(k string) error {
	oldVal := gf.Get(k)
	delete(gf.Values, k)
	if err := gf.Save(); err != nil {
		gf.Values[k] = oldVal
		return err
	}
	return nil
}
