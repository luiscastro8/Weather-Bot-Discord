package myerrors

type AddressNotFoundError struct {
	error            string
	UnmatchedAddress string
}

func (e AddressNotFoundError) Error() string {
	return e.error
}

func NewAddressNotFoundError(message, address string) AddressNotFoundError {
	return AddressNotFoundError{error: message, UnmatchedAddress: address}
}
