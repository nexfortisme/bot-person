package persistance

import (
	persistance "main/pkg/persistance/eums"
)

func GetUserAttribute(userId string, attribute persistance.Attribute) (string, error) {
	
	query := `SELECT * FROM UserAttributes WHERE UserId = ? AND Attribute = ?`

	var userAttribute UserAttribute
	err := RunQuery(query, &userAttribute, userId, attribute.String())
	if err != nil {
		return "", err
	}

	return userAttribute.Value, nil
}

func SetUserAttribute(userId string, attribute persistance.Attribute, value string) error {
	// Check if the attribute already exists for this user
	existingValue, err := GetUserAttribute(userId, attribute)
	
	if err == nil && existingValue != "" {
		// Attribute exists, update it
		query := `UPDATE UserAttributes SET Value = ? WHERE UserId = ? AND Attribute = ?`
		err = RunQuery(query, nil, value, userId, attribute.String())
		if err != nil {
			return err
		}
	} else {
		// Attribute doesn't exist, create it
		query := `INSERT INTO UserAttributes (UserId, Attribute, Value) VALUES (?, ?, ?)`
		err = RunQuery(query, nil, userId, attribute.String(), value)
		if err != nil {
			return err
		}
	}

	return nil
}