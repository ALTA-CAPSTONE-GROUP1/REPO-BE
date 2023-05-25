package usecase_test

// import (
// 	"errors"
// 	"testing"

// 	"github.com/ALTA-CAPSTONE-GROUP1/e-proposal-BE/feature/admin/user"
// 	"github.com/ALTA-CAPSTONE-GROUP1/e-proposal-BE/feature/admin/user/mocks"
// 	"github.com/ALTA-CAPSTONE-GROUP1/e-proposal-BE/feature/admin/user/usecase"
// 	"github.com/stretchr/testify/assert"
// )

// func TestDeleteUser(t *testing.T) {
// 	repo := mocks.NewRepository(t)
// 	ul := usecase.New(repo)

// 	t.Run("Success", func(t *testing.T) {
// 		userID := "12345"

// 		repo.On("DeleteUser", userID).Return(nil).Once()
// 		err := ul.DeleteUser(userID)

// 		assert.NoError(t, err)
// 		repo.AssertExpectations(t)
// 	})

// 	t.Run("User Not Found", func(t *testing.T) {
// 		userID := "12345"
// 		errMsg := "finding user error"

// 		repo.On("DeleteUser", userID).Return(errors.New(errMsg)).Once()
// 		err := ul.DeleteUser(userID)

// 		assert.EqualError(t, err, "bad request, user not found")
// 		repo.AssertExpectations(t)
// 	})

// 	t.Run("Cannot Delete User", func(t *testing.T) {
// 		userID := "12345"
// 		errMsg := "cannot delete user error"

// 		repo.On("DeleteUser", userID).Return(errors.New(errMsg)).Once()
// 		err := ul.DeleteUser(userID)

// 		assert.EqualError(t, err, "internal server error, cannot delete user")
// 		repo.AssertExpectations(t)
// 	})

// 	t.Run("Other Error", func(t *testing.T) {
// 		userID := "12345"
// 		errMsg := "other error"

// 		repo.On("DeleteUser", userID).Return(errors.New(errMsg)).Once()
// 		err := ul.DeleteUser(userID)

// 		assert.EqualError(t, err, errMsg)
// 		repo.AssertExpectations(t)
// 	})
// }

// func TestUpdateUser(t *testing.T) {
// 	repo := mocks.NewRepository(t)
// 	ul := usecase.New(repo)

// 	t.Run("Success", func(t *testing.T) {
// 		userID := "12345"
// 		updateUser := user.Core{
// 			ID:   userID,
// 			Name: "Adi Yuda",
// 		}

// 		repo.On("UpdateUser", userID, updateUser).Return(nil).Once()
// 		err := ul.UpdateUser(userID, updateUser)

// 		assert.NoError(t, err)
// 		repo.AssertExpectations(t)
// 	})

// 	t.Run("Invalid Password", func(t *testing.T) {
// 		userID := "12345"
// 		updateUser := user.Core{
// 			ID:       userID,
// 			Name:     "Adi Yuda",
// 			Password: "weak",
// 		}

// 		errMsg := "hashing password error"
// 		repo.On("UpdateUser", userID, updateUser).Return(errors.New(errMsg)).Once()
// 		err := ul.UpdateUser(userID, updateUser)

// 		assert.EqualError(t, err, "is invalid")
// 		repo.AssertExpectations(t)
// 	})

// 	t.Run("No Rows Affected", func(t *testing.T) {
// 		userID := "12345"
// 		updateUser := user.Core{
// 			ID:   userID,
// 			Name: "Adi Yuda",
// 		}

// 		errMsg := "no rows affected on update user"
// 		repo.On("UpdateUser", userID, updateUser).Return(errors.New(errMsg)).Once()
// 		err := ul.UpdateUser(userID, updateUser)

// 		assert.EqualError(t, err, "data is up to date")
// 		repo.AssertExpectations(t)
// 	})

// 	t.Run("Other Error", func(t *testing.T) {
// 		userID := "12345"
// 		updateUser := user.Core{
// 			ID:   userID,
// 			Name: "Adi Yuda",
// 		}

// 		errMsg := "other error"
// 		repo.On("UpdateUser", userID, updateUser).Return(errors.New(errMsg)).Once()
// 		err := ul.UpdateUser(userID, updateUser)

// 		assert.EqualError(t, err, errMsg)
// 		repo.AssertExpectations(t)
// 	})
// }

// func TestGetUserById(t *testing.T) {
// 	repo := mocks.NewRepository(t)
// 	ul := usecase.New(repo)

// 	t.Run("Success", func(t *testing.T) {
// 		userID := "12345"
// 		expectedUser := user.Core{
// 			ID:   userID,
// 			Name: "Adi Yuda",
// 		}

