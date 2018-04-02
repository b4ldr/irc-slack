# IRC-to-Slack gateway

`irc-slack` is an IRC-to-Slack gateway. Run it locally, and it will spawn an IRC
server that will let you use your Slack teams via IRC.

Slack is ending support for IRC and XMPP gateway on the 15th of April 2018. So
what's left to do for people like me, who want to still be able to log in via
IRC? Either you use [wee-slack](https://github.com/wee-slack/wee-slack) (but I
don't use WeeChat), or you implement your own stuff.

The code quality is currently at the `works-for-me` level.

## How to use it

```
go build
./irc-slack # by default on port 6666
irssi
  /network add SlackYourTeamName
  /server add -auto SlackYourTeamName localhost 6666 xoxp-<your-slack-token>
```

Get you Slack legacy token at https://api.slack.com/custom-integrations/legacy-tokens .


You can also specify the port to listen on, and the server name, e.g.
your-team-name.slack.com.

```
$ ./irc-slack -h
Usage of ./irc-slack:
  -p int
        Local port to listen on (default 6666)
  -s string
        IRC server name (i.e. the host name to send to clients)
```


## TODO

A lot of things. Want to help? Grep "TODO", "FIXME" and "XXX" in the code and send me a PR :)

This currently "works for me", but I published it in the hope that someone would use it so we can find and fix bugs.

## BUGS

Plenty of them. I wrote this project while on a plane (like many other projects of mine) so this is hack-level quality - no proper design, no RFC compliance, no testing. I just fired up an IRC client until I could reasonably chat on a few Slack teams. Please report all the bugs you find on the Github issue tracker, or privately to me.

## How do I contact you?

Find my contacts on https://insomniac.slackware.it .
