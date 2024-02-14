
# Alternative shell 
Written in golang

### Default configuration file
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
prompt: 'ASH>'
envs:
  - >-
    $PATH=/opt/homebrew/bin:/opt/homebrew/sbin:/usr/local/bin:/bin:/usr/sbin:/sbin:/var/run/com.apple.security.cryptexd/codex.system/bootstrap/usr/local/bin:/opt/homebrew/sbin:  

colors:
  defaultText: 0
  defaultBackground: 0

autocomplete:
  showFileInformation: true
```

### Prompt configuration
Its a simple json array:
```
 [{"value": "exe", "color": 0, "bold": true,"under": true }]
```
where:
```
"value": "exe", // text, $variable (see variables) or system binary exec "%(git log --pretty=format:"%s"  | head -n 1)"
"color": 0, //color 
"bold": true, // bold font
"under": true, //underline font
```

###Roadmap
```
internal eventbus
plugins system
>> > commands
prompt custom
```
