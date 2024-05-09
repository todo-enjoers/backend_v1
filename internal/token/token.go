package token

type DataInToken struct {
	ID int64
}

func CreateTokenWithData(token DataInToken) string {
	panic("not implemented")
}

func GetDataFromToken(token string) DataInToken {
	panic("not implemented")
}
