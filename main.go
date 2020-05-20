package main

import (
	"fmt"
	"log"
	"os"

	"github.com/cruisechang/dbServer/handler"
	"github.com/cruisechang/dbServer/middleware"
	"github.com/cruisechang/dbServer/util"
	"github.com/cruisechang/dbex"
)

func main() {

	dbx, err := dbex.NewDBEX("dbexConfig.json")
	if err != nil {
		exit(1, fmt.Errorf("dbex.NewDBEX() error=%s", err.Error()))
	}

	uniqueIDProvider, err := util.CreateUniqueIDProvider()
	if err != nil {
		exit(2, fmt.Errorf("CreateUniqueIDProvider() error=%s", err.Error()))
	}

	//middleware
	mw := middleware.NewMiddleware(dbx.Logger)

	//http server middleware
	router := dbx.HTTPServer.GetRouter()
	//router.Use(mw.LogRequestURI)
	router.Use(mw.CheckHead)

	//router

	router.Handle("/rooms/{id:[0-9]+}/historyResult", handler.NewRoomIDHandler(handler.NewBaseHandler(dbx.DB, dbx.Logger, uniqueIDProvider))).Methods( "PATCH")
	router.Handle("/rooms/{id:[0-9}+}/name", handler.NewRoomIDHandler(handler.NewBaseHandler(dbx.DB, dbx.Logger, uniqueIDProvider))).Methods("PATCH")
	router.Handle("/rooms/{id:[0-9}+}/active", handler.NewRoomIDHandler(handler.NewBaseHandler(dbx.DB, dbx.Logger, uniqueIDProvider))).Methods("PATCH")
	router.Handle("/rooms/{id:[0-9}+}/hlsURL", handler.NewRoomIDHandler(handler.NewBaseHandler(dbx.DB, dbx.Logger, uniqueIDProvider))).Methods("PATCH")
	router.Handle("/rooms/{id:[0-9}+}/boot", handler.NewRoomIDHandler(handler.NewBaseHandler(dbx.DB, dbx.Logger, uniqueIDProvider))).Methods("PATCH")
	router.Handle("/rooms/{id:[0-9}+}/round", handler.NewRoomIDHandler(handler.NewBaseHandler(dbx.DB, dbx.Logger, uniqueIDProvider))).Methods("PATCH")
	router.Handle("/rooms/{id:[0-9}+}/status", handler.NewRoomIDHandler(handler.NewBaseHandler(dbx.DB, dbx.Logger, uniqueIDProvider))).Methods("PATCH")
	router.Handle("/rooms/{id:[0-9}+}/betCountdown", handler.NewRoomIDHandler(handler.NewBaseHandler(dbx.DB, dbx.Logger, uniqueIDProvider))).Methods("PATCH")
	router.Handle("/rooms/{id:[0-9}+}/dealerID", handler.NewRoomIDHandler(handler.NewBaseHandler(dbx.DB, dbx.Logger, uniqueIDProvider))).Methods("PATCH")
	router.Handle("/rooms/{id:[0-9}+}", handler.NewRoomIDHandler(handler.NewBaseHandler(dbx.DB, dbx.Logger, uniqueIDProvider))).Methods("GET", "DELETE", "PATCH")
	router.Handle("/rooms", handler.NewRoomsHandler(handler.NewBaseHandler(dbx.DB, dbx.Logger, uniqueIDProvider))).Methods("GET", "POST")



	router.Handle("/users/{account}/password", handler.NewUserAccountHandler(handler.NewBaseHandler(dbx.DB, dbx.Logger, uniqueIDProvider))).Methods("GET")
	router.Handle("/users/{account}/id", handler.NewUserAccountHandler(handler.NewBaseHandler(dbx.DB, dbx.Logger, uniqueIDProvider))).Methods("GET")
	router.Handle("/users/{account}/accessToken", handler.NewUserAccountHandler(handler.NewBaseHandler(dbx.DB, dbx.Logger, uniqueIDProvider))).Methods("GET")
	router.Handle("/users/{accessToken}/tokenData", handler.NewUserAccessTokenHandler(handler.NewBaseHandler(dbx.DB, dbx.Logger, uniqueIDProvider))).Methods("GET")
	router.Handle("/users/{id:[0-9]+}/credit", handler.NewUserIDHandler(handler.NewBaseHandler(dbx.DB, dbx.Logger, uniqueIDProvider))).Methods("GET", "PATCH")
	//router.Handle("/users/{id:[0-9]+}/login", handler.NewUserIDHandler(handler.NewBaseHandler(dbx.DB, dbx.Logger, uniqueIDProvider))).Methods("GET", "PATCH")
	router.Handle("/users/{id:[0-9]+}/active", handler.NewUserIDHandler(handler.NewBaseHandler(dbx.DB, dbx.Logger, uniqueIDProvider))).Methods("GET", "PATCH")
	router.Handle("/users/{id:[0-9]+}", handler.NewUserIDHandler(handler.NewBaseHandler(dbx.DB, dbx.Logger, uniqueIDProvider))).Methods("GET")
	router.Handle("/users", handler.NewUsersHandler(handler.NewBaseHandler(dbx.DB, dbx.Logger, uniqueIDProvider))).Methods("GET", "POST")

	router.Handle("/partners/{account}/password", handler.NewPartnerAccountHandler(handler.NewBaseHandler(dbx.DB, dbx.Logger, uniqueIDProvider))).Methods("GET")
	router.Handle("/partners/{account}/id", handler.NewPartnerAccountHandler(handler.NewBaseHandler(dbx.DB, dbx.Logger, uniqueIDProvider))).Methods("GET")
	router.Handle("/partners/{account}/login", handler.NewPartnerAccountHandler(handler.NewBaseHandler(dbx.DB, dbx.Logger, uniqueIDProvider))).Methods("GET")
	router.Handle("/partners/{id:[0-9]+}/aesKey", handler.NewPartnerIDHandler(handler.NewBaseHandler(dbx.DB, dbx.Logger, uniqueIDProvider))).Methods("GET", "PATCH")
	router.Handle("/partners/{id:[0-9]+}/active", handler.NewPartnerIDHandler(handler.NewBaseHandler(dbx.DB, dbx.Logger, uniqueIDProvider))).Methods("GET", "PATCH")
	router.Handle("/partners/{id:[0-9]+}/apiBindIP", handler.NewPartnerIDHandler(handler.NewBaseHandler(dbx.DB, dbx.Logger, uniqueIDProvider))).Methods("GET", "PATCH")
	router.Handle("/partners/{id:[0-9]+}/cmdBindIP", handler.NewPartnerIDHandler(handler.NewBaseHandler(dbx.DB, dbx.Logger, uniqueIDProvider))).Methods("GET", "PATCH")
	router.Handle("/partners/{id:[0-9]+}/login", handler.NewPartnerIDHandler(handler.NewBaseHandler(dbx.DB, dbx.Logger, uniqueIDProvider))).Methods("PATCH")
	router.Handle("/partners/{id:[0-9]+}", handler.NewPartnerIDHandler(handler.NewBaseHandler(dbx.DB, dbx.Logger, uniqueIDProvider))).Methods("GET")
	router.Handle("/partners", handler.NewPartnersHandler(handler.NewBaseHandler(dbx.DB, dbx.Logger, uniqueIDProvider))).Methods("GET", "POST")

	//router.Handle("/halls/{id:[0-9]+}/name", handler.NewHallIDHandler(handler.NewBaseHandler(dbx.DB, dbx.Logger,uniqueIDProvider))).Methods("PATCH")
	//router.Handle("/halls/{id:[0-9]+}/active", handler.NewHallIDHandler(handler.NewBaseHandler(dbx.DB, dbx.Logger,uniqueIDProvider))).Methods("PATCH")
	router.Handle("/halls/{id:[0-9]+}", handler.NewHallIDHandler(handler.NewBaseHandler(dbx.DB, dbx.Logger, uniqueIDProvider))).Methods("GET", "DELETE", "PATCH")
	router.Handle("/halls", handler.NewHallsHandler(handler.NewBaseHandler(dbx.DB, dbx.Logger, uniqueIDProvider))).Methods("GET", "POST")


	router.Handle("/limitations", handler.NewLimitationsHandler(handler.NewBaseHandler(dbx.DB, dbx.Logger, uniqueIDProvider))).Methods("GET")

	//router.Handle("/dealers/{account}/password", handler.NewDealerIDHandler(handler.NewBaseHandler(dbx.DB, dbx.Logger,uniqueIDProvider))).Methods("GET")
	router.Handle("/dealers/{account}/login", handler.NewDealerIDHandler(handler.NewBaseHandler(dbx.DB, dbx.Logger, uniqueIDProvider))).Methods("GET")
	router.Handle("/dealers/{account}/id", handler.NewDealerIDHandler(handler.NewBaseHandler(dbx.DB, dbx.Logger, uniqueIDProvider))).Methods("GET")
	//router.Handle("/dealers/{id:[0-9]+}/active", handler.NewDealerIDHandler(handler.NewBaseHandler(dbx.DB, dbx.Logger,uniqueIDProvider))).Methods("PATCH")
	//router.Handle("/dealers/{id:[0-9]+}/portraitURL", handler.NewDealerIDHandler(handler.NewBaseHandler(dbx.DB, dbx.Logger,uniqueIDProvider))).Methods("PATCH")
	router.Handle("/dealers/{id:[0-9]+}", handler.NewDealerIDHandler(handler.NewBaseHandler(dbx.DB, dbx.Logger, uniqueIDProvider))).Methods("GET", "DELETE", "PATCH")
	router.Handle("/dealers", handler.NewDealersHandler(handler.NewBaseHandler(dbx.DB, dbx.Logger, uniqueIDProvider))).Methods("GET", "POST")

	router.Handle("/transfers/ptID/{partnerTransferID}", handler.NewTransferPartnerTransferIDHandler(handler.NewBaseHandler(dbx.DB, dbx.Logger, uniqueIDProvider))).Methods("GET")
	router.Handle("/transfers/{id:[0-9}+}/status", handler.NewTransferIDHandler(handler.NewBaseHandler(dbx.DB, dbx.Logger, uniqueIDProvider))).Methods("PATCH")
	router.Handle("/transfers/{id:[0-9]+}", handler.NewTransferIDHandler(handler.NewBaseHandler(dbx.DB, dbx.Logger, uniqueIDProvider))).Methods("GET")
	router.Handle("/transfers", handler.NewTransfersHandler(handler.NewBaseHandler(dbx.DB, dbx.Logger, uniqueIDProvider))).Methods("GET", "POST")

	router.Handle("/users/{id:[0-9]+}/log", handler.NewUserIDHandler(handler.NewBaseHandler(dbx.DB, dbx.Logger, uniqueIDProvider))).Methods("GET")
	router.Handle("/partners/{id:[0-9]+}/log", handler.NewPartnerIDHandler(handler.NewBaseHandler(dbx.DB, dbx.Logger, uniqueIDProvider))).Methods("GET")

	router.Handle("/bets/{id:[0-9]+}/status", handler.NewBetIDHandler(handler.NewBaseHandler(dbx.DB, dbx.Logger, uniqueIDProvider))).Methods("PATCH")
	router.Handle("/bets/{id:[0-9]+}", handler.NewBetIDHandler(handler.NewBaseHandler(dbx.DB, dbx.Logger, uniqueIDProvider))).Methods("GET")
	router.Handle("/bets", handler.NewBetsHandler(handler.NewBaseHandler(dbx.DB, dbx.Logger, uniqueIDProvider))).Methods("GET", "POST")

	router.Handle("/rounds/{id:[0-9]+}/patch", handler.NewRoundIDHandler(handler.NewBaseHandler(dbx.DB, dbx.Logger, uniqueIDProvider))).Methods("PATCH")
	router.Handle("/rounds/{id:[0-9]+}", handler.NewRoundIDHandler(handler.NewBaseHandler(dbx.DB, dbx.Logger, uniqueIDProvider))).Methods("GET")
	router.Handle("/rounds", handler.NewRoundsHandler(handler.NewBaseHandler(dbx.DB, dbx.Logger, uniqueIDProvider))).Methods("GET", "POST")

	router.Handle("/officialCMSManagers/{account}/login", handler.NewOfficialCMSManagerAccountHandler(handler.NewBaseHandler(dbx.DB, dbx.Logger, uniqueIDProvider))).Methods("GET")
	router.Handle("/officialCMSManagers/{id:[0-9]+}", handler.NewOfficialCMSManagerIDHandler(handler.NewBaseHandler(dbx.DB, dbx.Logger, uniqueIDProvider))).Methods("GET", "DELETE", "PATCH")
	router.Handle("/officialCMSManagers", handler.NewOfficialCMSManagersHandler(handler.NewBaseHandler(dbx.DB, dbx.Logger, uniqueIDProvider))).Methods("GET", "POST")

	router.Handle("/officialCMSRoles/{id:[0-9]+}", handler.NewOfficialCMSRoleIDHandler(handler.NewBaseHandler(dbx.DB, dbx.Logger, uniqueIDProvider))).Methods("GET", "PATCH")
	router.Handle("/officialCMSRoles", handler.NewOfficialCMSRolesHandler(handler.NewBaseHandler(dbx.DB, dbx.Logger, uniqueIDProvider))).Methods("GET", "POST")

	router.Handle("/broadcasts/{id:[0-9]+}", handler.NewBroadcastIDHandler(handler.NewBaseHandler(dbx.DB, dbx.Logger, uniqueIDProvider))).Methods("DELETE", "PATCH")
	router.Handle("/broadcasts", handler.NewBroadcastsHandler(handler.NewBaseHandler(dbx.DB, dbx.Logger, uniqueIDProvider))).Methods("GET", "POST")

	router.Handle("/banners/{id:[0-9]+}", handler.NewBannerIDHandler(handler.NewBaseHandler(dbx.DB, dbx.Logger, uniqueIDProvider))).Methods("DELETE", "PATCH")
	router.Handle("/banners", handler.NewBannersHandler(handler.NewBaseHandler(dbx.DB, dbx.Logger, uniqueIDProvider))).Methods("GET", "POST")

	router.Handle("/", handler.NewNotFoundHandler(handler.NewBaseHandler(dbx.DB, dbx.Logger, uniqueIDProvider)))

	//sub := router.Host("{subdomain}.domain.com").Subrouter()
	//sub := router.Host("localhost").Subrouter()
	//sub.Path("/users").HandlerFunc(handler.UserHandler).Methods("POST").Name("users")

	if err := dbx.HTTPServer.Start(); err != nil {
		exit(3, fmt.Errorf("StartHTTPServer error:%s", err.Error()))
	}

	//run forever (主程式沒有此功能才需要)
	dbx.HTTPServer.Hold()
}

func exit(id int, err error) {
	fmt.Println(err)
	log.Println(err)
	os.Exit(id)
}
