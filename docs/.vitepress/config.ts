import { defineConfig } from 'vitepress'

export default defineConfig({
  title: 'Go Pip SDK',
  description: 'A comprehensive Go SDK for managing Python pip operations, virtual environments, and Python projects',
  
  locales: {
    root: {
      label: 'English',
      lang: 'en',
      title: 'Go Pip SDK',
      description: 'A comprehensive Go SDK for managing Python pip operations, virtual environments, and Python projects',
      themeConfig: {
        nav: [
          { text: 'Home', link: '/' },
          { text: 'Guide', link: '/guide/getting-started' },
          { text: 'API Reference', link: '/api/' },
          { text: 'Examples', link: '/examples/' },
          { text: 'GitHub', link: 'https://github.com/scagogogo/go-pip-sdk' }
        ],
        sidebar: {
          '/guide/': [
            {
              text: 'Guide',
              items: [
                { text: 'Getting Started', link: '/guide/getting-started' },
                { text: 'Installation', link: '/guide/installation' },
                { text: 'Configuration', link: '/guide/configuration' },
                { text: 'Package Management', link: '/guide/package-management' },
                { text: 'Virtual Environments', link: '/guide/virtual-environments' },
                { text: 'Project Management', link: '/guide/project-management' },
                { text: 'Error Handling', link: '/guide/error-handling' },
                { text: 'Logging', link: '/guide/logging' }
              ]
            }
          ],
          '/api/': [
            {
              text: 'API Reference',
              items: [
                { text: 'Overview', link: '/api/' },
                { text: 'Manager', link: '/api/manager' },
                { text: 'Package Operations', link: '/api/package-operations' },
                { text: 'Virtual Environments', link: '/api/virtual-environments' },
                { text: 'Project Management', link: '/api/project-management' },
                { text: 'Types', link: '/api/types' },
                { text: 'Errors', link: '/api/errors' },
                { text: 'Logger', link: '/api/logger' },
                { text: 'Installer', link: '/api/installer' }
              ]
            }
          ],
          '/examples/': [
            {
              text: 'Examples',
              items: [
                { text: 'Overview', link: '/examples/' },
                { text: 'Basic Usage', link: '/examples/basic-usage' },
                { text: 'Package Management', link: '/examples/package-management' },
                { text: 'Virtual Environments', link: '/examples/virtual-environments' },
                { text: 'Project Initialization', link: '/examples/project-initialization' },
                { text: 'Advanced Usage', link: '/examples/advanced-usage' }
              ]
            }
          ]
        },
        socialLinks: [
          { icon: 'github', link: 'https://github.com/scagogogo/go-pip-sdk' }
        ],
        footer: {
          message: 'Released under the MIT License.',
          copyright: 'Copyright © 2024 Go Pip SDK Contributors'
        }
      }
    },
    zh: {
      label: '简体中文',
      lang: 'zh-CN',
      title: 'Go Pip SDK',
      description: '用于管理 Python pip 操作、虚拟环境和 Python 项目的综合 Go SDK',
      themeConfig: {
        nav: [
          { text: '首页', link: '/zh/' },
          { text: '指南', link: '/zh/guide/getting-started' },
          { text: 'API 参考', link: '/zh/api/' },
          { text: '示例', link: '/zh/examples/' },
          { text: 'GitHub', link: 'https://github.com/scagogogo/go-pip-sdk' }
        ],
        sidebar: {
          '/zh/guide/': [
            {
              text: '指南',
              items: [
                { text: '快速开始', link: '/zh/guide/getting-started' },
                { text: '安装', link: '/zh/guide/installation' },
                { text: '配置', link: '/zh/guide/configuration' },
                { text: '包管理', link: '/zh/guide/package-management' },
                { text: '虚拟环境', link: '/zh/guide/virtual-environments' },
                { text: '项目管理', link: '/zh/guide/project-management' },
                { text: '错误处理', link: '/zh/guide/error-handling' },
                { text: '日志记录', link: '/zh/guide/logging' }
              ]
            }
          ],
          '/zh/api/': [
            {
              text: 'API 参考',
              items: [
                { text: '概述', link: '/zh/api/' },
                { text: '管理器', link: '/zh/api/manager' },
                { text: '包操作', link: '/zh/api/package-operations' },
                { text: '虚拟环境', link: '/zh/api/virtual-environments' },
                { text: '项目管理', link: '/zh/api/project-management' },
                { text: '类型', link: '/zh/api/types' },
                { text: '错误', link: '/zh/api/errors' },
                { text: '日志', link: '/zh/api/logger' },
                { text: '安装器', link: '/zh/api/installer' }
              ]
            }
          ],
          '/zh/examples/': [
            {
              text: '示例',
              items: [
                { text: '概述', link: '/zh/examples/' },
                { text: '基本用法', link: '/zh/examples/basic-usage' },
                { text: '包管理', link: '/zh/examples/package-management' },
                { text: '虚拟环境', link: '/zh/examples/virtual-environments' },
                { text: '项目初始化', link: '/zh/examples/project-initialization' },
                { text: '高级用法', link: '/zh/examples/advanced-usage' }
              ]
            }
          ]
        },
        socialLinks: [
          { icon: 'github', link: 'https://github.com/scagogogo/go-pip-sdk' }
        ],
        footer: {
          message: '基于 MIT 许可证发布。',
          copyright: 'Copyright © 2024 Go Pip SDK 贡献者'
        }
      }
    }
  },
  
  themeConfig: {
    search: {
      provider: 'local'
    }
  },

  ignoreDeadLinks: true,
  
  head: [
    ['link', { rel: 'icon', href: '/favicon.ico' }],
    ['meta', { name: 'theme-color', content: '#3c82f6' }],
    ['meta', { property: 'og:type', content: 'website' }],
    ['meta', { property: 'og:locale', content: 'en' }],
    ['meta', { property: 'og:title', content: 'Go Pip SDK | A comprehensive Go SDK for Python pip operations' }],
    ['meta', { property: 'og:site_name', content: 'Go Pip SDK' }],
    ['meta', { property: 'og:image', content: 'https://scagogogo.github.io/go-pip-sdk/og-image.png' }],
    ['meta', { property: 'og:url', content: 'https://scagogogo.github.io/go-pip-sdk/' }]
  ]
})
