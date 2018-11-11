package config

import (
	"bytes"
	"encoding/json"
	"io"
	"log"
	"os"
	"os/user"
	"path/filepath"
)

var (
	configDir  string
	configFile string

	//configDescriptor os.File
)

var cfg config
var exist bool

type config struct {
	Projects map[string]Project `json:"projects"`
}

// Project full information
type Project struct {
	TrackerURL         string `json:"tracker_url"`
	TrackerAccessToken string `json:"tracker_access_token"`
	TrackerType        string `json:"tracker_type"`
	DefaultProject     string `json:"default_project"`
	LocalPath          string `json:"local_path"`
}

// paths: set and return config directory, config file path
func paths() (dir, file string) {
	u, err := user.Current()
	if err != nil {
		return "", ""
	}
	configDir = u.HomeDir + string(filepath.Separator) + ".toad"
	configFile = configDir + string(filepath.Separator) + ".config"

	return configDir, configFile
}

// create config file if not exist
func create() error {
	log.Println("config.create")
	dir, file := paths()

	if err := os.Mkdir(dir, os.ModePerm); err != nil {
		if !os.IsExist(err) {
			return err
		}
	}
	f, err := os.Create(file)
	if err != nil {
		log.Printf("could not create config file in user home directory %v", err)
		return err
	}
	defer f.Close()

	exist = true
	return nil
}

// if exist config file
func exists() (bool, error) {
	log.Println("config.exists")
	if exist {
		return true, nil
	}

	if _, err := os.Stat(configFile); os.IsNotExist(err) {
		return false, nil
	} else {
		return false, err
	}

	exist = true
	return exist, nil
}

func createIfNotExist() error {
	log.Println("createIfNotExist")
	exist, err := exists()
	if err != nil {
		return err
	}
	if exist {
		return nil
	}
	if err := create(); err != nil {
		return err
	}
	return nil
}

// Read projects configuration
func Read() (map[string]Project, error) {
	log.Println("Read")
	if err := createIfNotExist(); err != nil {
		return nil, err
	}

	//todo if file just create - return nil, nil

	f, err := os.OpenFile(configFile, os.O_RDONLY, os.ModeExclusive)
	if err != nil {
		log.Println("Read.OpenFile ", err)
		return nil, err
	}
	defer f.Close()

	dec := json.NewDecoder(f)
	if err := dec.Decode(&cfg); err != nil {
		log.Println("Read.Decode ", err)
		if err != io.EOF {
			return nil, err
		}
	}

	return cfg.Projects, nil
}

// Append project to config
func Append(project Project) error {
	log.Println("Append")
	cfg.Projects[project.LocalPath] = project
	return nil
}

// Save config file //todo don't forget to do it on app exit
func Save() error {
	log.Println("Save")
	if err := createIfNotExist(); err != nil {
		return err
	}
	f, err := os.OpenFile(configFile, os.O_WRONLY, os.ModeExclusive)
	if err != nil {
		log.Println("config.Save.OpenFile", err)
		return err
	}
	defer f.Close()

	buffer := bytes.NewBuffer([]byte{})
	enc := json.NewEncoder(buffer)
	if err := enc.Encode(cfg); err != nil {
		log.Println("config.Save.Encode", err)
		return err
	}
	log.Println(buffer.String())
	_, err = f.Write(buffer.Bytes())
	if err != io.EOF {
		log.Println(err)
		return err
	}

	return nil
}
