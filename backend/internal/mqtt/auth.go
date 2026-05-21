package mqtt

import "fmt"

type AuthConfig struct {
	Mode      string
	Username  string
	Password  string
	PublicKey string
	Token     string
}

func (a AuthConfig) Validate() error {
	switch a.Mode {
	case "", "none":
		return nil
	case "subscriber":
		if a.Username == "" || a.Password == "" {
			return fmt.Errorf("subscriber auth requires MQTT_USERNAME and MQTT_PASSWORD")
		}
		return nil
	case "jwt":
		if a.PublicKey == "" || a.Token == "" {
			return fmt.Errorf("jwt auth requires public key and generated token")
		}
		return nil
	default:
		return fmt.Errorf("unsupported auth mode %q", a.Mode)
	}
}
