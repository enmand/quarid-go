package bot

func (b *Bot) Start() error {
	return b.adapter.Start()
}

func (b *Bot) Stop() error {
	return b.adapter.Stop()
}
