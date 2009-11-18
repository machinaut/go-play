// IRC Messages in EBNF
// Translated to Go's EBNF from RFC 2812 by Alex Ray <ajray@ncsu.edu>
//    Note: this is less strict than the format required by RFC 2812
//    So don't expect this thing to check for _all_ the proper formats

message    =  [ ":" prefix " " ] command [ params ] "\r\n" .
prefix     =  servername | ( nickname [ [ "!" user ] "@" host ] ) .
command    =  letter {letter} | digit digit digit .
params     =  { " " middle } [ " " [ ":" ] trailing ] .

nospcrlfcl =  "\x01"..."\x09" | "\x0B"..."\x0C" | "\x0E"..."\x1F"
           |  "\x21"..."\x39" | "\x3B"..."\xFF" .
              // any octet except NUL, CR, LF, " " and ":" 
middle     =  nospcrlfcl { ":" | nospcrlfcl } .
trailing   =  { ":" | " " | nospcrlfcl } .

servername =  hostname .
host       =  hostname | hostaddr .
hostname   =  shortname { "." shortname } .
shortname  =  ( letter | digit ) { letter | digit | "-" }
              { letter | digit } .
              // as specified in RFC 1123 [HNAME]
hostaddr   =  ip4addr | ip6addr .
ipdigit    =  digit [digit [digit] ] .
ip4addr    =  ipdigit "." ipdigit "." ipdigit "." ipdigit .
hexnum     =  hexdigit {hexdigit} .
              // TODO: This is a hell of a kludge for IPv6. Make it right.
ip6addr    =  hexnum { ":" [hexnum]  } 
           |  "0:0:0:0:0:" ( "0" | "FFFF" ) ":" ip4addr .
nickname   =  ( letter | special ) { letter | digit | special | "-" } .
letter     =  "A"..."Z" | "a"..."z" .
digit      =  "0"..."9" .
hexdigit   =  digit | "A"..."F" | "a"..."f" .
userbit    =  "\x01"..."\x09" | "\x0B"..."\x0C" | "\x0E"..."\x1F"
           |  "\x21"..."\x3F" | "\x41"..."\xFF" .
              // any octet except NUL, CR, LF, " " and "@"
user       =  userbit {userbit} .
special    =  "\x5B"..."\x60" | "\x7B"..."\x7D" .
              // "[", "]", "\", "`", "_", "^", "{", "|", "}"
