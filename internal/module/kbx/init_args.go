// Package kbx provides utilities for working with initialization arguments.
package kbx

import (
	"os"
	"path/filepath"
	"time"

	"github.com/google/uuid"

	kbxGet "github.com/kubex-ecosystem/kbx/get"
)

var (
	CLIArgs *InitArgs
)

// Reference is the internal struct that holds the server startup unique identifier with name.
type Reference struct {
	// refID is the unique identifier for this context.
	ID uuid.UUID
	// refName is the name of the context.
	Name string
}

// GetReference returns the reference.
func (r *Reference) GetReference() *Reference {
	return r
}

// GetID returns the reference ID.
func (r *Reference) GetID() uuid.UUID {
	return r.ID
}

// GetName returns the reference name.
func (r *Reference) GetName() string {
	return r.Name
}

// GetType returns the reference type.
func (r *Reference) GetType() string {
	return "kbx_init_args"
}

// NewReference creates a new Reference with the given name.
func NewReference(name string) *Reference {
	return &Reference{
		ID:   uuid.New(),
		Name: name,
	}
}

// InitArgs holds the initialization arguments for the server.
type InitArgs struct {
	*Reference `yaml:"-" json:"-" mapstructure:"-"`

	// Basic options
	HideBanner     bool     `yaml:"hide_banner" json:"hide_banner" mapstructure:"hide_banner"`
	Debug          bool     `yaml:"debug" json:"debug" mapstructure:"debug"`
	ReleaseMode    *bool    `yaml:"release_mode" json:"release_mode" mapstructure:"release_mode"`
	IsConfidential bool     `yaml:"is_confidential" json:"is_confidential" mapstructure:"is_confidential"`
	CORSEnabled    *bool    `yaml:"enable_cors" json:"enable_cors" mapstructure:"enable_cors"`
	TrustedProxies []string `yaml:"trusted_proxies" json:"trusted_proxies" mapstructure:"trusted_proxies"`
	UIEnabled      bool     `yaml:"ui_enabled" json:"ui_enabled" mapstructure:"ui_enabled"`

	// Paths and files

	Cwd                  string `yaml:"cwd,omitempty" json:"cwd,omitempty" mapstructure:"cwd,omitempty"`
	LogFile              string `yaml:"log_file,omitempty" json:"log_file,omitempty" mapstructure:"log_file,omitempty"`
	EnvFile              string `yaml:"env_file,omitempty" json:"env_file,omitempty" mapstructure:"env_file,omitempty"`
	ConfigFile           string `yaml:"config_file,omitempty" json:"config_file,omitempty" mapstructure:"config_file,omitempty"`
	DBConfigFile         string `yaml:"db_config_file,omitempty" json:"db_config_file,omitempty" mapstructure:"db_config_file,omitempty"`
	MailerConfigFile     string `yaml:"mail_config_file,omitempty" json:"mail_config_file,omitempty" mapstructure:"mail_config_file,omitempty"`
	GoogleAuthClientPath string `yaml:"google_auth_client_path,omitempty" json:"google_auth_client_path,omitempty" mapstructure:"google_auth_client_path,omitempty"`

	FirebaseAPIKey         string `yaml:"-" json:"-" mapstructure:"-"`                                                                                                          // FirebaseAPIKey will be loaded from the environment variable at runtime. (ANY ENVIRONMENT)
	FirebaseProjectID      string `yaml:"firebase_project_id,omitempty" json:"firebase_project_id,omitempty" mapstructure:"firebase_project_id,omitempty"`                      // FirebaseProjectID will be loaded from the environment variable at runtime. (PRODUCTION ONLY)
	FirebaseSDKAdmFilePath string `yaml:"firebase_sdk_adm_file_path,omitempty" json:"firebase_sdk_adm_file_path,omitempty" mapstructure:"firebase_sdk_adm_file_path,omitempty"` // FirebaseKeyPath is just the path to the file containing configuration for the Firebase Admin SDK. (DEVELOPMENT ONLY. NOT FOR PRODUCTION.)

	ProvidersConfig string `yaml:"providers_config,omitempty" json:"providers_config,omitempty" mapstructure:"providers_config,omitempty"`
	ScorecardPath   string `yaml:"scorecard_path,omitempty" json:"scorecard_path,omitempty" mapstructure:"scorecard_path,omitempty"`
	TemplatesDir    string `yaml:"template_dir,omitempty" json:"template_dir,omitempty" mapstructure:"template_dir,omitempty"`
	WebDir          string `yaml:"web_dir,omitempty" json:"web_dir,omitempty" mapstructure:"web_dir,omitempty"`

	// Runtime options

	Host            string        `yaml:"host,omitempty" json:"host,omitempty" mapstructure:"host,omitempty"`
	Port            string        `yaml:"port,omitempty" json:"port,omitempty" mapstructure:"port,omitempty"`
	Bind            string        `yaml:"bind,omitempty" json:"bind,omitempty" mapstructure:"bind,omitempty"`
	PubCertKeyPath  string        `yaml:"pub_cert_key_path,omitempty" json:"pub_cert_key_path,omitempty" mapstructure:"pub_cert_key_path,omitempty"`
	PubKeyPath      string        `yaml:"pub_key_path,omitempty" json:"pub_key_path,omitempty" mapstructure:"pub_key_path,omitempty"`
	PrivKeyPath     string        `yaml:"priv_key_path,omitempty" json:"priv_key_path,omitempty" mapstructure:"priv_key_path,omitempty"`
	AccessTokenTTL  time.Duration `yaml:"access_token_ttl,omitempty" json:"access_token_ttl,omitempty" mapstructure:"access_token_ttl,omitempty"`
	RefreshTokenTTL time.Duration `yaml:"refresh_token_ttl,omitempty" json:"refresh_token_ttl,omitempty" mapstructure:"refresh_token_ttl,omitempty"`
	Issuer          string        `yaml:"issuer,omitempty" json:"issuer,omitempty" mapstructure:"issuer,omitempty"`

	// Advanced options

	Context    string            `yaml:"context,omitempty" json:"context,omitempty" mapstructure:"context,omitempty"`
	Command    string            `yaml:"command,omitempty" json:"command,omitempty" mapstructure:"command,omitempty"`
	Subcommand string            `yaml:"subcommand,omitempty" json:"subcommand,omitempty" mapstructure:"subcommand,omitempty"`
	Args       string            `yaml:"args,omitempty" json:"args,omitempty" mapstructure:"args,omitempty"`
	EnvVars    map[string]string `yaml:"env_vars,omitempty" json:"env_vars,omitempty" mapstructure:"env_vars,omitempty"`

	// Flags

	FailFast  bool `yaml:"fail_fast,omitempty" json:"fail_fast,omitempty" mapstructure:"fail_fast,omitempty"`
	Verbose   bool `yaml:"verbose,omitempty" json:"verbose,omitempty" mapstructure:"verbose,omitempty"`
	BatchMode bool `yaml:"batch_mode,omitempty" json:"batch_mode,omitempty" mapstructure:"batch_mode,omitempty"`
	NoColor   bool `yaml:"no_color,omitempty" json:"no_color,omitempty" mapstructure:"no_color,omitempty"`
	TraceMode bool `yaml:"trace_mode,omitempty" json:"trace_mode,omitempty" mapstructure:"trace_mode,omitempty"`
	RootMode  bool `yaml:"root_mode,omitempty" json:"root_mode,omitempty" mapstructure:"root_mode,omitempty"`

	// Performance options

	MaxProcs  int    `yaml:"max_procs,omitempty" json:"max_procs,omitempty" mapstructure:"max_procs,omitempty"`
	TimeoutMS int    `yaml:"timeout_ms,omitempty" json:"timeout_ms,omitempty" mapstructure:"timeout_ms,omitempty"`
	Hash      string `yaml:"hash,omitempty" json:"hash,omitempty" mapstructure:"hash,omitempty"`
}

