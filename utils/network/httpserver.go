/**
有一个缺陷：
若有捕获根目录("/")，则需要在"/"对应的HandleFunc过滤没有注册过的目录，否则会出现访问(http://localhost/test)也会返回(http://localhost/)的结果
 */
package network

import (
	"context"
	"crypto/sha256"
	"crypto/subtle"
	"errors"
	"net/http"
	"os"
	"os/signal"
	"regexp"
	"syscall"
)

//自定义httpserver才可以shutdown，否则http.ListenAndServe无法关闭服务
type HttpServer struct {
	mux *http.ServeMux
}

func (s *HttpServer) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	s.mux.ServeHTTP(w, r)
}

func (s *HttpServer) HandleFunc(pattern string, handler func(http.ResponseWriter, *http.Request)) {
	s.mux.HandleFunc(pattern, handler)
}

func newHs() *HttpServer {
	s := &HttpServer{mux: http.NewServeMux()}
	return s
}

func NewHttpServer(port string) *http.Server {
	return &http.Server{
		Addr:    port,
		Handler: newHs(),
	}
}

func ShutdownHttpServer(hs *http.Server) {
	ctx, cancel := context.WithCancel(context.Background())
	hs.Shutdown(ctx)
	cancel()
}

func HandleFunc(hs *http.Server, pattern string, handler func(http.ResponseWriter, *http.Request)) {
	if hs != nil {
		if s, ok := hs.Handler.(*HttpServer); ok {
			s.HandleFunc(pattern, handler)
		}
	}
}

func HandleFuncAuth(hs *http.Server, pattern string, handler func(http.ResponseWriter, *http.Request), user, pwd string) {
	if hs != nil {
		if s, ok := hs.Handler.(*HttpServer); ok {
			s.HandleFunc(pattern, authHandler(handler, user, pwd))
		}
	}
}

func ListenAndServe(hs *http.Server) error {
	if hs != nil {
		return hs.ListenAndServe()
	}
	return errors.New("ListenAndServe failed")
}

func hasher(s string) []byte {
	val := sha256.Sum256([]byte(s))
	return val[:]
}

func authHandler(handler http.HandlerFunc, user, pwd string) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		userHast := hasher(user)
		passHash := hasher(pwd)
		user, pass, ok := r.BasicAuth()
		if !ok || subtle.ConstantTimeCompare(hasher(user), userHast) != 1 || subtle.ConstantTimeCompare(hasher(pass), passHash) != 1 {
			w.Header().Set("WWW-Authenticate", `Basic realm="Please enter username and password"`)
			http.Error(w, "Unauthorized.", http.StatusUnauthorized)
			return
		}
		handler(w, r)
	}
}

func GetIpAddrByRequest(request *http.Request) string {
	regex, _ := regexp.Compile("([^:]+)")
	return regex.FindString(request.RemoteAddr)
}

func ShutdownHttpServerWaitSignal(hs *http.Server) {
	signalChan := make(chan os.Signal, 1)
	signal.Notify(signalChan, syscall.SIGINT, syscall.SIGTERM)
	<-signalChan
	ShutdownHttpServer(hs)
}
