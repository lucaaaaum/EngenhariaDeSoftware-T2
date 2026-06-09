package auth

// LoginCommand são os dados de login recebidos pelo cliente.
// O cliente envia email e senha em JSON: {"email": "...", "password": "..."}
type LoginCommand struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

// LoginResponse é o que retornamos após login bem-sucedido.
// O token JWT deve ser enviado no header Authorization de todas as requisições protegidas.
// Exemplo: Authorization: Bearer eyJhbGc...
type LoginResponse struct {
	Token string `json:"token"`
}
