package main

import (
	"bufio"
	"flag"
	"fmt"
	"net/url"
	"os"
	"os/exec"
	"os/signal"
	"syscall"
	"strings"
	"time"

	"github.com/robertpelloni/hustle/hustle/content"
	"github.com/robertpelloni/hustle/hustle/curation"
	"github.com/robertpelloni/hustle/hustle/research"
	"github.com/robertpelloni/hustle/hustle/social"
	"github.com/robertpelloni/hustle/hustle/trading"
	"github.com/robertpelloni/hustle/orchestrator"
)

func main() {
	hustleTask := flag.String("hustle", "", "Run a specific hustle module")
	hustleURI := flag.String("uri", "", "Run a hustle module via URI")
	syncMode := flag.Bool("sync", false, "Run repository sync protocol")
	params := flag.String("params", "{}", "JSON parameters")
	dashboard := flag.Bool("dashboard", false, "Start the live dashboard")
	daemon := flag.Bool("daemon", false, "Run as background scheduler")
	interactive := flag.Bool("interactive", false, "Launch interactive menu")
	apiPort := flag.String("api", "", "Start the HTTP API")
	realPrices := flag.Bool("real-prices", false, "Use real-world prices")
	seedURL := flag.String("seed", "", "URL of a seed node")
	agentMode := flag.Bool("agent", false, "Run as autonomous agent loop")
	agentType := flag.String("agent-type", "general", "Agent focus")
	agentIter := flag.Int("agent-iterations", 20, "Max iterations")
	autoPlan := flag.Bool("autoplan", false, "LLM generates and executes a plan")
	dryRun := flag.Bool("dry-run", false, "Execute in dry-run mode")
	refreshRate := flag.Int("refresh", 1000, "Dashboard refresh rate (ms)")
	flag.Parse()

	versionData, err := os.ReadFile("VERSION.md")
	version := "unknown"
	if err == nil { version = strings.TrimSpace(string(versionData)) }

	fmt.Printf("=== AI Hustle Machine Orchestrator v%s ===\n", version)

	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)

	orch := orchestrator.NewOrchestrator()
	orch.Load("memory.json")
	orch.DryRun = *dryRun

	db, err := orchestrator.NewSQLiteStore("hustle.db")
	if err == nil { orch.DB = db }

	llmProvider := orchestrator.NewOpenAICompatProvider()
	if model, err := llmProvider.DetectModel(); err == nil {
		fmt.Printf("[LLM] ✅ Connected to local LLM: %s\n", model)
		orch.LLM = &orchestrator.CachedLLM{Provider: llmProvider, DB: orch.DB, Enabled: true}
		orch.Embedder = orchestrator.NewOpenAICompatEmbedder()
	} else {
		fmt.Printf("[LLM] ⚠️  Local LLM not available. Falling back to MockLLM.\n")
		orch.LLM = &orchestrator.MockLLM{}
	}

	protocol := orchestrator.NewHustleProtocol()
	chainManager := orchestrator.NewChainManager(orch, protocol)
	chainManager.LoadState("chains.json")
	discoverer := orchestrator.NewChainDiscoverer(orch, chainManager)
	broker := orchestrator.NewA2ABroker(orch)
	swarm := orchestrator.NewMemorySwarm(orch, broker)
	multiAgent := orchestrator.NewMultiAgentOrchestrator(orch, protocol, broker)

	traderModule := &trading.TradingModule{Orchestrator: orch, Broker: broker, Fetcher: &trading.MockPriceFetcher{}}
	if *realPrices {
		traderModule.Fetcher = &trading.CoinGeckoFetcher{APIKey: os.Getenv("COINGECKO_API_KEY")}
	}
	contentModule := content.NewContentModule(orch, "output/content")

	// ── Protocol Registration ──
	protocol.Register("research", func(p url.Values) error {
		searcher := research.NewResearchSearch(research.Tavily, orch, broker)
		_, err := searcher.Query(p.Get("query")); return err
	})
	protocol.Register("curation", func(p url.Values) error {
		c := &curation.CurationModule{Orchestrator: orch, Fetcher: curation.NewRSSFetcher(), Feeds: orch.RSSFeeds}
		return c.Curate(p.Get("topic"))
	})
	protocol.Register("social", func(p url.Values) error {
		var provider social.Provider
		if strings.ToLower(p.Get("platform")) == "linkedin" {
			provider = social.NewLinkedInProvider(os.Getenv("LINKEDIN_ACCESS_TOKEN"), os.Getenv("LINKEDIN_AUTHOR_URN"))
		} else {
			provider = social.NewTwitterProvider(os.Getenv("TWITTER_API_KEY"), os.Getenv("TWITTER_API_SECRET"), os.Getenv("TWITTER_ACCESS_TOKEN"), os.Getenv("TWITTER_ACCESS_SECRET"))
		}
		if *dryRun { provider.SetDryRun(true) }
		social.SchedulePost(orch, provider, p.Get("platform"), p.Get("topic")); return nil
	})
	protocol.Register("trading", func(p url.Values) error {
		traderModule.Symbol = p.Get("symbol"); return traderModule.ExecuteStrategy()
	})
	protocol.Register("content", func(p url.Values) error {
		req := content.ContentRequest{Topic: p.Get("topic"), Type: content.ContentType(p.Get("type")), Niche: p.Get("niche")}
		_, err := contentModule.Generate(req); return err
	})
	protocol.Register("swarm", func(p url.Values) error {
		action := p.Get("action")
		if action == "sync" { swarm.Sync() } else if action == "aggregate" { swarm.AggregateStatus() }
		return nil
	})
	protocol.Register("chain", func(p url.Values) error {
		if p.Get("action") == "discover" { _, err := discoverer.Discover(); return err }
		return chainManager.Execute(p.Get("name"))
	})
	protocol.Register("healer", func(p url.Values) error {
		h := orchestrator.NewHealer(orch); h.Loop(p.Get("issue")); return nil
	})
	protocol.Register("sync", func(p url.Values) error { return runSyncProtocol() })

	// ── Modes ──
	if *interactive {
		runInteractiveMenu(orch, protocol, broker, traderModule, contentModule, version, discoverer)
		return
	}
	if *dashboard {
		orchestrator.StartLiveDashboard(orch, time.Duration(*refreshRate)*time.Millisecond)
		return
	}
	if *apiPort != "" {
		api := orchestrator.NewAPI(orch, protocol, broker, chainManager, discoverer)
		api.MultiAgent = multiAgent
		api.Start(*apiPort)
		return
	}
	if *syncMode { runSyncProtocol(); return }
	if *agentMode {
		agent := orchestrator.NewAgentLoop(orch, protocol, broker, *agentType)
		agent.State.MaxIter = *agentIter
		agent.Run(); return
	}
	if *autoPlan {
		plans, _ := orchestrator.PlanHustles(orch)
		for _, p := range plans {
			if p.Priority == "high" { multiAgent.AddAgent(p.Category, 10) }
		}
		multiAgent.RunAll(); return
	}
	if *seedURL != "" {
		broker.RegisterPeer("seed", *seedURL)
		protocol.HandleURI("hustle://swarm?action=sync")
	}

	if *hustleURI != "" { protocol.HandleURI(*hustleURI) } else if *hustleTask != "" {
		protocol.HandleURI(fmt.Sprintf("hustle://%s?%s", *hustleTask, *params))
	}

	if *daemon {
		scheduler := orchestrator.NewScheduler(orch)
		scheduler.Register("Heartbeat", 5*time.Minute, func(o *orchestrator.Orchestrator) error { return nil })
		scheduler.Start()
	}
}

