package main

import (
	"context"
	"log"
	"os/exec"
	"time"

	"github.com/chromedp/cdproto/browser"
	"github.com/chromedp/cdproto/cdp"
	"github.com/chromedp/cdproto/network"
	"github.com/chromedp/chromedp"
)

func main() {

	URL := "https://class.kavano.org/class/azad-admin/test01"
	TOKEN := "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJpZCI6IjYwZjY0ZjI1YWZmYjdiMDAxZWQxYzEyYSIsInVzZXJuYW1lIjoiYXphZC1hc3Npc3RhbnQiLCJlbWFpbCI6bnVsbCwicm9sZSI6ImFkbWluIiwiZmlyc3ROYW1lIjoi2KLYstin2K8iLCJsYXN0TmFtZSI6Itiv2LPYqtuM2KfYsSIsImxhc3RMb2dvdXQiOiIyMDIyLTAyLTA1VDA3OjE4OjU5LjgyM1oiLCJsYXN0TG9naW4iOiIyMDIyLTAyLTA3VDA0OjUwOjAzLjA0OFoiLCJkeW5hbWljUm9sZSI6ImFkbWluIiwiaWF0IjoxNjQ0MjA5NDAzfQ.xJbnwEsK4z06cmZwCWGMParBYBGfAsvBO70K4qxEaYU"
	APARAT_URL := "rtmp://rtmp.cdn.asset.aparat.com:443/event/e4a035237879741166eb3db92b7044377?s=6587f3a42091cce0"

	log.Println("Opening Chrome")
	opts := append(chromedp.DefaultExecAllocatorOptions[:],
		chromedp.NoFirstRun,
		chromedp.NoDefaultBrowserCheck,
		// chromedp.DisableGPU,
		chromedp.IgnoreCertErrors,
		chromedp.Flag("autoplay-policy", "no-user-gesture-required"),
		chromedp.Flag("headless", false),
		chromedp.Flag("disable-gpu", true),
		chromedp.Flag("enable-automation", false),
		chromedp.Flag("disable-extensions", true),
		chromedp.Flag("kiosk", true),
		chromedp.Flag("window-size", "1504,846"),
		chromedp.Flag("ignore-certificate-errors", true),
		chromedp.Flag("disable-sync", true),
		chromedp.Flag("allow-file-access-from-files", true),
		chromedp.Flag("enable-usermedia-screen-capturing", true),
		chromedp.Flag("disable-gesture-requirement-for-media-playback", true),
		chromedp.Flag("use-fake-device-for-media-stream", true),
		chromedp.ExecPath("/usr/bin/google-chrome-stable"),
		chromedp.Flag("disable-infobars", true),
		chromedp.Flag("excludeSwitches", "enable-automation"),
		chromedp.Flag("disable-background-networking", true),
		chromedp.Flag("enable-features", "NetworkService,NetworkServiceInProcess"),
		chromedp.Flag("disable-background-timer-throttling", true),
		chromedp.Flag("disable-backgrounding-occluded-windows", true),
		chromedp.Flag("disable-breakpad", true),
		chromedp.Flag("disable-client-side-phishing-detection", true),
		chromedp.Flag("disable-default-apps", true),
		chromedp.Flag("disable-dev-shm-usage", true),
		chromedp.Flag("disable-extensions", true),
		chromedp.Flag("disable-features", "site-per-process,TranslateUI,BlinkGenPropertyTrees"),
		chromedp.Flag("disable-hang-monitor", true),
		chromedp.Flag("disable-ipc-flooding-protection", true),
		chromedp.Flag("disable-popup-blocking", true),
		chromedp.Flag("disable-prompt-on-repost", true),
		chromedp.Flag("disable-renderer-backgrounding", true),
		chromedp.Flag("disable-sync", true),
		chromedp.Flag("force-color-profile", "srgb"),
		chromedp.Flag("metrics-recording-only", true),
		chromedp.Flag("safebrowsing-disable-auto-update", true),
		chromedp.Flag("enable-automation", true),
		chromedp.Flag("password-store", "basic"),
		chromedp.Flag("use-mock-keychain", true),
	)

	allocCtx, cancel := chromedp.NewExecAllocator(context.Background(), opts...)
	defer cancel()

	ctx, cancel := chromedp.NewContext(allocCtx, chromedp.WithLogf(log.Printf))
	defer cancel()

	var res string
	go func() {
		if err := chromedp.Run(ctx,
			chromedp.Navigate(URL),
			browser.GrantPermissions([]browser.PermissionType{
				"audioCapture", "videoCapture",
			}),
			setcookies(URL, "class.kavano.org", &res, "auth", TOKEN),
			chromedp.Sleep(time.Hour),
		); err != nil {
			log.Fatal(err)
		}
	}()
	chromedp.Reload()
	log.Println("Waiting")
	time.Sleep(time.Second * 3)
	log.Println("Startinc FFMPGE command")
	chromedp.Reload()
	startFfmpegCommand(APARAT_URL)

}

func setcookies(host string, domain string, res *string, cookies ...string) chromedp.Tasks {
	if len(cookies)%2 != 0 {
		panic("length of cookies must be divisible by 2")
	}
	return chromedp.Tasks{
		chromedp.ActionFunc(func(ctx context.Context) error {
			// create cookie expiration
			expr := cdp.TimeSinceEpoch(time.Now().Add(180 * 24 * time.Hour))
			// add cookies to chrome
			for i := 0; i < len(cookies); i += 2 {
				err := network.SetCookie(cookies[i], cookies[i+1]).
					WithExpires(&expr).
					WithDomain(domain).
					WithHTTPOnly(false).
					Do(ctx)
				if err != nil {
					return err
				}
			}
			return nil
		}),
		chromedp.Navigate(host),
	}
}

func startFfmpegCommand(token string) {
	cmd := exec.Command(
		"ffmpeg",
		"-f",
		"pulse",
		"-ac",
		" 2",
		"-i",
		"default",
		"-f",
		"x11grab",
		"-video_size",
		"1504x846",
		"-framerate",
		"24",
		"-i",
		":44",
		"-codec:v",
		"libx264",
		"-pix_fmt",
		"yuv420p",
		"-profile",
		"high",
		"-preset",
		"veryfast",
		"-x264-params",
		"keyint=48:scenecut=0",
		"-b:v",
		"4000k",
		"-b:a",
		"4000k",
		"-bufsize",
		"3000k",
		"-maxrate",
		"6000k",
		"-minrate",
		"1500k",
		"-f",
		"flv",
		token,
	)

	// stderr, _ := cmd.StderrPipe()
	log.Println("Starting Command")
	cmd.Start()

	// scanner := bufio.NewScanner(stderr)
	// scanner.Split(bufio.ScanWords)
	// for scanner.Scan() {
	// 	m := scanner.Text()
	// 	fmt.Println(m)
	// }
	cmd.Wait()
}
