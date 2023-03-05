package astra

import (
	"fmt"
	"path/filepath"

	"github.com/bitsgofer/notebook-go/internal/render"
	klog "k8s.io/klog/v2"
)

// luaFilters are filters used when rendering Markdown pages using this theme.
var luaFilters = []string{
	"pandoc/lua-filters/standard-code.lua",
}

func (th *bitsgoferTheme) Render(outputDir string, contentDir string) error {
	if err := th.renderIndexPage(outputDir, contentDir); err != nil {
		return err
	}
	if err := th.renderAboutPage(outputDir, contentDir); err != nil {
		return err
	}
	if err := th.renderProjectsPage(outputDir, contentDir); err != nil {
		return err
	}
	if err := th.renderArticles(outputDir, contentDir); err != nil {
		return err
	}
	if err := th.renderDrafts(outputDir, contentDir); err != nil {
		return err
	}

	return nil
}

func (th *bitsgoferTheme) renderIndexPage(outputDir, contentDir string) error {
	dst := filepath.Join(outputDir, "index.html")
	src := filepath.Join(contentDir, "README.md")

	if err := render.RenderSinglePage(
		dst,
		src,
		themeFS,
		"pandoc/index.template.html",
		luaFilters...,
	); err != nil {
		return fmt.Errorf("cannot render index page; err= %w", err)
	}

	klog.V(2).InfoS("rendered index page", "src", src, "dst", dst)
	return nil
}

func (th *bitsgoferTheme) renderAboutPage(outputDir, contentDir string) error {
	dst := filepath.Join(outputDir, "about.html")
	src := filepath.Join(contentDir, "ABOUT.md")

	if err := render.RenderSinglePage(
		dst,
		src,
		themeFS,
		"pandoc/about.template.html",
		luaFilters...,
	); err != nil {
		return fmt.Errorf("cannot render about page; err= %w", err)
	}

	klog.V(2).InfoS("rendered about page", "src", src, "dst", dst)
	return nil
}

func (th *bitsgoferTheme) renderProjectsPage(outputDir, contentDir string) error {
	dst := filepath.Join(outputDir, "projects.html")
	src := filepath.Join(contentDir, "PROJECTS.md")

	if err := render.RenderSinglePage(
		dst,
		src,
		themeFS,
		"pandoc/projects.template.html",
		luaFilters...,
	); err != nil {
		return fmt.Errorf("cannot render projects page; err= %w", err)
	}

	klog.V(2).InfoS("rendered projects page", "src", src, "dst", dst)
	return nil
}

func (th *bitsgoferTheme) renderArticles(outputDir, contentDir string) error {
	dstDir := filepath.Join(outputDir, "articles")
	srcDir := filepath.Join(contentDir, "articles")

	if err := render.RenderMultiplePages(
		dstDir,
		srcDir,
		themeFS,
		"pandoc/single.template.html",
		luaFilters...,
	); err != nil {
		return fmt.Errorf("cannot render articles; err= %w", err)
	}

	klog.V(2).InfoS("rendered articles", "srcDir", srcDir, "dstDir", dstDir)
	return nil
}

func (th *bitsgoferTheme) renderDrafts(outputDir, contentDir string) error {
	dstDir := filepath.Join(outputDir, "drafts")
	srcDir := filepath.Join(contentDir, "drafts")

	if err := render.RenderMultiplePages(
		dstDir,
		srcDir,
		themeFS,
		"pandoc/googledoc.template.html",
		luaFilters...,
	); err != nil {
		return fmt.Errorf("cannot render drafts; err= %w", err)
	}

	klog.V(2).InfoS("rendered drafts", "srcDir", srcDir, "dstDir", dstDir)
	return nil
}
