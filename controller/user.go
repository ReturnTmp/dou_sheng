package controller

import (
	"fmt"
	"net/http"

	"gitee.com/Whitroom/imitate-tiktok/middlewares"
	"gitee.com/Whitroom/imitate-tiktok/sql/crud"
	"gitee.com/Whitroom/imitate-tiktok/sql/models"
	"github.com/gin-gonic/gin"
	"golang.org/x/crypto/bcrypt"
)

func hashEncode(str string) string {
	hash, err := bcrypt.GenerateFromPassword([]byte(str), bcrypt.DefaultCost)
	if err != nil {
		fmt.Println("failed to hash:%w", err)
	}
	return string(hash)
}

func comparePasswords(sourcePwd, hashPwd string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hashPwd), []byte(sourcePwd))
	return err == nil
}

type UserLoginResponse struct {
	Response
	UserId int64  `json:"user_id,omitempty"`
	Token  string `json:"token"`
}

type UserResponse struct {
	Response
	User User `json:"user"`
}

type RegisterRequest struct {
	Username string `json:"username" binding:"required" form:"username"`
	Password string `json:"password" binding:"required" form:"password"`
}

func Register(ctx *gin.Context) {
	var request RegisterRequest
	if !BindAndValid(ctx, &request) {
		return
	}

	if crud.GetUserByName(request.Username) == nil {
		ctx.JSON(http.StatusBadRequest, Response{
			StatusCode: 2,
			StatusMsg:  "存在用户姓名",
		})
		return
	}

	newUser := &models.User{
		Name:     request.Username,
		Password: hashEncode(request.Password),
		Content:  "",
	}

	newUser = crud.CreateUser(newUser)

	token, err := middlewares.Sign(newUser.ID)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, Response{
			StatusCode: 3,
			StatusMsg:  "token创建失败",
		})
		return
	}

	ctx.JSON(http.StatusOK, UserLoginResponse{
		Response: Response{
			StatusCode: 0,
			StatusMsg:  "用户创建成功",
		},
		UserId: int64(newUser.ID),
		Token:  token,
	})
}

func Login(ctx *gin.Context) {
	var request RegisterRequest
	if !BindAndValid(ctx, &request) {
		return
	}
	existedUser := crud.GetUserByName(request.Username)

	if existedUser == nil {
		ctx.JSON(http.StatusNotFound, Response{
			StatusCode: 2,
			StatusMsg:  "找不到用户",
		})
		return
	}

	pwdMatch := comparePasswords(request.Password, existedUser.Password)
	if !pwdMatch {
		ctx.JSON(http.StatusUnauthorized, UserLoginResponse{
			Response: Response{
				StatusCode: 3,
				StatusMsg:  "用户名或密码错误",
			},
		})
		return
	}

	token, err := middlewares.Sign(existedUser.ID)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, Response{
			StatusCode: 4,
			StatusMsg:  "token创建失败",
		})
		return
	}

	ctx.JSON(http.StatusOK, UserLoginResponse{
		Response: Response{StatusCode: 0},
		UserId:   int64(existedUser.ID),
		Token:    token,
	})
}

// 查询用户信息接口函数。
func UserInfo(ctx *gin.Context) {

	user_, _ := ctx.Get("User")
	user, _ := user_.(*models.User)

	toUserID := QueryIDAndValid(ctx, "user_id")
	if toUserID == 0 {
		return
	}

	toUser, err := crud.GetUserByID(toUserID)
	if err != nil {
		ctx.JSON(http.StatusNotFound, Response{
			StatusCode: 3,
			StatusMsg:  "找不到用户",
		})
		return
	}

	responseUser := UserModelChange(*toUser)
	if user != nil {
		responseUser.IsFollow = crud.IsUserFollow(user.ID, toUserID)
	} else {
		responseUser.IsFollow = false
	}

	ctx.JSON(http.StatusOK, UserResponse{
		Response: Response{
			StatusCode: 0,
			StatusMsg:  "已找到用户",
		},
		User: responseUser,
	})

}
