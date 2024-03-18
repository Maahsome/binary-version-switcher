package cmd

import (
	"fmt"
	"io"
	"os"
	"path/filepath"
	"plugin"
	"regexp"
	"strings"

	"binary-version-switcher/ask"
	"binary-version-switcher/cmd/plugins"
	"binary-version-switcher/config"

	glog "github.com/maahsome/golang-logger"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	// "gopkg.in/yaml.v3"
	// "sigs.k8s.io/yaml"
)

type (
	Project struct {
		ID     int    `json:"id"`
		Name   string `json:"name"`
		Path   string `json:"path"`
		SSHURL string `json:"ssh_url_to_repo"`
	}
)

var (
	cfgFile   string
	semVer    string
	gitCommit string
	gitRef    string
	buildDate string

	// semVerReg - gets the semVer portion only, cutting off any other release details
	semVerReg = regexp.MustCompile(`(v[0-9]+\.[0-9]+\.[0-9]+).*`)

	c = &config.Config{}
)

// RootCmd represents the base command when called without any subcommands
var RootCmd = &cobra.Command{
	Use:   "binary-version-switcher",
	Short: "",
	Long: `EXAMPLE:

  TODO: add description

  > binary-version-switcher

`,
	PersistentPreRun: func(cmd *cobra.Command, args []string) {

		c.VersionDetail.SemVer = semVer
		c.VersionDetail.BuildDate = buildDate
		c.VersionDetail.GitCommit = gitCommit
		c.VersionDetail.GitRef = gitRef
		c.VersionJSON = fmt.Sprintf("{\"SemVer\": \"%s\", \"BuildDate\": \"%s\", \"GitCommit\": \"%s\", \"GitRef\": \"%s\"}", semVer, buildDate, gitCommit, gitRef)
		if c.OutputFormat != "" {
			c.FormatOverridden = true
			c.NoHeaders = false
			c.OutputFormat = strings.ToLower(c.OutputFormat)
			switch c.OutputFormat {
			case "json", "gron", "yaml", "text", "table", "raw":
				break
			default:
				fmt.Println("Valid options for -o are [json|gron|text|table|yaml|raw]")
				os.Exit(1)
			}
		}

		// if os.Args[1] != "version" {
		// }
	},
}

func GetCmd(pluginPath, cmdName string) *cobra.Command {
	p, err := plugin.Open(pluginPath)
	if err != nil {
		panic(err)
	}
	b, err := p.Lookup(cmdName)
	if err != nil {
		panic(err)
	}

	// sld, err := p.Lookup("SymLinkDir")
	// if err != nil {
	// 	panic(err)
	// }
	// *sld.(*string) = "/usr/local/bin/"

	// bnd, err := p.Lookup("BinDir")
	// if err != nil {
	// 	panic(err)
	// }
	// *bnd.(*string) = "/usr/local/bvs/"

	f, err := p.Lookup("Init" + cmdName)
	if err == nil {
		f.(func(string, string, string))(c.SymLinkPath, c.BinPath, c.LogLevel)
	}
	return *b.(**cobra.Command)
}
func buildRootCmd() *cobra.Command {
	RootCmd.PersistentFlags().StringVar(&cfgFile, "config", "", "config file (default is $HOME/.splicectl/config.yml)")
	RootCmd.PersistentFlags().StringVarP(&c.OutputFormat, "output", "o", "", "output types: json, text, yaml, gron, raw")
	RootCmd.PersistentFlags().BoolVar(&c.NoHeaders, "no-headers", false, "Suppress header output in Text output")

	return RootCmd
}

func addSubCommands() {
	RootCmd.AddCommand(
		// from 'import binary-version-switcher/cmd/<subcommand:package>'
		// <package>.InitSubCommands(c),
		plugins.InitSubCommands(c),
	)

	plugIns, perr := os.ReadDir(c.PluginPath)
	if perr != nil {
		c.Logger.WithError(perr).Error("failed to get a list of plugins")
		return
	}
	if len(plugIns) == 0 {
		c.Logger.Warn("No application plugins detected, please run 'binary-version-switcher plugins get'")
	}
	for _, p := range plugIns {
		piPath := filepath.Join(c.PluginPath, p.Name())
		RootCmd.AddCommand(GetCmd(piPath, "MainCmd"))
	}
}

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
	if err := RootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

func setupLogging() {
	// log level enablement inclusions from left to right
	// trace, debug, info, warning, error, fatal, panic

	logLevel := "warning"
	if len(os.Getenv("BVS_LOG_LEVEL")) > 0 {
		logLevel = os.Getenv("BVS_LOG_LEVEL")
	}
	c.LogLevel = logLevel
	c.Logger = glog.CreateStandardLogger()
	c.Logger.Level = glog.LogLevelFromString(logLevel)

}

func init() {
	setupLogging()
	buildRootCmd()
	initConfig()
	cobra.OnInitialize(initConfig)
	addSubCommands()
}

