package jwt

import (
	"crypto/rsa"
	"fmt"
	"github.com/golang-jwt/jwt"
	"github.com/google/uuid"
	"github.com/todo-enjoers/backend_v1/internal/config"
	"github.com/todo-enjoers/backend_v1/internal/model"
	"github.com/todo-enjoers/backend_v1/internal/pkg/token"
	"go.uber.org/zap"
	"os"
	"time"
)

// Checking whether the interface "ProviderI" implements the structure "Provider"
var _ token.ProviderI = (*Provider)(nil)

type Provider struct {
	publicKey       *rsa.PublicKey
	privateKey      *rsa.PrivateKey
	accessLifetime  int
	refreshLifetime int
}

type CustomClaims struct {
	jwt.StandardClaims
	IsAccess bool `json:"access"`
}

func NewProvider(cfg *config.Config, log *zap.Logger) (*Provider, error) {
	//Read and Parsing Private Key
	privateKeyRaw, err := os.ReadFile(cfg.JWT.PrivateKeyPath)
	if err != nil {
		return nil, fmt.Errorf("failed to read private key file")
	}
	privateKey, err := jwt.ParseRSAPrivateKeyFromPEM(privateKeyRaw)
	if err != nil {
		return nil, fmt.Errorf("failed to parse private key file")
	}

	//Read and Parsing Public Key
	publicKeyRaw, err := os.ReadFile(cfg.JWT.PublicKeyPath)
	if err != nil {
		return nil, fmt.Errorf("failed to read public key file")
	}
	publicKey, err := jwt.ParseRSAPublicKeyFromPEM(publicKeyRaw)
	if err != nil {
		return nil, fmt.Errorf("failed to parse public key file")
	}

	log.Info("public key", zap.ByteString("publicKey", publicKeyRaw))

	provider := &Provider{
		publicKey:       publicKey,
		privateKey:      privateKey,
		accessLifetime:  cfg.JWT.AccessTokenLifeTime,
		refreshLifetime: cfg.JWT.RefreshTokenLifeTime,
	}

	return provider, nil
}

func (provider *Provider) readKeyFunc(token *jwt.Token) (interface{}, error) {
	// readKeyFunc is a reader of public key.
	if _, ok := token.Method.(*jwt.SigningMethodRSA); !ok {
		return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
	}

	return provider.publicKey, nil
}
func (provider *Provider) GetDataFromToken(token string) (*model.UserDataInToken, error) {
	parsed, err := jwt.ParseWithClaims(token, &CustomClaims{}, provider.readKeyFunc)
	if err != nil {
		return nil, fmt.Errorf("failed to parse token")
	}

	if !parsed.Valid {
		return nil, fmt.Errorf("invalid token: not valid")
	}

	claims, ok := parsed.Claims.(*CustomClaims)
	if !ok {
		return nil, fmt.Errorf("invalid token: can't parse claims")
	}

	var ParsedID uuid.UUID

	ParsedID, err = uuid.Parse(claims.Issuer)
	if err != nil {
		return nil, fmt.Errorf("invalid token: issuer is not UUID")
	}

	return &model.UserDataInToken{
		ID:       ParsedID,
		IsAccess: claims.IsAccess,
	}, nil
}

// CreateTokenForUser : Create a JWT Token for user
func (provider *Provider) CreateTokenForUser(userID uuid.UUID, isAccess bool) (string, error) {
	now := time.Now()

	var add time.Duration

	// checking token accessible
	if isAccess {
		add = time.Duration(provider.accessLifetime) * time.Minute
	} else {
		add = time.Duration(provider.refreshLifetime) * time.Minute
	}

	// creating payload part of JWT Token (header is self-creating)
	claims := &CustomClaims{
		StandardClaims: jwt.StandardClaims{
			Issuer:    userID.String(),
			IssuedAt:  now.Unix(),
			NotBefore: now.Unix(),
			ExpiresAt: now.Add(add).Unix(),
		},
		IsAccess: isAccess,
	}

	// JWT token is signed with claims
	token := jwt.NewWithClaims(jwt.SigningMethodRS256, claims)
	return token.SignedString(provider.privateKey)
}
