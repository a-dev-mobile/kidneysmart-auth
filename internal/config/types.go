// types.go
package config

type Environment string

const (
    Dev  Environment = "dev"
    Prod Environment = "prod"
)

type RotationPolicy string

const (
    Monthly RotationPolicy = "monthly"
    Weekly  RotationPolicy = "weekly"
    Daily   RotationPolicy = "daily"
)