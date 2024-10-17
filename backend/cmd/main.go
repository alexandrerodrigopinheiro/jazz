// main.go
package main

import "jazz/backend/pkg/logger"

func main() {
	// Inicializa o logger
	logger.InitializeLogger()

	// Usa o logger global para logar informações iniciais da aplicação
	logger.Logger.Info("Application has started")

	// Chama outros módulos
	runApplication()
}

// Função exemplo para rodar a aplicação
func runApplication() {
	logger.Logger.Info("Running application logic")
	// ... outras lógicas da aplicação
}
