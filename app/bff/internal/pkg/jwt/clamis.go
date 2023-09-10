package jwt

import (
	"context"
)

func GetClaims(ctx context.Context) (*CustomClaims, error) {
	if claims, exists := ctx.Value("claims").(*CustomClaims); !exists {
		return &CustomClaims{}, nil
	} else {
		return claims, nil
	}
}

func GetUserID(ctx context.Context) uint64 {
	if claims, exists := ctx.Value("claims").(*CustomClaims); !exists {
		return 0
	} else {
		return claims.BaseClaims.Id
	}
}

func GetUserName(ctx context.Context) string {
	if claims, exists := ctx.Value("claims").(*CustomClaims); !exists {
		return ""
	} else {
		return claims.BaseClaims.Name
	}
}
