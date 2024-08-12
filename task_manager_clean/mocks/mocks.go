package mocks

import (
	domain "task_manager/domain"
	"time"

	primitive "go.mongodb.org/mongo-driver/bson/primitive"
)

func GetNewTask() *domain.Task {
	return &domain.Task{
		ID:          GetPrimitiveID1(),
		Title:       "My First Task",
		Description: "This is some example description for the first task.",
		DueDate:     format(time.Now().AddDate(0, 0, 3)),
		Status:      "In Progress",
		UserID:      GetPrimitiveID2(),
	}
}

func GetCreateTaskData() *domain.CreateTaskData {
	return &domain.CreateTaskData{
		Title:       "My First Task",
		Description: "This is some example description for the first task.",
		DueDate:     format(time.Now().AddDate(0, 0, 3)),
		Status:      "In Progress",
	}
}

func GetReplaceTaskData() *domain.ReplaceTaskData {
	return &domain.ReplaceTaskData{
		Title:       "My New Task",
		Description: "This is some example description for the new task.",
		DueDate:     format(time.Now().AddDate(0, 0, 3)),
		Status:      "Pending",
	}
}

func GetUpdateTaskData() *domain.UpdateTaskData {
	return &domain.UpdateTaskData{
		Title:       "My New Task",
		Description: "This is some example description for the new task.",
		DueDate:     format(time.Now().AddDate(0, 0, 3)),
		Status:      "Pending",
	}
}

func GetTaskView(task *domain.Task) *domain.TaskView {
	return &domain.TaskView{
		ID:          task.ID.Hex(),
		Title:       task.Title,
		Description: task.Description,
		DueDate:     task.DueDate,
		Status:      task.Status,
	}
}

func GetNewTask2() *domain.Task {
	return &domain.Task{
		ID:          GetPrimitiveID2(),
		Title:       "My Second Task",
		Description: "This is some example description for the second task.",
		DueDate:     format(time.Now().AddDate(0, 0, 1)),
		Status:      "Completed",
		UserID:      GetPrimitiveID1(),
	}
}

func GetManyTasks() []domain.Task {
	return []domain.Task{
		*GetNewTask(),
		*GetNewTask2(),
		{
			ID:          GetPrimitiveID3(),
			Title:       "My Third Task",
			Description: "This is some example description for the third task.",
			DueDate:     format(time.Now().AddDate(0, 0, 5)),
			Status:      "In Progress",
			UserID:      GetPrimitiveID2(),
		},
	}
}

func GetAuthUserData() *domain.AuthUserData {
	return &domain.AuthUserData{
		Username: "user1",
		Password: "password1",
	}
}

func GetUser3(userData *domain.AuthUserData) *domain.User {
	return &domain.User{
		ID:       GetNextID(primitive.NewObjectID()),
		Username: userData.Username,
		Password: userData.Password,
		Role:     "user",
	}
}

func GetAuthData(user *domain.User) *domain.AuthUserData {
	return &domain.AuthUserData{
		Username: user.Username,
		Password: user.Password,
	}
}

func format(d time.Time) time.Time {
	return d.UTC().Truncate(time.Millisecond)
}

func GetNewUser() *domain.User {
	return &domain.User{
		ID:       primitive.NewObjectID(),
		Username: "user1",
		Password: "password1",
		Role:     "user",
	}
}

func GetNewUser2() *domain.User {
	return &domain.User{
		ID:       primitive.NewObjectID(),
		Username: "user2",
		Password: "password2",
		Role:     "user",
	}
}

func GetCreateUserData() *domain.CreateUserData {
	return &domain.CreateUserData{
		Username: "user1",
		Password: "password1",
		Role:     "user",
	}
}

func GetClaims() *domain.Claims {
	return &domain.Claims{
		ID:       GetPrimitiveID1(),
		Username: "user1",
		Password: "password1",
		Role:     "user",
	}
}

func GetClaims2() *domain.Claims {
	return &domain.Claims{
		ID:       GetPrimitiveID2(),
		Username: "user2",
		Password: "password2",
		Role:     "admin",
	}
}

func GetClaims3() *domain.Claims {
	return &domain.Claims{
		ID:       GetPrimitiveID3(),
		Username: "user3",
		Password: "password3",
		Role:     "root",
	}
}

