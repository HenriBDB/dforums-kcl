package configuration

import (
	"fmt"
	"testing"

	"github.com/spf13/viper"
)

func TestLoadingConfig(*testing.T) {
	fmt.Println("All settings:", GetJsonConfigs())
	for _, v := range GetNetworkSeeds() {
		fmt.Println(v)
	}
	// UpdateConfig("newPath", 101, 201, 21)
	UpdateConfig("database", 100, 200, 20)
	fmt.Println("All settings:", GetJsonConfigs())
}

func TestSendMessage(*testing.T) {
	viper.SetConfigName("config")
	viper.SetConfigType("yaml")
	viper.AddConfigPath("../")
	err := viper.ReadInConfig() // Find and read the config file
	if err != nil {             // Handle errors reading the config file
		panic(fmt.Errorf("Fatal error config file: %w \n", err))
	}
	fmt.Printf("config data: \n %#v\n", viper.GetInt("network.low"))
	viper.Set("network.low", 50)
	fmt.Printf("config data: \n %#v\n", viper.GetInt("network.low"))
	viper.WriteConfig()

}
