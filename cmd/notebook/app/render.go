package app

import (
	"fmt"

	"github.com/bitsgofer/notebook-go/internal/render"
	"github.com/bitsgofer/notebook-go/themes/bitsgofer"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var renderCmd = &cobra.Command{
	Use:   "render",
	Short: "Render articles, drafts and themes",
	RunE: func(cmd *cobra.Command, args []string) error {
		themeName := viper.GetString(cfgKeyTheme)
		contentDir := viper.GetString(cfgKeyContentDir)
		outputDir := viper.GetString(cfgKeyOutputDir)

		var theme render.Theme
		switch themeName {
		case "bitsgofer":
			theme = bitsgofer.New()
		default:
			return fmt.Errorf("theme does not exist")
		}

		fmt.Printf("Render:\n")
		fmt.Printf("- Theme  : %s\n", themeName)
		fmt.Printf("- Content: %s\n", contentDir)
		fmt.Printf("- Output : %s\n", outputDir)

		if err := theme.CompileAssets(outputDir); err != nil {
			return fmt.Errorf("cannot compile theme's assests; err= %w", err)
		}

		if err := theme.Render(outputDir, contentDir); err != nil {
			return fmt.Errorf("cannot render content; err= %w", err)
		}

		return nil
	},
}
