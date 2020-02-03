language swift(go);


lang = "swift"
package = "github.com/swift-land/llswift"                      # generate an event-based parser

eventBased = true
eventFields = true
eventAST = true
fileNode = "Module"
cancellable = true
recursiveLookaheads = true

reportTokens = [MultiLineComment, SingleLineComment, invalid_token]    # report these tokens as parser events

# ... lexer specification
:: lexer

# Accept end-of-input in all states.
<*> eoi: /{eoi}/

invalid_token:
error:

# LineTerminatorSequence
WhiteSpace: /[\n\r\u2028\u2029]|\r\n/ (space)

commentChars = /([^*]|\*+[^*\/])*\**/
MultiLineComment:  /\/\*{commentChars}\*\//     (space)
# Note: the following rule disables backtracking for incomplete multiline comments, which
# would otherwise be reported as '/', '*', etc.
invalid_token: /\/\*{commentChars}/
SingleLineComment: /\/\/[^\n\r\u2028\u2029]*/   (space)

# Shebang.
SingleLineComment: /#![^\n\r\u2028\u2029]*/   (space)

# Note: see http://unicode.org/reports/tr31/
IDStart = /\p{Lu}|\p{Ll}|\p{Lt}|\p{Lm}|\p{Lo}|\p{Nl}/
IDContinue = /{IDStart}|\p{Mn}|\p{Mc}|\p{Nd}|\p{Pc}/
JoinControl = /\u200c|\u200d/

hex = /[0-9a-fA-F]/
unicodeEscapeSequence = /u(\{{hex}+\}|{hex}{4})/
brokenEscapeSequence = /\\(u({hex}{0,3}|\{{hex}*))?/

identifierStart = /{IDStart}|$|_|\\{unicodeEscapeSequence}/
identifierPart =  /{identifierStart}|{IDContinue}|{JoinControl}/

Identifier: /{identifierStart}{identifierPart}*/    (class)
# Note: the following rule disables backtracking for incomplete identifiers.
invalid_token: /({identifierStart}{identifierPart}*)?{brokenEscapeSequence}/







# Literals.
'nil' : /nil/
'true':  /true/
'false': /false/

'Bool': /Bool/
'Double':  /Double/
'String':  /String/

# modifier
'internal':   /internal/
'private':      /private/
'fileprivate':    /fileprivate/
'public':       /public/
# open # not supported

'convenience':   /convenience/
'dynamic':   /dynamic/
'final':   /final/
'lazy' :   /lazy/
'mutating' :   /mutating/
'nonmutating':   /nonmutating/

'optional':   /optional/
'override' :   /override/
'required' :   /required/
'static':   /static/
'unowned':   /unowned/
'weak':   /weak/


# Keywords.
# 'await':      /await/
'As':     /As/
'Self': /Self/

# associativity  # Future-reserved. AdvancedOperators
# right # Future-reserved. AdvancedOperators
# none # Future-reserved. AdvancedOperators
# left # Future-reserved. AdvancedOperators

# operators (prefix, binary, postfix, assignment)
# operator 
# postfix 
# prefix 

# precedence 


'break':      /break/
'catch':      /catch/
'case':       /case/
'class':      /class/
'continue':   /continue/
'default':    /default/
'defer':    /defer/
'deinit':    /deinit/
'didSet': /didSet/
'do':         /do/
'enum':  /enum/

'extension' : /extension/
'else' : /else/
'fallthrough' : /fallthrough/
'for' : /for/
'func' : /func/
'get' : /get/
'guard' : /guard/
'if': /if/
'import' : /import/
'in' : /in/
'indirect' : /indirect/
'infix' : /infix/
'init' : /init/
'inout' : /inout/
'is' : /is/
'let': /let/




'protocol' : /protocol/


'repeat' : /repeat/
'rethrows' : /rethrows/
'return' : /return/



'self' : /self/
'set': /set/
'struct' : /struct/

'subscript' : /subscript/
'safe' : /safe/ # subscript in-Future

# unsafe : /unsafe/  #  in-Future

'super' : /super/
'switch' : /switch/
'throw' : /throw/
'throws' : /throws/
'try': /try/
'typealias' : /typealias/

'var' : /var/
'where' : /where/
'while' : /while/
'willSet' : /willSet/

# metatype-type
'Type': /Type/
'Protocol': /Protocol/
'Self': /Self/



# 'async':  /async/
# 'await':      /await/
# 'yield':      /yield/




# End of swift keywords.
# Punctuation

'{': /\{/
'}':          /* See below */
'(': /\(/
')': /\)/
'[': /\[/
']': /\]/
'.': /\./
invalid_token: /\.\./
'...': /\.\.\./
';': /;/
',': /,/
'<': /</
'>': />/
'<=': /<=/
'>=': />=/
'==': /==/
'!=': /!=/
'===': /===/
'!==': /!==/
'@': /@/
'+': /\+/
'-': /-/
'*': /\*/
'/':          /* See below */
'%': /%/
'++': /\+\+/
'--': /--/
# '<<': /<</
# '>>': />>/
# '>>>': />>>/
# '&': /&/
# '|': /\|/
# '^': /^/
'!': /!/
# '~': /~/
'&&': /&&/
'||': /\|\|/
'?': /\?/
'??': /\?\?/
invalid_token: /\?\.[0-9]/   { l.rewind(l.tokenOffset+1); token = QUEST }
# '?.': /\?\./
':': /:/
'=': /=/
'+=': /\+=/
'-=': /-=/
'*=': /\*=/
'/=':         /* See below */
'%=': /%=/
# '<<=': /<<=/
# '>>=': />>=/
# '>>>=': />>>=/
# '&=': /&=/
# '|=': /\|=/
# '^=': /^=/

'->': /->/



# ... parser
# :: parser

# # %input Module;

# %assert empty set(follow error & ~('}' | ')' | ',' | ';' | ']'));

# %generate afterErr = set(follow error);

# # Keywords
# | 'nil' |'true'|'false'|'Bool'|'Double'|'String'
# |'internal'|'private'|'fileprivate'|'public'
# |'convenience'|'dynamic'|'final'|'lazy' |'mutating' |'nonmutating'|'optional'|'override' |'required'|'static'|'unowned'|'weak'
# |'As'|'Self'|'break'|'catch'|'case'|'class'|'continue'|'default'|'defer'|'deinit'|'didSet'|'do'|'enum'|'extension' |'else' 
# |'fallthrough' |'for' |'func' |'get' |'guard' |'if'|'import' |'in' |'indirect' |'infix' |'init' |'inout' |'is' 
# |'let'|'protocol' |'repeat'|'rethrows'|'return' |'self' |'set'|'struct' |'subscript' |'safe' |'super' 
# |'switch'|'throw'|'throws' |'try'|'typealias'|'var'|'where' |'while'|'willSet'
# |'Type'|'Protocol'|'Self'
# ;

# IdentifierNameDecl :
#     IdentifierName                                    -> BindingIdentifier
# ;

# IdentifierNameRef :
#     IdentifierName                                    -> IdentifierReference
# ;
