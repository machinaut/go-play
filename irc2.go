package main

import (
    "bufio";
    "log";
    "net";
    "os";
    "regexp";
    "strings";
    "strconv";
    "time";
)

// IRC server connnection, has a simple logger and a TCP connection
type IRCConn struct {
    *net.Conn;                      // TCP Connection to write on
    *log.Logger;                    // Logging enabled
    *bufio.ReadWriter;              // Handy functions to read/write TCP conn
    Send               chan string; // Channel for sending messages to the server
    Rec                chan string; // Channel for recieving messages from the server
}

// Constructor for IRCConn
// TODO: would it be cleaner to just oneline this in the code itself?
func NewIRCConn(conn *net.Conn, logger *log.Logger, rw *bufio.ReadWriter, send, rec chan string) *IRCConn {
    return &IRCConn{conn, logger, rw, send, rec}
}

// Constructor for IRCConn uses the defaults for setup
func DialIRC(laddr, raddr string) (c *IRCConn, err os.Error) {
    logger := log.New(os.Stdout, nil, laddr+";"+raddr+";", log.Lok|log.Ltime|log.Ldate);
    conn, err := net.Dial("tcp", laddr, raddr);
    if err != nil {
        return nil, err
    }
    rw := bufio.NewReadWriter(bufio.NewReader(conn), bufio.NewWriter(conn));
    send := make(chan string); // Unbuffered
    rec := make(chan string);  // Unbuffered
    // Make new IRCConn from logger, tcp conn, and readwriter
    c = NewIRCConn(&conn, logger, rw, send, rec);
    // Start send/recieve handlers
    go c.HandleSend();
    go c.HandleRec();

    return;
}

// Handle sending messages (insert wait to avoid flood)
func (c *IRCConn) HandleSend() {
    for { // Loop forever
        next := <-c.Send; // Get next string to send
        n, err := c.Write(strings.Bytes(next));
        c.Log("SENT:", strings.TrimSpace(next));
        if err != nil {
            c.Logf("SEND ERROR (Wrote %d): %s", n, err)
        }
        time.Sleep(1e8); // Wait 1/10th of a second
    }
}

// Handle recieving messages (pass on anything not explicitly handled)
func (c *IRCConn) HandleRec() {
    for { // Loop forever
        line, err := c.ReadString('\n');
        c.Log("READ:", strings.TrimSpace(line));
        if err != nil {
            c.Log("READ ERROR:", err)
        }
        if line[0] != ':' { // not a private message
            words := strings.Split(line, " ", 2);
            switch words[0] { // Message Type
            case "PING":
                // TODO: find a better way to remove leading character in string
                pong := "PONG " + words[1] + "\r\n";
                c.Send <- pong;
                continue;
            case "NOTICE": // silently eat these
                continue
            }

        }
        c.Rec <- line; // If not handled here, pass it on
    }
}
// Send a message to the server
// Generic enough to be used to send any message
// TODO: remove all of the n, err stuff
func (c *IRCConn) Mesg(mesgType string, params []string) (n int, err os.Error) {
    message := mesgType + " " + strings.Join(params, " ") + "\r\n";
    c.Send <- message;
    return 1, nil;
    //    return c.Write(strings.Bytes(message));
}

// Send NICK message to server, securing a nickname
func (c *IRCConn) Nick(nick string) (n int, err os.Error) {
    params := []string{nick};
    return c.Mesg("NICK", params);
}

// Send USER message to server
// user : username requested
// name : real name
// mode : numeric(int) with 2 bit flags
//    'w'     bit 2 : recieve wallops
//    'i'     bit 3 : invisible
// (if in doubt use 8 as the mode)
func (c *IRCConn) User(user, name string, mode int) (n int, err os.Error) {
    params := []string{user, strconv.Itoa(mode), "*", ":" + name};
    return c.Mesg("USER", params);
}

// Login to an IRC server
// Uses the Nick and User methods to ident with some defaults
func (c *IRCConn) Login(nick string) (n int, err os.Error) {
    i, err := c.Nick(nick);
    if err != nil {
        return i, err
    }
    j, err := c.User(nick, "Go Programming Language", 8);
    return i + j, err;
}

// Join an IRC Chatroom
// Chat names are usually of the form "#name" (e.g. "#fedora")
func (c *IRCConn) Join(name string) (n int, err os.Error) {
    params := []string{name};
    return c.Mesg("JOIN", params);
}

// Say hi
func (c *IRCConn) Hi(name string) (n int, err os.Error) {
    params := []string{name, "Hi!"};
    return c.Mesg("PRIVMSG", params);
}

// Type to parse irc users into
type IRCUser struct {
    Nick, User string;
}

// Regex to use
var userExp = regexp.MustCompile(":([^!]+){!(.+)}?")

// Use regexp's to parse userstrings
func ParseUser(s string) (u *IRCUser) {
    u = new(IRCUser);
    log.Stdout("Got User", s);
    a := userExp.MatchStrings(s);
    log.Stdout("Parsed", len(a));
    if len(a) > 0 { // Got a nick
        u.Nick = a[1];
        u.User = a[2];
    }
    return u;
}

