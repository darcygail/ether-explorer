package schema

type Account struct {
	Address   string  `bson:"address"`
	Balance   string  `bson:"balance"`
	Assets    []Asset `bson:"assets"`
	UpdatedAt int64   `bson:"updated_at"`
}

type Asset struct {
	ContractAddress string          `bson:"contract_address"`
	Type            string          `bson:"type"` 
	TokenIDs        []TokenIDAmount `bson:"token_ids"`
}

type TokenIDAmount struct {
	ID string `bson:"id"`
}
