package main

import (
	"net";
	"fmt";
)

//!/usr/bin/perl -w
// irc.pl
// A simple IRC robot.
// Usage: perl irc.pl
/*
use strict;

// We will use a raw socket to connect to the IRC server.
use IO::Socket;

// The server to connect to and our details.
my $server = "irc.freenode.net";
my $nick = "Homecoming";
my $login = "christmas";

// The channel which the bot will join.
my $channel = "//DSotM";

// Connect to the IRC server.
my $sock = new IO::Socket::INET(PeerAddr => $server,
                                PeerPort => 6667,
                                Proto => 'tcp') or
                                    die "Can't connect\n";

// Log on to the server.
print $sock "NICK $nick\r\n";
print $sock "USER $login 8 * :Viva Pink Floyd!\r\n";

// Read lines from the server until it tells us we have connected.
print "Connecting to irc.freenode.org...";
while (my $input = <$sock>) {
    // Check the numerical responses from the server.
    if ($input =~ /004/) {
        // We are now logged in.
	print "Connected.\n";
        last;
    }
    elsif ($input =~ /433/) {
        die "Nickname is already in use.";
    }
}

// Join the channel.
print $sock "JOIN $channel\r\n";

// Keep reading lines from the server.
while (my $input = <$sock>) {
    chop $input;
    if ($input =~ /^PING(.*)$/i) {
        // We must respond to PINGs to avoid being disconnected.
        print $sock "PONG $1\r\n";
    }
    else {
        // Print the raw line received by the bot.
        print "viva viva viva$input\n";
    }
}
*/

func main() {
	fmt.Printf("Hello, world\n");
	a, b, c := net.LookupHost("machinaut.blogdns.net");
	fmt.Println(a, b, c);
}
