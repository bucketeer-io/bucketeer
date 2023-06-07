
È
proto/auth/token.protobucketeer.auth"¡
Token!
access_token (	RaccessToken

token_type (	R	tokenType#
refresh_token (	RrefreshToken
expiry (Rexpiry
id_token (	RidToken"B
IDTokenSubject
user_id (	RuserId
conn_id (	RconnIdB.Z,github.com/bucketeer-io/bucketeer/proto/authbproto3
Û
proto/auth/service.protobucketeer.authproto/auth/token.proto"P
GetAuthCodeURLRequest
state (	Rstate!
redirect_url (	RredirectUrl"*
GetAuthCodeURLResponse
url (	Rurl"M
ExchangeTokenRequest
code (	Rcode!
redirect_url (	RredirectUrl"D
ExchangeTokenResponse+
token (2.bucketeer.auth.TokenRtoken"]
RefreshTokenRequest#
refresh_token (	RrefreshToken!
redirect_url (	RredirectUrl"C
RefreshTokenResponse+
token (2.bucketeer.auth.TokenRtoken2§
AuthService_
GetAuthCodeURL%.bucketeer.auth.GetAuthCodeURLRequest&.bucketeer.auth.GetAuthCodeURLResponse\
ExchangeToken$.bucketeer.auth.ExchangeTokenRequest%.bucketeer.auth.ExchangeTokenResponseY
RefreshToken#.bucketeer.auth.RefreshTokenRequest$.bucketeer.auth.RefreshTokenResponseB.Z,github.com/bucketeer-io/bucketeer/proto/authbproto3