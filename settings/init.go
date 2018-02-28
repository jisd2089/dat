package settings

const (
	defaultSettingsPath = "settings.yaml"
)

func Initialize(settingsPath string) error {
	if settingsPath == "" {
		settingsPath = defaultSettingsPath
	}

	se, err := CreateSettingsFromYAML(settingsPath)
	if err != nil {
		return err
	}

	defaultSettings = se
	return nil
}
