package cmd

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"regexp"

	"github.com/PrashantRaj18198/parsedir/pkg/parser"
	"github.com/spf13/cobra"
	"gopkg.in/yaml.v3"

	"github.com/spf13/viper"
)

var cfgFile string

type ParseDirFlags struct {
	VariablesFile string
	TemplateDir   string
	OutputDir     string
	Replacements  map[string]string
}

var ParseDirFlagsVar = ParseDirFlags{}

func Completions() *cobra.Command {
	c := &cobra.Command{
		Use:       "completion [shell]",
		Short:     "Generate auto complete for given shell",
		Args:      cobra.OnlyValidArgs,
		ValidArgs: []string{"bash", "fish", "zsh", "powershell"},
		RunE: func(cmd *cobra.Command, args []string) error {
			if len(args) == 0 {
				return errors.New("atleast one argument is required")
			}
			cmd.SilenceUsage = true
			switch args[0] {
			case "bash":
				rootCmd.GenBashCompletion(os.Stdout)
			case "fish":
				rootCmd.GenFishCompletion(os.Stdout, true)
			case "zsh":
				rootCmd.GenZshCompletion(os.Stdout)
			case "powershell":
				rootCmd.GenPowerShellCompletion(os.Stdout)
			default:
				fmt.Fprintf(os.Stderr, "%s is not a supported shell", args[0])
			}
			return nil
		},
	}
	return c
}

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:   "parsedir",
	Short: "Parse a golang template folder using this tool",
	Long: `Parse a golang template folder using this tools.
Pass a directory location which has your templates, template in path is also valid.
To get a detailed doc on how to use this tools refer to the docs/ folder on github.


`,
	Example: `
Example folder structure:

|example
|---- {{.dog.name}}
|-------- {{dog.name}}.yaml
|---- {{range .pets}}{{.name}}
|-------- detail.txt

./parse --vars-file=config.yaml --template-dir example/ --out-dir result/
The filepath will be generated from the config.yaml and written to result/ dir.
`,
	RunE: func(cmd *cobra.Command, args []string) error {
		files, err := parser.RecurseThroughDir(ParseDirFlagsVar.TemplateDir)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Could not recurse through directory %s. Error: %v\n", ParseDirFlagsVar.TemplateDir, err)
			os.Exit(1)
		}
		fileinfos, err := parser.ReadAllFiles(files)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Could not read all files. Error: %v\n", err)
			os.Exit(1)
		}
		data, err := ioutil.ReadFile(ParseDirFlagsVar.VariablesFile)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Could not read config file. Error: %v\n", err)
			os.Exit(1)
		}
		string_data := string(data)
		var out interface{}
		if string_data[0] == '{' || string_data[0] == '[' {
			fmt.Fprintf(os.Stdout, "File passed is in json format.\n")
			err = json.Unmarshal(data, &out)
		} else {
			fmt.Fprintf(os.Stdout, "File passed is in yaml format.\n")
			err = yaml.Unmarshal(data, &out)
		}
		if err != nil {
			fmt.Fprintf(os.Stderr, "Could not convert file to yaml. Error: %v\n", err)
			os.Exit(1)
		}
		parsedFiles, err := parser.PopulateAllFiles(fileinfos, out)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error: %v\n", err)
			os.Exit(1)
		}
		for _, f := range parsedFiles {
			for key, val := range ParseDirFlagsVar.Replacements {
				re, err := regexp.Compile(key)
				if err != nil {
					return err
				}
				f.Path = string(re.ReplaceAll([]byte(f.Path), []byte(val)))
			}
			parser.WriteFile([]byte(f.Content), filepath.Join(ParseDirFlagsVar.OutputDir, f.Path))
		}
		return nil
	},
}

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
	cobra.CheckErr(rootCmd.Execute())
}

func init() {
	cobra.OnInitialize(initConfig)

	// Here you will define your flags and configuration settings.
	// Cobra supports persistent flags, which, if defined here,
	// will be global for your application.

	rootCmd.PersistentFlags().StringVar(&cfgFile, "config", "", "config file (default is $HOME/.parsedir.yaml)")

	// Cobra also supports local flags, which will only run
	// when this action is called directly.
	rootCmd.Flags().StringVar(&ParseDirFlagsVar.VariablesFile, "vars-file", "", "Pass a file to read the input data from")
	rootCmd.MarkFlagRequired("vars-file")
	rootCmd.Flags().StringVar(&ParseDirFlagsVar.TemplateDir, "template-dir", "", "The directory which has the templates")
	rootCmd.MarkFlagRequired("template-dir")
	rootCmd.Flags().StringVar(&ParseDirFlagsVar.OutputDir, "out-dir", "", "The base directory where the output needs to be written to")
	rootCmd.Flags().StringToStringVar(&ParseDirFlagsVar.Replacements, "replace", map[string]string{}, "Replace key with value. Key must be a valid regex")

}

// initConfig reads in config file and ENV variables if set.
func initConfig() {
	if cfgFile != "" {
		// Use config file from the flag.
		viper.SetConfigFile(cfgFile)
	} else {
		// Find home directory.
		home, err := os.UserHomeDir()
		cobra.CheckErr(err)

		// Search config in home directory with name ".parsedir" (without extension).
		viper.AddConfigPath(home)
		viper.SetConfigType("yaml")
		viper.SetConfigName(".parsedir")
	}

	viper.AutomaticEnv() // read in environment variables that match

	// If a config file is found, read it in.
	if err := viper.ReadInConfig(); err == nil {
		fmt.Fprintln(os.Stderr, "Using config file:", viper.ConfigFileUsed())
	}
}
