package repositories

// Проверяем, что пользователь существует по ID
func (r *UserRepository) isUserExists(userID int) (bool, error) {
	var userExists bool
	err := r.db.Get(&userExists, `SELECT EXISTS(SELECT 1 FROM users WHERE id = $1)`, userID)
	if err != nil {
		return false, err
	}

	return userExists, nil
}

// получаем id юзера по username
func (r *UserRepository) getUserIdByUsername(username string) (int, error) {
	var userID int
	err := r.db.Get(&userID, `SELECT id FROM users WHERE username = $1`, username)
	if err != nil {
		return 0, err
	}

	return userID, nil
}
