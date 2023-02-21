package app

// config keys + default values
const (
	cfgConfiFileFormat = "yaml"

	// root cmd
	cfgKeyConfigFile     = "config"
	cfgDefaultConfigFile = ".notebook"
	// render + develop + publish
	cfgKeyRender                  = "render"
	cfgKeyTheme                   = cfgKeyRender + ".theme"
	cfgKeyContentDir              = cfgKeyRender + ".contentDir"
	cfgKeyOutputDir               = cfgKeyRender + ".outputDir"
	cfgKeyDevServer               = "devServer"
	cfgKeyDevServerAddr           = cfgKeyDevServer + ".addr"
	cfgKeyDevServerDataDir        = cfgKeyDevServer + ".dir"
	cfgKeyPublish                 = "publish"
	cfgKeyPublishBranch           = cfgKeyPublish + ".branch"
	cfgKeyPublishDir              = cfgKeyPublish + ".dir"
	cfgDefaultPublishBranch       = "gh-pages"
	cfgDefaultOutputAndPublishDir = "_public_html"
)
