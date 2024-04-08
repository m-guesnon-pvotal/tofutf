package main

import (
	"context"
	"fmt"
	"io"
	"log/slog"
	"os"
	"os/signal"
	"syscall"

	"cloud.google.com/go/profiler"
	"github.com/pkg/errors"
	"github.com/spf13/cobra"
	cmdutil "github.com/tofutf/tofutf/cmd"
	"github.com/tofutf/tofutf/internal"
	"github.com/tofutf/tofutf/internal/agent"
	"github.com/tofutf/tofutf/internal/authenticator"
	"github.com/tofutf/tofutf/internal/daemon"
	"github.com/tofutf/tofutf/internal/github"
	"github.com/tofutf/tofutf/internal/gitlab"
	"github.com/tofutf/tofutf/internal/xslog"
)

const (
	defaultAddress  = ":8080"
	defaultDatabase = "postgres:///otf?host=/var/run/postgresql"
)

func main() {
	// Configure ^C to terminate program
	ctx, cancel := signal.NotifyContext(context.Background(), syscall.SIGTERM, syscall.SIGINT)
	go func() {
		<-ctx.Done()
		// Stop handling ^C; another ^C will exit the program.
		cancel()
	}()

	defer func() {
		fmt.Printf("Closing app")
		if x := recover(); x != nil {
			// recovering from a panic; x contains whatever was passed to panic()
			fmt.Printf("run time panic: %v", x)

			// if you just want to log the panic, panic again
			panic(x)
		}
	}()

	cfg := profiler.Config{
		Service:        "tofutf",
		ServiceVersion: "0.9.1",
		// ProjectID must be set if not running on GCP.
		// ProjectID: "my-project",

		// For OpenCensus users:
		// To see Profiler agent spans in APM backend,
		// set EnableOCTelemetry to true
		// EnableOCTelemetry: true,
	}

	// Profiler initialization, best done as early as possible.
	if err := profiler.Start(cfg); err != nil {
		cmdutil.PrintError(err)
	}
	fmt.Printf("Starting\n")

	if err := parseFlags(ctx, os.Args[1:], os.Stderr); err != nil {
		cmdutil.PrintError(err)
	}
	fmt.Printf("Exited\n")
}

