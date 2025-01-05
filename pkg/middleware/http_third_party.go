package middleware

import (
	chi_middleware "github.com/go-chi/chi/v5/middleware"
)

var HTTPRequestID = chi_middleware.RequestID
var HTTPRealIP = chi_middleware.RealIP
var HTTPCompress = chi_middleware.Compress
