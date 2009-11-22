// Copyright 2009 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// A library for PEG grammars. The input is text ([]byte) satisfying
// the following grammar (represented itself in PEG):
// (Based off Bryan Ford's PEG paper)

// # Hierarchical syntax
// Grammar     <- Spacing Definition+ EndOfFile    # The whole grammar file
// Definition  <- Identifier LEFTARROW Expression  # The Definition Definition
// 
// Expression  <- Sequence (SLASH Sequence)*       # Choices
// Sequence    <- Prefix*                          # Series of primaries w/ stuff
// Prefix      <- (AND / NOT)? Suffix               # Stuff before a primary
// Suffix      <- Primary (QUESTION / STAR / PLUS)? # Stuff after a primary
// Primary     <- Identifier !LEFTARROW            # Nonterminal
//              / OPEN Expression CLOSE            # Subexpression
//              / Literal / Class / DOT            # Terminal
// 
// # Lexical syntax
// Identifier  <- IdentStart IdentCont* Spacing # Name has to start with a
// IdentStart  <- [a-zA-Z_]                     # letter or an underscore
// IdentCont   <- IdentStart / [0-9]            # and can have numbers in it
// 
// Literal     <- ['] (!['] Char)* ['] Spacing  # 'lawl' 'nar' ''
//              / ["] (!["] Char)* ["] Spacing  # "lawl" "nar" ""
// Class       <- '[' (!']' Range)* ']' Spacing # [bar] []
// Range       <- Char '-' Char / Char          # a-z L
// Char        <- '\\' [nrt'"\[\]\\]            # \n \r \t \' \" \] \] \\
//              / '\\' ([0-2]?[0-7])?[0-7]      # \123 \50 \2
//              / !'\\' .                       # Anything else
// 
// LEFTARROW   <- '<-' Spacing
// SLASH       <- '/' Spacing
// AND         <- '&' Spacing
// NOT         <- '!' Spacing
// QUESTION    <- '?' Spacing
// STAR        <- '*' Spacing
// PLUS        <- '+' Spacing
// OPEN        <- '(' Spacing
// CLOSE       <- ')' Spacing
// DOT         <- '.' Spacing
// 
// Spacing     <- (Space / Comment)*
// Comment     <- '#' (!EndOfLine .)* EndOfLine
// Space       <- ' ' / '\t' / EndOfLine
// EndOfLine   <- '\r\n' / '\n' / '\r'
// EndOfFile   <- !.

// A name is a Go identifier, a token is a Go string, and comments
// and white space follow the same rules as for the Go language.
// Production names starting with an uppercase Unicode letter denote
// non-terminal productions (i.e., productions which allow white-space
// and comments between tokens); all other production names denote
// lexical productions.
//
package peg

import (
	"container/vector";
	"go/scanner";
	"go/token";
	"os";
	"unicode";
	"utf8";
)


// ----------------------------------------------------------------------------
// Internal representation

type (
	// An Expression node represents a production expression.
	Expression	interface {
		// Pos is the position of the first character of the syntactic construct
		Pos() token.Position;
	};

	// An Alternative node represents a non-empty list of alternative expressions.
	Alternative	[]Expression;	// x | y | z

	// A Sequence node represents a non-empty list of sequential expressions.
	Sequence	[]Expression;	// x y z

	// A Name node represents a production name.
	Name	struct {
		token.Position;
		String	string;
	};

	// A Token node represents a literal.
	Token	struct {
		token.Position;
		String	string;
	};

	// A List node represents a range of characters.
	Range	struct {
		Begin, End *Token;	// begin ... end
	};

	// A Group node represents a grouped expression.
	Group	struct {
		token.Position;
		Body	Expression;	// (body)
	};

	// An Option node represents an optional expression.
	Option	struct {
		token.Position;
		Body	Expression;	// [body]
	};

	// A Repetition node represents a repeated expression.
	Repetition	struct {
		token.Position;
		Body	Expression;	// {body}
	};

	// A Production node represents an EBNF production.
	Production	struct {
		Name	*Name;
		Expr	Expression;
	};

	// A Grammar is a set of EBNF productions. The map
	// is indexed by production name.
	//
	Grammar	map[string]*Production;
)


