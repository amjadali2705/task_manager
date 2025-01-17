package models

import "task_manager/config"

func ReadAvatar(uid int64) (*config.Avatar, error) {
	var avatar config.Avatar

	if err := config.DB.Where("user_id = ?", uid).First(&avatar).Error; err != nil {
		return nil, err
	}
	return &avatar, nil
}

func SaveAvatar(uid int64, content []byte, fileName string) error {
	avatar := config.Avatar{
		Data:   content,
		Name:   fileName,
		UserID: uid,
	}

	if err := config.DB.Create(&avatar).Error; err != nil {
		return err
	}

	return nil
}
