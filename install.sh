#!/bin/bash

sudo << EOF > /etc/netctl/hooks/update-iface-state.sh
#!/bin/bash

ExecUpPost="pkill -RTMIN+1 go-i3status"
ExecDownPre="pkill -RTMIN+1 go-i3status"
EOF

sudo chmod +x /etc/netctl/hooks/update-iface-state.sh


sudo << EOF > /etc/acpi/update-battery-state.sh
#!/bin/bash

pkill -RTMIN+2 go-i3status
EOF

sudo chmod +x /etc/acpi/update-battery-state.sh
