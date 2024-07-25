package cmd

import (
	"fmt"
	"log/slog"
	"os"
	"strings"

	"github.com/hxhac/docs-alfred/pkg/gh"

	"github.com/spf13/cobra"
)

// const ghTpl = `
// {{- each .}}
// ## {{.Type}}
// {{- each .Qs}}
// * {{.}}
// {{- end}}
// {{- each .Repos}}
// ### {{.URL}}
// {{- each .Qs}}
// * {{.}}
// {{- end}}
// {{- end}}
// `

// const ghTpl = `
// {{- range . -}}
// ## {{.Type}}
//
// {{- if .Qs -}}
// - {{ range .Qs }}{{ . }}
// {{- end -}}
// {{- end }}
//
// {{- range .Repos -}}
// ### {{.URL}}
//
// {{- range .Qs -}}
// - {{ . }}
// {{- end -}}
// {{- end -}}
// {{- end -}}
// `

// const ghTpl = `
// {{- range . -}}
// {{"\n"}}
// ## {{.Type}}
// {{"\n"}}
//
// {{- if .Qs -}}
// {{- range .Qs}}
// - {{.}}
// {{- end}}
// {{- end }}
//
// {{- range .Repos -}}
// {{"\n"}}
// ### {{.URL}}
// {{"\n"}}
//
// {{- range .Pix}}
// {{.}}
// {{- end}}
//
// {{- range .Qs}}
// - {{.}}
// {{- end}}
//
// {{- if .Qq -}}
// {{- range .Qq -}}
// {{"\n"}}
// #### {{.Topic}}
// {{- range .Qs}}
// - {{.}}
// {{- end}}
// {{- end}}
// {{- end }}
//
// {{- end -}}
// {{- end -}}
// `

// ghCmd represents the gh command
var ghCmd = &cobra.Command{
	Use:   "gh",
	Short: "A brief description of your command",
	Run: func(cmd *cobra.Command, args []string) {
		f, err := os.ReadFile(cfgFile)
		if err != nil {
			return
		}

		dfo := gh.NewConfigRepos(f)
		df := dfo.FilterReposMD()

		// 清理掉 Qs == nil 的 Type
		dfr := FilterRepos(df)

		var res strings.Builder

		for _, d := range dfr {

			res.WriteString(fmt.Sprintf("## %s \n", d.Type))
			if d.Qs != nil {
				res.WriteString(addMarkdownQsFormat(d.Qs))
			}

			for _, repo := range d.Repos {
				repoName, f := strings.CutPrefix(repo.URL, gh.GhURL)
				if !f {
					repoName = ""
				}

				res.WriteString(fmt.Sprintf("\n\n### [%s](%s)\n\n", repoName, repo.URL))
				if repo.Qs != nil {
					res.WriteString(addMarkdownQsFormat(repo.Qs))
				}
				if repo.Qq != nil {
					for _, s := range repo.Qq {
						if s.Qs != nil {
							res.WriteString(fmt.Sprintf("\n\n#### %s \n\n", s.Topic))
							res.WriteString(addMarkdownQsFormat(s.Qs))
						}
					}
				}
			}
		}

		tf, _ := strings.CutSuffix(cfgFile, ".yml")
		targetFile := fmt.Sprintf("%s.md", tf)
		err = os.WriteFile(fmt.Sprintf("%s.md", targetFile), []byte(res.String()), os.ModePerm)
		if err != nil {
			return
		}

		slog.Info("Markdown output has been written to", slog.String("File", targetFile))
	},
}

func addMarkdownQsFormat(qs gh.Qs) string {
	var builder strings.Builder
	// builder.WriteString("<dl>")
	for _, q := range qs {
		if q.X == "" {
			if q.U != "" {
				builder.WriteString(fmt.Sprintf("- [%s](%s)\n", q.Q, q.U))
			} else {
				builder.WriteString(fmt.Sprintf("- %s\n", q.Q))
			}
		} else {
			if q.U != "" {
				builder.WriteString(fmt.Sprintf("\n\n<details>\n<summary>[%s](%s)</summary>\n\n%s\n\n</details>\n\n", q.Q, q.U, q.X))
			} else {
				builder.WriteString(fmt.Sprintf("\n\n<details>\n<summary>%s</summary>\n\n%s\n\n</details>\n\n", q.Q, q.X))
			}
		}
	}
	// builder.WriteString("</dl>")

	return builder.String()
}

// FilterRepos 过滤掉Repo中Qs为nil的ConfigRepos
func FilterRepos(configRepos gh.ConfigRepos) (filteredRepos gh.ConfigRepos) {
	for _, repoGroup := range configRepos {
		// 过滤掉qs为nil的Repository
		filteredGroup := gh.ConfigRepo{
			Type:  repoGroup.Type,
			Repos: make([]gh.Repository, 0),
			Qs:    make(gh.Qs, 0),
		}
		filteredGroup.Type = repoGroup.Type
		for _, repo := range repoGroup.Repos {
			if repo.Qs != nil { // 假设Qs的零值是Qs{}
				filteredGroup.Repos = append(filteredGroup.Repos, repo)
			}
		}
		if len(filteredGroup.Repos) > 0 {
			// 只有当过滤后的Repositories不为空时，才添加到结果中
			filteredRepos = append(filteredRepos, filteredGroup)
		}
	}
	return filteredRepos
}

func init() {
	rootCmd.AddCommand(ghCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// ghCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// ghCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
