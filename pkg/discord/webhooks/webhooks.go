package webhooks

import (
  "context"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/bwmarrin/discordgo"
	"github.com/pkg/errors"
	"github.com/thompsonja/discord_bots_lib/pkg/discord/channellogger"
	"github.com/thompsonja/discord_bots_lib/pkg/discord/interactions"
	"github.com/thompsonja/discord_bots_lib/pkg/discord/session"
	"github.com/thompsonja/discord_bots_lib/pkg/logger"
)

type WebhookFunc func(*Client, *discordgo.Interaction, *http.Request) error

type Client struct {
	commands          []*discordgo.ApplicationCommand
	port              int
	epk               []byte
	fns               map[string]WebhookFunc
	secretKey         string
	projectID         string
	session           *discordgo.Session
	appID             string
	channelLogger     *channellogger.ChannelLogger
	interactionLogger *interactions.InteractionLogger
	logger            logger.Logger
	pool              chan func()
}

type ClientConfig struct {
	Commands  []*discordgo.ApplicationCommand
	Port      int
	Epk       []byte
	Fns       map[string]WebhookFunc
	SecretKey string
	ProjectID string
	AppID     string
	PoolSize  int
	Logger    logger.Logger
}

func NewClient(cfg ClientConfig) (*Client, error) {
	if len(cfg.Commands) == 0 {
		return nil, fmt.Errorf("empty command list passed to NewClient")
	}
	if len(cfg.Fns) == 0 {
		return nil, fmt.Errorf("empty function map passed to NewClient")
	}
	if cfg.SecretKey == "" {
		return nil, fmt.Errorf("config SecretKey must be set")
	}
	if cfg.ProjectID == "" {
		return nil, fmt.Errorf("config ProjectID must be set")
	}
	if cfg.AppID == "" {
		return nil, fmt.Errorf("config AppID must be set")
	}
	l := cfg.Logger
	if l == nil {
		l = &logger.StandardLogger{}
	}
	poolSize := 100
	if cfg.PoolSize > 0 {
		poolSize = cfg.PoolSize
	}
	port := 8080
	if cfg.Port != 0 {
		port = cfg.Port
	}
	pool := make(chan func(), poolSize)
	for i := 0; i < poolSize; i++ {
		go func() {
			for f := range pool {
				f()
			}
		}()
	}
	return &Client{
		commands:          cfg.Commands,
		port:              port,
		epk:               cfg.Epk,
		fns:               cfg.Fns,
		secretKey:         cfg.SecretKey,
		projectID:         cfg.ProjectID,
		appID:             cfg.AppID,
		interactionLogger: &interactions.InteractionLogger{},
		logger:            l,
		pool:              pool,
	}, nil
}

func (c *Client) handlePing(w http.ResponseWriter) error {
	c.logger.Info("Got a ping request")
	if _, err := w.Write([]byte(`{"type":1}`)); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return fmt.Errorf("w.Write: %v", err)
	}
	return nil
}

func (c *Client) sendResponse(w http.ResponseWriter, r *discordgo.InteractionResponse) error {
	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(r); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return fmt.Errorf("json.NewEncoder.Encode: %v", err)
	}
	return nil
}

func (c *Client) SendDeferredResponse(w http.ResponseWriter) error {
	return c.sendResponse(w, &discordgo.InteractionResponse{
		Type: discordgo.InteractionResponseDeferredChannelMessageWithSource,
		Data: &discordgo.InteractionResponseData{
			Content: "Command received, please wait...",
		},
	})
}

func (c *Client) SendFilesMessage(ctx context.Context, channel string, files []*discordgo.File) error {
	if err := c.updateSession(ctx); err != nil {
		return fmt.Errorf("c.updateSession: %v", err)
	}

	for _, f := range files {
		c.channelLogger.SendFile(channel, f.Name, f.Reader)
	}
	return nil
}

func (c *Client) SendFilesResponse(ctx context.Context, i *discordgo.Interaction, files []*discordgo.File) error {
	if err := c.updateSession(ctx); err != nil {
		return fmt.Errorf("c.updateSession: %v", err)
	}

	c.interactionLogger.SendEditedInteractionFiles(c.session, i, files)
	return nil
}

func (c *Client) SendStringMessage(ctx context.Context, channel, msg string) error {
	if err := c.updateSession(ctx); err != nil {
		return fmt.Errorf("c.updateSession: %v", err)
	}

	if channel == "" {
		return fmt.Errorf("channel is empty")
	}

	c.channelLogger.SendMessage(channel, msg)
	return nil
}

func (c *Client) SendStringResponse(ctx context.Context, i *discordgo.Interaction, msg string) error {
	if err := c.updateSession(ctx); err != nil {
		return fmt.Errorf("c.updateSession: %v", err)
	}

	if i == nil {
		return fmt.Errorf("interaction is nil")
	}

	c.interactionLogger.SendEditedInteractionMessage(c.session, i, msg)
	return nil
}

