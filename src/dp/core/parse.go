package dp

import (
	"dp/err"
	"fmt"
	"strings"
)

const (
	NOT_FOUND       rune = 0
	FOUND_ENDING    rune = 1
	FOUND_BEGINNING rune = 2
)

const (
	BRANCH_TOKEN  = "->"
	DIMENSION_SEP = ","
)

const CONDITION_ELSE = "else"

/// clean the expression
func Clean(str string) string {
	return strings.Replace(
		strings.Replace(
			strings.Replace(
				strings.Replace(
					strings.Replace(
						str, "and", "&&", -1),
					"or", "||", -1),
				"\n", "", -1),
			" ", "", -1),
		"\t", "", 1)
}

func newCondition(str string) string {
	var BRACEState = NOT_FOUND
	var endingIndex = -1
	var beginningIndex = -1
	for i := len(str) - 1; i >= 0; i-- {
		switch BRACEState {
		case NOT_FOUND:
			if str[i] == BRACE_SMALL_RIGHT {
				BRACEState = FOUND_ENDING
				endingIndex = i
			}
			break
		case FOUND_ENDING:
			if str[i] == BRACE_SMALL_LEFT {
				BRACEState = FOUND_BEGINNING
				beginningIndex = i + 1
				if endingIndex != -1 {
					return str[beginningIndex:endingIndex]
				} else {
					err.Raise("Ending BRACE not found")
					return CONDITION_ELSE
				}
			}
			break
		case FOUND_BEGINNING:
			err.Raise("Program has been mysteriously exited")
			return CONDITION_ELSE
		}
	}
	return CONDITION_ELSE
}

func newExpression(str string) string {
	for i := 0; i < len(str); i++ {
		if str[i] == BRACE_SMALL_LEFT || str[i] == BRACE_SMALL_RIGHT {
			return str[:i]
		}
	}
	return str
}

func newBranch(code string) *branch {
	ret := new(branch)
	ret.Condition = newCondition(code)
	ret.Expression = newExpression(code)
	ret.IsDefault = ret.Condition == CONDITION_ELSE
	return ret
}

func newState(str string) *state {
	ret := new(state)
	sep := strings.Index(str, string(BRACE_MIDDLE_LEFT))
	if sep == -1 {
		err.Raise("Main expression error")
	}
	ret.NameExpr = str
	ret.Name = str[:sep]
	ret.DimExpr = strings.Split(str[sep + 1:len(str) - 1], DIMENSION_SEP)
	ret.RelationExpr = make([]string, 0) /// TODO
	return ret
}

/// parse the given source code
func Parse(source string) *dyProInfo {
	checkBrace(source)
	ret := newStateEquation(source)
	checkSymbol(ret)
	checkDimension(ret)
	fmt.Errorf(err.GetErrors())
	return ret
}

func newStateEquation(code string) *dyProInfo {
	source := Clean(code)
	ret := new(dyProInfo)
	ret.Type = "int"
	ret.Detail = *NewImplDetail(101)
	split := strings.Split(source, BRANCH_TOKEN)
	if len(split) < 2 {
		err.Raise("Require branches!")
	} else {
		branches := make([]branch, len(split) - 1)
		for index, i := range split[1:] {
			branches[index] = *newBranch(i)
		}
		ret.Branches = branches
		ret.State = *newState(split[0])
	}
	return ret
}
