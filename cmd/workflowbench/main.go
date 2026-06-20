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
	actions := []game.Action{
		game.ActionOptimizeMetric,
		game.ActionComparePeers,
		game.ActionRest,
		game.ActionSwitchTrack,
		game.ActionJobTrack,
		game.ActionRefuseMetric,
	}

	fmt.Println("WorkflowBench CLI")
	fmt.Println()
	result := engine.InitialResult(*seed)
	printInitial(result.State)

	for turn := 0; !result.Ended && turn < 32; turn++ {
		if result.CurrentNode == nil {
			break
		}
		action := actions[turn%len(actions)]
		fmt.Printf("Turn %d: %s\n", result.State.Turn+1, result.CurrentNode.Title)
		fmt.Printf("Action: %s\n", game.ActionLabel(action))
		next, err := engine.Step(result.State, action)
		if err != nil {
			fmt.Fprintf(os.Stderr, "step failed: %v\n", err)
			os.Exit(1)
		}
		if len(next.State.EventLog) > 0 {
			fmt.Printf("系统提示：%s\n", next.State.EventLog[len(next.State.EventLog)-1])
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
