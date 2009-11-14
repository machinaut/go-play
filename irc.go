// IRC - Internet Relay Chat package provides a convenient method of writing IRC clients and chat bots (RFC 2812).
package main

import (
    "net";
    "strings";
    "log";
    "time";
)


func main() {
    // Connect to the server
    server := "irc.freenode.net:6667"; // TODO: make this not hardcoded
    // First resolve the name:port to an address
    addr, err := net.ResolveTCPAddr(server);
    if err != nil {
        log.Exitf("Error resolving server '%s': %s", server, err)
    }
    log.Stdoutf("Successfully resolved server '%s'", server);
    // Then Dial the address we just found to connect
    conn, err := net.DialTCP("tcp", nil, addr);
    if err != nil {
        log.Exitf("Error dialing server '%s': %s", server, err)
    }
    log.Stdoutf("Connected to server '%s'", server);

    // Get the server reply
    reply := make([]byte, 100000);
    if n, err := conn.Read(reply); err != nil {
        log.Exitf("Error reading server reply (Read %d): %s", n, err)
    }
    log.Stdoutf("Connection Response %s", reply);

    // Tell it a session password
    password := "foobarbazqux"; // TODO: make a real password
    conn.Write(strings.Bytes("PASS " + password + "\r\n"));

    // Get the server reply
    reply = make([]byte, 100000);
    if n, err := conn.Read(reply); err != nil {
        log.Exitf("Error reading server reply (Read %d): %s", n, err)
    }
    log.Stdoutf("Password Response %s", reply);

    // Tell it a nickname
    nickname := "go_bot1"; // TODO: make this a parameter
    conn.Write(strings.Bytes("NICK " + nickname + "\r\n"));

    // Get the server reply
    reply = make([]byte, 100000);
    if n, err := conn.Read(reply); err != nil {
        log.Exitf("Error reading server reply (read %d): %s", n, err)
    }
    log.Stdoutf("Nickname Response %s", reply);

    // Tell it a username
    username := "turing";        // TODO: make this a parameter
    realname := "Alonzo Church"; // TODO: make this a paramater
    conn.Write(strings.Bytes("USER " + username + " 8 * :" + realname + "\r\n"));

    // Get the server reply (all of it)
    log.Stdout("Username Response");
    for {
        reply = make([]byte, 1000);
        for i := 0; i < 1000; i++ {
            if n, err := conn.Read(reply[i : i+1]); err != nil {
                log.Exitf("Error reading server reply (read %d, %s): %s", n, reply, err)
            }
            switch string(reply[i]) {
            case "\n":
                break
            case "\x04":
                goto chans
            }
        }
        log.Stdoutf("%s", reply);
    }

    // Tell it a channame
chans:
    time.Sleep(int64(1e6));
    channame := "#ncsulug"; // TODO: make this a parameter
    conn.Write(strings.Bytes("JOIN " + channame + "\r\n"));

    // Get the server reply
    reply = make([]byte, 100000);
    if n, err := conn.Read(reply); err != nil {
        log.Exitf("Error reading server reply (read %d): %s", n, err)
    }
    log.Stdoutf("Channame Response %s", reply);

    // We're done with the connection, close it
    time.Sleep(int64(1));
    conn.Close();

}
//// We will use a raw socket to connect to the IRC server.
//use IO::Socket;
//
//// The server to connect to and our details.
//my $server = "irc.freenode.net";
//my $nick = "Homecoming";
//my $login = "christmas";
//
//// The channel which the bot will join.
//my $channel = "//DSotM";
//
//// Connect to the IRC server.
//my $sock = new IO::Socket::INET(PeerAddr => $server,
//                                PeerPort => 6667,
//                                Proto => 'tcp') or
//                                    die "Can't connect\n";
//
//// Log on to the server.
//print $sock "NICK $nick\r\n";
//print $sock "USER $login 8 * :Viva Pink Floyd!\r\n";
//
//// Read lines from the server until it tells us we have connected.
//print "Connecting to irc.freenode.org...";
//while (my $input = <$sock>) {
//    // Check the numerical responses from the server.
//    if ($input =~ /004/) {
//        // We are now logged in.
//	print "Connected.\n";
//        last;
//    }
//    elsif ($input =~ /433/) {
//        die "Nickname is already in use.";
//    }
//}
//
//// Join the channel.
//print $sock "JOIN $channel\r\n";
//
//// Keep reading lines from the server.
//while (my $input = <$sock>) {
//    chop $input;
//    if ($input =~ /^PING(.*)$/i) {
//        // We must respond to PINGs to avoid being disconnected.
//        print $sock "PONG $1\r\n";
//    }
//    else {
//        // Print the raw line received by the bot.
//        print "viva viva viva$input\n";
//    }
//}
//
