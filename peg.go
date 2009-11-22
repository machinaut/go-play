# Hierarchical syntax
Grammar     <- Spacing Definition+ EndOfFile
Definition  <- Identifier LEFTARROW Expression

Expression  <- Sequence (SLASH Sequence)*
Sequence    <- Prefix*
Prefix      <- (AND / NOT)? Suffix
Suffix      <- Primary (QUESTION / STAR / PLUS)?
Primary     <- Identifier !LEFTARROW
             / OPEN Expression CLOSE
             / Literal / Class / DOT


