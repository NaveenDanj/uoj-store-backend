package handlers

import (
	"net/http"
	"peer-store/db"
	"peer-store/dto"
	"peer-store/models"
	"peer-store/service/auth"
	"strings"

	"github.com/gin-gonic/gin"
	"golang.org/x/crypto/bcrypt"
)

func CreateNewUser(c *gin.Context) {

	var requestJSON dto.UserRequestDTO

	if err := c.ShouldBindJSON(&requestJSON); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Generate key pairs and store public key on the user table and private key on the backend

	user, err := auth.CreateNewUser(&requestJSON)

	if err != nil && user.Email != "" {
		c.JSON(http.StatusBadRequest, gin.H{
			"message": err.Error(),
		})
		return
	}

	authToken, err := auth.GenerateJWT(user.ID, user.Email, false)

	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"message": err.Error(),
		})
		return
	}

	newFolder := models.Folder{
		Name:          "root",
		UserId:        user.ID,
		ParentID:      nil,
		SpecialFolder: "ROOT_FOLDER",
	}

	sessionFolder := models.Folder{
		Name:          "session",
		UserId:        user.ID,
		ParentID:      &newFolder.ID,
		SpecialFolder: "SESSION_FOLDER",
	}

	if err := db.GetDB().Create(&newFolder).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": "Error while creating folder DB record"})
		return
	}

	if err := db.GetDB().Create(&sessionFolder).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": "Error while creating folder DB record"})
		return
	}

	user.SessionFolder = sessionFolder.ID
	user.RootFolder = newFolder.ID
	if err := db.GetDB().Save(&user).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": "Error while updating user"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message":   "User account created successfully",
		"user":      requestJSON,
		"authToken": authToken,
	})

}

func UserSignIn(c *gin.Context) {
	var requestJSON dto.UserSignInDTO

	if err := c.ShouldBindJSON(&requestJSON); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	user, err := auth.GetUserByUsername(requestJSON.Username)

	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"message": err.Error(),
		})
		return
	}

	err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(requestJSON.Password))

	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"message": "Username or password is in-correct, Please try again.",
		})
		return
	}

	if !user.IsVerified {
		c.JSON(http.StatusBadRequest, gin.H{
			"message": "user account is not verified",
		})
		return
	}

	if !user.IsActive {
		c.JSON(http.StatusBadRequest, gin.H{
			"message": "user account is not activated",
		})
		return
	}

	authToken, err := auth.GenerateJWT(user.ID, user.Username, false)

	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"message": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message":   "Login success",
		"user":      user,
		"authToken": authToken,
	})

}

func PrivateSessionSignIn(c *gin.Context) {
	var requestJSON dto.PrivateSessionSignInDTO

	if err := c.ShouldBindJSON(&requestJSON); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	user, err := auth.GetUserBySessionId(requestJSON.SessionId)

	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"message": err.Error(),
		})
		return
	}

	if !user.IsVerified {
		c.JSON(http.StatusBadRequest, gin.H{
			"message": "user account is not verified",
		})
		return
	}

	if !user.IsActive {
		c.JSON(http.StatusBadRequest, gin.H{
			"message": "user account is not activated",
		})
		return
	}

	authToken, err := auth.GenerateJWT(user.ID, user.Username, true)

	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"message": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message":   "Login success",
		"user":      user,
		"authToken": authToken,
	})

}

func GetCurrentUser(c *gin.Context) {

	user, _ := c.Get("currentUser")

	c.JSON(http.StatusOK, gin.H{
		"user": user,
	})

}

func Logout(c *gin.Context) {

	authToken := c.GetHeader("Authorization")

	if authToken == "" || !strings.HasPrefix(authToken, "Bearer ") {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid token format"})
		c.Abort()
		return
	}

	authToken = strings.TrimPrefix(authToken, "Bearer ")

	err := auth.BlockToken(authToken)

	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"message": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Logout success!",
	})

}

