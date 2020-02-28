// Copyright 2019 Tobias Guggenmos
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
// http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package cache

import (
	"fmt"
	"go/token"
	"reflect"
	"testing"

	promql "github.com/prometheus/prometheus/promql/parser"
)

func TestSmallestSurroundingNode(t *testing.T) { //nolint:funlen
	shouldMatchFull := []struct {
		input string
		pos   token.Pos
	}{
		{
			input: "1",
			pos:   1,
		}, {
			input: "+1 + -2 * 1",
			pos:   4,
		},
	}

	for _, test := range shouldMatchFull {
		parseResult, err := promql.ParseExpr(test.input)
		if err != nil {
			panic("Parser should not have failed on " + test.input)
		}

		node := getSmallestSurroundingNode(&CompiledQuery{Ast: parseResult}, test.pos)

		if !reflect.DeepEqual(node, parseResult) {
			panic("Whole Expression should have been matched for " + test.input)
		}
	}

	var testExpressions = []string{
		"1",
		" 1",
		"-1",
		"+Inf",
		"-Inf",
		".5",
		"5.",
		"123.4567",
		"5e-3",
		"5e3",
		"0xc",
		"0755",
		"+5.5e-3",
		"-0755",
		"1 + 1",
		"1 - 1",
		"1 * 1",
		"1 % 1",
		"1 / 1",
		"1 == bool 1",
		"1 != bool 1",
		"1 > bool 1",
		"1 >= bool 1",
		"1 < bool 1",
		"1 <= bool 1",
		"+1 + -2 * 1",
		"1 + 2/(3*1)",
		"1 < bool 2 - 1 * 2",
		"-some_metric",
		"+some_metric",
		"",
		"# just a comment\n\n",
		"1+",
		".",
		"2.5.",
		"100..4",
		"0deadbeef",
		"1 /",
		"*1",
		"(1))",
		"((1)",
		"999999999999999999999999999999999999999999999999999999999999999999999999999999999999999999999999999999999999999999999999999999999999999999999999999999999999999999999999999999999999999999999999999999999999999999999999999999999999999999999999999999999999999999999999999999999999999999999999999999999999999999999",
		"(",
		"1 and 1",
		"1 == 1",
		"1 or 1",
		"1 unless 1",
		"1 !~ 1",
		"1 =~ 1",
		`-"string"`,
		`-test[5m]`,
		`*test`,
		"1 offset 1d",
		"a - on(b) ignoring(c) d",
		"foo * bar",
		"foo == 1",
		"foo == bool 1",
		"2.5 / bar",
		"foo and bar",
		"foo unless bar",
		"foo + bar or bla and blub",
		"foo and bar unless baz or qux",
		"bar + on(foo) bla / on(baz, buz) group_right(test) blub",
		"foo * on(test,blub) bar",
		"foo * on(test,blub) group_left bar",
		"foo and on(test,blub) bar",
		"foo and on() bar",
		"foo and ignoring(test,blub) bar",
		"foo and ignoring() bar",
		"foo unless on(bar) baz",
		"foo / on(test,blub) group_left(bar) bar",
		"foo / ignoring(test,blub) group_left(blub) bar",
		"foo - on(test,blub) group_right(bar,foo) bar",
		"foo - ignoring(test,blub) group_right(bar,foo) bar",
		"foo and 1",
		"1 and foo",
		"foo or 1",
		"1 or foo",
		"foo unless 1",
		"1 unless foo",
		"1 or on(bar) foo",
		"foo == on(bar) 10",
		"foo and on(bar) group_left(baz) bar",
		"foo and on(bar) group_right(baz) bar",
		"foo or on(bar) group_left(baz) bar",
		"foo or on(bar) group_right(baz) bar",
		"foo unless on(bar) group_left(baz) bar",
		"foo unless on(bar) group_right(baz) bar",
		`http_requests{group="production"} + on(instance) group_left(job,instance) cpu_count{type="smp"}`,
		"foo + bool bar",
		"foo + bool 10",
		"foo and bool 10",
		"foo",
		"foo offset 5m",
		`foo:bar{a="bc"}`,
		`foo{NaN='bc'}`,
		`foo{a="b", foo!="bar", test=~"test", bar!~"baz"}`,
		`{`,
		`}`,
		`some{`,
		`some}`,
		`some_metric{a=b}`,
		`some_metric{a:b="b"}`,
		`foo{a*"b"}`,
		`foo{a>="b"}`,
		"some_metric{a=\"\xff\"}",
		`foo{gibberish}`,
		`foo{1}`,
		`{}`,
		`{x=""}`,
		`{x=~".*"}`,
		`{x!~".+"}`,
		`{x!="a"}`,
		`foo{__name__="bar"}`,
		"test[5s]",
		"test[5m]",
		"test[5h] OFFSET 5m",
		"test[5d] OFFSET 10s",
		"test[5w] offset 2w",
		`test{a="b"}[5y] OFFSET 3d`,
		`foo[5mm]`,
		`foo[0m]`,
		`foo[5m30s]`,
		`foo[5m] OFFSET 1h30m`,
		`foo["5m"]`,
		`foo[]`,
		`foo[1]`,
		`some_metric[5m] OFFSET 1`,
		`some_metric[5m] OFFSET 1mm`,
		`some_metric[5m] OFFSET`,
		`some_metric OFFSET 1m[5m]`,
		`(foo + bar)[5m]`,
		"sum by (foo)(some_metric)",
		"avg by (foo)(some_metric)",
		"max by (foo)(some_metric)",
		"sum without (foo) (some_metric)",
		"sum (some_metric) without (foo)",
		"stddev(some_metric)",
		"stdvar by (foo)(some_metric)",
		"sum by ()(some_metric)",
		"topk(5, some_metric)",
		"count_values(\"value\", some_metric)",
		"sum without(and, by, avg, count, alert, annotations)(some_metric)",
		"sum without(==)(some_metric)",
		`sum some_metric by (test)`,
		`sum (some_metric) by test`,
		`sum (some_metric) by test`,
		`sum () by (test)`,
		"MIN keep_common (some_metric)",
		"MIN (some_metric) keep_common",
		`sum (some_metric) without (test) by (test)`,
		`sum without (test) (some_metric) by (test)`,
		`topk(some_metric)`,
		`topk(some_metric, other_metric)`,
		`count_values(5, other_metric)`,
		"time()",
		`floor(some_metric{foo!="bar"})`,
		"rate(some_metric[5m])",
		"round(some_metric)",
		"round(some_metric, 5)",
		"floor()",
		"floor(some_metric, other_metric)",
		"floor(1)",
		"non_existent_function_far_bar()",
		"rate(some_metric)",
		"label_replace(a, `b`, `c\xff`, `d`, `.*`)",
		"-=",
		"++-++-+-+-<",
		"e-+=/(0)",
		`"double-quoted string \" with escaped quote"`,
		`'single-quoted string \' with escaped quote'`,
		"`backtick-quoted string`",
		`"\a\b\f\n\r\t\v\\\" - \xFF\377\u1234\U00010111\U0001011111☺"`,
		`'\a\b\f\n\r\t\v\\\' - \xFF\377\u1234\U00010111\U0001011111☺'`,
		"`" + `\a\b\f\n\r\t\v\\\"\' - \xFF\377\u1234\U00010111\U0001011111☺` + "`",
		"`\\``",
		`"\`,
		`"\c"`,
		`"\x."`,
		`foo{bar="baz"}[10m:6s]`,
		`foo[10m:]`,
		`min_over_time(rate(foo{bar="baz"}[2s])[5m:5s])`,
		`min_over_time(rate(foo{bar="baz"}[2s])[5m:])[4m:3s]`,
		`min_over_time(rate(foo{bar="baz"}[2s])[5m:] offset 4m)[4m:3s]`,
		"sum without(and, by, avg, count, alert, annotations)(some_metric) [30m:10s]",
		`some_metric OFFSET 1m [10m:5s]`,
		`(foo + bar{nm="val"})[5m:]`,
		`(foo + bar{nm="val"})[5m:] offset 10m`,
		"test[5d] OFFSET 10s [10m:5s]",
		`(foo + bar{nm="val"})[5m:][10m:5s]`,
		`{} 1 2 3`,
		`{a="b"} -1 2 3`,
		`my_metric 1 2 3`,
		`my_metric{} 1 2 3`,
		`my_metric{a="b"} 1 2 3`,
		`my_metric{a="b"} 1 2 3-10x4`,
		`my_metric{a="b"} 1 2 3-0x4`,
		`my_metric{a="b"} 1 3 _ 5 _x4`,
		`my_metric{a="b"} 1 3 _ 5 _a4`,
		`my_metric{a="b"} 1 -1`,
		`my_metric{a="b"} 1 +1`,
		`my_metric{a="b"} 1 -1 -3-10x4 7 9 +5`,
		`my_metric{a="b"} 1 +1 +4 -6 -2 8`,
		`my_metric{a="b"} 1 2 3    `,
		`my_metric{a="b"} -3-3 -3`,
		`my_metric{a="b"} -3 -3-3`,
		`my_metric{a="b"} -3 _-2`,
		`my_metric{a="b"} -3 3+3x4-4`,
	}

	for _, test := range testExpressions {
		parseResult, _ := promql.ParseExpr(test)

		for pos := 1; pos <= len(test)+1; pos++ {
			node := getSmallestSurroundingNode(&CompiledQuery{Ast: parseResult}, token.Pos(pos))

			// If we are outside the outermost Expression, nothing should be matched
			if parseResult == nil || int(parseResult.PositionRange().Start) > pos || int(parseResult.PositionRange().End) < pos {
				if node != nil {
					panic("nothing should have been matched")
				}

				continue
			}

			if node == nil || int(node.PositionRange().Start) > pos || int(node.PositionRange().End) < pos {
				panic("The smallestSurroundingNode is not actually surrounding for input " + test + " and pos " + fmt.Sprintln(pos))
			}
		}
	}
}
