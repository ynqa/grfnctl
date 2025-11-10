package cmd

var (
	// version is overridden at build time via -ldflags "-X github.com/ynqa/grfnctl/cmd.version=vX.Y.Z".
	version = "dev"
)

func init() {
	rootCmd.Version = version
	rootCmd.SetVersionTemplate(`{{printf "%s version %s\n" .Name .Version}}`)
}
