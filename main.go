package main

import (
    "log"
    "os"
    "async-service/internal/app/handler"
    "async-service/internal/app/service"
    "github.com/gin-gonic/gin"
    "github.com/joho/godotenv"
)

func main() {
    // Загружаем переменные окружения
    if err := godotenv.Load(); err != nil {
        log.Printf("Warning: .env file not found")
    }
    
    // Создаем сервис
    svc := service.NewService()
    
    // Создаем обработчик
    h := handler.NewHandler(svc)
    
    // Настраиваем Gin
    r := gin.Default()
    
    // Добавляем обработчики
    r.POST("/calculate_population", h.CalculatePopulation)
    r.GET("/health", h.HealthCheck)
    
    // Запускаем сервер
    port := os.Getenv("GO_PORT")
    if port == "" {
        port = "8081"
    }
    
    log.Printf("Starting async service on port %s", port)
    if err := r.Run(":" + port); err != nil {
        log.Fatal(err)
    }
}