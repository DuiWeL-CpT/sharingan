package cli

import (
	"bufio"
	"fmt"
	"log"
	"os"

	"github.com/leobeosab/sharingan/internal/helpers"
	"github.com/leobeosab/sharingan/internal/models"
	"github.com/leobeosab/sharingan/pkg/dns"
	"github.com/leobeosab/sharingan/pkg/storage"
)

func RunDNSRecon(settings *models.ScanSettings) {
	if settings.Target == "" {
		log.Fatal("Target needs to be defined")
	}
	if settings.DNSWordlist == "" {
		log.Fatalf("No program found - DNS Wordlist needs to be defined")
	}

	_, p := storage.RetrieveOrCreateProgram(settings.Store, settings.Target)
	subs := dns.DNSBruteForce(settings.RootDomain, settings.DNSWordlist, settings.Threads)

	// Pesky progressbars not ending their lines
	fmt.Printf("\n")
	if settings.Rescan {
		AddSubsToProgram(&p, &subs)
	} else {
		ReplaceSubsInProgram(&p, &subs)
	}

	storage.UpdateOrCreateProgram(settings.Store, &p)
}

func AddSubsFromInput(settings *models.ScanSettings) {

	info, err := os.Stdin.Stat()
	if err != nil {
		panic(err)
	}

	e, p := storage.RetrieveOrCreateProgram(settings.Store, settings.Target)
	if !e {
		fmt.Println(settings.Target + " not found in store, creating new entry")
		p.Hosts = make(map[string]models.Host)
	}

	if info.Mode()&os.ModeNamedPipe == 0 {
		log.Println("DNS addsubs is intended to work with pipies.")
		log.Println("Usage: cat subs | sharingancli --target program dns addsubs")
		return
	}

	reader := bufio.NewScanner(os.Stdin)
	var subdomains []string

	for reader.Scan() {
		input := reader.Text()
		subdomains = append(subdomains, input)
	}
	subdomains = helpers.RemoveDuplicatesInSlice(subdomains)

	cliOut := fmt.Sprintf("Added %v subdomains to %s \n", len(subdomains), settings.Target)

	if settings.ReplaceSubs {
		ReplaceSubsInProgram(&p, &subdomains)
		cliOut = fmt.Sprintf("Replacing subdomains for %s \n", settings.Target)
	} else {
		AddSubsToProgram(&p, &subdomains)
	}

	storage.UpdateOrCreateProgram(settings.Store, &p)

	log.Printf(cliOut)
}

func ReplaceSubsInProgram(p *models.Program, subs *[]string) {
	p.Hosts = make(map[string]models.Host)
	for _, s := range *subs {
		fmt.Println(s)
		p.Hosts[s] = models.Host{
			Subdomain: s,
		}
	}
}

func AddSubsToProgram(p *models.Program, subs *[]string) {
	for _, s := range *subs {
		if _, ok := p.Hosts[s]; ok {
			continue
		}

		fmt.Println(s)
		p.Hosts[s] = models.Host{
			Subdomain: s,
		}
	}
}
