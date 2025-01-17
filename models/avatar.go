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

func UpdateAvatar(uid int64, content []byte, fileName string) error {
	return config.DB.Model(&config.Avatar{}).
		Where("user_id = ?", uid).
		Updates(map[string]interface{}{
			"data": content,
			"name": fileName,
		}).Error
}

func DeleteAvatar(uid int64) error {
	var avatar config.Avatar

	result := config.DB.Where("user_id = ?", uid).Delete(&avatar)
	if result.RowsAffected == 0 {
		return nil
	}

	if result.Error != nil {
		return result.Error
	}
	return nil
}