func (x Alternative) Pos() token.Position {
	return x[0].Pos()	// the parser always generates non-empty Alternative
}


func (x Sequence) Pos() token.Position {
	return x[0].Pos()	// the parser always generates non-empty Sequences
}


func (x Range) Pos() token.Position	{ return x.Begin.Pos() }


func (p *Production) Pos() token.Position	{ return p.Name.Pos() }


// ----------------------------------------------------------------------------
// Grammar verification

func isLexical(name string) bool {
	ch, _ := utf8.DecodeRuneInString(name);
	return !unicode.IsUpper(ch);
}


type verifier struct {
	scanner.ErrorVector;
	worklist	vector.Vector;
	reached		Grammar;	// set of productions reached from (and including) the root production
	grammar		Grammar;
}


func (v *verifier) push(prod *Production) {
	name := prod.Name.String;
	if _, found := v.reached[name]; !found {
		v.worklist.Push(prod);
		v.reached[name] = prod;
	}
}


func (v *verifier) verifyChar(x *Token) int {
	s := x.String;
	if utf8.RuneCountInString(s) != 1 {
		v.Error(x.Pos(), "single char expected, found "+s);
		return 0;
	}
	ch, _ := utf8.DecodeRuneInString(s);
	return ch;
}


func (v *verifier) verifyExpr(expr Expression, lexical bool) {
	switch x := expr.(type) {
	case nil:
		// empty expression
	case Alternative:
		for _, e := range x {
			v.verifyExpr(e, lexical)
		}
	case Sequence:
		for _, e := range x {
			v.verifyExpr(e, lexical)
		}
	case *Name:
		// a production with this name must exist;
		// add it to the worklist if not yet processed
		if prod, found := v.grammar[x.String]; found {
			v.push(prod)
		} else {
			v.Error(x.Pos(), "missing production "+x.String)
		}
		// within a lexical production references
		// to non-lexical productions are invalid
		if lexical && !isLexical(x.String) {
			v.Error(x.Pos(), "reference to non-lexical production "+x.String)
		}
	case *Token:
		// nothing to do for now
	case *Range:
		i := v.verifyChar(x.Begin);
		j := v.verifyChar(x.End);
		if i >= j {
			v.Error(x.Pos(), "decreasing character range")
		}
	case *Group:
		v.verifyExpr(x.Body, lexical)
	case *Option:
		v.verifyExpr(x.Body, lexical)
	case *Repetition:
		v.verifyExpr(x.Body, lexical)
	default:
		panic("unreachable")
	}
}


func (v *verifier) verify(grammar Grammar, start string) {
	// find root production
	root, found := grammar[start];
	if !found {
		var noPos token.Position;
		v.Error(noPos, "no start production "+start);
		return;
	}

	// initialize verifier
	v.ErrorVector.Init();
	v.worklist.Init(0);
	v.reached = make(Grammar);
	v.grammar = grammar;

	// work through the worklist
	v.push(root);
	for v.worklist.Len() > 0 {
		prod := v.worklist.Pop().(*Production);
		v.verifyExpr(prod.Expr, isLexical(prod.Name.String));
	}

	// check if all productions were reached
	if len(v.reached) < len(v.grammar) {
		for name, prod := range v.grammar {
			if _, found := v.reached[name]; !found {
				v.Error(prod.Pos(), name+" is unreachable")
			}
		}
	}
}


// Verify checks that:
//	- all productions used are defined
//	- all productions defined are used when beginning at start
//	- lexical productions refer only to other lexical productions
//
func Verify(grammar Grammar, start string) os.Error {
	var v verifier;
	v.verify(grammar, start);
	return v.GetError(scanner.Sorted);
}
