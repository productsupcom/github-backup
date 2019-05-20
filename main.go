package main

import (
	"context"
	"flag"
	"fmt"
	"log"
	"os"

	"github.com/google/go-github/github"
	"golang.org/x/oauth2"
	"gopkg.in/src-d/go-git.v4"
)

var (
	org string
	token string
	path string
)

func main() {
	flag.StringVar(&org, "org", "productsupcom", "GitHub org name")
	flag.StringVar(&token, "token", "", "GitHub token")
	flag.StringVar(&path, "path", "/tmp/backup", "Git local clone path")

	flag.Parse()
	ctx := context.Background()
	ts := oauth2.StaticTokenSource(
		&oauth2.Token{AccessToken: token},
	)
	tc := oauth2.NewClient(ctx, ts)

	client := github.NewClient(tc)

	opt := &github.RepositoryListByOrgOptions{
		ListOptions: github.ListOptions{PerPage: 100},
	}

	var allRepos []*github.Repository

	fmt.Println("Querying the Github API")
	for {
		repos, resp, err := client.Repositories.ListByOrg(ctx, org, opt)
		if err != nil {
			log.Fatalln(err)
		}
		allRepos = append(allRepos, repos...)
		if resp.NextPage == 0 {
			break
		}
		opt.Page = resp.NextPage
	}
	for _, repo := range allRepos {
		fmt.Println(*repo.SSHURL)
		if *repo.Archived == true {
			fmt.Printf("Excluding archived repository %s from backup\n", *repo.Name)
		} else {
			_, err := git.PlainClone(path+"/"+*repo.Name, false, &git.CloneOptions{
				URL:      *repo.SSHURL,
				Progress: os.Stdout,
			})
			if err != nil {
				fmt.Println(err)
			}
		}
	}
}
