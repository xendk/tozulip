# tozulip

Create a Zulip bot and send messages to streams from the command line.
Handy for deployment messages, or other things you might want to send
from the CLI.

# Usage

```
  tozulip [message to send] [flags]

Flags:
  -k, --apikey string   API key of bot
  -c, --config string   config file (default is $HOME/.tozulip.yaml)
  -h, --help            help for tozulip
  -H, --host string     hostname of Zulip server
  -m, --mail string     bot email
  -s, --stream string   stream to send message to
  -t, --topic string    topic of message
```

All flags (apart from `--config`) can be specified in
`$HOME/.tozulip.yaml` or by as `TOZULIP_<name in uppercase>`.

# License

[MIT](https://xendk.mit-license.org/)
