package seeders

import (
	"log"
	"time"

	"jobtrack-backend/config"
	"jobtrack-backend/models"

	"golang.org/x/crypto/bcrypt"
)

// SeedAll menjalankan seluruh proses seeding data awal
func SeedAll() {
	log.Println("========================================")
	log.Println("Memulai proses seeding database...")
	log.Println("========================================")

	// 1. Seed User Demo
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte("password123"), bcrypt.DefaultCost)
	if err != nil {
		log.Printf("Gagal hash password: %v\n", err)
		return
	}

	demoUser := models.User{
		Name:     "Demo User",
		Email:    "demo@example.com",
		Password: string(hashedPassword),
	}

	// Cek jika user demo sudah ada
	var user models.User
	if err := config.DB.Where("email = ?", demoUser.Email).First(&user).Error; err != nil {
		if err := config.DB.Create(&demoUser).Error; err != nil {
			log.Printf("Gagal seeding user demo: %v\n", err)
			return
		}
		user = demoUser
		log.Printf("[+] User dibuat: %s (Password: password123)\n", user.Email)
	} else {
		log.Printf("[*] User %s sudah ada, menggunakan ID: %d\n", user.Email, user.ID)
	}

	// 2. Seed Data Lowongan Pekerjaan (Jobs)
	jobs := []models.Job{
		{
			Title:       "Frontend Developer",
			Company:     "Google Indonesia",
			Link:        "https://careers.google.com",
			Type:        "Full-time",
			Description: "Jadwal interview teknis tahap 2 via Google Meet",
			Status:      "Interview",
			UserID:      user.ID,
			CreatedAt:   time.Now().AddDate(0, 0, -10),
		},
		{
			Title:       "Backend Golang Engineer",
			Company:     "Tokopedia",
			Link:        "https://tokopedia.com/careers",
			Type:        "Full-time",
			Description: "Lamaran dikirim via LinkedIn, menunggu respon recruiter",
			Status:      "Terkirim",
			UserID:      user.ID,
			CreatedAt:   time.Now().AddDate(0, 0, -7),
		},
		{
			Title:       "Software Engineer Intern",
			Company:     "Traveloka",
			Link:        "https://traveloka.com/careers",
			Type:        "Internship",
			Description: "Offering letter sudah diterima! Jadwal onboarding awal bulan",
			Status:      "Diterima",
			UserID:      user.ID,
			CreatedAt:   time.Now().AddDate(0, 0, -14),
		},
		{
			Title:       "UI/UX Designer",
			Company:     "Shopee",
			Link:        "https://careers.shopee.co.id",
			Type:        "Full-time",
			Description: "Belum lolos screening portfolio",
			Status:      "Ditolak",
			UserID:      user.ID,
			CreatedAt:   time.Now().AddDate(0, 0, -20),
		},
	}

	for _, j := range jobs {
		var existingJob models.Job
		err := config.DB.Where("title = ? AND company = ? AND user_id = ?", j.Title, j.Company, user.ID).First(&existingJob).Error
		if err != nil {
			if err := config.DB.Create(&j).Error; err != nil {
				log.Printf("[-] Gagal membuat data job %s: %v\n", j.Title, err)
			} else {
				log.Printf("[+] Job dibuat: %s di %s (Status: %s)\n", j.Title, j.Company, j.Status)
			}
		} else {
			log.Printf("[*] Job %s di %s sudah ada, melewati.\n", j.Title, j.Company)
		}
	}

	log.Println("========================================")
	log.Println("Seeding database selesai!")
	log.Println("========================================")
}
