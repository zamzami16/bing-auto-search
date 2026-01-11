package main

import (
	"bufio"
	"fmt"
	"math/rand"
	"os"
	"os/exec"
	"path/filepath"
	"strconv"
	"strings"
	"time"

	"bing-auto-search/internal/app"
	"bing-auto-search/internal/words"
)

func main() {
	fmt.Println()
	fmt.Println("HUMAN-LIKE BING SEARCH AUTOMATION")
	fmt.Println()

	// Load or create config
	exePath, err := os.Executable()
	if err != nil {
		fmt.Printf("Failed to get executable path: %v\n", err)
		return
	}
	exeDir := filepath.Dir(exePath)
	configPath := filepath.Join(exeDir, "config", "mobile-config.json")

	configPath, err = app.EnsureMobileConfigFile(configPath)
	if err != nil {
		fmt.Printf("Failed to ensure config file: %v\n", err)
		return
	}

	cfg, err := app.LoadMobileConfig(configPath)
	if err != nil {
		fmt.Printf("Failed to load config: %v\n", err)
		return
	}

	fmt.Printf("Config loaded from: %s\n", configPath)
	fmt.Println()

	// Initialize random word generator
	w := words.New()

	// Get random keywords
	keywords := w.GetMany(cfg.Picks)
	if len(keywords) == 0 {
		fmt.Println("Failed to get keywords!")
		return
	}

	fmt.Printf("Selected %d random keywords from dictionary.\n", len(keywords))
	fmt.Println()

	// Show selected keywords
	for i, kw := range keywords {
		fmt.Printf("%d = %s\n", i+1, kw)
	}
	fmt.Println()
	fmt.Println("✓ Finished selecting random keywords!")
	fmt.Println()

	// Choose app
	fmt.Println("Choose which app to use:")
	fmt.Println("  [1] Bing       (com.microsoft.bing)")
	fmt.Println("  [2] Bing News  (com.microsoft.amp.apps.bingnews)")

	var appChoice string
	for {
		fmt.Print("Enter 1 or 2: ")
		reader := bufio.NewReader(os.Stdin)
		appChoice, _ = reader.ReadString('\n')
		appChoice = strings.TrimSpace(appChoice)
		if appChoice == "1" || appChoice == "2" {
			break
		}
	}

	pkg := "com.microsoft.bing"
	if appChoice == "2" {
		pkg = "com.microsoft.amp.apps.bingnews"
	}

	fmt.Println()
	fmt.Printf("Work-profile selected - manually open \"%s\" now (%d seconds)...\n", pkg, cfg.WaitAppOpen)
	time.Sleep(time.Duration(cfg.WaitAppOpen) * time.Second)

	// Start automation
	totalSearch := len(keywords)
	var prevKeyword string

	for count := 1; count <= totalSearch; count++ {
		fmt.Println()
		fmt.Printf("Cycle #%d\n", count)

		// Occasional swipe starting from configured cycle
		if count >= cfg.ScrollStartFrom {
			if rand.Intn(cfg.Scroll.Chance) == 0 { // 1/N chance
				fmt.Println("Scrolling...")
				runADB("shell", "input", "swipe", "500", "1200", "500", "600", "300")
				time.Sleep(1 * time.Second)
			}
		}

		// Calculate tap coordinates with jitter
		tapX := cfg.SearchArea.X + rand.Intn(cfg.SearchArea.DeltaX+1)
		tapY := cfg.SearchArea.Y + rand.Intn(cfg.SearchArea.DeltaY+1)

		// Add jitter ±5px
		jitterX := rand.Intn(11) - 5 // -5 to +5
		jitterY := rand.Intn(11) - 5
		tapX += jitterX
		tapY += jitterY

		fmt.Printf("Tap at X=%d Y=%d\n", tapX, tapY)
		runADB("shell", "input", "tap", strconv.Itoa(tapX), strconv.Itoa(tapY))

		// Random delay from config
		rndDelay := rand.Intn(cfg.TapDelay.Max-cfg.TapDelay.Min+1) + cfg.TapDelay.Min
		fmt.Printf("Waiting %ds...\n", rndDelay)
		time.Sleep(time.Duration(rndDelay) * time.Second)

		// Clear previous keyword if any
		if prevKeyword != "" {
			keyLen := len(prevKeyword)
			extraClear := rand.Intn(5) + 1
			clearCount := keyLen + extraClear
			fmt.Printf("Clearing %d chars from \"%s\" (len=%d + %d)\n",
				clearCount, prevKeyword, keyLen, extraClear)

			// Bangun string dengan multiple KEYCODE_DEL dalam 1 command
			deleteKeys := strings.Repeat("67 ", clearCount)
			deleteKeys = strings.TrimSpace(deleteKeys)
			runADB("shell", "input", "keyevent", deleteKeys)
			time.Sleep(time.Duration(cfg.ClearDelay) * time.Millisecond)
		}

		// Get current keyword
		currentKeyword := keywords[count-1]
		fmt.Printf("Typing \"%s\"...\n", currentKeyword)

		// Type keyword (escape single quotes for shell)
		sendText := strings.ReplaceAll(currentKeyword, "'", "\\'")
		runADB("shell", fmt.Sprintf("input text '%s'", sendText))

		// Pause before Enter from config
		rndEnter := rand.Intn(cfg.EnterDelay.Max-cfg.EnterDelay.Min+1) + cfg.EnterDelay.Min
		time.Sleep(time.Duration(rndEnter) * time.Second)

		fmt.Println("Press Enter")
		runADB("shell", "input", "keyevent", "66") // KEYCODE_ENTER
		time.Sleep(time.Duration(cfg.AfterEnter) * time.Second)

		prevKeyword = currentKeyword
	}

	fmt.Println()
	fmt.Printf("Completed %d searches!\n", totalSearch)
	fmt.Print("Press Enter to exit...")
	bufio.NewReader(os.Stdin).ReadString('\n')
}

func runADB(args ...string) {
	cmd := exec.Command("adb", args...)
	cmd.Stdout = nil
	cmd.Stderr = nil
	_ = cmd.Run()
}