// NewInitArgs creates a new InitArgs with the given parameters.
func NewInitArgs(
	configFile string,
	configType string,
	configDBFile string,
	configDBType string,
	googleAuthClientPath string,
	envFile string,
	logFile string,
	name string,
	debug bool,
	releaseMode bool,
	isConfidential bool,
	port string,
	bind string,
	pubCertKeyPath string,
	pubKeyPath string,
	pwd string,
	templatesDir string,
	uiEnabled bool,
) *InitArgs {
	// Resiliency layer: use os.ExpandEnv to resolve environment variables in paths
	// Define default values, if not provided
	ref := NewReference(name)

	logFile = os.ExpandEnv(kbxGet.ValOrType(kbxGet.EnvOr("KUBEX_XTUI_LOG_FILE_PATH", logFile), DefaultXTUILogPath))
	configFile = os.ExpandEnv(kbxGet.ValOrType(kbxGet.EnvOr("KUBEX_XTUI_CONFIG_FILE_PATH", configFile), DefaultXTUIConfigPath))
	configDBFile = os.ExpandEnv(kbxGet.ValOrType(kbxGet.EnvOr("KUBEX_XTUI_DOMUS_CONFIG_FILE_PATH", configDBFile), DefaultKubexDomusConfigPath))
	googleAuthClientPath = os.ExpandEnv(kbxGet.ValOrType(kbxGet.EnvOr("KUBEX_GOOGLE_CREDENTIALS_PATH", googleAuthClientPath), DefaultGoogleAuthClientPath))
	envFile = os.ExpandEnv(kbxGet.ValOrType(kbxGet.EnvOr("KUBEX_XTUI_ENV_FILE_PATH", envFile), filepath.Join("$PWD", ".env")))
	port = kbxGet.ValOrType(kbxGet.EnvOr("KUBEX_XTUI_PORT", port), "5000")
	bind = kbxGet.ValOrType(kbxGet.EnvOr("KUBEX_XTUI_BIND", bind), DefaultServerBind)
	templatesDir = os.ExpandEnv(kbxGet.ValOrType(kbxGet.EnvOr("KUBEX_XTUI_TEMPLATES_DIR", templatesDir), DefaultTemplatesDir))
	webDir := os.ExpandEnv(kbxGet.ValOrType(kbxGet.EnvOr("KUBEX_XTUI_WEB_DIR", ""), ""))

	if CLIArgs == nil {
		CLIArgs = &InitArgs{
			Reference: ref.GetReference(),

			Debug:          kbxGet.ValOrType(debug, false),
			ReleaseMode:    BoolPtr(kbxGet.ValOrType(releaseMode, false)),
			IsConfidential: kbxGet.ValOrType(isConfidential, false),
			PubCertKeyPath: os.ExpandEnv(kbxGet.ValOrType(pubCertKeyPath, DefaultXTUIKeyPath)),
			PubKeyPath:     os.ExpandEnv(kbxGet.ValOrType(pubKeyPath, DefaultXTUICertPath)),
			Cwd:            os.ExpandEnv(kbxGet.ValOrType(pwd, "")),

			ConfigFile:           configFile,
			DBConfigFile:         configDBFile,
			GoogleAuthClientPath: googleAuthClientPath,
			EnvFile:              envFile,
			LogFile:              logFile,
			Port:                 port,
			Bind:                 bind,

			TemplatesDir: templatesDir,
			WebDir:       webDir,
		}
	} else {
		return &InitArgs{
			Reference: ref.GetReference(),

			Debug:          kbxGet.ValOrType(debug, false),
			ReleaseMode:    BoolPtr(kbxGet.ValOrType(releaseMode, false)),
			IsConfidential: kbxGet.ValOrType(isConfidential, false),
			PubCertKeyPath: os.ExpandEnv(kbxGet.ValOrType(pubCertKeyPath, DefaultXTUIKeyPath)),
			PubKeyPath:     os.ExpandEnv(kbxGet.ValOrType(pubKeyPath, DefaultXTUICertPath)),
			Cwd:            os.ExpandEnv(kbxGet.ValOrType(pwd, "")),

			ConfigFile:           configFile,
			DBConfigFile:         configDBFile,
			GoogleAuthClientPath: googleAuthClientPath,
			EnvFile:              envFile,
			LogFile:              logFile,
			Port:                 port,
			Bind:                 bind,

			TemplatesDir: templatesDir,
			WebDir:       webDir,
		}
	}

	return CLIArgs
}

func GetCLIArgs() *InitArgs {
	if CLIArgs == nil {
		return &InitArgs{}
	}

	return CLIArgs
}
