Dev Environment
------

Quixotic attempt at a minimal development environment.
Updated as new simplifications or better tools are discovered.

* windows
  * [git for windows](https://git-scm.com/downloads/win) for `git`, `vim`, `bash`
  * windows Terminal with bash shell.
  * [renderdoc](https://renderdoc.org)
  * [glslc](https://github.com/google/shaderc) for compiling shaders.
  * prefer syscall over c-go.
* macos (+ios)
  * xcode developer tools for `git`, `vim`, `zsh`, `clang`
  * macos Terminal.
  * c-go required (clang) until [purego](https://github.com/ebitengine/purego) or equivalent is ready.
* all
  * [hack](https://github.com/source-foundry/Hack) font
  * [ripgrep](https://github.com/BurntSushi/ripgrep) search
  * [tre](https://github.com/dduan/tre) list
  * [tokei](https://github.com/XAMPPRocky/tokei) code count
* essential
  * [go](https://go.dev/doc/install)
  * [vulkan](https://vulkan.lunarg.com/sdk/home)
  * [openal](https://openal.org) + [openal-soft](https://openal-soft.org/openal-binaries)

Add, as needed, the essential open source game asset creation tools.
  * [blender](https://www.blender.org) 3D
  * [krita](https://krita.org) 2D
  * [audacity](https://www.audacityteam.org) audio

Capture the, ideally small, shell and vim config scripts.
Update as new simplifications are discovered.

Config
------

`~/.vimrc`
```vim
"prevent backup files
set nobackup nowritebackup

" want some nice colors
syntax on
set background=dark
colorscheme solarized

" turn on nice indenting with 4 spaces per tab.
set tabstop=4 shiftwidth=4 expandtab
set laststatus=2
" turn on the status line so line numbers are there too.
set statusline=%t\ %y\ format:\ %{&ff};\ [%c,%l]
" allow lots of tabs.
set tabpagemax=100

" vim-tabs
set tabpagemax=50
nnoremap <S-l> :tabnext<CR>
nnoremap <S-h> :tabprevious<CR>
nnoremap tt :tabedit<Space>
nnoremap tn :tabm<Space>

" save with Control-X
nnoremap <C-x> :w<CR>

" clear trailing spaces with \s
nnoremap \s ::%s/\s\+$//e<CR>

" highlight trailing whitespace
highlight ExtraWhitespace ctermbg=red guibg=red
match ExtraWhitespace /\s\+$/

" glsl syntax highlighting
" from https://www.vim.org/scripts/script.php?script_id=1002
augroup glsl_syntax
	autocmd!
	autocmd BufNewFile,BufRead *.frag,*.vert,*.glsl setf glsl
augroup END

" gofmt based on https://stackoverflow.com/questions/72135274/run-gofmt-on-vim-without-plugin
function! GoFmt()
	cexpr system('gofmt -e -w ' . expand('%'))
	edit!
endfunction
command! GoFmt call GoFmt()
augroup go_autocmd
	autocmd BufWritePost *.go GoFmt
augroup END

" go compile
autocmd Filetype go set makeprg=go\ build
autocmd QuickFixCmdPost [^l]* cwindow
nnoremap <F5> :silent make<CR><C-L><CR>

" use existing tabs or switch to a new tab
set switchbuf+=usetab,newtab
```

`~/.vim`
```
.vim
├── colors
│   └── solarized.vim
└── syntax
    └── glsl.vim
```

`~/.bashrc`
```bash
# needed for bash shared by all users.
cd $HOME

# edit files that match a ripgrep pattern
function rgv {
	vi -p $(rg -l $1)
}

# open vim files in tabs.
alias vi='vim -p'

# golang
alias gb='go build'
alias gd='go build --tags debug'
# build windows app without console popping up.
alias gbwin='go build -ldflags "-H windowsgui"'
alias gbrel='go build -ldflags "-H windowsgui -s -w"'

# put the following in .inputrc
#set bell-style none

# git
alias gl='git log --date=short --pretty=format:"%C(yellow)%h%Cblue%>(12)%ad %Cgreen%<(7)%aN %Cred%d%Creset %s"'

# git aware prompt
export PS1="\[\e[32m\][\w]\[\e[89m\]\$(GIT_PS1_SHOWDIRTYSTATE=1 __git_ps1)\[\033[00m\] $ "
# add local binaries
export PATH=$PATH:~/bin
# add windows SDK for packaging windows store apps
export PATH=$PATH:"/c/Program Files (x86)/Windows Kits/10/bin/10.0.26100.0/x64/"
```

`~/.zshrc`
```zsh
# colour the terminal
export CLICOLOR=1
export GREP_OPTIONS='--color=auto'

# use tabbed vim
export EDITOR=vim
alias vi='vim -v -p'

# add local commands
path+=('/usr/local/go/bin')
path+=('/Users/rust/bin')
export PATH

# customize the prompt : https://jonasjacek.github.io/colors/
# 0-15 are the solorized colors in terminal preferences.
export PROMPT='%(?.%F{3}√.%F{1}?%?)%f %F{3}%~%f %# '
# make it git friendly.
autoload -Uz vcs_info
precmd_vcs_info() { vcs_info }
precmd_functions+=( precmd_vcs_info )
setopt prompt_subst
RPROMPT=\$vcs_info_msg_0_
zstyle ':vcs_info:git:*' formats '%F{3}(%b)%f'

# Ignore duplicate commands in history
export HISTORY_IGNORE='(ls|fg|bg|exit)'
setopt hist_ignore_all_dups
setopt hist_find_no_dups
setopt hist_ignore_space
# Dont add ignored commands to history.
zshaddhistory() {
  emulate -L zsh
  ## uncomment if HISTORY_IGNORE
  ## should use EXTENDED_GLOB syntax
  setopt extendedglob
  [[ $1 != ${~HISTORY_IGNORE} ]]
}

# golang
alias gb='go build'
alias gd='go build --tags debug'

# git
alias gl='git log --date=short --pretty=format:"%C(yellow)%h%Cblue%>(12)%ad %Cgreen%<(7)%aN %Cred%d%Creset %s"'
```
