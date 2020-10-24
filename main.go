package main

import (
	"github.com/lestrrat-go/jwx/jwa"
	log "github.com/sirupsen/logrus"
	"github.com/urfave/cli/v2"
	"gopkg.in/natefinch/lumberjack.v2"
	"gopkg.in/yaml.v3"
	"net/url"
	"os"
	"strings"
	"time"
)

type AuthServerConfig struct {
	ListenAddr  string `yaml:"listenAddr"`
	TokenReqPath string `yaml:"tokenReqPath,omitempty"`
	Http2       bool   `yaml:"http2"`
	TlsCertFile string `yaml:"tlsCertFile,omitempty"`
	TlsKeyFile  string `yaml:"tlsKeyFile,omitempty"`
	InstanceId  string `yaml:"instanceId"`
	TokenExpire int64  `yaml:"tokenExpire"`
	Signature   struct {
		Algorithm string
		KeyFile   string `yaml:"keyFile"`
	}
}

func loadYamlConfig(fileName string, intf interface{}) error {
	f, err := os.Open(fileName)
	if err != nil {
		return err
	}

	defer f.Close()

	decoder := yaml.NewDecoder(f)
	err = decoder.Decode(intf)

	if err != nil {
		return err
	}

	return nil
}

func toYaml(intf interface{}) (string, error) {
	b, err := yaml.Marshal(intf)
	return string(b), err
}

func isQuoteDisabled() bool {
	disableQuoteEnv := strings.ToLower(os.Getenv("DISABLE_QUOTE"))

	return disableQuoteEnv == "true" || disableQuoteEnv == "t" || disableQuoteEnv == "yes" || disableQuoteEnv == "y"
}

func initLog(logFile string, strLevel string, logSize int, backups int) {
	level, err := log.ParseLevel(strLevel)
	if err != nil {
		level = log.InfoLevel
	}
	log.SetLevel(level)
	if len(logFile) <= 0 {
		log.SetOutput(os.Stdout)
	} else {
		log.SetFormatter(&log.TextFormatter{DisableColors: true,
			FullTimestamp: true,
			ForceQuote:    false,
			DisableQuote:  isQuoteDisabled()})
		log.SetOutput(&lumberjack.Logger{Filename: logFile,
			LocalTime:  true,
			MaxSize:    logSize,
			MaxBackups: backups})
	}
}

func loadAuthServerConfig(fileName string) (*AuthServerConfig, error) {
	r := &AuthServerConfig{}
	err := loadYamlConfig(fileName, r)
	if err == nil {
		return r, err
	} else {
		return nil, err
	}

}
func startAuthServer(c *cli.Context) error {
	config, err := loadAuthServerConfig(c.String("config"))
	if err != nil {
		return err
	}
	strLevel := c.String("log-level")
	fileName := c.String("log-file")
	logSize := c.Int("log-size")
	backups := c.Int("log-backups")
	initLog(fileName, strLevel, logSize, backups)
	b, _ := toYamlBytes(config)
	log.Info("Load configuration:", string(b))
	alg := jwa.SignatureAlgorithm(config.Signature.Algorithm)
	key, err := loadSignatureKeyFromFile(config.Signature.KeyFile)
	if err != nil {
		return err
	}
	return NewOAuthServer(config.TokenReqPath,
		config.InstanceId,
		time.Duration(config.TokenExpire)*time.Second,
		config.Http2,
		config.TlsCertFile,
		config.TlsKeyFile,
		alg,
		key).Start(config.ListenAddr)
}

type AuthProxyConfig struct {
	Proxies []struct {
		ListenAddr string `yaml:"listenAddr"`
		AuthServer struct {
			Fqdn       string `yaml:"fqdn,omitempty"`
			Http2      bool   `yaml:"http2"`
			Url        string `yaml:"url"`
			CaCertFile string `yaml:"caCertFile,omitempty"`
			CertFile   string `yaml:"certFile,omitempty"`
			KeyFile    string `yaml:"keyFile,omitempty"`
		} `yaml:"authServer"`
		TokenReqPath         string `yaml:"tokenReqPath"`
		TokenVerifyPath      string `yaml:"tokenVerifyPath"`
		TokenVerifyAlgorithm string `yaml:"tokenVerifyAlgorithm,omitempty"`
		TokenVerifyKeyFile   string `yaml:"tokenVerifyKeyFile,omitempty"`
	}
}