func runSyncProtocol() error {
	cmd := exec.Command("./sync.sh"); cmd.Stdout = os.Stdout; cmd.Stderr = os.Stderr; return cmd.Run()
}

func runMemoryMenu(orch *orchestrator.Orchestrator, reader *bufio.Reader) {
	for {
		fmt.Printf("\n--- MEMORY EXPLORATION ---\n L1:%d, L2:%d, L3:%d\n 1. Keyword Search\n 2. Semantic Search\n r. Return\nSelect: ", len(orch.L1.Entries), len(orch.L2.Entries), len(orch.L3.Entries))
		input, _ := reader.ReadString('\n'); input = strings.TrimSpace(input)
		switch input {
		case "1":
			fmt.Print("Query: "); q, _ := reader.ReadString('\n'); results := orch.L1.Search(strings.TrimSpace(q))
			for _, r := range results { fmt.Printf(" [%s] %s\n", r.Timestamp.Format("15:04"), r.Content) }
		case "2":
			fmt.Print("Concept: "); q, _ := reader.ReadString('\n')
			if orch.Embedder != nil && orch.DB != nil {
				vec, _ := orch.Embedder.Embed(strings.TrimSpace(q))
				res, _ := orch.DB.VectorSearch("L2", vec, 5)
				for _, r := range res { fmt.Printf(" - %s\n", r.Content) }
			}
		case "r": return
		}
	}
}

