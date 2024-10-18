package logger

import (
	"context"
	"os"
	"testing"

	"gorm.io/gorm/logger"
)

// Função que simula o comportamento padrão do Fatal durante os testes para evitar o encerramento do processo
var FatalFn = func(msg string) {
	Logger.sugarLogger.Fatal(msg)
}

func TestInitializeLogger(t *testing.T) {
	// Testando a inicialização do logger em modo de desenvolvimento
	os.Setenv("APP_ENV", "development")
	InitializeLogger()
	if Logger == nil {
		t.Fatal("Logger não deve ser nulo após a inicialização")
	}

	// Testando a inicialização do logger em modo de produção
	os.Setenv("APP_ENV", "production")
	InitializeLogger()
	if Logger == nil {
		t.Fatal("Logger não deve ser nulo após a inicialização")
	}
}

func TestInfo(t *testing.T) {
	InitializeLogger()
	Logger.Info("Esta é uma mensagem de informação para fins de teste")
}

func TestWarn(t *testing.T) {
	InitializeLogger()
	Logger.Warn("Esta é uma mensagem de aviso para fins de teste")
}

func TestError(t *testing.T) {
	InitializeLogger()
	Logger.Error("Esta é uma mensagem de erro para fins de teste")
}

func TestLogMode(t *testing.T) {
	InitializeLogger()
	gormLogger := NewGormLogger(logger.Info)
	gormLogger.Info(context.Background(), "Testando nível de informação do Gorm Logger")
	gormLogger.Warn(context.Background(), "Testando nível de aviso do Gorm Logger")
	gormLogger.Error(context.Background(), "Testando nível de erro do Gorm Logger")
}

func TestFatal(t *testing.T) {
	defer func() {
		if r := recover(); r != nil {
			t.Log("Recuperado de erro fatal")
		}
	}()

	// Substituir o comportamento padrão do Fatal para evitar o encerramento do processo
	originalFatalFn := FatalFn
	FatalFn = func(msg string) {
		t.Log("Simulação de Fatal: ", msg)
		panic("simulando os.Exit")
	}
	defer func() { FatalFn = originalFatalFn }()

	// Logger.Fatal("Esta é uma mensagem fatal para fins de teste")
}
