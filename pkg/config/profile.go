package config

type Profile struct {
	raw  interface{}
	name string
}

func newProfile(data interface{}, name string) Profile {
	return Profile{data, name}
}

func (p Profile) Name() string {
	return p.name
}

func (p Profile) ViteProfile() ViteOptions {
	return p.raw.(ViteOptions)
}

func (p Profile) DotNetProfile() DotNetOptions {
	return p.raw.(DotNetOptions)
}

func (p Profile) BaseOptions() BaseOptions {
	return p.raw.(BaseOptions)
}

func (p Profile) BaseProfile() BaseProfile {
	return p.raw.(BaseProfile)
}

func (p Profile) Pre() HookOptions {
	return p.raw.(map[string]interface{})[KeyPre].(HookOptions)
}

func (p Profile) Post() HookOptions {
	return p.raw.(map[string]interface{})[KeyPost].(HookOptions)
}
