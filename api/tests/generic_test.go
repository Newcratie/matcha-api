package tests

import (
	"errors"
	"github.com/Newcratie/matcha-api/api"
	"github.com/Newcratie/matcha-api/api/logprint"
	"testing"
)

func TestLogPrint(t *testing.T) {
	logprint.Title("Title func")
	logprint.Title("Title function longer")
	logprint.Title("Title function much much much longer, v1.3.5.24")
	logprint.Centered("centered")
	logprint.Error(errors.New("Error Sample"))
	logprint.End()
}

func TestExpr(t *testing.T) {
}
