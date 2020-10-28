package cmdlr

import (
	"github.com/karrick/tparse/v2"
	"regexp"
	"strconv"
	"strings"
	"time"
)

var (
	RegexArguments = regexp.MustCompile("(\"[^\"]+\"|[^\\s]+)")

	RegexUserMention = regexp.MustCompile("<@!?(\\d+)>")

	RegexRoleMention = regexp.MustCompile("<@&(\\d+)>")

	RegexChannelMention = regexp.MustCompile("<#(\\d+)>")

	RegexBigCodeblock = regexp.MustCompile("(?s)\\n*```(?:([\\w.\\-]*)\\n)?(.*)```")

	RegexSmallCodeblock = regexp.MustCompile("(?s)\\n*`(.*)`")

	CodeblockLanguages = []string{
		"as",
		"1c",
		"abnf",
		"accesslog",
		"actionscript",
		"ada",
		"ado",
		"adoc",
		"apache",
		"apacheconf",
		"applescript",
		"arduino",
		"arm",
		"armasm",
		"asciidoc",
		"aspectj",
		"atom",
		"autohotkey",
		"autoit",
		"avrasm",
		"awk",
		"axapta",
		"bash",
		"basic",
		"bat",
		"bf",
		"bind",
		"bnf",
		"brainfuck",
		"c",
		"c++",
		"cal",
		"capnp",
		"capnproto",
		"cc",
		"ceylon",
		"clean",
		"clj",
		"clojure-repl",
		"clojure",
		"cls",
		"cmake.in",
		"cmake",
		"cmd",
		"coffee",
		"coffeescript",
		"console",
		"coq",
		"cos",
		"cpp",
		"cr",
		"craftcms",
		"crm",
		"crmsh",
		"crystal",
		"cs",
		"csharp",
		"cson",
		"csp",
		"css",
		"d",
		"dart",
		"dcl",
		"delphi",
		"dfm",
		"diff",
		"django",
		"dns",
		"do",
		"docker",
		"dockerfile",
		"dos",
		"dpr",
		"dsconfig",
		"dst",
		"dts",
		"dust",
		"ebnf",
		"elixir",
		"elm",
		"erb",
		"erl",
		"erlang-repl",
		"erlang",
		"excel",
		"f90",
		"f95",
		"feature",
		"fix",
		"flix",
		"fortran",
		"freepascal",
		"fs",
		"fsharp",
		"gams",
		"gauss",
		"gcode",
		"gemspec",
		"gherkin",
		"glsl",
		"gms",
		"go",
		"golang",
		"golo",
		"gradle",
		"graph",
		"groovy",
		"gss",
		"gyp",
		"h",
		"h++",
		"haml",
		"handlebars",
		"haskell",
		"haxe",
		"hbs",
		"hpp",
		"hs",
		"hsp",
		"html.handlebars",
		"html.hbs",
		"html",
		"htmlbars",
		"http",
		"https",
		"hx",
		"hy",
		"hylang",
		"i7",
		"iced",
		"icl",
		"inform7",
		"ini",
		"instances",
		"irb",
		"irpf90",
		"java",
		"javascript",
		"jboss-cli",
		"jinja",
		"js",
		"json",
		"jsp",
		"jsx",
		"julia",
		"k",
		"kdb",
		"kotlin",
		"lasso",
		"lassoscript",
		"lazarus",
		"ldif",
		"leaf",
		"less",
		"lfm",
		"lisp",
		"livecodeserver",
		"livescript",
		"llvm",
		"lpr",
		"ls",
		"lsl",
		"lua",
		"m",
		"mak",
		"makefile",
		"markdown",
		"mathematica",
		"matlab",
		"maxima",
		"md",
		"mel",
		"mercury",
		"mips",
		"mipsasm",
		"mizar",
		"mk",
		"mkd",
		"mkdown",
		"ml",
		"mm",
		"mma",
		"mojolicious",
		"monkey",
		"moo",
		"moon",
		"moonscript",
		"n1ql",
		"nc",
		"nginx",
		"nginxconf",
		"nim",
		"nimrod",
		"nix",
		"nixos",
		"nsis",
		"obj-c",
		"objc",
		"objectivec",
		"ocaml",
		"openscad",
		"osascript",
		"oxygene",
		"p21",
		"parser3",
		"pas",
		"pascal",
		"patch",
		"pb",
		"pbi",
		"pcmk",
		"perl",
		"pf.conf",
		"pf",
		"php",
		"php3",
		"php4",
		"php5",
		"php6",
		"pl",
		"plist",
		"pm",
		"podspec",
		"pony",
		"powershell",
		"pp",
		"processing",
		"profile",
		"prolog",
		"protobuf",
		"ps",
		"puppet",
		"purebasic",
		"py",
		"python",
		"q",
		"qml",
		"qt",
		"r",
		"rb",
		"rib",
		"roboconf",
		"rs",
		"rsl",
		"rss",
		"ruby",
		"ruleslanguage",
		"rust",
		"scad",
		"scala",
		"scheme",
		"sci",
		"scilab",
		"scss",
		"sh",
		"shell",
		"smali",
		"smalltalk",
		"sml",
		"sqf",
		"sql",
		"st",
		"stan",
		"stata",
		"step",
		"step21",
		"stp",
		"styl",
		"stylus",
		"subunit",
		"sv",
		"svh",
		"swift",
		"taggerscript",
		"tao",
		"tap",
		"tcl",
		"tex",
		"thor",
		"thrift",
		"tk",
		"toml",
		"tp",
		"ts",
		"twig",
		"typescript",
		"v",
		"vala",
		"vb",
		"vbnet",
		"vbs",
		"vbscript-html",
		"vbscript",
		"verilog",
		"vhdl",
		"vim",
		"wildfly-cli",
		"x86asm",
		"xhtml",
		"xjb",
		"xl",
		"xls",
		"xlsx",
		"xml",
		"xpath",
		"xq",
		"xquery",
		"xsd",
		"xsl",
		"yaml",
		"yml",
		"zep",
		"zephir",
		"zone",
		"zsh",
	}
)