func loadAuthProxyConfig(fileName string) (*AuthProxyConfig, error) {
	f, err := os.Open(fileName)
	if err != nil {
		return nil, err
	}

	defer f.Close()
	r := &AuthProxyConfig{}

	decoder := yaml.NewDecoder(f)
	err = decoder.Decode(r)

	if err != nil {
		return nil, err
	}

	return r, nil
}
func startAuthProxy(c *cli.Context) error {
	authProxyConfig, err := loadAuthProxyConfig(c.String("config"))
	if err != nil {
		return err
	}
	strLevel := c.String("log-level")
	fileName := c.String("log-file")
	logSize := c.Int("log-size")
	backups := c.Int("log-backups")
	initLog(fileName, strLevel, logSize, backups)

	for _, item := range authProxyConfig.Proxies {
		key, err := loadSignatureKeyFromFile(item.TokenVerifyKeyFile)
		if err != nil {
			log.Error("Fail to load key file ", item.TokenVerifyKeyFile, " with error:", err)
			return err
		}

		tlsConfig, err := loadCertFile(item.AuthServer.CaCertFile, item.AuthServer.CertFile, item.AuthServer.KeyFile)
		if err != nil {
			log.Error("Fail to load the certificate file ", item.AuthServer.CaCertFile)
			return err
		}
		serverFqdn := item.AuthServer.Fqdn
		if len(serverFqdn) <= 0 {
			u, err := url.Parse(item.AuthServer.Url)
			if err == nil {
				serverFqdn = u.Host
			}
		}
		if len(serverFqdn) > 0 && tlsConfig != nil {
			tlsConfig.ServerName = serverFqdn
			log.Info("tlsConfig.ServerName is ", tlsConfig.ServerName)
		}
		tlsConfig.BuildNameToCertificate()
		go func() {
			NewProxy(item.TokenReqPath,
				item.TokenVerifyPath,
				item.AuthServer.Url,
				tlsConfig,
				item.AuthServer.Http2,
				jwa.SignatureAlgorithm(item.TokenVerifyAlgorithm),
				key).Listen(item.ListenAddr)
		}()
	}

	for {
		time.Sleep(time.Duration(5 * time.Second))
	}
	return nil
}

func main() {
	serverCommand := &cli.Command{
		Name:  "server",
		Usage: "oauth server",
		Flags: []cli.Flag{
			&cli.StringFlag{
				Name:     "config",
				Aliases:  []string{"c"},
				Required: true,
				Usage:    "Load configuration from `FILE`",
			},
			&cli.StringFlag{
				Name:  "log-file",
				Usage: "log file name",
			},
			&cli.StringFlag{
				Name:  "log-level",
				Usage: "one of following level: Trace, Debug, Info, Warn, Error, Fatal, Panic",
			},
			&cli.IntFlag{
				Name:  "log-size",
				Usage: "size of log file in Megabytes",
				Value: 50,
			},
			&cli.IntFlag{
				Name:  "log-backups",
				Usage: "number of log rotate files",
				Value: 10,
			},
		},
		Action: startAuthServer,
	}
	proxyCommand := &cli.Command{
		Name:  "proxy",
		Usage: "oauth proxy",
		Flags: []cli.Flag{
			&cli.StringFlag{
				Name:     "config",
				Aliases:  []string{"c"},
				Required: true,
				Usage:    "Load configuration from `FILE`",
			},
			&cli.StringFlag{
				Name:  "log-file",
				Usage: "log file name",
			},
			&cli.StringFlag{
				Name:  "log-level",
				Usage: "one of following level: Trace, Debug, Info, Warn, Error, Fatal, Panic",
			},
			&cli.IntFlag{
				Name:  "log-size",
				Usage: "size of log file in Megabytes",
				Value: 50,
			},
			&cli.IntFlag{
				Name:  "log-backups",
				Usage: "number of log rotate files",
				Value: 10,
			},
		},
		Action: startAuthProxy,
	}
	app := &cli.App{
		Name:     "rest-oauth-proxy",
		Usage:    "oauth-proxy with rest interface",
		Commands: []*cli.Command{serverCommand, proxyCommand},
	}
	err := app.Run(os.Args)
	if err != nil {
		log.WithFields(log.Fields{"error": err}).Error("Fail to start application")
	}
}
