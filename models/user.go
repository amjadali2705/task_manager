package models

import (
	"errors"
	"task_manager/config"
	"task_manager/utils"
	"time"

	"go.uber.org/zap"
)

type User struct {
	ID               int64
	Name             string `json:"name"`
	Mobile_No        int64  `json:"mob_no"`
	Gender           string `json:"gender"`
	Email            string `json:"email"`
	Password         string `json:"password"`
	Confirm_Password string `json:"confirm_password"`
}

type UserResponse struct {
	Name      string `json:"name"`
	Mobile_No int64  `json:"mob_no"`
	Gender    string `json:"gender"`
	Email     string `json:"email"`
	User_id   int64  `json:"user_id"`
}

type Login struct {
	ID       int64
	Email    string `json:"username"`
	Password string `json:"password"`
	UserID   int64
}

type UpdateUserRequest struct {
	Name      string `json:"name" binding:"required"`
	Mobile_No int64  `json:"mobile_no" binding:"required"`
}

type UpdatePasswordRequest struct {
	OldPassword string `json:"old_password" binding:"required"`
	NewPassword string `json:"new_password" binding:"required"`
}

func (u *User) Save() (int64, error) {
	user := config.User{Name: u.Name, MobileNo: u.Mobile_No, Gender: u.Gender, Email: u.Email}

	if err := config.DB.Create(&user).Error; err != nil {
		return 0, err
	}

	u.ID = user.ID

	hashedpassword, err := utils.HashPassword(u.Password)
	if err != nil {
		return 0, err
	}

	login := Login{Email: u.Email, Password: hashedpassword, UserID: u.ID}
	if err := config.DB.Create(&login).Error; err != nil {
		return 0, err
	}
	return u.ID, nil
}

func (u *User) SaveToken(user_token, refresh_token string) error {
	token := config.Token{UserToken: user_token, RefreshToken: refresh_token, Timestamp: time.Now()}

	if err := config.DB.Create(&token).Error; err != nil {
		utils.Logger.Error("Failed to save token", zap.Error(err))
		return err
	}

	utils.Logger.Info("Token saved successfully")
	return nil
}

func (u *Login) ValidateCredentials() error {
	var login Login
	if err := config.DB.Where("email = ?", u.Email).First(&login).Error; err != nil {
		utils.Logger.Warn("Invalid credentials provided")
		return errors.New("invalid credentials")
	}

	passwordIsValid := utils.CheckPasswordHash(u.Password, login.Password)
	if !passwordIsValid {
		utils.Logger.Warn("Password mismatch")
		return errors.New("invalid credentials")
	}

	utils.Logger.Info("User credentials validated successfully")
	u.ID = login.ID
	return nil
}

func GetUserByEmail(email string) (*User, error) {
	var user User
	if err := config.DB.Where("email = ?", email).First(&user).Error; err != nil {
		utils.Logger.Error("Failed to fetch user by email", zap.Error(err))
		return nil, err
	}

	utils.Logger.Info("User fetched by email successfully")
	return &user, nil
}

func DeleteToken(tokenString string) error {
	var token config.Token
	if err := config.DB.Where("user_token = ?", tokenString).First(&token).Error; err != nil {
		utils.Logger.Error("Token not found for deletion", zap.Error(err))
		return err
	}

	if err := config.DB.Delete(&token).Error; err != nil {
		utils.Logger.Error("Failed to delete token", zap.Error(err))
		return err
	}

	utils.Logger.Info("Token deleted successfully")
	return nil
}

func DeleteRefreshToken(tokenString string) error {
	var token config.Token
	if err := config.DB.Where("refresh_token = ?", tokenString).First(&token).Error; err != nil {
		utils.Logger.Error("Refresh Token not found for deletion", zap.Error(err))
		return err
	}

	if err := config.DB.Delete(&token).Error; err != nil {
		utils.Logger.Error("Failed to delete refresh token", zap.Error(err))
		return err
	}

	utils.Logger.Info("Refresh Token deleted successfully")
	return nil
}

func (u *User) UpdateUserTable() error {
	if err := config.DB.Model(&config.User{}).Where("email = ?", u.Email).Updates(map[string]interface{}{"name": u.Name, "mobile_no": u.Mobile_No, "gender": u.Gender}).Error; err != nil {
		utils.Logger.Error("Failed to update user table", zap.Error(err))
		return err
	}

	utils.Logger.Info("User table updated successfully")
	return nil
}

func UpdateUserDetails(uid int64, req UpdateUserRequest) error {
	result := config.DB.Model(&config.User{}).Where("id = ?", uid).Updates(map[string]interface{}{"name": req.Name, "mobile_no": req.Mobile_No})
	if result.Error != nil {
		return result.Error
	}

	return nil
}

func GetUserByIdPassChng(uid int64) (*config.Login, error) {
	var login config.Login

	if err := config.DB.Where("id = ?", uid).First(&login).Error; err != nil {
		return nil, err
	}
	return &login, nil
}

func UpdatePassById(uid int64, password string) error {
	if err := config.DB.Model(&config.Login{}).Where("user_id = ?", uid).Update("password", password).Error; err != nil {
		return err
	}

	return nil
}

func DeleteTokenById(uid int64, tokenString string) error {
	var token config.Token

	result := config.DB.Where("user_id = ? AND user_token != ?", uid, tokenString).Delete(&token)
	if result.RowsAffected == 0 {
		return nil
	}

	if result.Error != nil {
		return result.Error
	}
	return nil
}
