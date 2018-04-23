package service

import (
	"compress/gzip"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"path"
	"runtime"
	"strconv"

	"github.com/gorilla/mux"
	"github.com/k0kubun/pp"
	"jvole.com/dsp/config"
	"jvole.com/dsp/dsperror"
	"jvole.com/dsp/util"
	// "golang.org/x/time/rate"

	kitlog "github.com/go-kit/kit/log"
	"github.com/go-kit/kit/log/level"
	kithttp "github.com/go-kit/kit/transport/http"
)

var (
	errBadRoute = errors.New("bad route")
	errNotBid   = errors.New("notbid")
)

var User, Password string

func MakeHandler(bs DSPBidder, logger kitlog.Logger) http.Handler {

	opts := []kithttp.ServerOption{
		kithttp.ServerBefore(kithttp.PopulateRequestContext),
		kithttp.ServerErrorLogger(logger),
		kithttp.ServerErrorEncoder(encodeError),
	}
	r := mux.NewRouter()

	smaatoBidder := makeBidderEndpoint(bs)
	// smaatoBidder = basic.AuthMiddleware(User, Password, "")(smaatoBidder)
	// smaatoBidder = middleware.ValidMiddleware()(smaatoBidder)
	// smaatoBidder = ratelimit.NewErroringLimiter(rate.NewLimiter(rate.Every(time.Second), config.QPS))(smaatoBidder)
	smaatoBidderHandler := kithttp.NewServer(
		smaatoBidder,
		decodeBidderRequest,
		encodeResponse,
		opts...,
	)

	smaatoADXNotify := makeADXNotifyEndpoint(bs)
	// smaatoADXNotify = basic.AuthMiddleware(User, Password, "")(smaatoADXNotify)
	// smaatoADXNotify = middleware.ValidMiddleware()(smaatoADXNotify)
	// smaatoADXNotify = ratelimit.NewErroringLimiter(rate.NewLimiter(rate.Every(time.Second), config.QPS))(smaatoADXNotify)
	smaatoADXNotifyHandler := kithttp.NewServer(
		smaatoADXNotify,
		decodeADXNotifyRequest,
		encodeResponse,
		opts...,
	)

	smaatoStopBidByCID := makeStopBidByCIDEndpoint(bs)
	// smaatoBidder = basic.AuthMiddleware(User, Password, "")(smaatoBidder)
	// smaatoBidder = middleware.ValidMiddleware()(smaatoBidder)
	// smaatoBidder = ratelimit.NewErroringLimiter(rate.NewLimiter(rate.Every(time.Second), config.QPS))(smaatoBidder)
	smaatoStopBidByCIDHandler := kithttp.NewServer(
		smaatoStopBidByCID,
		decodeStopBidByCIDRequest,
		encodeResponse,
		opts...,
	)

	smaatoSyncIndex := makeSyncIndexEndpoint(bs)
	// smaatoBidder = basic.AuthMiddleware(User, Password, "")(smaatoBidder)
	// smaatoBidder = middleware.ValidMiddleware()(smaatoBidder)
	// smaatoBidder = ratelimit.NewErroringLimiter(rate.NewLimiter(rate.Every(time.Second), config.QPS))(smaatoBidder)
	smaatoSyncIndexHandler := kithttp.NewServer(
		smaatoSyncIndex,
		decodeSyncIndexRequest,
		encodeResponse,
		opts...,
	)
	//不同adx通过二级域名区别
	r.Handle("/v1/Bidder", smaatoBidderHandler).Methods("POST")
	r.Handle("/v1/ADXNotify", smaatoADXNotifyHandler).Methods("GET")
	r.Handle("/v1/SyncIndex", smaatoSyncIndexHandler).Methods("POST")
	r.Handle("/v1/StopBid", smaatoStopBidByCIDHandler).Methods("POST")

	return r
}

func decodeSyncIndexRequest(ctx context.Context, r *http.Request) (interface{}, error) {
	var body struct {
		Server uint32 `json:"server"`
	}
	if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
		return nil, err
	}
	defer r.Body.Close()
	return SyncIndexRequest{}, nil
}

func decodeStopBidByCIDRequest(ctx context.Context, r *http.Request) (interface{}, error) {
	var body struct {
		CID uint32 `json:"cid"`
	}
	if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
		return nil, err
	}
	defer r.Body.Close()

	return StopBidByCIDRequest{body.CID}, nil
}

func decodeBidderRequest(ctx context.Context, r *http.Request) (interface{}, error) {
	encoding := r.Header.Get("Content-Encoding")
	var body []byte

	if encoding == "gzip" { //判断gzip
		resp, err := gzip.NewReader(r.Body)
		if err != nil {
			pc, file, line, _ := runtime.Caller(1)
			f := runtime.FuncForPC(pc)
			level.Error(util.KitLogger).Log(
				"method", f.Name(),
				"file", path.Base(file),
				"line", line,
				"msg", "gzip 解压错误",
			)
		}
		body, err = ioutil.ReadAll(resp)
		defer r.Body.Close()
		if err != nil {
			pc, file, line, _ := runtime.Caller(1)
			f := runtime.FuncForPC(pc)
			level.Error(util.KitLogger).Log(
				"method", f.Name(),
				"file", path.Base(file),
				"line", line,
				"msg", "读取body错误",
			)
		}
	} else {
		body, _ = ioutil.ReadAll(r.Body)
	}

	return BidderRequest{body, r.Host}, nil
}

