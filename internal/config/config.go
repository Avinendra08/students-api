package config

import (
	"flag"
	"log"
	"os"

	"github.com/ilyakaznacheev/cleanenv"
)

//we are starting things with capital letters, bcs it makes them public for export and import outside package
//below the things are representing local.yaml in root config folder,
//basically we are serializing the env thing to be used globally
//we are using cleanenv package for serializing
//we are using struct tags ``

type HTTPServer struct {
	Addr string
}
//env-default:"production"
type Config struct {
	Env         string `yaml:"env" env:"Env" env-required:"true"`
	StoragePath string `yaml:"storage_path" env-required:"true"`
	HTTPServer `yaml:"http_server"`
}

//now we come to logic part, how to parse
func MustLoad() *Config{
	//so sabse phle we check for the config path
	var configPath string 
	
	configPath = os.Getenv("CONFIG_PATH") //we can take from env if available
	if configPath == ""{
		//if not in env, we can check in cli flags
		flags := flag.String("config","","path to the configuration file")
		flag.Parse()

		configPath = *flags //dereferencing

		if configPath==""{
			log.Fatal("config path is not set")
		}
	}
	
	//now we check if file is avaialble on the path given or not
	if _,err := os.Stat(configPath); os.IsNotExist(err){
		log.Fatalf("config file does not exist: %s", configPath)
	}

	//if everything is going well till now, we move ahead for serializing logic
	var cfg Config

	err := cleanenv.ReadConfig(configPath, &cfg)
	if err!=nil{
		log.Fatalf("cannot read config file: %s",err.Error())
	}

	return &cfg
}