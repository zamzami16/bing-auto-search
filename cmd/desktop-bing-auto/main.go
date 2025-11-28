//go:build windows

package main

import (
	"bing-auto-search/internal/app"
	"bing-auto-search/internal/uiwin"
	"bing-auto-search/internal/words"
	"flag"
	"fmt"
	"math/rand"
	"os"
	"os/exec"
	"path/filepath"
	"time"
)

var version = "dev" // injected via ldflags during build

func main() {
	random := rand.New(rand.NewSource(time.Now().UnixNano()))
	cfgPath := flag.String("config", "config/config.json", "Path config JSON")
	shutdown := flag.Bool("shutdown", false, "Auto shutdown computer setelah selesai")
	showVersion := flag.Bool("version", false, "Tampilkan versi aplikasi")
	flag.Parse()

	if *showVersion {
		fmt.Printf("Bing Auto Search %s\n", version)
		return
	}

	// Make config path relative to executable location
	exePath, err := os.Executable()
	if err != nil {
		fmt.Printf("Gagal mendapatkan path exe: %v\n", err)
		return
	}
	exeDir := filepath.Dir(exePath)
	absCfgPath := *cfgPath
	if !filepath.IsAbs(absCfgPath) {
		absCfgPath = filepath.Join(exeDir, absCfgPath)
	}

	path, _ := app.EnsureConfigFile(absCfgPath)
	cfg, err := app.LoadConfig(path)
	if err != nil {
		fmt.Printf("Load config pakai default (err: %v)", err)
	}

	gen := words.New()
	delayMin := cfg.GlobalSetting.Delay.Min
	delayMax := cfg.GlobalSetting.Delay.Max
	scrollMin := cfg.GlobalSetting.Scroll.Min
	scrollMax := cfg.GlobalSetting.Scroll.Max
	totalScrollMin := cfg.GlobalSetting.TotalScroll.Min
	totalScrollMax := cfg.GlobalSetting.TotalScroll.Max
	totalSearch := cfg.GlobalSetting.TotalSearch

	// Create UIWin instance with 250ms delay
	ui := uiwin.NewUIWin(250)

	fmt.Println("Fokuskan jendela target. Mulai sebentar lagi…")
	delayStart := 6
	for i := delayStart; i > 0; i-- {
		fmt.Printf("Mulai dalam %ds…\n", i)
		time.Sleep(1 * time.Second)
	}

	nDesktop := len(cfg.Data)
	for searchCount := range totalSearch {
		for dIdx, desktop := range cfg.Data {
			for vIdx, view := range desktop.Configs {
				fmt.Printf("Search %d/%d On Desktop %d (%s) * View %d\n", searchCount+1, totalSearch, dIdx+1, desktop.Name, vIdx+1)
				// Move and click
				x, y := GetCoordinate(view, random)
				ui.MoveSmooth(x, y, 1.0)
				ui.DoubleClick()
				ui.CtrlA()
				ui.TypeString(gen.GetOne(), 30)
				ui.PressEnter()
				ui.Sleep(GetDelay(delayMin, delayMax, random))

				// Scroll
				ui.MoveSmooth(x, y+200, 1.0)
				scrollAmount := scrollMin
				if scrollMax > scrollMin {
					scrollAmount = scrollMin + gen.Rng().Intn(scrollMax-scrollMin+1)
				}
				scrollAmount = -scrollAmount
				totalScroll := totalScrollMin
				if totalScrollMax > totalScrollMin {
					totalScroll = totalScrollMin + gen.Rng().Intn(totalScrollMax-totalScrollMin+1)
				}
				for s := 0; s < totalScroll; s++ {
					ui.ScrollNotches(scrollAmount)
					ui.Sleep(0.2)
				}

				ui.Sleep(GetDelay(delayMin, delayMax, random))
			}

			// Switch desktop only if not last desktop
			if dIdx < nDesktop-1 {
				ui.SwitchDesktop(uiwin.DesktopRight)
			} else if nDesktop > 1 {
				// After last desktop, go back to first
				for k := 0; k < nDesktop-1; k++ {
					ui.SwitchDesktop(uiwin.DesktopLeft)
					ui.Sleep(0.8)
				}
			}
			ui.Sleep(GetDelay(delayMin, delayMax, random))
		}
	}

	fmt.Println("Selesai:", time.Now().Format(time.RFC3339))

	if *shutdown {
		fmt.Println("Shutdown komputer dalam 10 detik...")
		time.Sleep(10 * time.Second)
		cmd := exec.Command("shutdown", "/s", "/t", "0")
		if err := cmd.Run(); err != nil {
			fmt.Printf("Gagal shutdown: %v\n", err)
		}
	}
}

func GetDelay(min int, max int, random *rand.Rand) float64 {
	return float64(random.Intn(max-min+1) + min)
}

func GetCoordinate(cfg app.ViewConfig, random *rand.Rand) (int, int) {
	dx := random.Intn(cfg.DX*2+1) - cfg.DX
	dy := random.Intn(cfg.DY*2+1) - cfg.DY
	return cfg.PosX + dx, cfg.PosY + dy
}
