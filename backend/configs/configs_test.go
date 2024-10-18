// configs/configs_test.go
package configs

import (
	"os"
	"testing"
)

func TestLoadConfig(t *testing.T) {
	// Remove any existing .env values to ensure the test starts fresh
	os.Clearenv()

	// Carrega as configurações do arquivo .env
	config := LoadConfig()

	if config == nil {
		t.Fatal("Configurações não devem ser nulas após a inicialização")
	}

	// Verifica se a variável APP_ENV foi carregada
	if config["APP_ENV"] != "development" {
		t.Errorf("Esperado APP_ENV='development', mas obteve %v", config["APP_ENV"])
	}

	// Verifica se a variável TEST_KEY foi carregada
	if config["TEST_KEY"] != "test_value" {
		t.Errorf("Esperado TEST_KEY='test_value', mas obteve %v", config["TEST_KEY"])
	}
}

func TestGet(t *testing.T) {
	os.Setenv("APP_ENV", "testing")
	value := Get("APP_ENV")
	if value != "testing" {
		t.Errorf("Esperado APP_ENV='testing', mas obteve %v", value)
	}
}

func TestGetWithDefault(t *testing.T) {
	// Verifica a recuperação de um valor existente
	os.Setenv("APP_ENV", "production")
	value := GetWithDefault("APP_ENV", "development")
	if value != "production" {
		t.Errorf("Esperado APP_ENV='production', mas obteve %v", value)
	}

	// Verifica a recuperação de um valor padrão quando a chave não existe
	value = GetWithDefault("NON_EXISTENT_KEY", "default_value")
	if value != "default_value" {
		t.Errorf("Esperado valor padrão 'default_value', mas obteve %v", value)
	}
}

func TestAll(t *testing.T) {
	// Define algumas variáveis de ambiente para teste
	os.Setenv("APP_ENV", "staging")
	os.Setenv("TEST_VAR", "test_value")

	config := All()

	if config["APP_ENV"] != "staging" {
		t.Errorf("Esperado APP_ENV='staging', mas obteve %v", config["APP_ENV"])
	}

	if config["TEST_VAR"] != "test_value" {
		t.Errorf("Esperado TEST_VAR='test_value', mas obteve %v", config["TEST_VAR"])
	}
}
