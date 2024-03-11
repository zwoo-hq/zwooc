package model

const (
	ModeRun   = "run"
	ModeBuild = "build"
	ModeWatch = "watch"
)

const (
	AdapterViteYarn  = "vite-yarn"
	AdapterViteNpm   = "vite-npm"
	AdapterVitePnpm  = "vite-pnpm"
	AdapterTauriYarn = "tauri-yarn"
	AdapterTauriNpm  = "tauri-npm"
	AdapterTauriPnpm = "tauri-pnpm"
	AdapterDotnet    = "dotnet"
)

const (
	KeyDefault   = "$default"
	KeyAdapter   = "$adapter"
	KeyDirectory = "$dir"
	KeyFragment  = "$fragments"
	KeyCompound  = "$compounds"
	KeyPre       = "$pre"
	KeyPost      = "$post"
)
