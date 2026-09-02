package enum

type SecretType string

func (st SecretType) String() string {
	return string(st)
}

const (
	LoginPass SecretType = "login_pass"
	Binary    SecretType = "binary"
	Card      SecretType = "card"
)
