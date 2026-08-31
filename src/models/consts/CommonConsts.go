package consts

type ContextKey string

const (
	OwnerType            ContextKey = "owner-type"
	OwnerId              ContextKey = "owner-id"
	NationalIdentityCode ContextKey = "national-identity-code"
	Username             ContextKey = "username"
	RequestUuid          ContextKey = "request-uuid"
)
