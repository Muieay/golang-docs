import { defineConfig } from 'vitepress'

// https://vitepress.dev/reference/site-config
export default defineConfig({
  title: "Golang 技术文档",
  description: "从 Java/C++ 转向 Golang 的完整学习指南",
  lastUpdated: true,
  themeConfig: {
    // https://vitepress.dev/reference/default-theme-config
    nav: [
      { text: '首页', link: '/' },
      { text: '教程', link: '/golang-basics' },
      { text: '关于', link: '/about' }
    ],

    sidebar: [
      {
        text: 'Golang 学习路线',
        collapsed: false,
        items: [
          { text: '基础教程', link: '/golang-basics' },
          { text: '进阶教程', link: '/golang-advanced' },
          { text: '并发编程', link: '/golang-concurrency' },
          { text: 'Web开发', link: '/golang-web' },
          { text: '数据库操作', link: '/golang-database' },
        ]
      },
      {
        text: '基础篇',
        collapsed: false,
        items: [
          { text: '环境安装与配置', link: '/golang-basics#环境安装与配置' },
          { text: 'Hello World', link: '/golang-basics#hello-world' },
          { text: '基础语法', link: '/golang-basics#基础语法' },
          { text: '错误处理', link: '/golang-basics#错误处理' }
        ]
      },
      {
        text: '进阶篇',
        collapsed: false,
        items: [
          { text: '指针与内存管理', link: '/golang-advanced#指针与内存管理' },
          { text: '结构体与方法', link: '/golang-advanced#结构体与方法' },
          { text: '接口与多态', link: '/golang-advanced#接口与多态' },
          { text: '反射', link: '/golang-advanced#反射' },
          { text: '泛型', link: '/golang-advanced#泛型' },
          { text: '测试', link: '/golang-advanced#测试' }
        ]
      },
      {
        text: '并发篇',
        collapsed: false,
        items: [
          { text: 'Goroutine基础', link: '/golang-concurrency#goroutine基础' },
          { text: 'Channel通信', link: '/golang-concurrency#channel通信' },
          { text: 'Select语句', link: '/golang-concurrency#select语句' },
          { text: '并发模式', link: '/golang-concurrency#并发模式' },
          { text: '并发安全', link: '/golang-concurrency#并发安全' }
        ]
      },
      {
        text: 'Web开发篇',
        collapsed: false,
        items: [
          { text: 'HTTP基础', link: '/golang-web#http基础' },
          { text: '使用Gin框架', link: '/golang-web#使用gin框架' },
          { text: '中间件', link: '/golang-web#中间件' },
          { text: '静态文件服务', link: '/golang-web#静态文件服务' },
          { text: '模板渲染', link: '/golang-web#模板渲染' },
          { text: '错误处理', link: '/golang-web#错误处理' }
        ]
      },
      {
        text: '数据库篇',
        collapsed: false,
        items: [
          { text: '标准库database/sql', link: '/golang-database#数据库sql标准库' },
          { text: '使用GORM', link: '/golang-database#使用gorm' },
          { text: '关联关系', link: '/golang-database#关联关系' },
          { text: '事务处理', link: '/golang-database#事务处理' },
          { text: '连接池配置', link: '/golang-database#连接池配置' },
          { text: 'Redis操作', link: '/golang-database#redis操作' }
        ]
      }
    ],

    socialLinks: [
      { icon: 'github', link: 'https://github.com/muieay/golang-docs' },
      { icon: 'wechat', link: 'https://muieay.github.io/WeChat/' },
    ],

    footer: {
      message: 'Released under the MIT License.',
      copyright: 'Copyright © 2025 Golang 技术文档'
    },

    search: {
      provider: 'local'
    }
  }
})
