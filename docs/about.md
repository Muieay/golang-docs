# 关于本站

## 文档简介

本网站是一个专为Java/C++开发者设计的Golang技术文档站点，旨在帮助有编程经验的开发者快速掌握Golang语言。

<script setup>
import { VPTeamMembers } from 'vitepress/theme'

const members = [
    {
    avatar: 'https://www.github.com/muieay.png',
    name: 'Muieay',
    title: 'Creator',
    links: [
      { icon: 'github', link: 'https://github.com/muieay' },
      { icon: 'wechat', link: 'https://muieay.github.io/WeChat/' },
    ]
    },
    {
    avatar: 'https://www.github.com/golang.png',
    name: 'Golang',
    title: 'Official',
    links: [
      { icon: 'github', link: 'https://github.com/golang/go' },
      { icon: 'twitter', link: 'https://twitter.com/golang' }
    ]
  }
]
</script>

<VPTeamMembers size="small" :members />


## 目标受众

- **Java开发者**：希望了解Golang与Java的差异和优势
- **C++开发者**：想要体验Golang的简洁和高效
- **后端工程师**：寻求在微服务架构中使用Golang
- **全栈开发者**：扩展技术栈，学习Golang全栈开发

## 内容特色

### 📚 系统化学习路径
- **基础篇**：环境安装、语法基础、错误处理
- **进阶篇**：指针、接口、反射、泛型、测试
- **并发篇**：Goroutine、Channel、并发模式
- **Web篇**：HTTP服务、Gin框架、RESTful API
- **数据库篇**：SQL操作、GORM、Redis缓存

### 🔍 对比式学习
- Java vs Golang 语法对比
- C++ vs Golang 内存管理对比
- 设计模式在Golang中的实现

### 💡 实战导向
- 每个章节都包含完整代码示例
- 提供可直接运行的项目案例
- 涵盖真实业务场景

## 技术栈

- **构建工具**：VitePress
- **主题**：VitePress默认主题
- **部署**：支持GitHub Pages、Vercel等平台
- **搜索**：本地全文搜索

## 更新计划

- [ ] 添加微服务开发章节
- [ ] 补充Kubernetes部署指南
- [ ] 增加性能优化技巧
- [ ] 添加更多实战项目案例

## 贡献指南

如果你发现文档中的错误或想要添加新内容，欢迎通过以下方式贡献：

1. Fork本项目
2. 创建功能分支
3. 提交你的修改
4. 创建Pull Request

## 联系方式

- **文档维护**：[muieay/golang-docs](https://github.com/muieay/golang-docs)
- **GitHub**：[golang/go](https://github.com/golang/go)
- **官方文档**：[golang.org](https://golang.org)
