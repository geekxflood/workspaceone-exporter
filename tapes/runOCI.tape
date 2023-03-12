Output assets/runOCI.gif

Set Theme "Dracula"
Set FontSize 15
Set Padding 20
Set Height 210

Type "docker build -t workspaceone-prometheus-exporter ."
Enter
Sleep 15s

Type "docker run -d -p 9740:9740 workspaceone-prometheus-exporter"
Enter
Sleep 3s