// 		repo.On("GetUserById", userID).Return(expectedUser, nil).Once()
// 		user, err := ul.GetUserById(userID)

// 		assert.NoError(t, err)
// 		assert.Equal(t, expectedUser, user)
// 		repo.AssertExpectations(t)
// 	})

// 	t.Run("Error", func(t *testing.T) {
// 		userID := "12345"
// 		errMsg := "failed to find user error"

// 		repo.On("GetUserById", userID).Return(user.Core{}, errors.New(errMsg)).Once()
// 		userRes, err := ul.GetUserById(userID)

// 		assert.EqualError(t, err, "internal server error")
// 		assert.Equal(t, user.Core{}, userRes)
// 		repo.AssertExpectations(t)
// 	})
// }

// func TestGetAllUser(t *testing.T) {
// 	repo := mocks.NewRepository(t)
// 	ul := usecase.New(repo)

// 	t.Run("Success", func(t *testing.T) {
// 		limit := 10
// 		offset := 0
// 		name := "John"
// 		expectedUsers := []user.Core{
// 			{ID: "1", Name: "Adi Yuda"},
// 			{ID: "2", Name: "John Smith"},
// 		}

// 		repo.On("SelectAllUser", limit, offset, name).Return(expectedUsers, nil).Once()
// 		users, err := ul.GetAllUser(limit, offset, name)

// 		assert.NoError(t, err)
// 		assert.Equal(t, expectedUsers, users)
// 		repo.AssertExpectations(t)
// 	})

// 	t.Run("Error", func(t *testing.T) {
// 		limit := 10
// 		offset := 0
// 		name := "John"
// 		errMsg := "failed to find all user error"

// 		repo.On("SelectAllUser", limit, offset, name).Return([]user.Core{}, errors.New(errMsg)).Once()
// 		users, err := ul.GetAllUser(limit, offset, name)

// 		assert.EqualError(t, err, "internal server error")
// 		assert.Equal(t, []user.Core{}, users)
// 		repo.AssertExpectations(t)
// 	})
// }

// func TestRegisterUser(t *testing.T) {
// 	repo := mocks.NewRepository(t)
// 	ul := usecase.New(repo)

// 	t.Run("Success", func(t *testing.T) {
// 		newUser := user.Core{
// 			ID:       "1",
// 			Name:     "Adi Yuda",
// 			Email:    "johndoe@example.com",
// 			Password: "password",
// 		}

// 		repo.On("InsertUser", newUser).Return(nil).Once()
// 		err := ul.RegisterUser(newUser)

// 		assert.NoError(t, err)
// 		repo.AssertExpectations(t)
// 	})

// 	t.Run("Column Error", func(t *testing.T) {
// 		newUser := user.Core{
// 			ID:       "1",
// 			Name:     "Adi Yuda",
// 			Email:    "johndoe@example.com",
// 			Password: "password",
// 		}

// 		errMsg := "column error"
// 		repo.On("InsertUser", newUser).Return(errors.New(errMsg)).Once()
// 		err := ul.RegisterUser(newUser)

// 		assert.EqualError(t, err, "server error")
// 		repo.AssertExpectations(t)
// 	})

// 	t.Run("Value Error", func(t *testing.T) {
// 		newUser := user.Core{
// 			ID:       "1",
// 			Name:     "Adi Yuda",
// 			Email:    "johndoe@example.com",
// 			Password: "password",
// 		}

// 		errMsg := "value error"
// 		repo.On("InsertUser", newUser).Return(errors.New(errMsg)).Once()
// 		err := ul.RegisterUser(newUser)

// 		assert.EqualError(t, err, "invalid value")
// 		repo.AssertExpectations(t)
// 	})

// 	t.Run("Invalid Password Length", func(t *testing.T) {
// 		newUser := user.Core{
// 			ID:       "1",
// 			Name:     "Adi Yuda",
// 			Email:    "johndoe@example.com",
// 			Password: "short",
// 		}

// 		errMsg := "password too short error"
// 		repo.On("InsertUser", newUser).Return(errors.New(errMsg)).Once()
// 		err := ul.RegisterUser(newUser)

// 		assert.EqualError(t, err, "invalid password length")
// 		repo.AssertExpectations(t)
// 	})

// 	t.Run("Other Error", func(t *testing.T) {
// 		newUser := user.Core{
// 			ID:       "1",
// 			Name:     "Adi Yuda",
// 			Email:    "johndoe@example.com",
// 			Password: "password",
// 		}

// 		errMsg := "other error"
// 		repo.On("InsertUser", newUser).Return(errors.New(errMsg)).Once()
// 		err := ul.RegisterUser(newUser)

// 		assert.EqualError(t, err, "server error")
// 		repo.AssertExpectations(t)
// 	})
// }
