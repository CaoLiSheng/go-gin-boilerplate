package srv

import (
	"errors"
	"go-gin-boilerplate/cmd"
	"go-gin-boilerplate/db"
	"log"
	"time"

	jwt "github.com/appleboy/gin-jwt/v2"
	"github.com/gin-gonic/gin"
)



const BearerID = "__bearer-id__"
const SuperUserType = "super_admin"

// 登录校验
type jwtAuthenticator func(*gin.Context) (interface{}, error)
type loginForm struct {
	Type     string `form:"type" binding:"required"`
	UserName string `form:"name" binding:"required"`
	Password string `form:"pwd" binding:"required"`
}
type User struct {
	Type  string `json:"type"`
	ID    int64  `json:"id"`
	Name  string `json:"name"`
	Title string `json:"title,omitempty"`
}
type authenticatorJob func(string, string) (func (*db.Core) *User)
var authenticators = func() (avs map[string]authenticatorJob) {
	avs = make(map[string]authenticatorJob)
	avs[SuperUserType] = func(name, pwd string) func(*db.Core) *User {
		return func(c *db.Core) *User {
			if name == cmd.TheFlags.SUN && pwd == cmd.TheFlags.SUP {
				return &User{Title: name}
			}
			panic(errors.New("登录验证失败"))
		}
	}
	return
}()
func authenticator(tipe string) jwtAuthenticator {
	return func(c *gin.Context) (interface{}, error) {
		form := new(loginForm)
		err := c.ShouldBind(form)

		if err != nil {
			return nil, errors.New("登录参数错误【" + err.Error() + "】")
		}

		if form.Type != tipe {
			return nil, errors.New("登录参数错误【用户类型错误】")
		}

		name := form.UserName
		pwd := form.Password

		user := &User{
			Type: tipe,
			Name: name,
		}
		defer func() {
			c.Set(BearerID, user)
		}()

		var dbErr error
		authenticator := authenticators[tipe]
		MustGet(c).DoSimple(&db.JobOptions{Timeout: 5*time.Second}, func(c *db.Core) {
			dbUser := authenticator(name, pwd)(c)
			user.ID = dbUser.ID
			user.Title = dbUser.Title
		}, func(e error) { dbErr = e })

		if dbErr != nil {
			return nil, dbErr
		}
		return user, nil
	}
}

func payload(data interface{}) jwt.MapClaims {
	if u, ok := data.(*User); ok {
		return jwt.MapClaims{"type": u.Type, "name": u.Name, "id": u.ID}
	}
	return jwt.MapClaims{}
}

type jwtLoginResponse func(*gin.Context, int, string, time.Time)
func loginResponse(tipe string) jwtLoginResponse {
	return func(c *gin.Context, code int, token string, expire time.Time) {
		if u, ok := c.Get(BearerID); ok {
			if user, ok := u.(*User); ok {
				user.ID = -1
			}
			(&Result{
				Code: code,
				Results: gin.H{
					"expire": expire.Format("2006-01-02T15:04:05-07:00"),
					"token":  token,
					"user":   u,
				},
			}).Send(c)
			return
		}

		(&Result{
			Code: code,
			Results: gin.H{
				"expire": expire.Format("2006-01-02T15:04:05-07:00"),
				"token":  token,
			},
		}).Send(c)
	}
}

// 定位用户
func identify(c *gin.Context) interface{} {
	claims := jwt.ExtractClaims(c)
	return &User{
		Type: claims["type"].(string),
		Name: claims["name"].(string),
		ID:   int64(claims["id"].(float64)),
	}
}

// 验证用户信息
type jwtAuthorizator func(interface{}, *gin.Context) bool
type authorizatorJob func(*User) func(*db.Core)
var authorizators = func() (avs map[string]authorizatorJob) {
	avs = make(map[string]authorizatorJob)
	avs[SuperUserType] = func(u *User) func(*db.Core) {
		return func(c *db.Core) {}
	}
	return
}()
func authorizator(tipe string) jwtAuthorizator {
	return func(data interface{}, c *gin.Context) bool {
		if u, ok := data.(*User); ok && u.Type == tipe {
			var dbErr error
			authorizator := authorizators[tipe]
			MustGet(c).DoSimple(&db.JobOptions{Timeout: 5*time.Second}, authorizator(u), func(e error) { dbErr = e })

			return dbErr == nil
		}

		return false
	}
}

func unauthorized(c *gin.Context, code int, message string) {
	c.JSON(code, gin.H{
		"code":    code,
		"message": message,
	})
}

func BearerAuth(tipe string) *jwt.GinJWTMiddleware {
	auth, err := jwt.New(&jwt.GinJWTMiddleware{
		Authenticator:   authenticator(tipe),
		PayloadFunc:     payload,
		LoginResponse:   loginResponse(tipe),
		// LogoutResponse:  , // 需要前端自行删除 Token 数据以退出登录
		IdentityHandler: identify,
		Authorizator:    authorizator(tipe),
		Unauthorized:    unauthorized,
		TimeFunc:        time.Now,
		TokenLookup:     "header: Authorization",
		TokenHeadName:   "Bearer",
		Realm:           "HS256",
		Key:             []byte("詳しい情報は下記リンク先をご覧"),
		Timeout:         time.Hour * 12,
		MaxRefresh:      time.Hour * 36,
		IdentityKey:     BearerID,
	})

	if err != nil {
		log.Fatal("JWT Error:", err)
	}

	return auth
}