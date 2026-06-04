// Package kbx has default configuration values
package kbx

import (
	"context"
	"net"
	"strings"

	"github.com/kubex-ecosystem/ethyr/tools/flow"
	"github.com/kubex-ecosystem/ethyr/tools/flow/control"
	"github.com/kubex-ecosystem/ethyr/tools/flow/fsm"
	"github.com/spf13/cobra"

	kbxGet "github.com/kubex-ecosystem/kbx/get"
	gl "github.com/kubex-ecosystem/logz"
	// "github.com/kubex-ecosystem/xtui/internal/module/version"
)

// Estados do Ciclo de Vida do XTUI (Atômicos)
const (
	XTUIIdle fsm.State = 1 << iota
	XTUIBootstrapping
	XTUIConfiguring

	XTUIReading
	XTUIProcessing
	XTUIWaiting

	XTUIError
	XTUIStopping
	XTUIUnknown
)

// Eventos do XTUI
const (
	EvBoot fsm.Event = iota
	EvConfigure
	EvRead
	EvProcess
	EvWait
	EvOutput
	EvStop
	EvDone
)

type VersionService interface {
	// GetLatestVersion retrieves the latest version from the Git repository.
	GetLatestVersion() (string, error)
	// GetCurrentVersion returns the current version of the service.
	GetCurrentVersion() string
	// IsLatestVersion checks if the current version is the latest version.
	IsLatestVersion() (bool, error)
	// GetName returns the name of the service.
	GetName() string
	// GetVersion returns the current version of the service.
	GetVersion() string
	// GetRepository returns the Git repository URL of the service.
	GetRepository() string
}

// KbxModuleWithCtlr is a module that has a control manager and a flow instance
// using generics [T any] to define the type of data that the module will process
// T is the type of data that the module will process
// NOTE: If you don't need unsafe exposture and/or don't have any idea about what this
// interface means, I strongly recommend you not to use it. Trust me.
type KbxModuleWithCtlr[T any] interface {
	KbxModuleWithManager

	// FlowCtlr returns the FSM control manager pointer
	FlowCtlr(ctx context.Context) *fsm.FSM

	// Flow returns the flow instance pointer
	Flow(ctx context.Context) *flow.Flow

	// FlowPassCtlr returns the FlowPass instance pointer
	FlowPassCtlr(ctx context.Context) flow.Pass[T]
}

type KbxModuleWithService interface {
	KbxModule

	// Run runs a module as a service
	Run(opts ...any) error
}

type KbxModuleWithManager interface {
	KbxModule

	// Control returns the control manager
	Control() control.ManagerControl[control.StepFlag]

	// Status returns the current state of the module
	Status() (uint32, uint32, uint32, error) // (fsm.StateFlag, control.StepFlag, control.JobFlag, error)

	// Done channel is used to listen for state changes
	Done() chan<- uint32 // Only send updates to listeners, not receive

	// String returns a string representation of the manager
	String() string
}

type KbxModule interface {
	// Init initializes the module
	Init(ctx context.Context, subCmdsList []*cobra.Command) error

	// Alias returns the alias of the module
	Alias() string

	// ShortDescription returns a short description of the module
	ShortDescription() string

	// LongDescription returns a long description of the module
	LongDescription() string

	// Usage returns the usage of the module
	Usage() string

	// Examples returns examples of the module
	Examples() []string

	// Active returns true if the module is active
	Active() bool

	// Module returns the name of the module
	Module() string

	// Execute executes the module CLI root command
	Execute() error

	// Version returns the version service
	Version() VersionService
}

// Default configuration constants
const (
	DefaultXTUIName = "xtui"

	DefaultKubexConfigDir = "$HOME/.kubex/xtui"

	DefaultXTUICAPath   = "$HOME/.kubex/xtui/ca-cert.pem"
	DefaultXTUIKeyPath  = "$HOME/.kubex/xtui/xtui.key" // Priv
	DefaultXTUICertPath = "$HOME/.kubex/xtui/xtui.crt"

	DefaultXTUIEnvPath    = "$HOME/.kubex/xtui/config/.env"
	DefaultXTUIConfigPath = "$HOME/.kubex/xtui/config/config.json"
	DefaultXTUILogPath    = "$HOME/.kubex/xtui/logs/xtui_process.log.txt"

	DefaultKubexDomusConfigPath = "$HOME/.kubex/domus/config/config.json"

	DefaultMailConfigPath       = "$HOME/.kubex/xtui/config/mail_config.json"
	DefaultGoogleAuthClientPath = "$HOME/.kubex/xtui/config/xtui_google_client.json"

	DefaultFirebaseAPIKeyEnv         = "KUBEX_FIREBASE_API_KEY"
	DefaultFirebaseProjectIDEnv      = "KUBEX_FIREBASE_PROJECT_ID"
	DefaultFirebaseSDKAdmFilePathEnv = "$HOME/.kubex/xtui/config/xtui_firebase_service_account.json"

	DefaultVaultDir = "$HOME/.kubex/xtui/secrets"

	DefaultVaultKey = "kubex_kubex-jwt_secret.secret"

	DefaultTemplatesDir = "templates"

	DefaultProvidersConfig = "$HOME/.kubex/xtui/config/providers.yaml"
)

// Default General Rate Limiting Settings
const (
	DefaultRateLimitLimit  = 100
	DefaultRateLimitBurst  = 100
	DefaultRateLimitJitter = 0.1
	DefaultRequestWindow   = 1 * 60 * 1000 // 1 minute
)

