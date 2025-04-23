package routes

import (
	"pg_sandbox/controllers/auth"
	"pg_sandbox/controllers/collection"
	"pg_sandbox/controllers/dashboard"
	"pg_sandbox/controllers/disbursement"
	hostedcheckout "pg_sandbox/controllers/hosted_checkout"
	"pg_sandbox/controllers/mail"

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

	//emails
	r.POST("/v1/request/code", mail.SendMailHandler)
	r.POST("/v1/verify/code", mail.VerifyCode)

	// Authorized routes group
	au := r.Group("/v1")
	au.Use(middleware.AuthorizationMiddleWare())

	// Auth routes
	r.POST("/v1/auth/register", auth.RegisterHandler)
	r.POST("/v1/auth/login", auth.LoginHandler)

	// User routes
	au.POST("/users/get", auth.GetUsersHandler)
	au.POST("/user/edit", auth.EditUserHandler)
	au.POST("/merchants/get", auth.GetMerchantsHandler)
	au.GET("/user/get/:id", users.GetUserProfileHandler)

	r.POST("/v1/mobile-money/collect", collection.MakeCollectionHandler)
	r.POST("/v1/make-disbursement", disbursement.MakeDisbursementHandler)
	r.POST("/v1/oauth/token", auth.AuthorizationHandler)

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
	r.POST("/v1/checkout/respond/:condition", hostedcheckout.HostedCheckoutResponseHandler)
	r.POST("/v1/reset/password", users.ResetPasswordHandler)
	r.GET("/v1/mobile-money/disburse/balance", disbursement.CheckDisbursementBalance)

	//dashboard
	au.GET("/overview/cards", dashboard.GetOverviewCardsInfoHandler)
	au.POST("/overview/activity", dashboard.GetActivitiesHandler)

	//dashboard-users-tab
	au.GET("/dashboard/users/info", dashboard.GetUserStatisticsHandler)
	au.POST("/dashboard/users", dashboard.GetUsersHandler)

	//dashboard-merchants-tab
	au.GET("/dashboard/merchants/info", dashboard.GetMerchantStatisticsHandler)
	au.POST("/dashboard/merchants", dashboard.GetMerchantsHandler)
	au.GET("/dashboard/merchants/top", dashboard.GetTopMerchantsHandler)

	//dashboard-requests-tab
	au.GET("/dashboard/requests/info", dashboard.GetAPIRequestsInfoHandler)
	au.POST("/dashboard/requests", dashboard.GetAPIRequestsHandler)
	au.GET("/dashboard/request/response", dashboard.GetAPIResponeTimeHandler)

	//dashboard-transactions-tab
	au.GET("/dashboard/transactions/info", dashboard.GetTransactionInfoHandler)
	au.POST("/dashboard/transactions", dashboard.GetTransactionsHandler)
	au.GET("/dashboard/transactions/channels", dashboard.GetTransactionsChannelHandler)

}
