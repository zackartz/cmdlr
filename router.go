package cmdlr

import (
	"fmt"
	"github.com/bwmarrin/discordgo"
	"gorm.io/gorm"
	"sort"
	"strings"
)

type Router struct {
	Prefixes         []string
	IgnorePrefixCase bool
	BotsAllowed      bool
	Commands         []*Command
	Middlewares      []Middleware
	PingHandler      ExecutionHandler
	jobs             chan job
	MaxThreads       int
	Database         *gorm.DB
	Storage          map[string]*ObjectsMap
}

type job struct {
	Ctx     *Ctx
	Command *Command
}

func (r *Router) newJob(ctx *Ctx, cmd *Command) {
	r.jobs <- job{
		Ctx:     ctx,
		Command: cmd,
	}
}

func Create(router *Router) *Router {
	router.Storage = map[string]*ObjectsMap{}
	return router
}

func (r *Router) RegisterCmd(command *Command) {
	r.Commands = append(r.Commands, command)
}

func (r *Router) GetCmd(name string) *Command {
	sort.Slice(r.Commands, func(i, j int) bool {
		return len(r.Commands[i].Name) > len(r.Commands[j].Name)
	})

	for _, command := range r.Commands {
		toCheck := make([]string, len(command.Aliases)+1)
		toCheck = append(toCheck, command.Name)
		toCheck = append(toCheck, command.Aliases...)
		sort.Slice(toCheck, func(i, j int) bool {
			return len(toCheck[i]) > len(toCheck[j])
		})

		if StringArrayContains(toCheck, name, command.IgnoreCase) {
			return command
		}
	}
	return nil
}

func (r *Router) RegisterMiddleware(middleware Middleware) {
	r.Middlewares = append(r.Middlewares, middleware)
}

func (r *Router) InitializeStorage(name string) {
	r.Storage[name] = NewObjectsMap()
}

func (r *Router) Initialize(session *discordgo.Session) {
	session.AddHandler(r.Handler())

	r.jobs = make(chan job, 1000)

	for i := 0; i < r.MaxThreads; i++ {
		go worker(r.jobs)
	}

	fmt.Printf("Starting bot with %d threads...", r.MaxThreads)
}

func (r *Router) Handler() func(session *discordgo.Session, message *discordgo.MessageCreate) {
	return func(session *discordgo.Session, event *discordgo.MessageCreate) {
		message := event.Message
		content := event.Content

		if message.Author.Bot && !r.BotsAllowed {
			return
		}

		if (content == "<@!"+session.State.User.ID+">" || content == "<@"+session.State.User.ID+">") && r.PingHandler != nil {
			r.PingHandler(&Ctx{
				Session:  session,
				Event:    event,
				Database: r.Database,
				Args:     ParseArguments(""),
				Router:   r,
			})
			return
		}

		hasPrefix, content := StringHasPrefix(content, r.Prefixes, r.IgnorePrefixCase)
		if !hasPrefix {
			return
		}

		content = strings.Trim(content, " ")
		if content == "" {
			return
		}

		for _, m := range r.Middlewares {
			m.Trigger(Ctx{
				Session:  session,
				Event:    event,
				Database: r.Database,
				Router:   r,
			})
		}

		for _, cmd := range r.Commands {
			toCheck := BuildCheckPrefixes(cmd)

			isCommand, content := StringHasPrefix(content, toCheck, cmd.IgnoreCase)
			if !isCommand {
				continue
			}

			isValid, content := StringHasPrefix(content, []string{" ", "\n"}, false)
			if content == "" || isValid {
				ctx := &Ctx{
					Session:  session,
					Event:    event,
					Database: r.Database,
					Args:     ParseArguments(content),
					Router:   r,
					Command:  cmd,
				}

				r.newJob(ctx, cmd)
			}
		}
	}
}

func worker(jobs <-chan job) {
	for j := range jobs {
		j.Command.Trigger(j.Ctx)
	}
}