// initConfig reads in config file and ENV variables if set.
func initConfig() {
	homedir := ""
	if cfgFile != "" {
		// Use config file from the flag.
		viper.SetConfigFile(cfgFile)
	} else {
		// Find home directory.
		home, err := os.UserHomeDir()
		cobra.CheckErr(err)
		homedir = home

		workDir := fmt.Sprintf("%s/.config/binary-version-switcher", home)
		// if _, err := os.Stat(workDir); err != nil {
		// 	if os.IsNotExist(err) {
		// 		mkerr := os.MkdirAll(workDir, os.ModePerm)
		// 		if mkerr != nil {
		// 			c.Logger.Fatal("Error creating ~/.config/binary-version-switcher directory", mkerr)
		// 		}
		// 	}
		// }
		createConfigSubDirectory(home, fmt.Sprintf("plugins/%s", semVer))
		createConfigSubDirectory(home, "bin")
		createConfigSubDirectory(home, "bvs")

		if stat, err := os.Stat(workDir); err == nil && stat.IsDir() {
			configFile := fmt.Sprintf("%s/%s", workDir, "config.yaml")
			createRestrictedConfigFile(configFile)
			viper.SetConfigFile(configFile)
		} else {
			c.Logger.Info("The ~/.config/binary-version-switcher path is a file and not a directory, please remove the 'binary-version-switcher' file.")
			os.Exit(1)
		}
	}

	viper.AutomaticEnv() // read in environment variables that match

	// If a config file is found, read it in.
	if err := viper.ReadInConfig(); err != nil {
		c.Logger.Warn("Failed to read viper config file.")
	}

	// Check Config Items, prompt for defaults
	a := ask.New(c.LogLevel)
	// Symlink Path
	if viper.IsSet("symlinkPath") {
		c.SymLinkPath = viper.GetString("symlinkPath")
	} else {
		// Ask for the path
		paths := strings.Split(os.Getenv("PATH"), ":")
		paths = append(paths, fmt.Sprintf("%s/.config/binary-version-switcher/bin", homedir))
		p := a.PromptForPath(paths, "Choose a path to create symlinks in:")
		if len(p) == 0 {
			c.Logger.Fatal("you must select a path to create symlinks in")
		}
		// Set the default
		viper.Set("symlinkPath", p)
		c.SymLinkPath = p
		verr := viper.WriteConfig()
		if verr != nil {
			c.Logger.WithError(verr).Info("Failed to write config")
		}
	}

	// Symlink Path
	if viper.IsSet("binPath") {
		c.BinPath = viper.GetString("binPath")
	} else {
		// Ask for the path
		paths := []string{
			fmt.Sprintf("%s/.config/binary-version-switcher/bvs", homedir),
			"/usr/local/bvs",
		}
		p := a.PromptForPath(paths, "Choose a path to store downloaded binaries in:")
		if len(p) == 0 {
			c.Logger.Fatal("you must select a path to store binaries in")
		}
		// Set the default
		viper.Set("binPath", p)
		c.BinPath = p
		verr := viper.WriteConfig()
		if verr != nil {
			c.Logger.WithError(verr).Info("Failed to write config")
		}
	}

	c.PluginPath = fmt.Sprintf("%s/.config/binary-version-switcher/plugins/%s", homedir, semVer)
	// // Have Plugins?
	// plugIns, perr := os.ReadDir(fmt.Sprintf("%s/.config/binary-version-switcher/plugins/", homedir))
	// if perr != nil {
	// 	c.Logger.WithError(perr).Error("failed to get a list of plugins")
	// 	return
	// }
	// if len(plugIns) == 0 {
	// 	// prompt for plugins to download
	// 	// TODO: get this list from the repository
	// 	plugs := []string{
	// 		"kubectl",
	// 		"terraform",
	// 	}
	// 	p := a.PromptForMultipleString(plugs, "Select plugins to install")
	// 	for _, i := range p {
	// 		srcFile := fmt.Sprintf("/Users/christopher.maahs/src/cmaahsProjects/Maahsome/binary-version-switcher-plugins/%s/bvs-%s.so", i, i)
	// 		dstFile := fmt.Sprintf("%s/.config/binary-version-switcher/plugins/bvs-%s.so", homedir, i)
	// 		_, err := copy(srcFile, dstFile)
	// 		if err != nil {
	// 			c.Logger.WithError(err).Fatalf("failed to copy plugin %s", i)
	// 		}
	// 	}
	// }

}

func copy(src, dst string) (int64, error) {
	sourceFileStat, err := os.Stat(src)
	if err != nil {
		return 0, err
	}

	if !sourceFileStat.Mode().IsRegular() {
		return 0, fmt.Errorf("%s is not a regular file", src)
	}

	source, err := os.Open(src)
	if err != nil {
		return 0, err
	}
	defer source.Close()

	destination, err := os.Create(dst)
	if err != nil {
		return 0, err
	}
	defer destination.Close()
	nBytes, err := io.Copy(destination, source)
	return nBytes, err
}

func createConfigSubDirectory(homedir string, dirname string) {
	workDir := fmt.Sprintf("%s/.config/binary-version-switcher/%s", homedir, dirname)
	if _, err := os.Stat(workDir); err != nil {
		if os.IsNotExist(err) {
			mkerr := os.MkdirAll(workDir, os.ModePerm)
			if mkerr != nil {
				c.Logger.Fatal("Error creating ~/.config/binary-version-switcher directory", mkerr)
			}
		}
	}
}

func createRestrictedConfigFile(fileName string) {
	if _, err := os.Stat(fileName); err != nil {
		if os.IsNotExist(err) {
			file, ferr := os.Create(fileName)
			if ferr != nil {
				c.Logger.Info("Unable to create the configfile.")
				os.Exit(1)
			}
			mode := int(0600)
			if cherr := file.Chmod(os.FileMode(mode)); cherr != nil {
				c.Logger.Info("Chmod for config file failed, please set the mode to 0600.")
			}
		}
	}
}

// ClientSemVer - returns the full semVer as the first string and the numerical
// portion as the second string, they may be identical. One example where they
// would not be is:
//
//	semVer: v0.1.1-cacert -> (v0.1.1-cacert, v0.1.1).
func ClientSemVer() (string, string) {
	submatches := semVerReg.FindStringSubmatch(semVer)
	if submatches == nil || len(submatches) < 2 {
		c.Logger.Fatalf("the semver in the current build is not valid: %s", semVer)
	}
	return submatches[0], submatches[1]
}
