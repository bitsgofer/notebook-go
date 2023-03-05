## Components

### minicss

- The [original project](https://github.com/Chalarangelo/mini.css) is no longer maintained.
- We have [a fork at](https://github.com/bitsgofer/mini.css).
- After downloading v3.0.0, modify the `mini-default.css` based on these changes:

```diff
--- mini.css-3.0.0/dist/mini-default.css	2018-05-28 03:22:17.000000000 -0700
+++ mini-default.css	2022-10-06 21:03:25.098204307 -0700
@@ -18,7 +18,7 @@
   --pre-color: #1565c0;
   --border-color: #aaa;
   --secondary-border-color: #ddd;
-  --heading-ratio: 1.19;
+  --heading-ratio: 1.15;
   --universal-margin: 0.5rem;
   --universal-padding: 0.5rem;
   --universal-border-radius: 0.125rem;
@@ -85,7 +85,11 @@
 }

 h1 {
-  font-size: calc(1rem * var(--heading-ratio) * var(--heading-ratio) * var(--heading-ratio) * var(--heading-ratio));
+  margin-top: 4rem;
+  font-size: calc(1rem * var(--heading-ratio) * var(--heading-ratio) * var(--heading-ratio) * var(--heading-ratio) * var(--heading-ratio));
+  padding-bottom: 0.5rem;
+  border-bottom-style: solid;
+  border-bottom-width: 0.01rem;
 }

 h2 {
@@ -161,18 +165,23 @@
   white-space: pre;
 }

+/*
 code, kbd, pre, samp {
   font-family: Menlo, Consolas, monospace;
   font-size: 0.85em;
 }
+*/
+

-code {
+:not(pre) > code {
   background: var(--secondary-back-color);
   border-radius: var(--universal-border-radius);
   padding: calc(var(--universal-padding) / 4) calc(var(--universal-padding) / 2);
+  color: var(--pre-color);
 }

-kbd {
+/*
+ kbd {
   background: var(--fore-color);
   color: var(--back-color);
   border-radius: var(--universal-border-radius);
@@ -188,6 +197,7 @@
   border-left: 0.25rem solid var(--pre-color);
   border-radius: 0 var(--universal-border-radius) var(--universal-border-radius) 0;
 }
+*/

 sup, sub, code, kbd {
   line-height: 0;
@@ -1042,6 +1052,7 @@
   --header-hover-back-color: #f0f0f0;
   --header-fore-color: #444;
   --header-border-color: #ddd;
+  --header-height: 3.1875rem;
   --nav-back-color: #f8f8f8;
   --nav-hover-back-color: #f0f0f0;
   --nav-fore-color: #444;
@@ -1058,7 +1069,7 @@
 }

 header {
-  height: 3.1875rem;
+  height: var(--header-height);
   background: var(--header-back-color);
   color: var(--header-fore-color);
   border-bottom: 0.0625rem solid var(--header-border-color);
@@ -1272,6 +1283,37 @@
 }

 /*
+  Definitions for embedded iframe, mainly for GoogleDoc drafts.
+*/
+
+#doc-content .googledoc {
+	/* 100% view port height - var(--header-height). But the CSS variable didn't work */
+	height: calc(100vh - 3.1875rem); /* CSS3 */
+	height: -o-calc(100vh - 3.1875rem); /* opera */
+	height: -webkit-calc(100vh - 3.1875rem); /* google, safari */
+	height: -moz-calc(100vh - 3.1875rem); /* firefox */
+	width: 100%;
+}
+
+
+/*
+  Definitions for homepage's table of contents.
+*/
+
+.homepage ul li {
+}
+
+#doc-content .googledoc {
+	/* 100% view port height - var(--header-height). But the CSS variable didn't work */
+	height: calc(100vh - 3.1875rem); /* CSS3 */
+	height: -o-calc(100vh - 3.1875rem); /* opera */
+	height: -webkit-calc(100vh - 3.1875rem); /* google, safari */
+	height: -moz-calc(100vh - 3.1875rem); /* firefox */
+	width: 100%;
+}
+
+
+/*
   Definitions for the responsive table component.
 */
 /* Table module CSS variable definitions. */
```

### Prism.js

- Download the original items from <prismjs.com>, then make some modifications.
- [Download link](https://prismjs.com/download.html#themes=prism-tomorrow&languages=clike+awk+bash+diff+docker+git+go+go-module+hcl+http+log+makefile+plant-uml+promql+regex+rego+rust+sql+vim+yaml&plugins=line-numbers+autolinker+autoloader+toolbar+copy-to-clipboard+download-button+match-braces+diff-highlight)
- Theme/color scheme: Tomorrow Night
- Plugins:
  - **LineNumbers**:
    Show line numbers.
  - **Autolinker**:
    Create automatic links for URL and emails.
    Also handle custom links using Markdown's link syntax.
  - **Autoloader**:
    Automatically load grammars for non-default languages.
    This is useful for languages not often shown in the blog.
  - **Toolbar**:
    Base plugin to allow multiple buttons (used below).
  - **Copy to Clipboard Button**:
    Add a `Copy` button to code snippet.
    This is added by default.
  - **Download button**:
    Add a `Download` button to code snippet. Usage is
	`<pre data-src="code.file" data-download-link data-download-link-label="Download"></pre>`.
  - **Match braces**:
    Highlight matching braches when hovering.
  - **DiffHighlight**:
    - Enable language-specific diff when using `<pre><code class="language-diff-LANGUAGE"></code></pre>`.
	- Also use background color to highlight diff using `<pre><code class="language-diff-LANGUAGE diff-highlight"></code></pre>`.
	- For similar effect to Github PR review UI, use both.
- Languages:
  - Non-core languages are supported via asynchronous autoload.
  - Core (TBD: update usable language name):
    - C-like:
	- AWK + GAWK: `awk`, `gawk`.
	- Bash + Shell + Shell: `sh`, `shell`, `bash`.
	- Diff: `diff`.
	- Docker: `dockerfile`, `docker`.
	- Git: `gitignore`.
	- Go: `go`, `golang`
	- Go module: `go-mod`, `go-module`.
	- HCL: `hcl`.
	- HTTP: `http`.
	- Log file: `log`.
	- Makefile:
	- PlantUML: `plantuml`
	- PromQL: `promql`.
	- Regex:
	- Rego: `rego`.
	- Rust: `rust`.
	- SQL: `sql`.
	- vim: `vim`.
	- YAML: `yaml`.
- After downloading both CSS and JS, patch the `prism.js`:
  - We do this so some of the paths that are hard-coded in the library
    can be changed (to clearer names).
  - After that, make changes based on this patch:

```diff
diff --git a/src-go/internal/theme/bitsgofer/assets/prism/prism.js b/src-go/internal/theme/bitsgofer/assets/prism/prism.js
index dfd3a08..0c63ce6 100644
--- a/src-go/internal/theme/bitsgofer/assets/prism/prism.js
+++ b/src-go/internal/theme/bitsgofer/assets/prism/prism.js
@@ -3445,7 +3445,7 @@ Prism.languages.vim = {
        var lang_data = {};

        var ignored_language = 'none';
-       var languages_path = 'components/';
+       var languages_path = 'prism/languages/';

        var script = Prism.util.currentScript();
        if (script) {
@@ -3460,10 +3460,10 @@ Prism.languages.vim = {
                        var src = script.src;
                        if (autoloaderFile.test(src)) {
                                // the script is the original autoloader script in the usual Prism project structure
-                               languages_path = src.replace(autoloaderFile, 'components/');
+                               languages_path = src.replace(autoloaderFile, 'prism/languages/');
                        } else if (prismFile.test(src)) {
                                // the script is part of a bundle like a custom prism.js from the download page
-                               languages_path = src.replace(prismFile, '$1components/');
+                               languages_path = src.replace(prismFile, '$1prism/languages/');
                        }
                }
        }
```
