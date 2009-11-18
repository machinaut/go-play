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
