" Vim syntax file
" Filenames: *.gxs

if exists("b:current_syntax")
  finish
end

syn match gxsComment '\s*#.*$'
syn match gxsDirective '\v^(palette|action|mode|pattern|offset)' skipwhite
syn match gxsOperator '=>'
syn match gxsMode '\<\(commit\|leftedge\|bottomedge\|topedge\|rightedge\|xstitch\|tlbrline\|trblline\|hline\|vline\)\>'
syn match gxsOffset '\d\+x\d\+'

hi link gxsComment Comment
hi link gxsDirective Identifier
hi link gxsOperator Statement
hi link gxsMode Constant
hi link gxsOffset Operator

let b:current_syntax = "gxs"
