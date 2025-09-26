package web

import (
	"embed"
)

const StaticBaseDirName = "static"

//go:embed static/*
var static embed.FS

const TemplateBaseDirName = "tmpl"

//go:embed tmpl/*.html
var tmpl embed.FS

// GetStaticFS returns static assets as an [embed.FS]
func GetStaticFS() embed.FS {
	return static
}

// GetTmplFS returns template assets as an [embed.FS]
func GetTmplFS() embed.FS {
	return tmpl
}
