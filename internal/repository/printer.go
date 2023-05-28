package repository

type Printer interface {
	Print(code []byte) error
}
