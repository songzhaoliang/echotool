package echotool

import (
	"os"
	"path/filepath"
	"strings"

	"github.com/popeyeio/handy"
)

const (
	runtimeEnv    = "_ENV_"
	runtimeRegion = "_REGION_"
	runtimeTag    = "_TAG_"

	EnvDev     = "dev"
	EnvTest    = "test"
	EnvProduct = "prod"

	confPath      = "/conf"
	regionDefault = "cn"
)

func IsProduct() bool {
	return Env() == EnvProduct
}

func IsTest() bool {
	return Env() == EnvTest
}

func GetConfDir() string {
	dir, _ := os.Getwd()
	return dir + confPath
}

func Env() (env string) {
	switch os.Getenv(runtimeEnv) {
	case EnvTest:
		env = EnvTest
	case EnvProduct:
		env = EnvProduct
	default:
		env = EnvDev
	}
	return
}

func Reg() string {
	reg := regionDefault
	if r := os.Getenv(runtimeRegion); !handy.IsEmptyStr(r) {
		reg = strings.ToLower(r)
	}
	return reg
}

func Tag() string {
	return strings.ToLower(os.Getenv(runtimeTag))
}

func ParseByEnv(suffix string) string {
	return filepath.Join(Env(), suffix)
}

func ParseByReg(suffix string) string {
	return filepath.Join(Env(), Reg(), suffix)
}

func ParseByTag(suffix string) string {
	return filepath.Join(Env(), Tag(), suffix)
}
