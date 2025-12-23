package service

import (
    "encoding/json"
    "fmt"
    "math/rand"
    "net/http"
    "os"
    "strconv"
    "time"
    "bytes"
)

type Service struct {
    djangoURL    string
    authToken    string
    minDelay     int
    maxDelay     int
}

func NewService() *Service {
    return &Service{
        djangoURL: getEnv("DJANGO_URL", "http://localhost:8000"),
        authToken: getEnv("ASYNC_SERVICE_TOKEN", "my-secret-token-12345"),
        minDelay:  getEnvAsInt("DELAY_MIN", 5),
        maxDelay:  getEnvAsInt("DELAY_MAX", 10),
    }
}

func (s *Service) ValidateToken(token string) bool {
    // Простая проверка токена
    return token == s.authToken
}

func (s *Service) CalculatePopulationAsync(applicationID int, territoryArea float64, orders []interface{}) {
    // Имитируем долгий расчет (5-10 секунд)
    delay := s.minDelay + rand.Intn(s.maxDelay-s.minDelay+1)
    fmt.Printf("Starting calculation for application %d, will take %d seconds\n", applicationID, delay)
    
    time.Sleep(time.Duration(delay) * time.Second)
    
    // Выполняем расчет с добавлением случайного фактора
    totalPopulation := s.calculateTotalPopulation(territoryArea, orders)
    
    // Добавляем случайное отклонение (±10%)
    rand.Seed(time.Now().UnixNano())
    deviation := 0.9 + rand.Float64()*0.2 // от 0.9 до 1.1
    finalPopulation := int(float64(totalPopulation) * deviation)
    
    // Отправляем результат обратно в Django
    s.sendResultToDjango(applicationID, finalPopulation)
}

func (s *Service) calculateTotalPopulation(territoryArea float64, orders []interface{}) int {
    total := 0
    
    for _, order := range orders {
        // Преобразуем interface{} в map
        if orderMap, ok := order.(map[string]interface{}); ok {
            buildingDensity, _ := orderMap["building_density"].(float64)
            peoplePerBuilding, _ := orderMap["people_per_building"].(float64)
            
            // Расчет для одной услуги: площадь × плотность × человек
            population := territoryArea * buildingDensity * peoplePerBuilding
            total += int(population)
        }
    }
    
    return total
}

func (s *Service) sendResultToDjango(applicationID int, population int) error {
    // Формируем URL для отправки результата
    url := fmt.Sprintf("%s/api/applications/%d/async_result/", s.djangoURL, applicationID)
    
    // Формируем данные
    data := map[string]interface{}{
        "async_population": population,
        "calculation_status": "completed",
        "calculation_method": "async_refined",
        "auth_token": s.authToken,
    }
    
    jsonData, err := json.Marshal(data)
    if err != nil {
        fmt.Printf("Error marshaling data: %v\n", err)
        return err
    }
    
    // Отправляем PUT запрос
    req, err := http.NewRequest("PUT", url, bytes.NewBuffer(jsonData))
    if err != nil {
        fmt.Printf("Error creating request: %v\n", err)
        return err
    }
    
    req.Header.Set("Content-Type", "application/json")
    
    client := &http.Client{Timeout: 10 * time.Second}
    resp, err := client.Do(req)
    if err != nil {
        fmt.Printf("Error sending request: %v\n", err)
        return err
    }
    defer resp.Body.Close()
    
    if resp.StatusCode != http.StatusOK {
        fmt.Printf("Unexpected status code: %d\n", resp.StatusCode)
        return fmt.Errorf("status code: %d", resp.StatusCode)
    }
    
    fmt.Printf("Successfully sent async result for application %d: %d people\n", applicationID, population)
    return nil
}

// Вспомогательные функции
func getEnv(key, defaultValue string) string {
    if value := os.Getenv(key); value != "" {
        return value
    }
    return defaultValue
}

func getEnvAsInt(key string, defaultValue int) int {
    if value := os.Getenv(key); value != "" {
        if intValue, err := strconv.Atoi(value); err == nil {
            return intValue
        }
    }
    return defaultValue
}