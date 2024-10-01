// Anti-Bruteforce API.
//
//	 Schemes: http
//	 BasePath: /
//	 Version: 1.0.0
//	 Host: localhost:8080
//
//	 Consumes:
//	 - application/json
//
//	 Produces:
//	 - application/json
//
//	 Security:
//	 - basic
//
//	SecurityDefinitions:
//	basic:
//	  type: basic
//
// swagger:meta
package docs

import "github.com/loginovm/anti-bruteforce/internal/server/http/models"

// swagger:route PUT /check-login-attempt Check checkLoginAttempt
// Check login attempt.
// responses:
//   200: CheckLoginAttemptResponse
//	 400: BadRequestResponse

// swagger:route GET /blacklist IPBlacklist getIPBlackList
// Get IP blacklist.
// responses:
//   200: BWListResponse

// swagger:route GET /whitelist IPWhitelist getIPWhiteList
// Get IP whitelist.
// responses:
//   200: BWListResponse

// swagger:route POST /blacklist IPBlacklist addIPToBlackList
// Add IP to blacklist.
// responses:
//   200: OkResponse
//	 400: BadRequestResponse

// swagger:route DELETE /blacklist IPBlacklist deleteIPFromBlackList
// Delete IP from blacklist.
// responses:
//   200: OkResponse
//	 400: BadRequestResponse

// swagger:route POST /whitelist IPWhitelist addIPToWhiteList
// Add IP to whitelist.
// responses:
//   200: OkResponse
//	 400: BadRequestResponse

// swagger:route DELETE /whitelist IPWhitelist deleteIPFromWhiteList
// Delete IP from whitelist.
// responses:
//   200: OkResponse
//	 400: BadRequestResponse

// swagger:route PUT /reset-ip Reset resetIPBucket
// Reset IP.
// responses:
//   200: OkResponse
//	 400: BadRequestResponse

// swagger:route PUT /reset-login Reset resetLoginBucket
// Reset Login.
// responses:
//   200: OkResponse
//	 400: BadRequestResponse

// swagger:route GET /settings Settings getSettings
// Get settings.
// responses:
//   200: SettingsResponse

// swagger:route PUT /settings Settings updateSettings
// Update settings.
// responses:
//   200: OkResponse
//	 400: BadRequestResponse

// swagger:parameters updateSettings
type SettingsParamsWrapper struct {
	// in:body
	Body models.Setting
}

// swagger:parameters checkLoginAttempt
type CheckLoginAttemptParamsWrapper struct {
	// in:body
	Body models.CheckLoginAttemptRequest
}

// swagger:parameters addIPToBlackList deleteIPFromBlackList addIPToWhiteList deleteIPFromWhiteList
type IPParamsWrapper struct {
	// in:body
	Body models.CidrRequest
}

// swagger:parameters resetIPBucket
type ResetIPParamsWrapper struct {
	// in:body
	Body models.IPRequest
}

// swagger:parameters resetLoginBucket
type ResetLoginParamsWrapper struct {
	// in:body
	Body models.LoginRequest
}

// swagger:response CheckLoginAttemptResponse
type CheckLoginAttemptResponseWrapper struct {
	// in:body
	Body models.CheckAttemptResponse
}

// swagger:response BWListResponse
type BWListResponseWrapper struct {
	// in:body
	Body models.BWListResponse
}

// swagger:response SettingsResponse
type SettingsResponseWrapper struct {
	// in:body
	Body models.Setting
}

// No content
// swagger:response OkResponse
type OkResponse struct{}

// Model for BadRequest response
// swagger:response BadRequestResponse
type BadRequestResponse struct {
	// in:body
	Body struct {
		Message string `json:"message"`
	} `json:"body"`
}
