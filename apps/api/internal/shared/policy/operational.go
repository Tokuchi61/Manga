package policy

// RateLimitSurface defines canonical surfaces with baseline rate-limit expectations.
type RateLimitSurface string

const (
	RateLimitSurfaceAuthLogin       RateLimitSurface = "auth.login"
	RateLimitSurfaceCommentWrite    RateLimitSurface = "comment.write"
	RateLimitSurfaceSupportIntake   RateLimitSurface = "support.intake"
	RateLimitSurfaceSocialMessaging RateLimitSurface = "social.messaging"
	RateLimitSurfacePaymentCallback RateLimitSurface = "payment.callback"
	RateLimitSurfaceAdsClickIntake  RateLimitSurface = "ads.click_intake"
)

var AllRateLimitSurfaces = []RateLimitSurface{
	RateLimitSurfaceAuthLogin,
	RateLimitSurfaceCommentWrite,
	RateLimitSurfaceSupportIntake,
	RateLimitSurfaceSocialMessaging,
	RateLimitSurfacePaymentCallback,
	RateLimitSurfaceAdsClickIntake,
}

const (
	TraceFieldRequestID     = "request_id"
	TraceFieldCorrelationID = "correlation_id"
)

var CallbackAdditionalTraceReference = "provider_event_reference"

func IsValidRateLimitSurface(value RateLimitSurface) bool {
	switch value {
	case RateLimitSurfaceAuthLogin,
		RateLimitSurfaceCommentWrite,
		RateLimitSurfaceSupportIntake,
		RateLimitSurfaceSocialMessaging,
		RateLimitSurfacePaymentCallback,
		RateLimitSurfaceAdsClickIntake:
		return true
	default:
		return false
	}
}
