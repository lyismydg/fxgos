package router

import (
	"net/http"
	"strconv"

	"github.com/fidelfly/fxgo"
	"github.com/fidelfly/fxgo/httprxr"
	"github.com/fidelfly/fxgo/logx"
	"github.com/fidelfly/fxgo/routex"

	"github.com/fidelfly/fxgos/cmd/api"
	"github.com/fidelfly/fxgos/cmd/service/iam"
	"github.com/fidelfly/fxgos/cmd/utilities/auth"
	"github.com/fidelfly/fxgos/cmd/utilities/system"
)

func setupRoute(rr *fxgo.RootRouter) {
	setupAPI(rr.PathPrefix("/api").Subrouter())
}

func setupAPI(router *routex.Router) {
	router.Register(
		api.ProgressRoute,
		api.UserRoute,
		api.FileRoute,
		api.QueryRoute,
		api.RoleRoute,
	)
}

var rr *fxgo.RootRouter

func GetRootRouter() *fxgo.RootRouter {
	return rr
}

func Initialize() (err error) {
	ti, err := auth.SetupTokenIssuer(api.Token, system.OAuth2.Client...)
	if err != nil {
		return
	}
	rr = fxgo.NewRouter()

	rr.EnableAudit(logx.StandardLogger())
	rr.EnableRecover()

	rr.AttachPlugins(ti)
	//rr.EnableAuthFilter(ti.AuthFilter)

	rr.Use(AccessMiddleware)

	setupRoute(rr)

	return nil
}

func AccessMiddleware(handler http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if userKey := api.GetUserKey(r); len(userKey) > 0 {
			if userId, err := strconv.ParseInt(userKey, 10, 64); err == nil {
				if config, ok := rr.CurrentRouteConfig(r); ok {
					if ac := config.GetProps(api.AccessConfigKey); ac != nil {
						if ap, ok := ac.(iam.AccessPremise); ok {
							ok, err := iam.ValidateAccess(userId, ap)
							if err != nil || !ok {
								httprxr.ResponseJSON(w, http.StatusForbidden, nil)
								return
							}
						}
					}
				}
			}
		}

		handler.ServeHTTP(w, r)
	})
}