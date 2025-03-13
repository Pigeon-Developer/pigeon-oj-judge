package types

const (
	Language_c          = 0
	Language_cpp        = 1
	Language_pascal     = 2
	Language_java       = 3
	Language_ruby       = 4
	Language_bash       = 5
	Language_python     = 6
	Language_php        = 7
	Language_perl       = 8
	Language_csharp     = 9
	Language_objectivec = 10
	Language_freebasic  = 11
	Language_scheme     = 12
	Language_clang      = 13
	Language_clangpp    = 14
	Language_lua        = 15
	Language_javascript = 16
	Language_golang     = 17
	Language_sql        = 18
	Language_fortran    = 19
	Language_matlab     = 20
	Language_cobol      = 21
	Language_r          = 22
	Language_scratch3   = 23
	Language_cangjie    = 24
)

var (
	LanguageMap = make(map[int]string)
	LangMap     = make(map[string]int)
)

func init() {
	// 1
	LanguageMap[Language_c] = "c"
	LanguageMap[Language_cpp] = "cpp"
	LanguageMap[Language_pascal] = "pascal"
	LanguageMap[Language_java] = "java"
	LanguageMap[Language_ruby] = "ruby"
	// 2
	LanguageMap[Language_bash] = "bash"
	LanguageMap[Language_python] = "python"
	LanguageMap[Language_php] = "php"
	LanguageMap[Language_perl] = "perl"
	LanguageMap[Language_csharp] = "csharp"
	// 3
	LanguageMap[Language_objectivec] = "objectivec"
	LanguageMap[Language_freebasic] = "freebasic"
	LanguageMap[Language_scheme] = "scheme"
	LanguageMap[Language_clang] = "clang"
	LanguageMap[Language_clangpp] = "clangpp"
	// 4
	LanguageMap[Language_lua] = "lua"
	LanguageMap[Language_javascript] = "javascript"
	LanguageMap[Language_golang] = "golang"
	LanguageMap[Language_sql] = "sql"
	LanguageMap[Language_fortran] = "fortran"
	// 5
	LanguageMap[Language_matlab] = "matlab"
	LanguageMap[Language_cobol] = "cobol"
	LanguageMap[Language_r] = "r"
	LanguageMap[Language_scratch3] = "scratch3"
	LanguageMap[Language_cangjie] = "cangjie"

	// 1
	LangMap["c"] = Language_c
	LangMap["cpp"] = Language_cpp
	LangMap["pascal"] = Language_pascal
	LangMap["java"] = Language_java
	LangMap["ruby"] = Language_ruby
	// 2
	LangMap["bash"] = Language_bash
	LangMap["python"] = Language_python
	LangMap["php"] = Language_php
	LangMap["perl"] = Language_perl
	LangMap["csharp"] = Language_csharp
	// 3
	LangMap["objectivec"] = Language_objectivec
	LangMap["freebasic"] = Language_freebasic
	LangMap["scheme"] = Language_scheme
	LangMap["clang"] = Language_clang
	LangMap["clangpp"] = Language_clangpp
	// 4
	LangMap["lua"] = Language_lua
	LangMap["javascript"] = Language_javascript
	LangMap["golang"] = Language_golang
	LangMap["sql"] = Language_sql
	LangMap["fortran"] = Language_fortran
	// 5
	LangMap["matlab"] = Language_matlab
	LangMap["cobol"] = Language_cobol
	LangMap["r"] = Language_r
	LangMap["scratch3"] = Language_scratch3
	LangMap["cangjie"] = Language_cangjie
}
