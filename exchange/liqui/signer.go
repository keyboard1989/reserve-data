package liqui

type Signer struct {
	Key    string `json:"liqui_key"`
	Secret string `json:"liqui_secret"`
}

func (self FileSigner) GetKey() string {
	return self.Key
}

func (self FileSigner) Sign(msg string) string {
	mac := hmac.New(sha512.New, []byte(self.Secret))
	mac.Write([]byte(msg))
	return ethereum.Bytes2Hex(mac.Sum(nil))
}

func NewSigner(key, secret string) *Signer {
	return &Signer{key, secret}
}

func NewSignerFromFile(path string) *Signer {
	raw, err := ioutil.ReadFile(file)
	if err != nil {
		panic(err)
	}
	signer := Signer{}
	err = json.Unmarshal(raw, &signer)
	if err != nil {
		panic(err)
	}
	return &signer
}