func GetTask(taskData *domain.CreateTaskData, claims *domain.Claims) *domain.Task {
	return &domain.Task{
		ID:          claims.ID,
		Title:       taskData.Title,
		Description: taskData.Description,
		DueDate:     taskData.DueDate,
		Status:      taskData.Status,
		UserID:      claims.ID,
	}
}

func GetTask2(taskData *domain.ReplaceTaskData, claims *domain.Claims) *domain.Task {
	return &domain.Task{
		ID:          claims.ID,
		Title:       taskData.Title,
		Description: taskData.Description,
		DueDate:     taskData.DueDate,
		Status:      taskData.Status,
		UserID:      claims.ID,
	}
}

func GetTask3(taskData *domain.UpdateTaskData, claims *domain.Claims) *domain.Task {
	return &domain.Task{
		ID:          claims.ID,
		Title:       taskData.Title,
		Description: taskData.Description,
		DueDate:     taskData.DueDate,
		Status:      taskData.Status,
		UserID:      claims.ID,
	}
}

func GetUser(userData *domain.CreateUserData) *domain.User {
	return &domain.User{
		ID:       GetNextID(primitive.NewObjectID()),
		Username: userData.Username,
		Password: userData.Password,
		Role:     userData.Role,
	}
}

func GetUser2(claims *domain.Claims) *domain.User {
	return &domain.User{
		ID:       claims.ID,
		Username: claims.Username,
		Password: claims.Password,
		Role:     claims.Role,
	}
}

func GetPrimitiveID1() primitive.ObjectID {
	id, err := primitive.ObjectIDFromHex("60f1b3b3b3f3b3f3b3f3b3f3")
	if err != nil {
		panic(err)
	}

	return id
}

func GetPrimitiveID2() primitive.ObjectID {
	id, err := primitive.ObjectIDFromHex("60f1b3b3b3f3b3f3b3f3b3f4")
	if err != nil {
		panic(err)
	}

	return id
}

func GetPrimitiveID3() primitive.ObjectID {
	id, err := primitive.ObjectIDFromHex("60f1b3b3b3f3b3f3b3f3b3f5")
	if err != nil {
		panic(err)
	}

	return id
}

func GetNextID(id primitive.ObjectID) primitive.ObjectID {
	bytes := []byte(id.Hex())

	for i := len(bytes) - 1; i >= 0; i-- {
		if bytes[i] == 'f' {
			bytes[i] = '0'
		} else {
			if bytes[i] == '9' {
				bytes[i] = 'a'
			} else {
				bytes[i]++
			}

			break
		}
	}

	newID, err := primitive.ObjectIDFromHex(string(bytes))
	if err != nil {
		panic(err)
	}

	return newID
}

func GetManyUsers() []domain.User {
	return []domain.User{
		*GetNewUser(),
		*GetNewUser2(),
		{
			ID:       GetPrimitiveID3(),
			Username: "user3",
			Password: "password3",
			Role:     "root",
		},
	}
}

func GetUpdateUserData() *domain.UpdateUserData {
	return &domain.UpdateUserData{
		Username: "user2",
		Password: "password2",
		Role:     "user",
	}
}

func GetUser4(userData *domain.UpdateUserData) *domain.User {
	return &domain.User{
		ID:       GetPrimitiveID2(),
		Username: userData.Username,
		Password: userData.Password,
		Role:     userData.Role,
	}
}

func GetView(taskData *domain.CreateTaskData, claims *domain.Claims) *domain.TaskView {
	return &domain.TaskView{
		ID:          claims.ID.Hex(),
		Title:       taskData.Title,
		Description: taskData.Description,
		DueDate:     taskData.DueDate,
		Status:      taskData.Status,
	}
}

func GetView2(taskData *domain.ReplaceTaskData, claims *domain.Claims) *domain.TaskView {
	return &domain.TaskView{
		ID:          claims.ID.Hex(),
		Title:       taskData.Title,
		Description: taskData.Description,
		DueDate:     taskData.DueDate,
		Status:      taskData.Status,
	}
}

func GetView3(taskData *domain.UpdateTaskData, claims *domain.Claims) *domain.TaskView {
	return &domain.TaskView{
		ID:          claims.ID.Hex(),
		Title:       taskData.Title,
		Description: taskData.Description,
		DueDate:     taskData.DueDate,
		Status:      taskData.Status,
	}
}
