package config

import (
	"flag"
	"log"
	"os"

	"github.com/ilyakaznacheev/cleanenv"
)

//we are starting things with capital letters, bcs it makes them public for export and import outside package
//below the things are representing local.yaml in root config folder,
//basically we are serializing the env thing to be used globally, 
//it's like we used process.env in all files in node js, here we are doing at one place to be used globally
//cleanenv maps values from YAML config file + environment variables + struct tags into your Go struct in one go. 
//It’s not just “like process.env”, it also validates (env-required) and provides defaults (env-default) if defined.
//we are using struct tags ``

type HTTPServer struct {
	Addr string `yaml:"address" env-required:"true"`
}
//env-default:"production"
type Config struct {
	Env         string `yaml:"env" env:"Env" env-required:"true"`
	StoragePath string `yaml:"storage_path" env-required:"true"`
	HTTPServer `yaml:"http_server"`
}

//now we come to logic part, how to parse
//The function panics (log.Fatal) if something goes wrong instead of returning an error.
//➝ That’s why the name is MustLoad (convention in Go for "load or die").
//Useful for config, since without config the app cannot run.
func MustLoad() *Config{
	//so sabse phle we check for the config path
	var configPath string 
	
	configPath = os.Getenv("CONFIG_PATH") //we can take from env if available
	if configPath == ""{
		//if not in env, we can check in cli flags i.e --config
		//like in our case we dont have env, we have everything in local.yaml
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


// Config represents the application configuration loaded from YAML and env vars.
// Fields are exported (capitalized) so they can be accessed outside the package.
//
// cleanenv is used to populate this struct from:
//   1. A YAML file (path provided via CONFIG_PATH env or --config flag).
//   2. Environment variables (with validation, defaults, etc).
// This gives us a centralized, type-safe config (like process.env in Node.js, but stricter).
