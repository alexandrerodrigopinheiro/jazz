// main.go
package main

import (
	"fmt"
	"jazz/backend/configs"
	"jazz/backend/pkg/cache"
	"jazz/backend/pkg/database"
	"jazz/backend/pkg/logger"

	"gorm.io/gorm"
)

func main() {
	// Inicializa o logger antes de qualquer coisa
	logger.InitializeLogger()

	// Usa o logger global para logar informações iniciais da aplicação
	logger.Logger.Info("Application has started")

	// Carrega as configurações
	if err := configs.LoadConfig(); err != nil {
		panic(fmt.Sprintf("Error loading config: %v", err))
	}

	// Inicializa o banco de dados (Singleton)
	db := database.InitializeDatabase()

	// Inicializa o cache (Singleton)
	cacheManager := cache.NewCacheManager()

	// Testando os serviços
	runApplication(db, cacheManager)
}

func runApplication(_ *gorm.DB, cache cache.Cache) {
	// Executa a lógica principal da aplicação
	logger.Logger.Info("Running application logic")

	// Exemplo de uso do cache e do banco de dados
	key := "example_key"
	value := "example_value"

	// Salvando no cache
	cache.Set(key, value, 60)
	logger.Logger.Infof("Value '%s' set in cache with key '%s'", value, key)

	// Consultando o banco de dados (exemplo fictício)
	// logger.Logger.Infof("Database connection details: %v", db)
	// Outras lógicas da aplicação...
}
