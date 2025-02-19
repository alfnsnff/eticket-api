package utils

// Get config path for local or docker
func GetConfEnv(configPath string) string {
	switch configPath {
	case "development":
		return "./config/development-conf.yaml"
	case "staging":
		return "./config/staging-conf.yaml"
	case "production":
		return "./config/production-conf.yaml"
	default:
		return "./config/development-conf.yaml"
	}
}
