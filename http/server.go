package http

import (
	"github.com/valyala/fasthttp"
	"net/http"
	log "github.com/sirupsen/logrus"
	"github.com/didip/tollbooth"
	"time"
	"github.com/didip/tollbooth_fasthttp"
	"github.com/buaazp/fasthttprouter"
)

func StartServer(addr string) error {

	fastRouter := fasthttprouter.New()

	fastRouter.GET("/search", cors(SearchHandler))
	fastRouter.POST("/search", cors(SearchHandler))

	fastRouter.GET("/nearby",cors(NearbyEventsHandler))

	fastRouter.MethodNotAllowed = func(ctx *fasthttp.RequestCtx) {
		ctx.Error(http.StatusText(http.StatusMethodNotAllowed), http.StatusMethodNotAllowed)
	}

	fastRouter.NotFound = func(ctx *fasthttp.RequestCtx) {
		ctx.Error(http.StatusText(http.StatusNotFound), http.StatusNotFound)
	}

	fastRouter.PanicHandler = func(ctx *fasthttp.RequestCtx, i interface{}) {
		log.Errorln("something broke")
		log.Errorln("error: ", i)
		ctx.Error(http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
	}

	fastRouter.POST("/add-wishlist", cors(AddToWishlistHandler))
	fastRouter.POST("/remove-wishlist", cors(RemoveFromWishlistHandler))
	fastRouter.GET("/get-wishlist", cors(GetWishlistHandler))
	fastRouter.POST("/stats", cors(StatsHandler))
	fastRouter.POST("/get-sales", cors(GetSalesHandler))
	fastRouter.GET("/featured-events", cors(GetFeaturedEventsHandler))
	fastRouter.GET("/general-report", cors(GeneralReportByUserIdHandler))
	fastRouter.GET("/event-report", cors(GetReportByEventIdHandler))
	fastRouter.POST("/publish-event",cors(PublishEventHandler))
	fastRouter.GET("/ping", cors(func(ctx *fasthttp.RequestCtx) {
		ctx.WriteString("PONG\n")
	}))

	return newServer(fastRouter).ListenAndServe(addr)
}

func newServer(fastRouter *fasthttprouter.Router) *fasthttp.Server {

	limiter := tollbooth.NewLimiter(100, time.Second)

	return &fasthttp.Server{
		Name:              "injunx",
		Handler:           tollbooth_fasthttp.LimitHandler(fastRouter.Handler, limiter),
		ReduceMemoryUsage: true,
		LogAllErrors:      true,
	}
}

func cors(next fasthttp.RequestHandler) fasthttp.RequestHandler {
	return func(ctx *fasthttp.RequestCtx) {

		ctx.Response.Header.Set("Access-Control-Allow-Origin", "*")
		ctx.Response.Header.Set("Access-Control-Allow-Methods", "POST, GET, OPTIONS, PUT, DELETE")
		ctx.Response.Header.Set("Access-Control-Allow-Headers", "Accept, Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization")
		ctx.Response.Header.Set("Content-Type", "application/json; charset=UTF-8")

		next(ctx)
	}
}

func httpError(ctx *fasthttp.RequestCtx, err error) {
	log.Error("internal server error: " + err.Error())
	ctx.Error(http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
}