func parseFlags(ctx context.Context, args []string, out io.Writer) error {
	cfg := daemon.Config{}
	daemon.ApplyDefaults(&cfg)

	var logger *slog.Logger

	cmd := &cobra.Command{
		Use:           "tofutfd",
		Short:         "tofutf daemon",
		Long:          "tofutfd is the daemon component of the opentofu tuft framework.",
		SilenceUsage:  true,
		SilenceErrors: true,
		Version:       internal.Version,
		RunE: func(cmd *cobra.Command, args []string) error {

			// Confer superuser privileges on all calls to service endpoints
			ctx := internal.AddSubjectToContext(cmd.Context(), &internal.Superuser{Username: "app-user"})

			d, err := daemon.New(ctx, logger, cfg)
			if err != nil {
				return err
			}
			// block until ^C received
			return d.Start(ctx, make(chan struct{}))
		},
	}
	cmd.SetOut(out)

	// TODO: rename --address to --listen
	cmd.Flags().StringVar(&cfg.Address, "address", defaultAddress, "Listening address")
	cmd.Flags().StringVar(&cfg.Database, "database", defaultDatabase, "Postgres connection string")
	cmd.Flags().StringVar(&cfg.Host, "hostname", "", "User-facing hostname for otf")
	cmd.Flags().StringVar(&cfg.SiteToken, "site-token", "", "API token with site-wide unlimited permissions. Use with care.")
	cmd.Flags().StringSliceVar(&cfg.SiteAdmins, "site-admins", nil, "Promote a list of users to site admin.")
	cmd.Flags().BytesHexVar(&cfg.Secret, "secret", nil, "Hex-encoded 16 byte secret for cryptographic work. Required.")
	cmd.Flags().Int64Var(&cfg.MaxConfigSize, "max-config-size", cfg.MaxConfigSize, "Maximum permitted configuration size in bytes.")
	cmd.Flags().StringVar(&cfg.WebhookHost, "webhook-hostname", "", "External hostname for otf webhooks")

	cmd.Flags().StringVar(&cfg.CacheConfig.Address, "cache-address", "localhost:6379", "Redis address")
	cmd.Flags().DurationVar(&cfg.CacheConfig.TTL, "cache-expiry", internal.DefaultCacheTTL, "Cache entry TTL.")

	cmd.Flags().BoolVar(&cfg.SSL, "ssl", false, "Toggle SSL")
	cmd.Flags().StringVar(&cfg.CertFile, "cert-file", "", "Path to SSL certificate (required if enabling SSL)")
	cmd.Flags().StringVar(&cfg.KeyFile, "key-file", "", "Path to SSL key (required if enabling SSL)")
	cmd.Flags().BoolVar(&cfg.EnableRequestLogging, "log-http-requests", false, "Log HTTP requests")
	cmd.Flags().BoolVar(&cfg.DevMode, "dev-mode", false, "Enable developer mode.")

	cmd.Flags().StringVar(&cfg.GithubHostname, "github-hostname", github.DefaultHostname, "github hostname")
	cmd.Flags().StringVar(&cfg.GithubClientID, "github-client-id", "", "github client ID")
	cmd.Flags().StringVar(&cfg.GithubClientSecret, "github-client-secret", "", "github client secret")

	cmd.Flags().StringVar(&cfg.GitlabHostname, "gitlab-hostname", gitlab.DefaultHostname, "gitlab hostname")
	cmd.Flags().StringVar(&cfg.GitlabClientID, "gitlab-client-id", "", "gitlab client ID")
	cmd.Flags().StringVar(&cfg.GitlabClientSecret, "gitlab-client-secret", "", "gitlab client secret")

	cmd.Flags().StringVar(&cfg.BitbucketServerHostname, "bitbucketserver-hostname", cfg.BitbucketServerHostname, "bitbucket server hostname")

	cmd.Flags().StringVar(&cfg.OIDC.Name, "oidc-name", "", "User friendly OIDC name")
	cmd.Flags().StringVar(&cfg.OIDC.IssuerURL, "oidc-issuer-url", "", "OIDC issuer URL")
	cmd.Flags().StringVar(&cfg.OIDC.ClientID, "oidc-client-id", "", "OIDC client ID")
	cmd.Flags().StringVar(&cfg.OIDC.ClientSecret, "oidc-client-secret", "", "OIDC client secret")
	cmd.Flags().StringSliceVar(&cfg.OIDC.Scopes, "oidc-scopes", authenticator.DefaultOIDCScopes, "OIDC scopes")
	cmd.Flags().StringVar(&cfg.OIDC.UsernameClaim, "oidc-username-claim", string(authenticator.DefaultUsernameClaim), "OIDC claim to be used for username (name, email, or sub)")

	cmd.Flags().BoolVar(&cfg.RestrictOrganizationCreation, "restrict-org-creation", false, "Restrict organization creation capability to site admin role")

	cmd.Flags().StringVar(&cfg.GoogleIAPConfig.Audience, "google-jwt-audience", "", "The Google JWT audience claim for validation. If unspecified then validation is skipped")

	cmd.Flags().StringVar(&cfg.ProviderProxy.URL, "provider-proxy-url", "", "The URL of the provider registry to proxy provider registry requests to")
	cmd.Flags().BoolVar(&cfg.ProviderProxy.IsArtifactory, "provider-proxy-is-artifactory", false, "Set to true if using artifactory as the backing provider registry")

	loggerConfig := xslog.NewConfigFromFlags(cmd.Flags())
	cfg.AgentConfig = agent.NewConfigFromFlags(cmd.Flags())

	if err := cmdutil.SetFlagsFromEnvVariables(cmd.Flags()); err != nil {
		return errors.Wrap(err, "failed to populate config from environment vars")
	}
	logger, err := xslog.New(loggerConfig)
	if err != nil {
		return err
	}

	cmd.SetArgs(args)
	err = cmd.ExecuteContext(ctx)

	if err != nil {
		logger.Error("command exited with an error", "err", err)
	}
	return err
}
