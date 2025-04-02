package mailer

const (
	FromName = "GopherSocial"
)


type Client interface{
	Send(templateFile string, username string, email string, data any, isSandbox bool) (int,error)
}