func (c *Client) updateSession(ctx context.Context) error {
	if c.session != nil {
		return nil
	}
	s, err := session.GetSession(ctx, c.secretKey, c.projectID)
	if err != nil {
		return errors.Wrap(err, "session.GetSession")
	}

	c.channelLogger = channellogger.New(s)
	c.session = s

	return nil
}

func (c *Client) GetSession(ctx context.Context) (*discordgo.Session, error) {
	if err := c.updateSession(ctx); err != nil {
		return nil, fmt.Errorf("c.updateSession: %v", err)
	}
	return c.session, nil
}

func (c *Client) DeleteCommands(ctx context.Context) error {
	if err := c.updateSession(ctx); err != nil {
		return fmt.Errorf("c.updateSession: %v", err)
	}
	cmds, err := c.session.ApplicationCommands(c.appID, "")
	if err != nil {
		return fmt.Errorf("c.session.ApplicationCommands: %v", err)
	}
	for _, cmd := range cmds {
		c.logger.Infof("Deleting command %s (name %s) for app ID %s\n", cmd.ID, cmd.Name, c.appID)
		if err := c.session.ApplicationCommandDelete(c.appID, "", cmd.ID); err != nil {
			return fmt.Errorf("c.session.ApplicationCommandDelete(%s): %v", cmd.Name, err)
		}
	}
	return nil
}

func (c *Client) UpdateCommands(ctx context.Context) error {
	if err := c.updateSession(ctx); err != nil {
		return fmt.Errorf("c.updateSession: %v", err)
	}
	for _, v := range c.commands {
		c.logger.Infof("Creating command %s (name %s) for app ID %s\n", v.ID, v.Name, c.appID)
		_, err := c.session.ApplicationCommandCreate(c.appID, "", v)
		if err != nil {
			return fmt.Errorf("cannot create '%v' command: %v", v.Name, err)
		}
	}
	return nil
}

func (c *Client) handlePost(w http.ResponseWriter, i *discordgo.Interaction, r *http.Request) error {
	if i.Data.Type() != discordgo.InteractionApplicationCommand {
		return fmt.Errorf("invalid interaction type %v", i.Data.Type())
	}

	cmdData, ok := i.Data.(discordgo.ApplicationCommandInteractionData)
	if !ok {
		return fmt.Errorf("couldn't assert application command data")
	}

	fn, ok := c.fns[cmdData.Name]
	if !ok {
		return fmt.Errorf("unsupported command: %v", cmdData.Name)
	}
	return fn(c, i, r)
}

func (c *Client) ListenAndServe() error {
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		if !discordgo.VerifyInteraction(r, c.epk) {
			w.WriteHeader(http.StatusUnauthorized)
			fmt.Fprint(w, "Unauthorized")
			return
		}
		c.logger.Info("Got an authorized request")
		if r.Method != "POST" {
			w.WriteHeader(http.StatusBadRequest)
			fmt.Fprintf(w, "Unsupported method %v", r.Method)
			return
		}

		var i discordgo.Interaction
		if err := json.NewDecoder(r.Body).Decode(&i); err != nil {
			c.logger.Errorf("json.NewDecoder: %v\n", err)
			if err := c.SendStringResponse(r.Context(), &i, fmt.Sprintf("Error decoding command: %v", err)); err != nil {
				c.logger.Errorf("c.SendStringResponse: %v", err)
			}
			return
		}

		if i.Type == discordgo.InteractionPing {
			if err := c.handlePing(w); err != nil {
				c.logger.Errorf("c.handlePing: %v", err)
			}
			return
		}

		if err := c.SendDeferredResponse(w); err != nil {
			c.logger.Errorf("c.SendDeferredResponse: %v", err)
		}
		go func() {
			c.pool <- func() {
				if err := c.handlePost(w, &i, r); err != nil {
					c.logger.Errorf("handlePost: %v\n", err)
				}
			}
		}()
	})

	return http.ListenAndServe(fmt.Sprintf(":%d", c.port), nil)
}

func (c *Client) Run(ctx context.Context, destroyCommands, updateCommands bool) error {
	if destroyCommands {
		c.logger.Info("Destroying commands...")
		if err := c.DeleteCommands(ctx); err != nil {
			return fmt.Errorf("c.DestroyCommands: %v", err)
		}
		c.logger.Info("Done")
	}

	if updateCommands {
		c.logger.Info("Adding commands...")
		if err := c.UpdateCommands(ctx); err != nil {
			return fmt.Errorf("c.UpdateCommands: %v", err)
		}
		c.logger.Info("Done")
	}

	c.logger.Infof("Starting server at port %d\n", c.port)
	if err := c.ListenAndServe(); err != nil {
		return fmt.Errorf("c.ListenAndServe: %v", err)
	}
	return nil
}
