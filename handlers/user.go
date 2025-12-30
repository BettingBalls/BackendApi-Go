package handlers

import (
	"bytes"
	"io"
	"log"

	"github.com/gin-gonic/gin"

	"task/database"
	"task/models"
)

// Membuat user baru
func CreateUser(c *gin.Context) {
	var user models.User

	// 400 — invalid JSON
	if err := c.ShouldBindJSON(&user); err != nil {
		log.Println("[CreateUser] Bind JSON error:", err)
		c.JSON(400, gin.H{
			"error":   "invalid request body",
			"details": err.Error(),
		})
		return
	}

	// 400 — missing username
	if user.Username == "" {
		c.JSON(400, gin.H{
			"error": "username is required",
		})
		return
	}

	resp, err := database.SupabaseRequest(
		"POST",
		"/user",
		user,
	)
	if err != nil {
		log.Println("[CreateUser] Supabase POST error:", err)
		c.JSON(500, gin.H{"error": "failed to create user"})
		return
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		log.Println("[CreateUser] Read body error:", err)
		c.JSON(500, gin.H{"error": "failed to read response"})
		return
	}

	// cek status code
	if resp.StatusCode >= 400 {

		// 400 — username sudah ada
		if bytes.Contains(body, []byte("23505")) {
			c.JSON(400, gin.H{
				"error": "username already exists",
			})
			return
		}

		log.Println("[CreateUser] Supabase error:", string(body))
		c.JSON(resp.StatusCode, gin.H{
			"error": "failed to create user",
		})
		return
	}

	// 200 — sukses
	c.Data(resp.StatusCode, "application/json", body)
}

// Mengambil daftar user
func GetUsers(c *gin.Context) {
	resp, err := database.SupabaseRequest(
		"GET",
		"/user?select=id,username,role,created_at",
		nil,
	)
	if err != nil {
		log.Println("[GetUsers] Supabase GET error:", err)
		c.JSON(500, gin.H{"error": "failed to fetch users"})
		return
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		log.Println("[GetUsers] Read body error:", err)
		c.JSON(500, gin.H{"error": "failed to read response"})
		return
	}

	c.Data(200, "application/json", body)
}

// Menghapus user
func DeleteUser(c *gin.Context) {
	id := c.Param("id")

	resp, err := database.SupabaseRequest(
		"DELETE",
		"/user?id=eq."+id,
		nil,
	)
	if err != nil {
		log.Println("[DeleteUser] Supabase DELETE error:", err)
		c.JSON(500, gin.H{"error": "failed to delete user"})
		return
	}
	defer resp.Body.Close()

	body, _ := io.ReadAll(resp.Body)

	// cek status code
	if string(body) == "[]" {
		log.Println("[DeleteUser] User not found:", id)
		c.JSON(404, gin.H{
			"error": "user not found",
		})
		return
	}

	log.Println("[DeleteUser] User deleted:", id)
	c.JSON(200, gin.H{
		"message": "user deleted",
	})
}