func GetUserPublicKey(c *gin.Context) {

	userId := c.Param("userId")

	if userId == "" {
		c.JSON(http.StatusBadRequest, gin.H{
			"message": "User ID is required",
		})
		return
	}

	pubKey, err := auth.GetPublicKey(userId)

	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"message": "Error while fetching user public key ",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Public key retrieved successfully",
		"pubKey":  pubKey,
	})

}

func VerifyAccount(c *gin.Context) {
	var requestJSON dto.VerfyAccountDTO

	if err := c.ShouldBindJSON(&requestJSON); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": err.Error()})
		return
	}

	res := auth.VerifyAccount(requestJSON.OTP, requestJSON.Email)

	if !res {
		c.JSON(http.StatusBadRequest, gin.H{"message": "Account verification failed"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Account verified successfully"})

}

func ResetPasswordSendMail(c *gin.Context) {
	var requestJSON dto.ResetPasswordSendMailDTO

	if err := c.ShouldBindJSON(&requestJSON); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": err.Error()})
		return
	}

	user, err := auth.GetUserByEmail(requestJSON.Email)

	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": err.Error()})
		return
	}

	if !user.IsVerified {
		c.JSON(http.StatusBadRequest, gin.H{"message": "Unable to send reset password link"})
		return
	}

	if !user.IsActive {
		c.JSON(http.StatusBadRequest, gin.H{"message": "Unable to send reset password link"})
		return
	}

	if _, err := auth.ResetPasswordGenerateLink(user.ID); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": "Unable to send reset password link"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Success"})
}

func ResetPasswordSetPassword(c *gin.Context) {
	var requestJSON dto.ResetPasswordNewPasswordDTO

	if err := c.ShouldBindJSON(&requestJSON); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": err.Error()})
		return
	}

	if err := auth.HandleResetPassword(requestJSON.Token, requestJSON.Password); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Password reseted successfully!"})

}

func CheckPassPhrase(c *gin.Context) {
	var requestJSON dto.PassPhraseDTO

	if err := c.ShouldBindJSON(&requestJSON); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": err.Error()})
		return
	}

	user, _ := c.Get("currentUser")

	err := bcrypt.CompareHashAndPassword([]byte(user.(models.User).PassPhrase), []byte(requestJSON.PassPhrase))

	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"message": "Passphrase mismatch. Invalid passphrase",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Passphrase is valid",
	})

}

func GetUserNotifications(c *gin.Context) {
	user, _ := c.Get("currentUser")
	var n []models.Notification
	db.GetDB().Model(models.Notification{}).Where("user_id = ?", user.(models.User).ID).Where("is_read", false).Find(&n)
	c.JSON(http.StatusOK, gin.H{"notifications": n})
}

func MarkNotificationAsRead(c *gin.Context) {
	user, _ := c.Get("currentUser")
	db.GetDB().Model(models.Notification{}).Where("user_id = ?", user.(models.User).ID).Where("is_read", false).Update("is_read", true)
	c.JSON(http.StatusOK, gin.H{"message": "Notifications set as read"})
}

func UpdateUserProfile(c *gin.Context) {
	var requestJSON dto.UpdateUserProfileRequestDTO

	if err := c.ShouldBindJSON(&requestJSON); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": err.Error()})
		return
	}

	user_, _ := c.Get("currentUser")
	user, err := auth.GetUserByUsername(user_.(models.User).Username)

	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"message": err.Error(),
		})
		return
	}

	if requestJSON.TimoutTime < 5 || requestJSON.TimoutTime > 30 {
		c.JSON(http.StatusBadRequest, gin.H{
			"message": "Timeout time should be within 5 miniutes to 30 minutes",
		})
		return
	}

	user.SessionTime = requestJSON.TimoutTime
	if err := db.GetDB().Save(&user).Error; err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"message": err.Error(),
		})
	}

}
