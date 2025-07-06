package codescore

import (
	"github.com/mgechev/revive/lint"
	"github.com/mgechev/revive/rule"
)

func lintConfig() *lint.Config {
	defer func() {
		if r := recover(); r != nil {
			out("[codereview] [external ast parser crash/panic]")
		}
	}()
	config := lint.Config{
		Confidence: 0.8,
		Severity:   lint.SeverityWarning,
		Rules:      map[string]lint.RuleConfig{},
	}
	for _, r := range lintRules {
		config.Rules[r.Name()] = lint.RuleConfig{}
	}
	for k, v := range config.Rules {
		if v.Severity == "" {
			v.Severity = lint.SeverityWarning
		}
		config.Rules[k] = v
	}
	return &config
}

var lintRules = []lint.Rule{
	&rule.VarDeclarationsRule{},
	&rule.DotImportsRule{},
	&rule.BlankImportsRule{},
	&rule.ExportedRule{},
	&rule.IndentErrorFlowRule{},
	&rule.RangeRule{},
	&rule.ErrorfRule{},
	&rule.ErrorStringsRule{},
	&rule.IncrementDecrementRule{},
	&rule.UnexportedReturnRule{},
	&rule.ErrorReturnRule{},
	&rule.ContextKeysType{},
	&rule.ContextAsArgumentRule{},
	&rule.EmptyBlockRule{},
	&rule.SuperfluousElseRule{},
	&rule.GetReturnRule{},
	&rule.ModifiesParamRule{},
	&rule.ConfusingResultsRule{},
	&rule.UnusedParamRule{},
	&rule.UnreachableCodeRule{},
	&rule.UnnecessaryStmtRule{},
	&rule.StructTagRule{},
	&rule.ConstantLogicalExprRule{},
	&rule.BoolLiteralRule{},
	&rule.RangeValInClosureRule{},
	&rule.RangeValAddress{},
	&rule.WaitGroupByValueRule{},
	&rule.AtomicRule{},
	&rule.CallToGCRule{},
	&rule.DuplicatedImportsRule{},
	&rule.ImportShadowingRule{},
	&rule.BareReturnRule{},
	&rule.StringOfIntRule{},
	&rule.StringFormatRule{},
	&rule.EarlyReturnRule{},
	&rule.UnconditionalRecursionRule{},
	&rule.IdenticalBranchesRule{},
	&rule.DeferRule{},
	&rule.NestedStructs{},
	&rule.UselessBreak{},
	&rule.TimeEqualRule{},
	&rule.OptimizeOperandsOrderRule{},
	&rule.UseAnyRule{},
	&rule.DataRaceRule{},
	&rule.TimeNamingRule{},
	&rule.ConfusingNamingRule{},
	&rule.UnexportedNamingRule{},
	&rule.ReceiverNamingRule{},
}

// rules that are disabled [not usable in single file process scenario]
// &rule.PackageCommentsRule{},

// rules that need validate|fine-tuning
// &rule.FlagParamRule{},
// &rule.ModifiesValRecRule{}, // panics unrecoverable
// &rule.UnhandledErrorRule{},
// &rule.IfReturnRule{},

// rules that are too strict or very individual opinion
// &rule.VarNamingRule{},
// &rule.DeepExitRule{},
// &rule.ErrorNamingRule{},
// &rule.AddConstantRule{},
// &rule.EmptyLinesRule{},

// rules that need individual configuration [up for discussion]
// &rule.RedefinesBuiltinIDRule{},
// &rule.FunctionResultsLimitRule{},
// &rule.BannedCharsRule{},
// &rule.FunctionLength{},
// &rule.FileHeaderRule{},
// &rule.MaxPublicStructsRule{},
// &rule.CyclomaticRule{},
// &rule.ArgumentsLimitRule{},
// &rule.CognitiveComplexityRule{},
// &rule.LineLengthLimitRule{},

// broken
// &rule.UnusedReceiverRule{},
