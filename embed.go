package main

import "embed"

var (
	//go:embed templates
	TemplateFS embed.FS
	//go:embed sql
	SQLFS embed.FS
	//go:embed static
	StaticFS embed.FS
	useEmbed bool
)