func runInteractiveMenu(orch *orchestrator.Orchestrator, protocol *orchestrator.HustleProtocol, broker *orchestrator.A2ABroker, traderModule *trading.TradingModule, contentModule *content.ContentModule, version string, discoverer *orchestrator.ChainDiscoverer) {
	reader := bufio.NewReader(os.Stdin)
	for {
		fmt.Println("\n--- AI HUSTLE MACHINE INTERACTIVE ---")
		fmt.Println(" 1. Research\n 2. Social\n 3. Curation\n 4. Trading\n 5. Full Chain\n 6. Swarm Sync\n 9. Dashboard\n10. Sync Repo\n16. Multi-Agent Control\n17. Workflow Discovery\n18. Content Gallery\n24. Council Debate\n25. Memory Exploration\n26. System Maintenance\n q. Quit")
		fmt.Print("Select: ")
		input, _ := reader.ReadString('\n'); input = strings.TrimSpace(input)
		switch input {
		case "1":
			fmt.Print("Query: "); q, _ := reader.ReadString('\n'); protocol.HandleURI("hustle://research?query=" + url.QueryEscape(strings.TrimSpace(q)))
		case "2":
			fmt.Print("Platform (Twitter/LinkedIn): "); p, _ := reader.ReadString('\n'); fmt.Print("Topic: "); t, _ := reader.ReadString('\n')
			protocol.HandleURI(fmt.Sprintf("hustle://social?platform=%s&topic=%s", strings.TrimSpace(p), url.QueryEscape(strings.TrimSpace(t))))
		case "3":
			fmt.Print("Topic: "); t, _ := reader.ReadString('\n'); protocol.HandleURI("hustle://curation?topic=" + url.QueryEscape(strings.TrimSpace(t)))
		case "4":
			fmt.Print("Symbol: "); s, _ := reader.ReadString('\n'); protocol.HandleURI("hustle://trading?symbol=" + strings.TrimSpace(s))
		case "5": protocol.HandleURI("hustle://chain")
		case "6": protocol.HandleURI("hustle://swarm?action=sync")
		case "9": orchestrator.ShowDashboard(orch); reader.ReadString('\n')
		case "10": runSyncProtocol()
		case "16":
			fmt.Print("Type (research/trading/content/social): "); t, _ := reader.ReadString('\n')
			agent := orchestrator.NewAgentLoop(orch, protocol, broker, strings.TrimSpace(t)); agent.State.MaxIter = 10; go agent.Run()
		case "17":
			fmt.Println("Discovering luxury chain..."); d, _ := discoverer.Discover(); fmt.Printf("Discovered: %s\n", d.Name)
		case "18": contentModule.GenerateGallery(); fmt.Println("Gallery refreshed.")
		case "24":
			fmt.Print("Topic: "); t, _ := reader.ReadString('\n'); c := orchestrator.NewCouncil(orch); res, _ := c.Debate(strings.TrimSpace(t))
			fmt.Printf("Score: %.2f\nConsensus: %s\n", res.WeightedScore, res.Consensus)
		case "25": runMemoryMenu(orch, reader)
		case "26":
			fmt.Println("Refreshing gallery..."); contentModule.GenerateGallery()
			fmt.Println("Triggering self-heal..."); protocol.HandleURI("hustle://healer?issue=Maintenance")
		case "q": orch.Shutdown(); return
		}
	}
}
