" Required dependencies: vim-plug and tokyonight
" Don't forget to run :PlugInstall to make the plugins working
set number
set numberwidth=2

call plug#begin()

Plug 'ghifarit53/tokyonight-vim'

call plug#end()

set termguicolors

let g:tokyonight_style = 'night' " available: night, storm
let g:tokyonight_enable_italic = 1

colorscheme tokyonight
