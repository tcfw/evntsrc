package evntsrc

//ClientOptions all possible configurable client options
type ClientOptions struct {
	IgnoreSelf bool
	Crypto     EncrypterDecrypter
}

//ClientOption single option caller type
type ClientOption func(*ClientOptions)

//WithCrypto adds client-side stream level encryption
func WithCrypto(crypter EncrypterDecrypter) ClientOption {
	return func(options *ClientOptions) {
		options.Crypto = crypter
	}
}

//WithOwnEvents will receive events which the client sends
func WithOwnEvents() ClientOption {
	return func(options *ClientOptions) {
		options.IgnoreSelf = false
	}
}
