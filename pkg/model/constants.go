package model

const (
	ModeRun   = "run"
	ModeBuild = "build"
	ModeWatch = "watch"
)

const (
	AdapterViteYarn  = "vite-yarn"
	AdapterTauriYarn = "tauri-yarn"
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
