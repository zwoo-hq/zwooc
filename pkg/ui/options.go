package ui

type ViewOptions struct {
	DisableTUI    bool
	QuiteMode     bool
	InlineOutput  bool
	CombineOutput bool
	DisablePrefix bool
	// todo: move out of here
	MaxConcurrency int
}
