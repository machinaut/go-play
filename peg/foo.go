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
    "strings";
    "unicode";
    "utf8";
)


// ----------------------------------------------------------------------------
// Internal representation

type (
    // A Grammar node represents a whole PEG
    // Grammar ←  Spacing Definition* EndOfFile
    Grammar map[string]Expression;

    // A Definition node represents a nonterminal definition
    // Definition ←  Identifier '←' Expression
    Definition struct {
        Identifier string;
        Expr       Expression;
    };

    // An Epression node is a series of choices (in order)
    // Expression ←  Sequence ('/' Sequence)*
    Expression []Sequence;

    // A Sequence node is a sequence of primarys
    // Sequence ←  Primary*
    Sequence []Primary;

    // A Primary node is the basic element of sequences
    // Primary ←  Prefix? Inner Suffix?
    // Inner ←  Identifier !'←' / '(' Expression ')' / Literal / Class / '.'
    Primary struct {
        Prefix string; // Prefix ←  '&' / '!'
        Suffix string; // Suffix ←  '?' / '*' / '+'
    };
    // Primary with Inner ← Identifier !'←'
    PrimaryIdentifier struct {
        *Primary;
        Identifier string;
    };
    // Primary with Inner ← '(' Expression ')'
    PrimaryExpression struct {
        *Primary;
        Expr Expression;
    };
    // Primary with Literal
    PrimaryLiteral struct {
        *Primary;
        Literal string;
    };
    // Primary with Range
    PrimaryRange struct {
        *Primary;
        Range unicode.Range;
    };
)

// ----------------------------------------------------------------------------
// Parsing functions
// By convention:
//  * function names start with a lowercase 'p'
//  * functions evaluate to a single boolean value
//  * parameters/return values are passed by reference so it just evaluates to a boolean
//  * backtracking is the responsibility of the callee
//     (so correctly store the old *i if you modify it, and always return the correct value)

// Parses a PEG into grammar from src returning sucess (true) or failure (false)
// Grammar ←  Spacing Definition* EndOfFile
func pGrammar(gram Grammar, src []byte, i *int) bool {
    d := new(Definition); // To store parse results
    for pDefinition(d, src, i) {
        gram[d.Identifier] = d.Expr
    }
    return true;
}

// Parses a definition in a PEG
// Definition ←  Identifier '←' Expression
func pDefinition(d *Definition, src []byte, i *int) bool {
    old := *i; // Store old *i for backtracking
    ident, expr := new(string), new(Expression);
    if pIdentifier(ident, src, i) && pLiteral("←", src, i) && pExpression(expr, src, i) {
        d.Identifier = *ident;
        d.Expr = *expr;
        return true;
    }
    // No match
    *i = old;
    return false;
}

// Parses a PEG Expression
// Expression  ←  Sequence ('/' Sequence)*
func pExpression(expr *Expression, src []byte, i *int) bool {
    old := *i; // Store old *i for backtracking
    fmt.Println(expr, src, old);
    return true;
}

// Parses an identifier in a PEG
// Identifier  ←  Letter ( Letter / Digit )* White_Space*
func pIdentifier(ident *string, src []byte, i *int) bool {
    rune := new(string); // container for pRange results, holds one rune
    // Letter
    if pRange(unicode.Letter, rune, src, i) {
        *ident += *rune;
        // ( Letter / Digit )*
        for {
            // Letter
            if pRange(unicode.Letter, rune, src, i) {
                *ident += *rune;
                continue;
            }
            // Digit
            if pRange(unicode.Digit, rune, src, i) {
                *ident += *rune;
                continue;
            }
            // No match
            break;
        }
        // White_Space*
        for {
            // White_Space
            if pRange(unicode.White_Space, rune, src, i) {
                continue
            }
            // No match
            break;
        }
        return true;
    }
    // No match
    return false;
}

// Parses the next rune and checks to see if its in a given range
func pRange(ranges []unicode.Range, result *string, src []byte, i *int) bool {
    rune, size := utf8.DecodeRune(src);
    if unicode.Is(ranges, rune) {
        buf := make([]byte, size);
        utf8.EncodeRune(rune, buf);
        *result = string(buf);    // return resulting rune
        *i += size;               // Update index
        src = src[size:len(src)]; // Update slice
        return true;
    }
    // No match
    return false;
}

// Parses the next rune and checks to see if it is a specific literal
func pLiteral(literal string, src []byte, i *int) bool {
    size := len(strings.Bytes(literal));
    if literal == string(src[0:size]) {
        *i += size;
        src = src[size:len(src)];
        return true;
    }
    // No match
    return false;
}

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
