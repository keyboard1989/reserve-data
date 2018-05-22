package common

type Token struct {
	ID                      string
	Name                    string
	Address                 string
	Decimal                 int64
	Active                  bool
	Internal                bool
	MinimalRecordResolution string
	MaxTotalImbalance       string
	MaxPerBlockImbalance    string
}

// NewToken creates a new Token.
func NewToken(id, name, address string, decimal int64, active bool, internal bool, miminalrr, maxti, maxpbi string) Token {
	return Token{
		ID:                      id,
		Address:                 address,
		Decimal:                 decimal,
		Active:                  active,
		Internal:                internal,
		MinimalRecordResolution: miminalrr,
		MaxTotalImbalance:       maxti,
		MaxPerBlockImbalance:    maxpbi,
	}
}

func (self Token) MarshalText() (text []byte, err error) {
	// return []byte(fmt.Sprintf(
	// 	"%s-%s", self.ID, self.Address,
	// )), nil
	return []byte(self.ID), nil
}

func (self Token) IsETH() bool {
	return self.ID == "ETH"
}

type TokenPair struct {
	Base  Token
	Quote Token
}

func (self *TokenPair) PairID() TokenPairID {
	return NewTokenPairID(self.Base.ID, self.Quote.ID)
}