func decodeADXNotifyRequest(ctx context.Context, r *http.Request) (interface{}, error) {
	res := r.FormValue(config.NURLParam["OID"]) //广告id
	if res == "" {
		return nil, errBadRoute
	}
	oid := res
	res = r.FormValue(config.NURLParam["Postion"]) //postion
	pp.Println(config.NURLParam["Postion"], res)
	if res == "" {
		return nil, errBadRoute
	}
	postion := res
	res = r.FormValue(config.NURLParam["Device"]) //device
	pp.Println(config.NURLParam["Device"], res)
	// if res == "" {
	// 	return nil, errBadRoute
	// }
	device := res
	res = r.FormValue(config.NURLParam["User"]) //user
	pp.Println(config.NURLParam["User"], res)
	// if res == "" {
	// 	return nil, errBadRoute
	// }
	user := res
	res = r.FormValue(config.NURLParam["Price"]) //price
	pp.Println(config.NURLParam["Price"], res)
	if res == "" {
		return nil, errBadRoute
	}
	price, err := strconv.ParseFloat(res, 8)
	if err != nil {
		pc, file, line, _ := runtime.Caller(1)
		f := runtime.FuncForPC(pc)
		level.Error(util.KitLogger).Log(
			"method", f.Name(),
			"file", path.Base(file),
			"line", line,
			"content", res,
			"msg", "price传输的值错误",
		)
	}
	res = r.FormValue(config.NURLParam["CID"]) //cid
	pp.Println(config.NURLParam["CID"], res)
	if res == "" {
		return nil, errBadRoute
	}
	cid, err := strconv.ParseUint(res, 10, 32)
	if err != nil {
		pc, file, line, _ := runtime.Caller(1)
		f := runtime.FuncForPC(pc)
		level.Error(util.KitLogger).Log(
			"method", f.Name(),
			"file", path.Base(file),
			"line", line,
			"content", res,
			"msg", "cid传输的值错误",
		)
	}
	res = r.FormValue(config.NURLParam["UID"]) //uid
	pp.Println(config.NURLParam["UID"], res)
	if res == "" {
		return nil, errBadRoute
	}
	uid, err := strconv.ParseUint(res, 10, 32)
	if err != nil {
		pc, file, line, _ := runtime.Caller(1)
		f := runtime.FuncForPC(pc)
		level.Error(util.KitLogger).Log(
			"method", f.Name(),
			"file", path.Base(file),
			"line", line,
			"content", res,
			"msg", "uid传输的值错误",
		)
	}
	res = r.FormValue(config.NURLParam["T"]) //time
	pp.Println(config.NURLParam["T"], res)
	if res == "" {
		return nil, errBadRoute
	}
	t, err := strconv.ParseInt(res, 10, 63)
	if err != nil {
		pc, file, line, _ := runtime.Caller(1)
		f := runtime.FuncForPC(pc)
		level.Error(util.KitLogger).Log(
			"method", f.Name(),
			"file", path.Base(file),
			"line", line,
			"content", res,
			"msg", "time传输的值错误",
		)
	}
	res = r.FormValue(config.IMPURL_ClickID) //uid
	pp.Println(config.IMPURL_ClickID, res)
	if res == "" {
		return nil, errBadRoute
	}
	clickID := res

	return ADXNotifyRequest{
		oid,
		price,
		uint32(cid),
		uint32(uid),
		postion,
		device,
		t,
		user,
		clickID,
	}, nil
}

func encodeResponse(ctx context.Context, w http.ResponseWriter, response interface{}) error {
	defer func() {
		pc, file, line, _ := runtime.Caller(1)
		f := runtime.FuncForPC(pc)
		level.Debug(util.KitLogger).Log(
			"method", f.Name(),
			"file", path.Base(file),
			"line", line,
			"msg", "返回内容",
			"response", response,
		)
	}()
	if e, ok := response.(errorer); ok && e.error() != nil {
		encodeError(ctx, e.error(), w)
		return nil
	}
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	// w.Header().Set("Content-Type", "text/plain; charset=utf-8")
	// return response.(BIDResponse).Result
	if v, ok := response.(string); ok { //bidresponse
		fmt.Fprint(w, v)
		return nil
	} else if v, ok := response.([]byte); ok { //syncindex
		w.Header().Set("Content-Type", "application/zip")
		w.Header().Set("Content-Disposition", fmt.Sprintf("attachment; filename=\"%s\"", config.CPIDXNAME))
		util.GetResponseZip(w, v)
		return nil
	} else if v, ok := response.(*CommResponse); ok { //CommResponse
		w.Header().Set("Content-Type", "application/zip")
		w.Header().Set("Content-Disposition", fmt.Sprintf("attachment; filename=\"%s\"", config.CPIDXNAME))
		return json.NewEncoder(w).Encode(v)
	} else {
		pc, file, line, _ := runtime.Caller(1)
		f := runtime.FuncForPC(pc)
		level.Debug(util.KitLogger).Log(
			"method", f.Name(),
			"file", path.Base(file),
			"line", line,
			"msg", "返回数据错误",
			"response", response,
		)
		return errNotBid
	}

	// return json.NewEncoder(w).Encode(response)
}

type errorer interface {
	error() error
}

// encode errors from business-logic
func encodeError(_ context.Context, err error, w http.ResponseWriter) {
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	switch err {
	case errNotBid:
		w.WriteHeader(http.StatusNoContent)
	case errBadRoute:
		w.WriteHeader(http.StatusNotFound)
	case errBadRoute:
		w.WriteHeader(http.StatusBadRequest)
	default:
		w.WriteHeader(http.StatusInternalServerError)
	}
	var errcode int
	var msg string
	errcode = dsperror.ErrorUnknown
	msg = dsperror.ErrorText(errcode) + ":" + err.Error()
	json.NewEncoder(w).Encode(
		&CommResponse{Errcode: errcode, Msg: msg, Data: ""},
	)
}
