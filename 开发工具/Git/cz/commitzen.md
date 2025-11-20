Commitizen主要有两个流行的版本：一个是用**Python**开发的，另一个是用**Node.js**开发的（通常在npm上安装）。它们的安装方式是截然不同的。

下面的表格帮你快速梳理它们的区别：

| 特性               | Python 版 Commitizen                   | Node.js 版 Commitizen          |
| :----------------- | :------------------------------------- | :----------------------------- |
| **安装命令** | `pip install -U Commitizen`          | `npm install -g commitizen`  |
| **使用命令** | `cz commit` 或 `cz c`              | `git cz`                     |
| **配置文件** | 通常配置在 `pyproject.toml` 等文件中 | 通常配置在 `package.json` 中 |
| **主要用途** | Python 项目                            | Node.js/前端 项目              |

### 🐍 安装Python版Commitizen

如果你确定要为**Python项目**使用Commitizen，可以按以下步骤操作：

1. **安装工具**
   在命令行中执行以下命令进行全局安装：

   ```bash
   pip install -U Commitizen
   ```

   或者，在你的Python项目中使用 `poetry`安装：

   ```bash
   poetry add commitizen --dev
   ```
2. **在项目中使用**
   安装完成后，你就可以在Git仓库中使用 `cz commit`或简写 `cz c`命令来替代 `git commit`，它会启动一个交互式界面来帮助你生成符合规范的提交信息。

### 🟢 安装Node.js版Commitizen

如果你其实是想为**前端或Node.js项目**使用Commitizen，那么正确的步骤如下：

1. **安装Node.js**：首先，你需要确保系统已经安装了Node.js（它自带了npm包管理器）。这与安装Python是不同的。
2. **全局安装Commitizen**：
   在命令行中运行：

   ```bash
   npm install -g commitizen
   ```

   请注意，这一步是在你的系统全局环境中安装 `commitizen`命令行工具。
3. **在项目中初始化适配器**
   接着，你需要在你想要使用规范提交的**项目根目录**（该目录下通常有 `package.json`文件）下，运行以下命令来指定使用哪种提交规范（常用的如 `cz-conventional-changelog`）：

   ```bash
   commitizen init cz-conventional-changelog --save --save-exact
   ```

   如果项目之前配置过，你可以加上 `--force`参数强制重置。
4. **使用**
   完成以上步骤后，在这个项目里你就可以使用 `git cz`命令来开始交互式提交了。

### 💡 如何选择与确认

- **回想你的项目类型**：你主要在用什么语言开发？是Python后端项目，还是JavaScript/TypeScript前端项目？
- **检查安装来源**：你可以回想一下当初是在哪个网站或文档里看到Commitizen的，这能帮你判断是哪个版本。
- **确认现有环境**：在命令行中分别执行 `cz --version` 和 `npm list -g commitizen`，看哪个命令能返回有效的版本号，就说明你安装的是哪个版本。