/*
2.3.1 Message format in Augmented BNF

   The protocol messages must be extracted from the contiguous stream of
   octets.  The current solution is to designate two characters, CR and
   LF, as message separators.  Empty messages are silently ignored,
   which permits use of the sequence CR-LF between messages without
   extra problems.

   The extracted message is parsed into the components <prefix>,
   <command> and list of parameters (<params>).

    The Augmented BNF representation for this is:

"
    message    =  [ ":" prefix SPACE ] command [ params ] crlf
    prefix     =  servername / ( nickname [ [ "!" user ] "@" host ] )
    command    =  1*letter / 3digit
    params     =  *14( SPACE middle ) [ SPACE ":" trailing ]
               =/ 14( SPACE middle ) [ SPACE [ ":" ] trailing ]

    nospcrlfcl =  %x01-09 / %x0B-0C / %x0E-1F / %x21-39 / %x3B-FF
                    ; any octet except NUL, CR, LF, " " and ":"
    middle     =  nospcrlfcl *( ":" / nospcrlfcl )
    trailing   =  *( ":" / " " / nospcrlfcl )

    SPACE      =  %x20        ; space character
    crlf       =  %x0D %x0A   ; "carriage return" "linefeed"
"
   NOTES:
      1) After extracting the parameter list, all parameters are equal
         whether matched by <middle> or <trailing>. <trailing> is just a
         syntactic trick to allow SPACE within the parameter.

      2) The NUL (%x00) character is not special in message framing, and
         basically could end up inside a parameter, but it would cause
         extra complexities in normal C string handling. Therefore, NUL
         is not allowed within messages.

   Most protocol messages specify additional semantics and syntax for
   the extracted parameter strings dictated by their position in the
   list.  For example, many server commands will assume that the first
   parameter after the command is the list of targets, which can be
   described with:

  target     =  nickname / server
  msgtarget  =  msgto *( "," msgto )
  msgto      =  channel / ( user [ "%" host ] "@" servername )
  msgto      =/ ( user "%" host ) / targetmask
  msgto      =/ nickname / ( nickname "!" user "@" host )
  channel    =  ( "#" / "+" / ( "!" channelid ) / "&" ) chanstring
                [ ":" chanstring ]
  servername =  hostname
  host       =  hostname / hostaddr
  hostname   =  shortname *( "." shortname )
  shortname  =  ( letter / digit ) *( letter / digit / "-" )
                *( letter / digit )
                  ; as specified in RFC 1123 [HNAME]
  hostaddr   =  ip4addr / ip6addr
  ip4addr    =  1*3digit "." 1*3digit "." 1*3digit "." 1*3digit
  ip6addr    =  1*hexdigit 7( ":" 1*hexdigit )
  ip6addr    =/ "0:0:0:0:0:" ( "0" / "FFFF" ) ":" ip4addr
  nickname   =  ( letter / special ) *8( letter / digit / special / "-" )
  targetmask =  ( "$" / "#" ) mask
                  ; see details on allowed masks in section 3.3.1
  chanstring =  %x01-07 / %x08-09 / %x0B-0C / %x0E-1F / %x21-2B
  chanstring =/ %x2D-39 / %x3B-FF
                  ; any octet except NUL, BELL, CR, LF, " ", "," and ":"
  channelid  = 5( %x41-5A / digit )   ; 5( A-Z / 0-9 )

  Other parameter syntaxes are:

  user       =  1*( %x01-09 / %x0B-0C / %x0E-1F / %x21-3F / %x41-FF )
                  ; any octet except NUL, CR, LF, " " and "@"
  key        =  1*23( %x01-05 / %x07-08 / %x0C / %x0E-1F / %x21-7F )
                  ; any 7-bit US_ASCII character,
                  ; except NUL, CR, LF, FF, h/v TABs, and " "
  letter     =  %x41-5A / %x61-7A       ; A-Z / a-z
  digit      =  %x30-39                 ; 0-9
  hexdigit   =  digit / "A" / "B" / "C" / "D" / "E" / "F"
  special    =  %x5B-60 / %x7B-7D
                   ; "[", "]", "\", "`", "_", "^", "{", "|", "}"

  NOTES:
      1) The <hostaddr> syntax is given here for the sole purpose of
         indicating the format to follow for IP addresses.  This
         reflects the fact that the only available implementations of
         this protocol uses TCP/IP as underlying network protocol but is
         not meant to prevent other protocols to be used.

      2) <hostname> has a maximum length of 63 characters.  This is a
         limitation of the protocol as internet hostnames (in
         particular) can be longer.  Such restriction is necessary
         because IRC messages are limited to 512 characters in length.
         Clients connecting from a host which name is longer than 63
         characters are registered using the host (numeric) address
         instead of the host name.

      3) Some parameters used in the following sections of this
         documents are not defined here as there is nothing specific
         about them besides the name that is used for convenience.
         These parameters follow the general syntax defined for
         <params>.
*/
func main() {
    // start a connection with the server
    irc, err := DialIRC("", "irc.freenode.net:6667");
    if err != nil {
        log.Exit("Error connecting:", err)
    }
    irc.Log("Hello, world!");

    // Login w/ a default name and user
    irc.Login("go_bot");
    // Join a chat. Yay LUG!
    irc.Join("#bottest");

    // Get messages, handle them
    for { // loop forever
        line := <-irc.Rec;
        words := strings.Split(line, " ", 3);
        irc.Log(ParseUser(words[0]));
        switch words[1] { // Message Type
        case "PRIVMSG":
            go irc.Hi(words[2]);
            continue;
        }

    }

    time.Sleep(6e10); // Wait a minute
    irc.Log("Closing");
    irc.Close();
}
