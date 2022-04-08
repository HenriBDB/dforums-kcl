package configuration

import (
	"fmt"

	"github.com/spf13/viper"
)

const (
	configName        = "dfd-config"
	networkConnMinKey = "network.connection-low"
	networkConnMaxKey = "network.connection-high"
	networkSeedsKey   = "network.seeds"
	networkPortKey    = "network.port"
	networkPeersKey   = "network.peers"
	dbPathKey         = "database.storage-path"
	powLevelKey       = "security.proofofwork-level"
)

var defaults = map[string]interface{}{
	networkConnMinKey: 100,
	networkConnMaxKey: 200,
	networkPortKey:    6870,
	networkSeedsKey:   []string{},
	networkPeersKey:   []string{},
	dbPathKey:         "database/",
	powLevelKey:       "24",
}

func InitConfigs(configPath string) {
	viper.SetConfigName(configName)
	viper.SetConfigType("yaml")
	viper.AddConfigPath(configPath)
	if err := viper.ReadInConfig(); err != nil { // Handle errors reading the config file
		if _, ok := err.(viper.ConfigFileNotFoundError); ok {
			// Config file not found; Create new
			Logger.Info("creating new config file...")
			viper.SafeWriteConfig()
		} else {
			// Config file was found but another error was produced
			Logger.Errorf("fatal error config file: %w", err)
			panic(fmt.Errorf("fatal error config file: %w", err))
		}
	}
	validateConfig()
}

func UpdateConfig(dbPath string, minConn int, maxConn int, powLvl int) {
	if dbPath != "" && dbPath[len(dbPath)-1:][0] != byte('/') {
		dbPath = dbPath + "/"
	}
	viper.Set(dbPathKey, dbPath)
	if powLvl >= 0 && powLvl <= 3 {
		viper.Set(powLevelKey, 16+(powLvl*4))
	}
	viper.Set(networkConnMinKey, minConn)
	viper.Set(networkConnMaxKey, maxConn)
	viperSave()
}

func GetConnectionLimits() (int, int) {
	return viper.GetInt(networkConnMinKey),
		viper.GetInt(networkConnMaxKey)
}

func GetMinNodeDifficulty() int {
	return viper.GetInt(powLevelKey)
}

func GetDatabasePath() string {
	return viper.GetString(dbPathKey)
}

func GetNetworkSeeds() []string {
	return viper.GetStringSlice(networkSeedsKey)
}

func GetNetworkPeers() []string {
	return viper.GetStringSlice(networkPeersKey)
}

func SetNetworkPeers(peers []string) {
	viper.Set(networkPeersKey, peers)
	viperSave()
}

func GetNetworkPort() int {
	return viper.GetInt(networkPortKey)
}

func GetJsonConfigs() map[string]interface{} {
	return viper.AllSettings()
}

func validateConfig() {
	for k, v := range defaults {
		if viper.Get(k) == nil {
			viper.Set(k, v)
		}
	}
	// Write any missing values with default
	viperSave()
}

func viperSave() {
	if err := viper.WriteConfig(); err != nil {
		Logger.Error("could not save configs:", err.Error())
	}
}
