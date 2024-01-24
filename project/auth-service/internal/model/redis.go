package model

type CreateRedisUserSingInDB7Req struct {
	Key      string
	Password bool
	Salt     string
	ID       int64
	IP       string
	DFA      int64
}

type CreateRedisAuthSingInDB8Req struct {
	Key      string
	AuthCode int
	ID       int64
	IP       string
	DFA      int64
}

type CreateRedisAuthSingUpDB7Req struct {
	Key      string
	AuthCode int
	Salt     string
	IP       string
}

type CreateRedisUserSingUpDB8Req struct {
	Key        string
	Name       string
	Surname    string
	Patronymic string
	Email      string
	Username   string
	Birthday   int64
	Phone      int64
}

type CreateRedisConfirmEmailDB7Req struct {
	Key  string
	Salt string
	IP   string
}

type CreateRedisGeneralRecoveryDB7Req struct {
	Key          string
	AuthCode     int
	ActiveStatus bool
	IP           string
	ID           int64
	DFA          int64
	Salt         string
}

type CreateRedisConfirmRecoveryDB8Req struct {
	Key      string
	AuthCode int
	ID       string
}
