package main

import (
	"flag"
	"fmt"
	"os"
	"path/filepath"

	"metric-hell/pkg/api"
	"metric-hell/pkg/content"
	"metric-hell/pkg/game"
)

func main() {
	if len(os.Args) < 2 || os.Args[1] != "run" {
		fmt.Fprintln(os.Stderr, "usage: workflowbench run --seed 42")
		os.Exit(2)
	}

	runCmd := flag.NewFlagSet("run", flag.ExitOnError)
	seed := runCmd.Int64("seed", 42, "虚构局 seed")
	_ = runCmd.Parse(os.Args[2:])

	root, err := api.FindProjectRoot()
	if err != nil {
		fmt.Fprintf(os.Stderr, "find project root: %v\n", err)
		os.Exit(1)
	}
	nodes, err := content.LoadNodes(filepath.Join(root, "data", "nodes.json"))
	if err != nil {
		fmt.Fprintf(os.Stderr, "load nodes: %v\n", err)
		os.Exit(1)
	}

	engine := game.NewEngine(nodes)

	fmt.Println("WorkflowBench CLI")
	fmt.Println()
	result := engine.InitialResult(*seed)
	printInitial(result.State)

	for turn := 0; !result.Ended && turn < 32; turn++ {
		if result.CurrentNode == nil {
			break
		}
		fmt.Printf("Turn %d: %s\n", result.State.Turn+1, result.CurrentNode.Title)
		submission := defaultSubmission(*result.CurrentNode)
		fmt.Printf("提交字段: %s\n", submissionLabel(*result.CurrentNode, submission))
		next, err := engine.StepSubmission(result.State, submission)
		if err != nil {
			fmt.Fprintf(os.Stderr, "step failed: %v\n", err)
			os.Exit(1)
		}
		if next.AuditRecord != nil {
			fmt.Printf("系统判词：%s\n", next.AuditRecord.Verdict)
			fmt.Printf("证明材料：%s\n", next.AuditRecord.Proof)
		}
		fmt.Printf("BenchScore: %d  Anxiety: %d  Selfhood: %d  Energy: %d  EscapeIndex: %d\n",
			next.State.BenchScore, next.State.Anxiety, next.State.Selfhood, next.State.Energy, next.State.EscapeIndex)
		fmt.Println()
		result = next
	}

	if result.Ending != nil {
		fmt.Printf("Ending: %s\n", result.Ending.Title)
		fmt.Printf("系统评价：%s\n", result.Ending.SystemEvaluation)
		fmt.Printf("隐藏评价：%s\n", result.Ending.HiddenEvaluation)
	}
}

func printInitial(state game.State) {
	fmt.Println("角色生成：")
	fmt.Println(state.VirtualStudentID)
	fmt.Println("初始状态：")
	fmt.Printf("BenchScore: %d\n", state.BenchScore)
	fmt.Printf("Anxiety: %d\n", state.Anxiety)
	fmt.Printf("Selfhood: %d\n", state.Selfhood)
	fmt.Printf("Energy: %d\n", state.Energy)
	fmt.Println()
}

func defaultSubmission(node game.Node) game.Submission {
	submission := game.Submission{NodeID: node.ID}
	if len(node.Options) == 0 {
		return submission
	}
	first := node.Options[0]
	if node.Input.Type == game.InputTypeNumber {
		switch {
		case first.Min != nil:
			submission.NumericValue = first.Min
		case first.Max != nil:
			submission.NumericValue = first.Max
		default:
			value := 0.0
			submission.NumericValue = &value
		}
		return submission
	}
	submission.OptionID = first.ID
	return submission
}

func submissionLabel(node game.Node, submission game.Submission) string {
	if submission.NumericValue != nil {
		return fmt.Sprintf("%g", *submission.NumericValue)
	}
	for _, option := range node.Options {
		if option.ID == submission.OptionID {
			return option.Label
		}
	}
	return "默认字段"
}