// Default HTTP Client Settings
const (
	DefaultTLSHandshakeTimeout   = 10 * 1000 // 10 seconds
	DefaultExpectContinueTimeout = 1 * 1000  // 1 second
	DefaultResponseHeaderTimeout = 5 * 1000  // 5 seconds

	DefaultTimeout         = 30 * 1000 // 30 seconds
	DefaultKeepAlive       = 30 * 1000 // 30 seconds
	DefaultMaxConnsPerHost = 100
)

// Default Generic Retry and Connection Settings
const (
	DefaultMaxRetries = 3
	DefaultRetryDelay = 1 * 1000 // 1 second

	DefaultMaxIdleConns        = 100
	DefaultMaxIdleConnsPerHost = 100
	DefaultIdleConnTimeout     = 90 * 1000 // 90 seconds
)

// Default LLM Settings
const (
	DefaultLLMProvider    = "gemini"
	DefaultLLMModel       = "gemini-flash-latest"
	DefaultLLMMaxTokens   = 1024
	DefaultLLMTemperature = 0.3
)

// Default LLM Environment Variables
const (
	DefaultLLMOpenAIKeyEnv       = "OPENAI_API_KEY"
	DefaultLLMGoogleKeyEnv       = "GOOGLE_API_KEY"
	DefaultLLMAzureKeyEnv        = "AZURE_API_KEY"
	DefaultLLMAnthropicKeyEnv    = "ANTHROPIC_API_KEY"
	DefaultLLMGeminiKeyEnv       = "GEMINI_API_KEY"
	DefaultLLMOllamaKeyEnv       = "OLLAMA_API_KEY"
	DefaultLLMChatGPTKeyEnv      = "CHATGPT_API_KEY"
	DefaultLLMDeepseekKeyEnv     = "DEEPSEEK_API_KEY"
	DefaultLLMCohereKeyEnv       = "COHERE_API_KEY"
	DefaultLLMGroqKeyEnv         = "GROQ_API_KEY"
	DefaultLLMGrokKeyEnv         = "GROK_API_KEY"
	DefaultLLMMistralKeyEnv      = "MISTRAL_API_KEY"
	DefaultLLMCustomKeyEnv       = "CUSTOM_API_KEY"
	DefaultLLMMetaKeyEnv         = "META_API_KEY"
	DefaultLLMClaudeKeyEnv       = "CLAUDE_API_KEY"
	DefaultLLMErnieKeyEnv        = "ERNIE_API_KEY"
	DefaultLLMCustomKeyEnvPrefix = "CUSTOM_"
	DefaultLLMCustomKeyEnvSuffix = "_KEY_ENV"
)

// DefaultXTUILoopbackIP is the default loopback IP for XTUI
var DefaultXTUILoopbackIP string

// Default Server Settings
const (
	DefaultServerPort = "5000"
	DefaultServerBind = "0.0.0.0"
	DefaultServerHost = "localhost"
)

func init() {
	hostIps, err := net.LookupHost("localhost")
	if err != nil {
		gl.Warnf("Failed to lookup host localhost: %v", err)
	}

	DefaultXTUILoopbackIP = strings.Join(kbxGet.ValOrType(hostIps, []string{"localhost", "::1"}), ",")
}

// Default HTTP Basic Header Security Keys
const (
	HeaderRequestIDKey = "X-Request-ID"
	CookieSessionIDKey = "session_id"
)

// Default Authentication Types
const (
	AuthTypeNone   = "none"
	AuthTypeOIDC   = "oidc"
	AuthTypeBasic  = "basic"
	AuthTypeBearer = "bearer"
	AuthTypeAPIKey = "api_key" // pragma: allowlist secret
)

// Default Database Settings

// DBNameKey is a type for database names
type DBNameKey string

// Default Database Names
const (
	ContextDBNameKey      = DBNameKey("postgres")
	DefaultVolumesDir     = "$HOME/.kubex/domus/volumes"
	DefaultMongoVolume    = "$HOME/.kubex/domus/volumes/mongo"
	DefaultRedisVolume    = "$HOME/.kubex/domus/volumes/redis"
	DefaultPostgresVolume = "$HOME/.kubex/domus/volumes/postgresql"
	DefaultRabbitMQVolume = "$HOME/.kubex/domus/volumes/rabbitmq"
)

// Default Sankhya Catalog Settings
const (
	DefaultSankhyaCatalogDir     = "$HOME/.kubex/xtui/references/data/catalogo_bi"
	DefaultSankhyaConfigDir      = "$HOME/.kubex/xtui/config/getl/sankhya_catalog"
	DefaultSankhyaSyncManifest   = "$HOME/.kubex/xtui/config/getl/sankhya_catalog/sync.manifest.json"
	DefaultSankhyaCatalogSchema  = "sankhya_catalog"
	DefaultSankhyaCatalogDomain  = "bi_catalog"
	DefaultSankhyaCatalogRefresh = "full_refresh"
)

// Default Sankhya Server Settings
const (
	DefaultSankhyaHost     = "localhost:8081"
	DefaultSankhyaService  = ""
	DefaultSankhyaSession  = ""
	DefaultSankhyaUser     = ""
	DefaultSankhyaPassword = ""
	DefaultSankhyaAppKey   = ""
)
