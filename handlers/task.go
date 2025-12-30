package handlers

import (
	"bytes"
	"io"
	"log"
	"strconv"

	"github.com/gin-gonic/gin"

	"task/database"
	"task/models"
)


// Membuat tugas baru
func CreateTask(c *gin.Context) {
	var task models.Task

	if err := c.ShouldBindJSON(&task); err != nil {
		log.Println("[CreateTask] Bind JSON error:", err)
		c.JSON(400, gin.H{
			"error":   "invalid request body",
			"details": err.Error(),
		})
		return
	}

	if task.ID == 0 || task.UserID == 0 {
		log.Println("[CreateTask] Missing ID or UserID")
		c.JSON(400, gin.H{
			"error": "id and user_id are required",
		})
		return
	}

	// cek user
	userResp, err := database.SupabaseRequest(
		"GET",
		"/user?id=eq."+strconv.FormatInt(task.UserID, 10),
		nil,
	)
	if err != nil {
		log.Println("[CreateTask] Supabase GET user error:", err)
		c.JSON(500, gin.H{"error": err.Error()})
		return
	}
	defer userResp.Body.Close()

	userBody, err := io.ReadAll(userResp.Body)
	if err != nil {
		log.Println("[CreateTask] Read user body error:", err)
		c.JSON(500, gin.H{"error": "failed to read user response"})
		return
	}
	// 400 - user not found
	if string(userBody) == "[]" {
		log.Println("[CreateTask] User not found:", task.UserID)
		c.JSON(404, gin.H{
			"error": "user not found",
		})
		return
	}

	resp, err := database.SupabaseRequest(
		"POST",
		"/task",
		task,
	)
	if err != nil {
		log.Println("[CreateTask] Supabase POST task error:", err)
		c.JSON(500, gin.H{"error": err.Error()})
		return
	}
	defer resp.Body.Close()

body, err := io.ReadAll(resp.Body)
if err != nil {
	log.Println("[CreateTask] Read task response error:", err)
	c.JSON(500, gin.H{"error": "failed to read task response"})
	return
}

// cek status code
if resp.StatusCode >= 400 {

	// 404 - foreign key violation (user not found)
	if bytes.Contains(body, []byte("23503")) {
		log.Println("[CreateTask] Foreign key violation (user not found)")
		c.JSON(404, gin.H{
			"error": "user not found",
		})
		return
	}

	// error lain
	log.Println("[CreateTask] Supabase error:", string(body))
	c.JSON(resp.StatusCode, gin.H{
		"error": "failed to create task",
	})
	return
}

// 200 - sukses
c.Data(resp.StatusCode, "application/json", body)
}

// Membuat task baru untuk user tertentu
func CreateTaskByUser(c *gin.Context) {
	userID := c.Param("id")
	var task models.Task

	// 400 — invalid body
	if err := c.ShouldBindJSON(&task); err != nil {
		log.Println("[CreateTaskByUser] Bind JSON error:", err)
		c.JSON(400, gin.H{
			"error":   "invalid request body",
			"details": err.Error(),
		})
		return
	}

	// set user_id dari URL
	uid, err := strconv.ParseInt(userID, 10, 64)
	if err != nil {
		c.JSON(400, gin.H{"error": "invalid user id"})
		return
	}
	task.UserID = uid

	// cek user ada atau tidak
	userResp, err := database.SupabaseRequest(
		"GET",
		"/user?id=eq."+userID,
		nil,
	)
	if err != nil {
		log.Println("[CreateTaskByUser] Supabase GET user error:", err)
		c.JSON(500, gin.H{"error": "failed to check user"})
		return
	}
	defer userResp.Body.Close()

	userBody, _ := io.ReadAll(userResp.Body)
	if string(userBody) == "[]" {
		c.JSON(404, gin.H{
			"error": "user not found",
		})
		return
	}

	// create task
	resp, err := database.SupabaseRequest(
		"POST",
		"/task",
		task,
	)
	if err != nil {
		log.Println("[CreateTaskByUser] Supabase POST error:", err)
		c.JSON(500, gin.H{"error": "failed to create task"})
		return
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		c.JSON(500, gin.H{"error": "failed to read response"})
		return
	}

	// handle FK error
	if resp.StatusCode >= 400 {
		if bytes.Contains(body, []byte("23503")) {
			c.JSON(404, gin.H{"error": "user not found"})
			return
		}

		c.JSON(resp.StatusCode, gin.H{
			"error": "failed to create task",
		})
		return
	}

	// 200 — sukses
	c.Data(resp.StatusCode, "application/json", body)
}

