// This package has the values used by postgres
// Based off of the PostgreSQL 8.4.1 Documentation
// Chapter 45. Frontend/Backend Protocol
// http://www.postgresql.org/docs/8.4/interactive/protocol-message-formats.html
package psql_constants

// MessageType specifies the type of message sent/recieved
// (B) indicates this is sent by the Postgres backend
// (F) indicated this is sent by the Postgres frontend
// (F & B) indicates it can be sent by both
type MessageType string

const (
    Authentication       = "R"; //(B) Specifies the type of authentication to use.
    BackendKeyData       = "K"; //(B) Identifies the message as cancellation key data. Frontend must save these values in order to issue CancelRequest messages later.
    Bind                 = "B"; //(F) Identifies the message as a Bind command.
    BindComplete         = "2"; //(B) Identifies the message as a Bind-complete indicator.
    Close                = "C"; //(F) Identifies the message as a Close command.
    CloseComplete        = "3"; //(B) Identifies the message as a Close-complete indicator.
    CommandComplete      = "C"; //(B) Identifies the message as a command-completed response.
    CopyData             = "d"; //(F & B) Identifies the message as COPY data.
    CopyDone             = "c"; //(F & B) Identifies the message as a COPY-complete indicator.
    CopyFail             = "f"; //(F) Identifies the message as a COPY-failure indicator.
    CopyInResponse       = "G"; //(B) Identifies the message as a Start Copy In response. Frontend must now send copy-in data (if not prepared to do so, send a CopyFail message).
    CopyOutResponse      = "H"; //(B) Identifies the message as a Start Copy Out response. This message will be followed by copy-out data.
    DataRow              = "D"; //(B) Identifies the message as a data row.
    Describe             = "D"; //(F) Identifies the message as a Describe command.
    EmptyQueryResponse   = "I"; //(B) Identifies the message as a response to an empty query string. (This substitutes for CommandComplete.)
    ErrorResponse        = "E"; //(B) Identifies the message as an error.
    Execute              = "E"; //(F) Identifies the message as an Execute command.
    Flush                = "H"; //(F) Identifies the message as a Flush command.
    FunctionCall         = "F"; //(F) Identifies the message as a function call.
    FunctionCallResponse = "V"; //(B) Identifies the message as a function call result.
    NoData               = "n"; //(B) Identifies the message as a no-data indicator.
    NoticeResponse       = "N"; //(B) Identifies the message as a notice.
    NotificationResponse = "A"; //(B) Identifies the message as a notification response.
    ParameterDescription = "t"; //(B) Identifies the message as a parameter description.
    ParameterStatus      = "S"; //(B) Identifies the message as a run-time parameter status report.
    Parse                = "P"; //(F) Identifies the message as a Parse command.
    ParseComplete        = "1"; //(B) Identifies the message as a Parse-complete indicator.
    PasswordMessage      = "p"; //(F) Identifies the message as a password response. Note that this is also used for GSSAPI and SSPI response messages (which is really a design error, since the contained data is not a null-terminated string in that case, but can be arbitrary binary data).
    PortalSuspended      = "s"; //(B) Identifies the message as a portal-suspended indicator. Note this only appears if an Execute message's row-count limit was reached.
    Query                = "Q"; //(F) Identifies the message as a simple query.
    ReadyForQuery        = "Z"; //(B) Identifies the message type. ReadyForQuery is sent whenever the backend is ready for a new query cycle.
    RowDescription       = "T"; //(B) Identifies the message as a row description.
    Sync                 = "S"; //(F) Identifies the message as a Sync command.
    Terminate            = "X"; //(F) Identifies the message as a termination.
)

// AuthenticationType specifies the type of authentication to use
type AuthenticationType int

const (
    AuthenticationOK                = 0; //Specifies that the authentication was successful.
    AuthenticationKerberosV5        = 2; //Specifies that Kerberos V5 authentication is required.
    AuthenticationCleartextPassword = 3; //Specifies that a clear-text password is required.
    AuthenticationMD5Password       = 5; //Specifies that an MD5-encrypted password is required.
    AuthenticationSCMCredential     = 6; //Specifies that an SCM credentials message is required.
    AuthenticationGSS               = 7; //Specifies that GSSAPI authentication is required.
    AuthenticationGSSContinue       = 8; //Specifies that this message contains GSSAPI or SSPI data.
    AuthenticationSSPI              = 9; //Specifies that SSPI authentication is required.
)
