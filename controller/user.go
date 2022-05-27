package controller

import (
	"fmt"
	"gitee.com/Whitroom/imitate-tiktok/sql"
	"gitee.com/Whitroom/imitate-tiktok/sql/models"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
	"net/http"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
)

// usersLoginInfo use map to store user info, and key is username+password for demo
// user data will be cleared every time the server starts
// test data: username=zhanglei, password=douyin
var usersLoginInfo = map[string]User{
	"zhangleidouyin": {
		Id:            1,
		Name:          "zhanglei",
		FollowCount:   10,
		FollowerCount: 5,
		IsFollow:      true,
	},
}

var userIdSequence = int64(1)

type UserLoginResponse struct {
	Response
	UserId int64  `json:"user_id,omitempty"`
	Token  string `json:"token"`
}

type UserResponse struct {
	Response
	User models.User `json:"user"`
}

func Register(c *gin.Context) {

	var user models.User
	if err := c.ShouldBindQuery(&user); err != nil {
		c.JSON(400, UserLoginResponse{
			Response: Response{StatusCode: 1, StatusMsg: "bind failed"},
		})
	}

	if err := sql.DB.Where("name = ?", user.Name).First(&user).Error; err == gorm.ErrRecordNotFound {
		hashcode := hashEncode(user.Password)
		var newUser = models.User{
			Name:     user.Name,
			Password: hashcode,
		}
		sql.DB.Create(&newUser)
		fmt.Println("创建成功！！！！")
		token := string(newUser.ID) + "+" + string(time.Now().Unix())
		c.Set("token", token)
		c.JSON(http.StatusOK, UserLoginResponse{
			Response: Response{StatusCode: 0},
			UserId:   userIdSequence,
			Token:    token,
		})
	} else if err == nil {
		c.JSON(403, UserLoginResponse{
			Response: Response{StatusCode: 1, StatusMsg: "User is exist"},
		})
	}

	//if _, exist := usersLoginInfo[token]; exist {
	//	c.JSON(http.StatusOK, UserLoginResponse{
	//		Response: Response{StatusCode: 1, StatusMsg: "User already exist"},
	//	})
	//} else {
	//	atomic.AddInt64(&userIdSequence, 1)
	//	newUser := User{
	//		Id:   userIdSequence,
	//		Name: username,
	//	}
	//	usersLoginInfo[token] = newUser
	//	c.JSON(http.StatusOK, UserLoginResponse{
	//		Response: Response{StatusCode: 0},
	//		UserId:   userIdSequence,
	//		Token:    username + password,
	//	})
	//}
}

func Login(c *gin.Context) {

	var user models.User
	if err := c.ShouldBindQuery(&user); err != nil {
		c.JSON(400, UserLoginResponse{
			Response: Response{StatusCode: 1, StatusMsg: "bind failed"},
		})
	}
	var existedUser models.User

	if err := sql.DB.Where("name = ?", user.Name).First(&existedUser).Error; err == gorm.ErrRecordNotFound {
		c.JSON(404, UserLoginResponse{
			Response: Response{StatusCode: 1, StatusMsg: "User doesn't exist"},
		})
	}
	pwdMatch := comparePasswords(user.Password, existedUser.Password)
	if pwdMatch != true {
		c.JSON(401, UserLoginResponse{
			Response: Response{StatusCode: 1, StatusMsg: "password is incorrect"},
		})
	}
	token := string(existedUser.ID) + "+" + string(time.Now().Unix())
	c.Set("token", token)
	c.JSON(200, UserLoginResponse{
		Response: Response{StatusCode: 0},
		UserId:   int64(existedUser.ID),
		Token:    token,
	})
	//if user, exist := usersLoginInfo[token]; exist {
	//	c.JSON(http.StatusOK, UserLoginResponse{
	//		Response: Response{StatusCode: 0},
	//		UserId:   user.Id,
	//		Token:    token,
	//	})
	//} else {
	//	c.JSON(http.StatusOK, UserLoginResponse{
	//		Response: Response{StatusCode: 1, StatusMsg: "User doesn't exist"},
	//	})
	//}
}

func UserInfo(c *gin.Context) {
	//token := c.Query("token")
	var userId uint64
	var err error
	str := c.Query("user_id")
	if userId, err = strconv.ParseUint(str, 0, len(str)); err != nil {
		c.JSON(400, gin.H{
			"error": err.Error(),
		})
	}
	var user models.User
	if err = sql.DB.First(&user, uint(userId)).Error; err == gorm.ErrRecordNotFound {
		c.JSON(404, UserResponse{
			Response: Response{StatusCode: 1, StatusMsg: "User doesn't exist"},
		})
	} else if err == nil {
		c.JSON(200, UserResponse{
			Response: Response{StatusCode: 0},
			User:     user,
		})
	}

	//if user, exist := usersLoginInfo[token]; exist {
	//	c.JSON(http.StatusOK, UserResponse{
	//		Response: Response{StatusCode: 0},
	//		User:     user,
	//	})
	//} else {
	//	c.JSON(http.StatusOK, UserResponse{
	//		Response: Response{StatusCode: 1, StatusMsg: "User doesn't exist"},
	//	})
	//}
}
func hashEncode(str string) string {
	hash, err := bcrypt.GenerateFromPassword([]byte(str), bcrypt.DefaultCost)
	if err != nil {
		fmt.Println("failed to hash:%w", err)
	}
	return string(hash)
}
func comparePasswords(sourcePwd, hashPwd string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hashPwd), []byte(sourcePwd))
	if err != nil {
		return false
	}
	return true
}
