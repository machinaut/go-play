# -*- mode: shell-script; sh-basic-offset: 8; indent-tabs-mode: t -*-
# ex: ts=8 sw=8 noet filetype=sh
#
# go completion by Alex Ray <ajray@ncsu.edu>
#
# put this in /etc/bash_completion.d/ or a folder for bash completion to 
# find it

complete -f -X '!*.8' 8l
complete -f -X '!*.6' 6l
complete -f -X '!*.5' 5l
complete -f -X '!*.go' 8g 6g 5g
