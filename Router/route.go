package router

import (
	BookingController "Redbus_backend/Controllers/Booking"
	BusController "Redbus_backend/Controllers/Bus"
	OfferController "Redbus_backend/Controllers/Offer"
	ReviewController "Redbus_backend/Controllers/Review"
	TicketController "Redbus_backend/Controllers/Ticket"
	UserController "Redbus_backend/Controllers/User"
	WalletController "Redbus_backend/Controllers/Wallet"
	auth "Redbus_backend/Helpers/Auth"
	middlewares "Redbus_backend/Helpers/Middlewares"
	"net/http"

	"github.com/gorilla/mux"
)

func GetRouter() *mux.Router {
	apiRouter := mux.NewRouter()
	userHandler(apiRouter)
	busHandler(apiRouter)
	bookingHandler(apiRouter)
	walletHandler(apiRouter)
	offerHandler(apiRouter)
	reviewHandler(apiRouter)
	ticketHandler(apiRouter)
	return apiRouter
}

func userHandler(router *mux.Router) {
	router.HandleFunc("/register", UserController.RegisterUser).Methods("POST", "OPTIONS")
	router.HandleFunc("/login", UserController.LoginUser).Methods("POST", "OPTIONS")
	router.HandleFunc("/users/{userID}", UserController.GetUserByID).Methods("GET")
	router.Handle("/logout", auth.VerifyJWT(http.HandlerFunc(UserController.Logout))).Methods("GET")
	router.HandleFunc("/forgotpassword", UserController.VerifyEmailAndSendKey).Methods("POST", "OPTIONS")
	router.HandleFunc("/resetpassword", UserController.ResetPassword).Methods("POST", "OPTIONS")
}

func busHandler(router *mux.Router) {
	router.Handle("/addbus", auth.VerifyAdmin(http.HandlerFunc(BusController.AddBus))).Methods("POST", "OPTIONS")
	router.HandleFunc("/buses/{busID}", BusController.GetBusByID).Methods("GET")
	router.HandleFunc("/buses/search/{page}", BusController.GetSearchedBus).Methods("GET")
}

func bookingHandler(router *mux.Router) {
	router.Handle("/bookseat", auth.VerifyJWT(http.HandlerFunc(BookingController.BookSeat))).Methods("POST", "OPTIONS")
	router.HandleFunc("/viewseats/{busID}", middlewares.FetchCommonData(BookingController.ViewSeats)).Methods("GET")
	router.HandleFunc("/selectseat/{seatID}/{busID}", middlewares.FetchCommonData(BookingController.SelectSeat)).Methods("GET")
}

func walletHandler(router *mux.Router) {
	router.Handle("/addMoney/{userID}", auth.VerifyJWT(http.HandlerFunc(WalletController.AddToWallet))).Methods("POST", "OPTIONS")
	router.Handle("/withdrawMoney/{userID}", auth.VerifyJWT(http.HandlerFunc(WalletController.WithdrawFromWallet))).Methods("POST", "OPTIONS")
	router.Handle("/getWalletBalance/{userID}", auth.VerifyJWT(http.HandlerFunc(WalletController.GetWalletBalance))).Methods("GET")
}

func offerHandler(router *mux.Router) {
	router.Handle("/addOffer", auth.VerifyAdmin(http.HandlerFunc(OfferController.AddOffer))).Methods("POST", "OPTIONS")
	router.HandleFunc("/getOffers", OfferController.GetOffers).Methods("GET")
	router.HandleFunc("/applyOffer", OfferController.ApplyOffer).Methods("POST", "OPTIONS")
}

func reviewHandler(router *mux.Router) {
	router.Handle("/addReview/{customerID}/{busID}", auth.VerifyJWT(http.HandlerFunc(ReviewController.AddReview))).Methods("POST", "OPTIONS")
	router.HandleFunc("/getReviews/{busID}", ReviewController.GetReviewsByBusID).Methods("GET")
}

func ticketHandler(router *mux.Router) {
	router.HandleFunc("/cancelTicket/{ticketID}", TicketController.CancelTicket).Methods("DELETE")
	router.HandleFunc("/getTickets/{userID}", TicketController.GetTicketByUserID).Methods("GET")
}
