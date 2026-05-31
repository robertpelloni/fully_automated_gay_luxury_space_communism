package main

import (
	"bufio"
	"flag"
	"fmt"
	"github.com/robertpelloni/hustle/hustle/curation"
	"github.com/robertpelloni/hustle/hustle/research"
	"github.com/robertpelloni/hustle/hustle/social"
	"github.com/robertpelloni/hustle/orchestrator"
	"os"
	"strings"
	"time"
)

func main() {
	hustleTask := flag.String("hustle", "", "Run a specific hustle module (research, social, curation)")
	syncMode := flag.Bool("sync", false, "Run repository synchronization protocol")
	params := flag.String("params", "{}", "JSON parameters for the hustle module")
	dashboard := flag.Bool("dashboard", false, "Start the live terminal dashboard")
	daemon := flag.Bool("daemon", false, "Run the Orchestrator as a background task scheduler")
	interactive := flag.Bool("interactive", false, "Launch the interactive command menu")
	flag.Parse()

	// Source version from VERSION.md
	versionData, err := os.ReadFile("VERSION.md")
	version := "unknown"
	if err == nil {
		version = strings.TrimSpace(string(versionData))
	}

	fmt.Printf("=== AI Hustle Machine Orchestrator v%s ===\n", version)

	orch := orchestrator.NewOrchestrator()

	if *interactive {
		runInteractiveMenu(orch, version)
		return
	}

	if *dashboard {
		orchestrator.StartLiveDashboard(orch)
		return
	}

	if *daemon {
		scheduler := orchestrator.NewScheduler(orch)

		// Register Hustle Tasks
		scheduler.Register("Research", 1*time.Hour, func(o *orchestrator.Orchestrator) error {
			searcher := &research.ResearchSearch{Orchestrator: o}
			_, err := searcher.Query("AI Agent Trends")
			return err
		})

		scheduler.Register("Curation", 2*time.Hour, func(o *orchestrator.Orchestrator) error {
			c := &curation.CurationModule{
				Orchestrator: o,
				Fetcher:      curation.NewRSSFetcher(),
				Feeds:        []string{"https://news.ycombinator.com/rss"},
			}
			return c.Curate("Technology")
		})

		scheduler.Register("Heartbeat", 5*time.Minute, func(o *orchestrator.Orchestrator) error {
			return orchestrator.WriteStatusReport(version, "Active", "Scheduler Heartbeat", o.Ledger)
		})

		fmt.Println("Orchestrator running in Daemon mode.")
		scheduler.Start()
		return
	}

	if *syncMode {
		fmt.Println("Triggering repository sync protocol...")
		// Internal call to sync logic would go here
	}

	if *hustleTask != "" {
		fmt.Printf("Launching hustle module: %s with params: %s\n", *hustleTask, *params)
		launchHustle(orch, *hustleTask)
	}

	// Real-time status reporting with financial metrics
	orchestrator.WriteStatusReport(version, "Running", "Orchestrator command processed", orch.Ledger)

	fmt.Printf("Orchestrator initialized with L1: %d, L2: %d, L3: %d entries. Profit: $%.2f\n",
		len(orch.L1.Entries), len(orch.L2.Entries), len(orch.L3.Entries), orch.Ledger.Profit())
}

func launchHustle(orch *orchestrator.Orchestrator, task string) {
	switch strings.ToLower(task) {
	case "research":
		searcher := &research.ResearchSearch{Orchestrator: orch}
		searcher.Query("AI Trends")
	case "curation":
		c := &curation.CurationModule{
			Orchestrator: orch,
			Fetcher:      curation.NewRSSFetcher(),
			Feeds:        []string{"https://news.ycombinator.com/rss"},
		}
		c.Curate("AI")
	case "social":
		provider := social.NewTwitterProvider()
		social.SchedulePost(orch, provider, "Twitter", "AI")
	default:
		fmt.Printf("Unknown hustle: %s\n", task)
	}
}

func runInteractiveMenu(orch *orchestrator.Orchestrator, version string) {
	reader := bufio.NewReader(os.Stdin)
	for {
		fmt.Println("\n--- INTERACTIVE COMMAND MENU ---")
		fmt.Println("1. Launch Research Hustle")
		fmt.Println("2. Launch Social Hustle")
		fmt.Println("3. Launch Curation Hustle")
		fmt.Println("4. View Dashboard")
		fmt.Println("5. Run Repository Sync")
		fmt.Println("q. Quit")
		fmt.Print("Select an option: ")

		input, _ := reader.ReadString('\n')
		input = strings.TrimSpace(input)

		switch input {
		case "1":
			launchHustle(orch, "research")
		case "2":
			launchHustle(orch, "social")
		case "3":
			launchHustle(orch, "curation")
		case "4":
			orchestrator.ShowDashboard(orch)
			fmt.Println("\nPress Enter to return to menu...")
			reader.ReadString('\n')
		case "5":
			fmt.Println("Running Sync Protocol...")
			// Logic to invoke sync
		case "q":
			fmt.Println("Exiting...")
			return
		default:
			fmt.Println("Invalid option, please try again.")
		}
	}
}
