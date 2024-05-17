package app

//
//import (
//	"github.com/go-chi/chi/v5"
//	"github.com/pttrulez/investor-go/config"
//	"github.com/pttrulez/investor-go/internal/controller"
//	"github.com/pttrulez/investor-go/internal/controller/http_controllers"
//	"github.com/pttrulez/investor-go/internal/service"
//)
//
//func initControllers(cfg config.Config) {
//	c := NewControllers()
//	r := chi.NewRouter()
//
//	//                                   Public routes
//
//	//                               		Protected routes
//
//}
//
//type Controllers struct {
//	Auth controller.AuthController
//	//Deposit   controller.DepositController
//	//Cashout   controller.CashoutController
//	//Deal      controller.DealController
//	//Expert    controller.ExpertController
//	//MoexBond  controller.MoexBondController
//	//MoexShare controller.MoexShareController
//	//Portfolio controller.PortfolioController
//}
//
//func NewControllers() *Controllers {
//	userService := service.NewUserService()
//	return &Controllers{
//		Auth: http_controllers.NewAuthController(repo, services),
//		//Cashout:   http_controllers.NewCashoutController(repo, services),
//		//Deal:      http_controllers.NewDealController(repo, services),
//		//Deposit:   http_controllers.NewDepositController(repo, services),
//		//Expert:    http_controllers.NewExpertController(repo, services),
//		//MoexBond:  http_controllers.NewMoexBondController(repo, services),
//		//MoexShare: http_controllers.NewMoexShareController(repo, services),
//		//Portfolio: http_controllers.NewPortfolioController(repo, services),
//	}
//}
