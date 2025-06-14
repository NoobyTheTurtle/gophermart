package password

import "golang.org/x/crypto/bcrypt"

const DefaultCost = bcrypt.DefaultCost

type PasswordService struct {
	bcryptCost int
}

func New(bcryptCost int) *PasswordService {
	return &PasswordService{
		bcryptCost: bcryptCost,
	}
}

func (p *PasswordService) Hash(password string) (string, error) {
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), p.bcryptCost)
	if err != nil {
		return "", err
	}
	return string(hashedPassword), nil
}

func (*PasswordService) Compare(hashedPassword, password string) error {
	return bcrypt.CompareHashAndPassword([]byte(hashedPassword), []byte(password))
}
