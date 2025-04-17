package model

import (
	"encoding/json"
	"fmt"
	"log/slog"
	"net/http"
)

type ViaCepResponse struct {
	Cep         string `json:"cep"`
	Logradouro  string `json:"logradouro"`
	Complemento string `json:"complemento"`
	Unidade     string `json:"unidade"`
	Bairro      string `json:"bairro"`
	Localidade  string `json:"localidade"`
	Uf          string `json:"uf"`
	Estado      string `json:"estado"`
	Regiao      string `json:"regiao"`
	Ibge        string `json:"ibge"`
	Gia         string `json:"gia"`
	Ddd         string `json:"ddd"`
	Siafi       string `json:"siafi"`
}

func (v *ViaCepResponse) Validate() error {
	if v.Cep == "" {
		return fmt.Errorf("CEP não pode ser vazio")
	}
	if v.Logradouro == "" {
		return fmt.Errorf("Logradouro não pode ser vazio")
	}
	if v.Bairro == "" {
		return fmt.Errorf("Bairro não pode ser vazio")
	}
	if v.Localidade == "" {
		return fmt.Errorf("Localidade não pode ser vazia")
	}
	if v.Uf == "" {
		return fmt.Errorf("UF não pode ser vazia")
	}
	return nil
}

func FetchViaCep(cep string) (*ViaCepResponse, error) {
	url := fmt.Sprintf("https://viacep.com.br/ws/%s/json/", cep)
	slog.Info("Buscando endereço no ViaCEP", "url", url)
	resp, err := http.Get(url)
	if err != nil {
		slog.Error("Erro ao fazer request para ViaCEP", "error", err)
		return nil, fmt.Errorf("erro ao fazer request para ViaCEP: %w", err)
	}
	defer resp.Body.Close()

	slog.Info("Recebendo resposta do ViaCEP", "status", resp.StatusCode)

	if resp.StatusCode != http.StatusOK {
		slog.Error("Erro ao buscar endereço no ViaCEP", "status", resp.StatusCode)
		return nil, fmt.Errorf("ViaCEP retornou status %d", resp.StatusCode)
	}

	var result ViaCepResponse
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		slog.Error("Erro ao decodificar resposta do ViaCEP", "error", err)
		return nil, fmt.Errorf("erro ao decodificar resposta do ViaCEP: %w", err)
	}

	if err := result.Validate(); err != nil {
		slog.Error("Erro ao validar resposta do ViaCEP", "error", err)
		return nil, fmt.Errorf("erro ao validar resposta do ViaCEP: %w", err)
	}

	slog.Info("Endereço encontrado", "cep", result.Cep, "logradouro", result.Logradouro, "bairro", result.Bairro, "localidade", result.Localidade, "uf", result.Uf)

	return &result, nil
}
