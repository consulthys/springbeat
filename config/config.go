// Config is put into a different package to prevent cyclic imports in case
// it is needed in several locations

package config

type Config struct {
	Springbeat SpringbeatConfig
}

type SpringbeatConfig struct {
	Period string `config:"period"`

	URLs []string

	Stats struct {
		Metrics  *bool
		Health	 *bool
	}
}