// Mengambil daftar tugas
func GetTasks(c *gin.Context) {
	status := c.Query("status")
	path := "/task?select=*"

	if status != "" {
		path += "&status=eq." + status
	}

	resp, err := database.SupabaseRequest("GET", path, nil)
	if err != nil {
		log.Println("[GetTasks] Supabase GET error:", err)
		c.JSON(500, gin.H{"error": err.Error()})
		return
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		log.Println("[GetTasks] Read body error:", err)
		c.JSON(500, gin.H{"error": "failed to read response"})
		return
	}
	// 200 - sukses
	c.Data(200, "application/json", body)
}

// Mengambil semua task milik user tertentu
func GetUserTasks(c *gin.Context) {
	userID := c.Param("id")

	// cek user ada atau tidak
	userResp, err := database.SupabaseRequest(
		"GET",
		"/user?id=eq."+userID,
		nil,
	)
	if err != nil {
		log.Println("[GetUserTasks] Supabase GET user error:", err)
		c.JSON(500, gin.H{"error": "failed to check user"})
		return
	}
	defer userResp.Body.Close()

	userBody, err := io.ReadAll(userResp.Body)
	if err != nil {
		c.JSON(500, gin.H{"error": "failed to read user response"})
		return
	}

	if string(userBody) == "[]" {
		c.JSON(404, gin.H{
			"error": "user not found",
		})
		return
	}

	// ambil task milik user
	taskResp, err := database.SupabaseRequest(
		"GET",
		"/task?user_id=eq."+userID,
		nil,
	)
	if err != nil {
		log.Println("[GetUserTasks] Supabase GET task error:", err)
		c.JSON(500, gin.H{"error": "failed to fetch tasks"})
		return
	}
	defer taskResp.Body.Close()

	taskBody, err := io.ReadAll(taskResp.Body)
	if err != nil {
		c.JSON(500, gin.H{"error": "failed to read task response"})
		return
	}

	// 200 — sukses (kosong = [])
	c.Data(200, "application/json", taskBody)
}

// Update task
func UpdateTask(c *gin.Context) {
	id := c.Param("id")
	var task models.Task

	// 400 - invalid body
	if err := c.ShouldBindJSON(&task); err != nil {
		log.Println("[UpdateTask] Bind JSON error:", err)
		c.JSON(400, gin.H{"error": "invalid body"})
		return
	}

	resp, err := database.SupabaseRequest(
		"PATCH",
		"/task?id=eq."+id,
		task,
	)
	if err != nil {
		log.Println("[UpdateTask] Supabase PATCH error:", err)
		c.JSON(500, gin.H{"error": err.Error()})
		return
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		log.Println("[UpdateTask] Read body error:", err)
		c.JSON(500, gin.H{"error": "failed to read response"})
		return
	}

	// 404 - task not found
	if string(body) == "[]" {
		log.Println("[UpdateTask] Task not found:", id)
		c.JSON(404, gin.H{
			"error": "task not found",
		})
		return
	}

	c.Data(200, "application/json", body)
}

// Delete task
func DeleteTask(c *gin.Context) {
	id := c.Param("id")

	resp, err := database.SupabaseRequest(
		"DELETE",
		"/task?id=eq."+id,
		nil,
	)
	if err != nil {
		log.Println("[DeleteTask] Supabase DELETE error:", err)
		c.JSON(500, gin.H{"error": "failed to delete task"})
		return
	}
	defer resp.Body.Close()

	body, _ := io.ReadAll(resp.Body)

	// 404 - task not found
	if string(body) == "[]" {
		log.Println("[DeleteTask] Task not found:", id)
		c.JSON(404, gin.H{
			"error": "task not found",
		})
		return
	}

	log.Println("[DeleteTask] Task deleted:", id)
	c.JSON(200, gin.H{
		"message": "task deleted",
	})
}