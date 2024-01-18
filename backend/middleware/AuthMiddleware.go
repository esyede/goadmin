package middleware

import (
	"github.com/esyede/goadmin/backend/common"
	"github.com/esyede/goadmin/backend/config"
	"github.com/esyede/goadmin/backend/model"
	"github.com/esyede/goadmin/backend/repository"
	"github.com/esyede/goadmin/backend/response"
	"github.com/esyede/goadmin/backend/util"
	"github.com/esyede/goadmin/backend/vo"
	"fmt"
	"time"

	jwt "github.com/appleboy/gin-jwt/v2"
	"github.com/gin-gonic/gin"
)

// Initialize jwt middleware
func InitAuth() (*jwt.GinJWTMiddleware, error) {
	authMiddleware, err := jwt.New(&jwt.GinJWTMiddleware{
		Realm:           config.Conf.Jwt.Realm,                                 // jwt realm
		Key:             []byte(config.Conf.Jwt.Key),                           // Server key
		Timeout:         time.Hour * time.Duration(config.Conf.Jwt.Timeout),    // Token expiration time
		MaxRefresh:      time.Hour * time.Duration(config.Conf.Jwt.MaxRefresh), // Token maximum refresh time (RefreshToken expiration time = Timeout + MaxRefresh)
		PayloadFunc:     payloadFunc,                                           // payload processing
		IdentityHandler: identityHandler,                                       // Parse Claims
		Authenticator:   login,                                                 // Verify the correctness of the token and process the login logic
		Authorizator:    authorizator,                                          // User login verification is successfully processed
		Unauthorized:    unauthorized,                                          // Handling user login verification failure
		LoginResponse:   loginResponse,                                         // Response after successful login
		LogoutResponse:  logoutResponse,                                        // Response after logging out
		RefreshResponse: refreshResponse,                                       // Response after refreshing token
		TokenLookup:     "header: Authorization, query: token, cookie: jwt",    // Automatically search for the token in the request in these places
		TokenHeadName:   "Bearer",                                              // header name
		TimeFunc:        time.Now,
	})
	return authMiddleware, err
}

// payload processing
func payloadFunc(data interface{}) jwt.MapClaims {
	if v, ok := data.(map[string]interface{}); ok {
		var user model.User
		// Convert user json to structure
		util.JsonI2Struct(v["user"], &user)
		return jwt.MapClaims{
			jwt.IdentityKey: user.ID,
			"user":          v["user"],
		}
	}
	return jwt.MapClaims{}
}

// Parse Claims
func identityHandler(c *gin.Context) interface{} {
	claims := jwt.ExtractClaims(c)
	// The return value type map[string]interface{} here must be consistent with the data type of payloadFunc and authorizator,
	// otherwise it will cause authorization failure and it is not easy to find the reason.
	return map[string]interface{}{
		"IdentityKey": claims[jwt.IdentityKey],
		"user":        claims["user"],
	}
}

// Verify the correctness of the token and process login logic
func login(c *gin.Context) (interface{}, error) {
	var req vo.RegisterAndLoginRequest
	// Request json binding
	if err := c.ShouldBind(&req); err != nil {
		return "", err
	}

	// Password decrypted via RSA
	decodeData, err := util.RSADecrypt([]byte(req.Password), config.Conf.System.RSAPrivateBytes)
	if err != nil {
		return nil, err
	}

	u := &model.User{
		Username: req.Username,
		Password: string(decodeData),
	}

	// Password verification
	userRepository := repository.NewUserRepository()
	user, err := userRepository.Login(u)
	if err != nil {
		return nil, err
	}
	// Write the user in json format, which will be used by payloadFunc/authorizator
	return map[string]interface{}{
		"user": util.Struct2Json(user),
	}, nil
}

// User login verification is successfully processed
func authorizator(data interface{}, c *gin.Context) bool {
	if v, ok := data.(map[string]interface{}); ok {
		userStr := v["user"].(string)
		var user model.User
		// Convert user json to structure
		util.Json2Struct(userStr, &user)
		// Save the user to the context, it is convenient to retrieve the data when calling the api
		c.Set("user", user)
		return true
	}
	return false
}

// Handling user login verification failure
func unauthorized(c *gin.Context, code int, message string) {
	common.Log.Debugf("JWT authentication failed, error code: %d, message: %s", code, message)
	response.Response(c, code, code, nil, fmt.Sprintf("JWT authentication failed, error code: %d, message: %s", code, message))
}

// Response after successful login
func loginResponse(c *gin.Context, code int, token string, expires time.Time) {
	response.Response(c, code, code,
		gin.H{
			"token":   token,
			"expires": expires.Format("2006-01-02 15:04:05"),
		},
		"Login successful")
}

// Response after logging out
func logoutResponse(c *gin.Context, code int) {
	response.Success(c, nil, "Log out succeeful")
}

// Response after refreshing token
func refreshResponse(c *gin.Context, code int, token string, expires time.Time) {
	response.Response(c, code, code,
		gin.H{
			"token":   token,
			"expires": expires,
		},
		"Refresh token successful")
}
