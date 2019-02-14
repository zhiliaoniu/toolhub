# ~/.bashrc: executed by bash(1) for non-login shells.
# see /usr/share/doc/bash/examples/startup-files (in the package bash-doc)
# for examples

# If not running interactively, don't do anything
[ -z "$PS1" ] && return

# don't put duplicate lines in the history. See bash(1) for more options
# don't overwrite GNU Midnight Commander's setting of `ignorespace'.
export HISTCONTROL=$HISTCONTROL${HISTCONTROL+,}ignoredups
# ... or force ignoredups and ignorespace
export HISTCONTROL=ignoreboth

export GOROOT=/usr/lib/go
export GOPATH=~
export GOBIN=~/bin
export CSCOPE_DB=/home/yangshengzhi1/dfs_fastdfs/cscope.out

PATH="$HOME/bin:$HOME/.local/bin:$PATH"
export PATH=~/bin:$PATH

export BLADE_AUTO_UPGRADE=no

# append to the history file, don't overwrite it
shopt -s histappend

# for setting history length see HISTSIZE and HISTFILESIZE in bash(1)

# check the window size after each command and, if necessary,
# update the values of LINES and COLUMNS.
shopt -s checkwinsize

# make less more friendly for non-text input files, see lesspipe(1)
[ -x /usr/bin/lesspipe ] && eval "$(SHELL=/bin/sh lesspipe)"

# set variable identifying the chroot you work in (used in the prompt below)
if [ -z "$debian_chroot" ] && [ -r /etc/debian_chroot ]; then
    debian_chroot=$(cat /etc/debian_chroot)
fi

# set a fancy prompt (non-color, unless we know we "want" color)
case "$TERM" in
    xterm-color) color_prompt=yes;;
esac

# uncomment for a colored prompt, if the terminal has the capability; turned
# off by default to not distract the user: the focus in a terminal window
# should be on the output of commands, not on the prompt
#force_color_prompt=yes

if [ -n "$force_color_prompt" ]; then
    if [ -x /usr/bin/tput ] && tput setaf 1 >&/dev/null; then
        # We have color support; assume it's compliant with Ecma-48
        # (ISO/IEC-6429). (Lack of such support is extremely rare, and such
        # a case would tend to support setf rather than setaf.)
        color_prompt=yes
    else
        color_prompt=
    fi
fi

if [ "$color_prompt" = yes ]; then
    PS1='${debian_chroot:+($debian_chroot)}\[\033[01;32m\]\u@\h\[\033[00m\]:\[\033[01;34m\]\w\[\033[00m\]\$ '
else
    PS1='${debian_chroot:+($debian_chroot)}\u@\h:\w\$ '
fi
unset color_prompt force_color_prompt

# If this is an xterm set the title to user@host:dir
case "$TERM" in
xterm*|rxvt*)
    PS1="\[\e]0;${debian_chroot:+($debian_chroot)}\u@\h: \w\a\]$PS1"
    ;;
*)
    ;;
esac

# Alias definitions.
# You may want to put all your additions into a separate file like
# ~/.bash_aliases, instead of adding them here directly.
# See /usr/share/doc/bash-doc/examples in the bash-doc package.

#if [ -f ~/.bash_aliases ]; then
#    . ~/.bash_aliases
#fi

# enable color support of ls and also add handy aliases
if [ -x /usr/bin/dircolors ]; then
    eval "`dircolors -b`"
    alias ls='ls --color=auto'
    #alias dir='dir --color=auto'
    #alias vdir='vdir --color=auto'

    #alias grep='grep --color=auto'
    #alias fgrep='fgrep --color=auto'
    #alias egrep='egrep --color=auto'
fi

# some more ls aliases
alias la='ls -A'
alias ls='ls --color=auto'
alias ll='ls -lh'
alias grep='grep --color=auto'
#alias go='ssh -p 32200'
alias go='~/bin/go'
alias rm='rm -i'
alias cp='cp -i'
alias mv='mv -i'
alias vi='/usr/bin/vim'
alias vim='/usr/bin/vim -O'
alias du1='du -h --max-depth=1 ./'
# by yangshengzhi1
alias cdf='cd /home/yangshengzhi1/dfs_fastdfs'
alias cdm='cd /home/yangshengzhi1/dfs_fastdfs/client/meta_cache'
alias cdd='cd /home/yangshengzhi1/dfs_fdfsproxy'
alias cdfp='cd /home/yangshengzhi1/dfs_fcproxy'
alias cdp='cd /home/yangshengzhi1/fdfsbrpc_server/'
alias bb='blade build'

alias mcscope='find `pwd` -name "*.[ch]" -o -name "*.cpp" > cscope.files; cscope -bR -i cscope.files'
alias mtag='ctags -R ./; mcscope'

alias gitpush="git push origin HEAD:refs/for/master"
alias help_fp_server="~/dfs_fmproxy/src/server/help_fp_server"
alias mf="make -f ~/dfs_fastdfs/Makefile"
alias mfp="make -f ~/dfs_fcproxy/Makefile"
alias mp="make -f ~/fdfsbrpc_server/Makefile"

# enable programmable completion features (you don't need to enable
# this, if it's already enabled in /etc/bash.bashrc and /etc/profile
# sources /etc/bash.bashrc).
if [ -f /etc/bash_completion ]; then
    . /etc/bash_completion
fi

test -s ~/bin/bladefunctions && . ~/bin/bladefunctions || true
