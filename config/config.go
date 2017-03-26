package config

import (
	"encoding/json"
	"fmt"
	"github.com/kabukky/httpscerts"
	"io/ioutil"
	"os"
)

const (
	DirNetworkConfig = "network_config.json"
)

const (
	DefaultDomain      = "recipemanager.io"
	DefaultPort        = "443"
	DefaultMongo       = "mongodb://127.0.0.1:32017"
	DefaultSSLCertPath = "/etc/letsencrypt/live/recipemanager.io/fullchain.pem"
	DefaultSSLKeyPath  = "/etc/letsencrypt/live/recipemanager.io/privkey.pem"
	TempSSLCertPath    = "temp_cert.pem"
	TempSSLKeyPath     = "temp_key.pem"
)

type NetworkConfig struct {
	Domain      string `json:"domain"`
	Port        string `json:"port"`
	Mongo       string `json:"mongo"`
	SSLCertPath string `json:"ssl_cert_path"`
	SSLKeyPath  string `json:"ssl_key_path"`
}

func GetNetworkConfig() (*NetworkConfig, error) {
	// Get path of network config file.
	networkConfigRoot := getRoot() + DirNetworkConfig

	// Create network config file with default values if not exist.
	if _, e := ioutil.ReadFile(networkConfigRoot); e != nil {
		c := NetworkConfig{
			Domain:      DefaultDomain,
			Port:        DefaultPort,
			Mongo:       DefaultMongo,
			SSLCertPath: DefaultSSLCertPath,
			SSLKeyPath:  DefaultSSLKeyPath,
		}
		data, _ := json.MarshalIndent(c, "", "\t")
		if e := ioutil.WriteFile(networkConfigRoot, data, os.ModePerm); e != nil {
			return nil, e
		}
	}

	// Read network config file and load values to struct.
	data, e := ioutil.ReadFile(networkConfigRoot)
	if e != nil {
		return nil, e
	}
	c := &NetworkConfig{}
	e = json.Unmarshal(data, c)
	if e != nil {
		return nil, e
	}

	// Make random certificates if specified don't exist.
	if httpscerts.Check(c.SSLCertPath, c.SSLKeyPath) != nil {
		fmt.Printf(`One of either "%s" or "%s" does not exist.\n`, c.SSLCertPath, c.SSLKeyPath)
		fmt.Printf(`Generating temporary certificates for "%s"...\n`, c.Domain+":"+c.Port)
		e = httpscerts.Generate(TempSSLCertPath, TempSSLKeyPath, c.Domain+":"+c.Port)
		if e != nil {
			return nil, e
		}
		// Link struct values to temporary certificates.
		c.SSLCertPath = TempSSLCertPath
		c.SSLKeyPath = TempSSLKeyPath
	} else {
		data, _ := ioutil.ReadFile(c.SSLKeyPath)
		fmt.Println(string(data))
		data, _ = ioutil.ReadFile(c.SSLCertPath)
		fmt.Println(string(data))
	}
	return c, nil
}

func getRoot() string {
	path := "/home/ubuntu/recipe-manager/"
	if _, e := os.Stat("/home/ubuntu/"); os.IsNotExist(e) {
		path = "/home/ubuntu/" + os.Getenv("USER") + "/recipe-manager/"
	}
	fmt.Println("Configuration file is located at:", path)
	os.MkdirAll(path, os.ModePerm)
	return path
}
