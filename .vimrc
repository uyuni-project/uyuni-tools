" SPDX-FileCopyrightText: 2023 SUSE LLC
"
" SPDX-License-Identifier: Apache-2.0

" Local vim configuration loaded by https://github.com/LucHermitte/local_vimrc
" For local_vimrc to use this file, ensure .vimrc is in the g:local_vimrc
" list. You can set it like the following in the vim or neovim config:
"
"     let g:local_vimrc = ['.vimrc']

" Set make command
set makeprg=go\ build\ ./...
