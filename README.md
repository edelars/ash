

# Alternative shell 
Written in golang

#### Example configuration file
```yaml
keybindings:
  - key: 13
    action: ':Execute'
  - key: 9
    action: ':Autocomplete'
  - key: 55
    action: ':Exit'
  - key: 127
    action: ':RemoveLeftSymbol'
  - key: 167
    action: ':Close'

aliases: 
  - short: lg
    full: lazygit
prompt: '[{"value": "ash> ", "color":"#8ec07c", "bold": true}]' #See section "Prompt configuration"
envs:
  - >-
    $PATH=/opt/homebrew/bin:/opt/homebrew/sbin:/usr/local/bin:/bin:/usr/sbin:/sbin:/var/run/com.apple.security.cryptexd/codex.system/bootstrap/usr/local/bin:/opt/homebrew/sbin:  

colors:
  defaultText: "#1d2021"
  defaultBackground: "#1d2021"

autocomplete:
  showFileInformation: true
```
##### Keybindings
##### Aliases
##### OS Envs
##### Prompt configuration
Its a simple json array:
```json
 [{"value": "exe", "color": "#8ec07c", "bold": true,"underline": true }]
```
Where:
**"value"** - text, $variable (see variables) or system binary exec "%(git log --pretty=format:"%s"  | head -n 1)"
**"color"** - color, string 
**"bold"** - bold font, bool
**"underline"**- underline font, bool

Example:
```json
[{"value": "ash> ", "color":"#8ec07c", "bold": true}]
```


#### Roadmap
```
internal eventbus
plugins system
more commands (>>,>, etc)
```
