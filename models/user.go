// package models

// import (
// 	"errors"
// 	"task_manager/config"
// 	"task_manager/utils"
// 	"time"
// )

// type User struct {
// 	ID        int64
// 	Name      string `binding:"required" json:"name"`
// 	Mobile_No int64  `binding:"required" json:"mob_no"`
// 	Gender    string `binding:"required" json:"gender"`
// 	Email     string `binding:"required" json:"email"`
// 	//Image []byte `binding:"required"`
// 	Password         string `binding:"required" json:"password"`
// 	Confirm_Password string `binding:"required" json:"confirm_password"`
// 	//Timestamp time.Time `binding:"required"`
// }

// type UserResponse struct {
// 	Name      string `binding:"required"`
// 	Mobile_No int64  `binding:"required"`
// 	Gender    string `binding:"required"`
// 	Email     string `binding:"required"`
// 	User_id   int64  `binding:"required"`
// 	//Image []byte `binding:"required"`
// 	//Password string `binding:"required"`
// }

// type Login struct {
// 	ID       int64
// 	Email    string `binding:"required" json:"username"`
// 	Password string `binding:"required" json:"password"`
// 	// User_id  int64  `binding:"required"`
// }

// type Token struct {
// 	ID            int64
// 	Refresh_Token string `binding:"required"`
// 	JWT_Token     string `binding:"required"`
// 	// User_id int64 `binding:"required"`
// 	Timestamp time.Time `binding:"required"`
// }

// func (u User) Save() (int64, error) {
// 	query := "INSERT INTO users(name, mobile_no, gender, email) VALUES (?,?,?,?)"

// 	stmt, err := config.DB.Prepare(query)
// 	if err != nil {
// 		return 0, err
// 	}
// 	defer stmt.Close()

// 	result, err := stmt.Exec(u.Name, u.Mobile_No, u.Gender, u.Email)
// 	if err != nil {
// 		return 0, err
// 	}

// 	userId, err := result.LastInsertId()
// 	if err != nil {
// 		return 0, err
// 	}
// 	u.ID = userId

// 	loginQuery := "INSERT INTO login (email, password, user_id) VALUES(?,?,?)"

// 	loginstmt, err := config.DB.Prepare(loginQuery)
// 	if err != nil {
// 		return 0, err
// 	}
// 	defer loginstmt.Close()

// 	hashedpassword, err := utils.HashPassword(u.Password)
// 	if err != nil {
// 		return 0, err
// 	}

// 	_, err = loginstmt.Exec(u.Email, hashedpassword, u.ID)
// 	if err != nil {
// 		return 0, err
// 	}

// 	return u.ID, err
// }

// func GetAllUsers() ([]UserResponse, error) {
// 	query := "SELECT * FROM users"

// 	rows, err := config.DB.Query(query)
// 	if err != nil {
// 		return nil, err
// 	}
// 	defer rows.Close()

// 	var usersR []UserResponse
// 	for rows.Next() {
// 		var userR UserResponse

// 		err := rows.Scan(&userR.User_id, &userR.Name, &userR.Mobile_No, &userR.Gender, &userR.Email)
// 		if err != nil {
// 			return nil, err
// 		}

// 		usersR = append(usersR, userR)
// 	}

// 	return usersR, nil
// }

// // func GetAllLOGIN() ([]Login, error) {
// // 	query := "SELECT * FROM login"

// // 	rows, err := config.DB.Query(query)
// // 	if err != nil {
// // 		return nil, err
// // 	}
// // 	defer rows.Close()

// // 	var logins []Login
// // 	for rows.Next() {
// // 		var login Login

// // 		err := rows.Scan(&login.ID, &login.Email, &login.Password, &login.User_id)
// // 		if err != nil {
// // 			return nil, err
// // 		}

// // 		logins = append(logins, login)
// // 	}

// // 	return logins, nil
// // }

// func (u User) SaveToken(jwt_token string, refresh_token string) error {
// 	query := "INSERT INTO token (refresh_token, jwt_token, timestamp) VALUES (?,?,?)"

// 	stmt, err := config.DB.Prepare(query)
// 	if err != nil {
// 		return err
// 	}
// 	defer stmt.Close()

// 	_, err = stmt.Exec(jwt_token, refresh_token, time.Now())
// 	if err != nil {
// 		return err
// 	}

// 	return err
// }

// // func GetAllToken() ([]Token, error) {
// // 	query := "SELECT * FROM token"

// // 	rows, err := config.DB.Query(query)
// // 	if err != nil {
// // 		return nil, err
// // 	}
// // 	defer rows.Close()

// // 	var tokens []Token
// // 	for rows.Next() {
// // 		var token Token

// // 		err := rows.Scan(&token.ID, &token.Refresh_Token, &token.JWT_Token, &token.Timestamp)
// // 		if err != nil {
// // 			return nil, err
// // 		}

// // 		tokens = append(tokens, token)
// // 	}

// // 	return tokens, nil
// // }

// func (u *Login) ValidateCredentials() error {
// 	query := "SELECT id, password FROM login WHERE email = ?"

// 	row := config.DB.QueryRow(query, u.Email)

// 	var retrievedpassword string
// 	err := row.Scan(&u.ID, &retrievedpassword)
// 	if err != nil {
// 		return errors.New("invalid credentials")
// 	}

// 	passwordIsValid := utils.CheckPasswordHash(u.Password, retrievedpassword)

// 	if !passwordIsValid {
// 		return errors.New("invalid credentials")
// 	}

// 	return nil
// }

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

// type Token struct {
// 	ID            int64     `gorm:"primaryKey"`
// 	Refresh_Token string    `gorm:"not null"`
// 	JWT_Token     string    `gorm:"not null"`
// 	Timestamp     time.Time `gorm:"not null"`
// }

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

// func GetAllUsers() ([]UserResponse, error) {
// 	var users []User
// 	if err := config.DB.Find(&users).Error; err != nil {
// 		return nil, err
// 	}
// 	var usersR []UserResponse
// 	for _, user := range users {
// 		usersR = append(usersR, UserResponse{Name: user.Name, Mobile_No: user.Mobile_No, Gender: user.Gender, Email: user.Email, User_id: user.ID})
// 	}
// 	return usersR, nil
// }

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
