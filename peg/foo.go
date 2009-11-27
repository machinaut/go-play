// Copyright 2009 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// A library for PEG grammars. The input is text ([]byte) satisfying
// the following grammar (represented itself in PEG):
// (Based off Bryan Ford's PEG paper)
//
// # Hierarchical syntax
// Grammar     ←  Spacing Definition* EndOfFile # The whole grammar file
// Definition  ←  Identifier '←' Expression     # The Definition Definition
//
// Expression  ←  Sequence ('/' Sequence)*      # Choices
// Sequence    ←  Primary*                      # Series of primaries
// Primary     ←  Prefix? Inner Suffix?         # Optional prefix/suffix
// Prefix      ←  '&' / '!'                     # Stuff before a primary
// Suffix      ←  '?' / '*' / '+'               # Stuff after a primary
// Inner       ←  Identifier !'←'               # Nonterminal
//             /  '(' Expression ')'            # Subexpression
//             /  Literal / Class / '.'         # Terminal
//
// # Lexical syntax
// # Identifiers must not start with a numeric character, but may contain them
// Identifier  ←  [a-zA-Z_] [a-zA-Z_0-9]* Spacing # Or unicode stuffs
//
// Literal     ←  "'" (!"'" Char)*  "'" Spacing # 'literal' ''
//             /  '"' (!'"' Char)*  '"' Spacing # "literal" ""
// Class       ←  '[' (!']' Range)* ']' Spacing # [class] []
// Range       ←  Char '-' Char / Char          # a-z L
// Char        ←  '\\' [nrt'"\[\]\\]            # \n \r \t \' \" \] \] \\
//             /  '\\' ([0-2]?[0-7])?[0-7]      # \123 \50 \2 (octets)
//             /  '\\x'[0-9a-fA-F][0-9a-fA-F]   # \x89 \x00 (hex escapes)
//             /  !'\\' .                       # Anything else
//
// Spacing     ←  (WhiteSpace / Comment)*
// Comment     ←  '//' (!'\n' .)* '\n'          # C-style comments
//             /  '/*' (Comment / !'*/' .)* ’*/’
// WhiteSpace  ←  ' ' / '\t' / '\n' / '\r'      # Or unicode whitespace
// EndOfFile   ←  !.
//
// A name is a Go identifier, a token is a Go string, and comments
// and white space follow the same rules as for the Go language.
// Production names starting with an uppercase Unicode letter denote
// non-terminal productions (i.e., productions which allow white-space
// and comments between tokens); all other production names denote
// lexical productions.
//
package main

import (
    "fmt";
    "io";
    "log";
    //	"unicode";
    "utf8";
)


// ----------------------------------------------------------------------------
// Internal representation

type (
    // A Grammar node represents a whole PEG
    //
    // Grammar ←  Spacing Definition* EndOfFile
    Grammar []Definition;

    // A Definition node represents a nonterminal definition
    //
    // Definition ←  Identifier '←' Expression
    Definition struct {
        Ident string;
        Expr  Expression;
    };

    // An Epression node is a series of choices (in order)
    //
    // Expression ←  Sequence ('/' Sequence)*
    Expression []Sequence;

    // A Sequence node is a sequence of primarys
    //
    // Sequence ←  Primary*
    Sequence []Primary;

    // A Primary node is the basic element of sequences
    //
    // Primary ←  Prefix? Inner Suffix?
    Primary struct {
        Prefix string; // Prefix ←  '&' / '!'
        Suffix string; // Suffix ←  '?' / '*' / '+'
        Inn    Inner;
    };

    // An Inner node is the inside of a primary statement
    //
    // Inner ←  Identifier !'←' / '(' Expression ')' / Literal / Class / '.'
    Inner interface {
        // Nothing in here yet
    };
)

//func Gram() {
//    for {
//        Defn()
//    }
//    return true;
//}
//
//func Defn() {
//    return 1 and 2 and 3

func main() {
    maths, err := io.ReadFile("math.peg");
    if err != nil {
        log.Exit("Error: ", err)
    }

    fmt.Println(maths);
    rune, size := utf8.DecodeRuneInString(string(maths));
    fmt.Println(rune, size);
    return;
}
