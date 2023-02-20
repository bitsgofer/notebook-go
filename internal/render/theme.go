package render

// Theme inspect the content dir and render HTML pages.
// It also provide different looks via CSS and additional features
// like syntax-highlighting with JS.
//
// The theme package will embed its content,
// so we can interact with them in code.
//
// NOTE: Technically, we can just work with the local file system,
// rather than having the theme embed its file. However, that approach
// will leak the location of the theme into the interface -> not as nice.
//
// => We can let the them embed its own file and just tell it where to
// compile assets into.
type Theme interface {
	CompileAssets(outputDir string) error
	Render(outputDir, contentDir string) error
}