type Argument struct {
	raw string
}

type Arguments struct {
	raw       string
	arguments []*Argument
}

type Codeblock struct {
	Language string
	Content  string
}

func ParseArguments(raw string) *Arguments {
	rawArguments := RegexArguments.FindAllString(raw, -1)
	arguments := make([]*Argument, len(rawArguments))

	for index, rawArgument := range rawArguments {
		rawArgument = StringTrimPreSuffix(rawArgument, "\"")
		arguments[index] = &Argument{
			raw: rawArgument,
		}
	}

	return &Arguments{
		raw:       raw,
		arguments: arguments,
	}
}

func (a *Arguments) Raw() string {
	return a.raw
}

func (a *Arguments) AsSingle() *Argument {
	return &Argument{
		raw: a.raw,
	}
}

func (a *Arguments) Amount() int {
	return len(a.arguments)
}

func (a *Arguments) Get(n int) *Argument {
	if a.Amount() <= n {
		return &Argument{
			raw: "",
		}
	}
	return a.arguments[n]
}

// Remove removes the n'th argument
func (a *Arguments) Remove(n int) {
	// Check if the given index is valid
	if a.Amount() <= n {
		return
	}

	// Set the new argument slice
	a.arguments = append(a.arguments[:n], a.arguments[n+1:]...)

	// Set the new raw string
	raw := ""
	for _, argument := range a.arguments {
		raw += argument.raw + " "
	}
	a.raw = strings.TrimSpace(raw)
}

// AsCodeblock parses the given arguments as a codeblock
func (a *Arguments) AsCodeblock() *Codeblock {
	raw := a.Raw()

	// Check if the raw string is a big codeblock
	matches := RegexBigCodeblock.MatchString(raw)
	if !matches {
		// Check if the raw string is a small codeblock
		matches = RegexSmallCodeblock.MatchString(raw)
		if matches {
			submatches := RegexSmallCodeblock.FindStringSubmatch(raw)
			return &Codeblock{
				Language: "",
				Content:  submatches[1],
			}
		}
		return nil
	}

	// Define the content and the language
	submatches := RegexBigCodeblock.FindStringSubmatch(raw)
	language := ""
	content := submatches[1] + submatches[2]
	if submatches[1] != "" && !StringArrayContains(CodeblockLanguages, submatches[1], false) {
		language = submatches[1]
		content = submatches[2]
	}

	// Return the codeblock
	return &Codeblock{
		Language: language,
		Content:  content,
	}
}

// Raw returns the raw string value of the argument
func (arg *Argument) Raw() string {
	return arg.raw
}

// AsBool parses the given argument into a boolean
func (arg *Argument) AsBool() (bool, error) {
	return strconv.ParseBool(arg.raw)
}

// AsInt parses the given argument into an int32
func (arg *Argument) AsInt() (int, error) {
	return strconv.Atoi(arg.raw)
}

// AsInt64 parses the given argument into an int64
func (arg *Argument) AsInt64() (int64, error) {
	return strconv.ParseInt(arg.raw, 10, 64)
}

// AsUserMentionID returns the ID of the mentioned user or an empty string if it is no mention
func (arg *Argument) AsUserMentionID() string {
	// Check if the arg is a user mention
	matches := RegexUserMention.MatchString(arg.raw)
	if !matches {
		return ""
	}

	// Parse the user ID
	userID := RegexUserMention.FindStringSubmatch(arg.raw)[1]
	return userID
}

// AsRoleMentionID returns the ID of the mentioned role or an empty string if it is no mention
func (arg *Argument) AsRoleMentionID() string {
	// Check if the arg is a role mention
	matches := RegexRoleMention.MatchString(arg.raw)
	if !matches {
		return ""
	}

	// Parse the role ID
	roleID := RegexRoleMention.FindStringSubmatch(arg.raw)[1]
	return roleID
}

// AsChannelMentionID returns the ID of the mentioned channel or an empty string if it is no mention
func (arg *Argument) AsChannelMentionID() string {
	// Check if the arg is a channel mention
	matches := RegexChannelMention.MatchString(arg.raw)
	if !matches {
		return ""
	}

	// Parse the channel ID
	channelID := RegexChannelMention.FindStringSubmatch(arg.raw)[1]
	return channelID
}

// AsDuration parses the given argument into a duration
func (arg *Argument) AsDuration() (time.Duration, error) {
	return tparse.AbsoluteDuration(time.Now(), arg.raw)
}
