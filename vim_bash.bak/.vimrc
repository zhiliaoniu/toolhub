set nocompatible "关闭兼容模式，不在兼容vi，可以正常使用vim的扩展功能

syntax on "语法高亮
filetype plugin indent on
syntax enable
set number "设置行号
set ts=4

set history=1000

"自动对齐；智能选择对齐方式，使用与C编程
set autoindent
set smartindent
set autowrite       
set linebreak        " 整词换行
set whichwrap=b,s,<,>,[,] " 光标从行首和行末时可以跳到另一行去
set splitright "split windows on right

set hlsearch "hight light search word
set incsearch
set cindent

set shiftwidth=4 "设置tab键为4个空格
set tabstop=4 "设置当行之间交错时使用4个空格
set softtabstop=4
set expandtab

set showmatch "设置匹配模式，类似当输入一个左括号时会匹配相应的那个右括号
set laststatus=2 "显示文件名

set cinoptions={0,1s,t0,n-2,p2s,(03s,=.5s,>1s,=1s,:1s
if &term=="xterm"
    set t_Co=256
    set t_Sb=^[[4%dm
    set t_Sf=^[[3%dm
endif

set enc=utf-8
set fencs=utf-8,GB18030,ucs-bom,shift-jis,gbk,gb2312,cp936,default
let NERDTreeWinPos="right"


"------ctags setting begin------
" 按下F5重新生成tag文件，并更新taglist
"map <F5> :!ctags -R --c++-kinds=+p --fields=+iaS --extra=+q .<CR><CR> :TlistUpdate<CR>
"imap <F5> <ESC>:!ctags -R --c++-kinds=+p --fields=+iaS --extra=+q .<CR><CR> :TlistUpdate<CR>
set tags=tags
" set autochdir
set tags+=./tags "add current directory's generated tags file
set tags+=~/.vim/tags/cpp_src/tags
source ~/.vim/bundle/cscope_maps.vim
"------ctags setting end--------


" pathongen
call pathogen#infect()

let g:go_disable_autoinstall = 0
let g:go_version_warning = 0

" taglist
let Tlist_Show_One_File=1
let Tlist_Exit_OnlyWindow=1
let Tlist_Use_Right_Window=1
let Tlist_Sort_Type="name"

" omnicppcomplete
set completeopt=longest,menu
let OmniCpp_NamespaceSearch = 2     " search namespaces in the current
"buffer and in included files
let OmniCpp_ShowPrototypeInAbbr = 1 " 显示函数参数列表
let OmniCpp_MayCompleteScope = 1    " 输入 :: 后自动补全
let OmniCpp_DefaultNamespaces = ["std", "_GLIBCXX_STD"]

"cpp-enhanced-highlight https://github.com/octol/vim-cpp-enhanced-highlight
"高亮类，成员函数，标准库和模板
let g:cpp_class_scope_highlight = 1
let g:cpp_member_variable_highlight = 1
let g:cpp_class_decl_highlight = 1
let g:cpp_concepts_highlight = 1
let g:cpp_experimental_simple_template_highlight = 1
"文件较大时使用下面的设置高亮模板速度较快，但会有一些小错误
let g:cpp_experimental_template_highlight = 1"



"==============有关Vundle.vim的配置==================
"set the runtime path to include Vundle and initialize
set rtp+=~/.vim/bundle/Vundle.vim
call vundle#begin()

" let Vundle manage Vundle, required
Bundle 'gmarik/Vundle.vim'
Plugin 'L9'
Plugin 'rstacruz/sparkup', {'rtp': 'vim/'}
"Plugin 'fatih/vim-go'
"All of your Plugins must be added before the following line
call vundle#end()
"==============有关Vundle.vim的配置==================


"augroup filetype
"    autocmd! BufRead,BufNewFile BUILD set filetype=blade
"augroup end

"high light for blade BUILD file
autocmd! BufRead,BufNewFile BUILD set filetype=blade
