package routes

import (
	"pg_sandbox/controllers/auth"
	"pg_sandbox/controllers/collection"
	"pg_sandbox/controllers/disbursement"
	hostedcheckout "pg_sandbox/controllers/hosted_checkout"
	"pg_sandbox/controllers/transactions"
	"pg_sandbox/controllers/users"
	"pg_sandbox/middleware"

	"github.com/gin-gonic/gin"
)

func SetupRoutes(r *gin.Engine) {
	// Static file routes
	r.Static("/v1/images", "./Files/Business")
	r.Static("/v1/files", "./Files/Business")

	// Dynamic routes initialization
	//dynamiccontrollers.LoadDynamicRoutes(r)

	// Authorized routes group
	au := r.Group("/v1")
	au.Use(middleware.AuthorizationMiddleWare())

	// Auth routes
	r.POST("/v1/auth/register", auth.RegisterHandler)
	r.POST("/v1/auth/login", auth.LoginHandler)

	// User routes
	au.POST("/users/get", auth.GetUsersHandler)

	r.POST("/v1/make-collection", collection.MakeCollectionHandler)
	r.POST("/v1/make-disbursement", disbursement.MakeDisbursementHandler)
	r.POST("/v1/token-generate", auth.AuthorizationHandler)

	au.POST("/secret/generate", auth.GenerateSecretHandler)
	au.POST("/signature/generate", auth.GenerateOAuthSignatureHandler)
	au.POST("/pin/create", auth.SetPinHandler)
	au.GET("/user/credentials/get/:id", auth.GetAPICredentialsHandler)
	r.GET("/v1/mobile-money/check-status/:id", transactions.TransactionQueryHandler)
	r.GET("/v1/mobile-money/name-lookup/:number", users.NameLookUpHandler)
	r.POST("/v1/mobile-money/disburse", disbursement.MakeDisbursementHandler)
	r.GET("/v1/mobile-money/disburse/status/:reference", disbursement.QueryDisbursementHandler)
	au.POST("/float/update", users.SetFloatBalanceHander)
	r.POST("/v1/checkout/session", hostedcheckout.HostedCheckOutHandler)
	r.GET("/v1/checkout/get/:id", hostedcheckout.GetHostedCheckoutDetailsHandler)
	r.POST("/callback", users.CallbackHandler)
	r.POST("/v1/checkout/respond", hostedcheckout.HostedCheckoutResponseHandler)
}
