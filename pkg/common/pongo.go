package common

import (
	"github.com/cnmade/bsmi-mail-kernel/pkg/version"
	"github.com/flosch/pongo2/v4"
)
import "github.com/imdario/mergo"

func Pongo2ContextAppend(ctx pongo2.Context, ctxAddition pongo2.Context) pongo2.Context {
	currentMethod := "Pongo2ContextAppend"
	err := mergo.Map(&ctx, ctxAddition);

	if err != nil {
		Sugar.Errorf(currentMethod + " merge pnogo2 context err: %v", err)
	}
	return ctx
}

func Pongo2ContextWithVersion(ctx pongo2.Context) pongo2.Context {
	outCtx := Pongo2ContextAppend(ctx, pongo2.Context{
		"BuildTag": version.BuildTag,
		"BuildNum": version.BuildNum,
		"BsmiKbVersion": BsmiKbVersion,
		"SiteKey": Config.HCaptchaSiteKey,
		"CaptchaEnabled": Config.CaptchaEnabled,
	})

	if Config.TongjiConfig.TongjiEnabled == 1 {
		outCtx = Pongo2ContextAppend(outCtx, pongo2.Context{
			"TongjiCode": Config.TongjiConfig.TongjiCode,
		})
	}
	return outCtx
}
