#### 文章目录

*   [CI/CD](#CICD_2)
*   *   [持续集成（CI）](#CI_5)
    *   *   [概念](#_7)
        *   [流程](#_10)
        *   [关键组件](#_13)
        *   [作用](#_23)
    *   [持续交付（CD）](#CD_33)
    *   *   [概念](#_35)
        *   [流程](#_38)
        *   [关键组件](#_41)
        *   [作用](#_51)
*   [GitHub Action](#GitHub_Action_58)
*   *   [概念](#_60)
    *   [基本要点](#_63)
    *   *   [工作流（Workflows）](#Workflows_65)
        *   [任务（Jobs）](#Jobs_68)
        *   [步骤（Steps）](#Steps_71)
    *   [实践场景](#_74)
    *   *   [需求](#_76)
        *   [步骤](#_79)
        *   *   [创建 pull_request 工作流文件](#pull_request_81)
            *   [提交 PR 并观察](#PR_120)
            *   [创建用于部署的工作流文件](#_125)
            *   [部署文件](#_202)
            *   [提交到部署分支并观察](#_226)
    *   [其他概念介绍](#_239)
    *   *   [概念功能](#_241)
        *   *   [触发器（Triggers）](#Triggers_243)
            *   [环境（Environments）](#Environments_257)
            *   [矩阵构建（Matrix Builds）](#Matrix_Builds_267)
            *   [缓存（Caching）](#Caching_275)
            *   [自定义操作（Custom Actions）](#Custom_Actions_299)
            *   [部署（Deployment）](#Deployment_302)
            *   [Secrets](#Secrets_305)
        *   [语法](#_318)
        *   *   [工作流程文件结构](#_320)
            *   [步骤（Steps）](#Steps_350)
            *   [操作（Actions）](#Actions_363)
            *   [矩阵构建（Matrix Builds）](#Matrix_Builds_371)
            *   [环境变量](#_379)
            *   [条件（Conditions）](#Conditions_391)
            *   [超时和重试](#_400)
            *   [缓存（Caching）](#Caching_410)

CI/CD
-----

CI/CD 是**持续集成（Continuous Integration）**和**持续交付（Continuous Delivery）**的缩写，它旨在通过**自动化**的流程和工具，提高[软件开发](https://so.csdn.net/so/search?q=%E8%BD%AF%E4%BB%B6%E5%BC%80%E5%8F%91&spm=1001.2101.3001.7020)的效率、质量和交付速度。  

### [持续集成](https://so.csdn.net/so/search?q=%E6%8C%81%E7%BB%AD%E9%9B%86%E6%88%90&spm=1001.2101.3001.7020)（CI）

#### 概念

持续集成是开发团队通过将代码的不同部分集成到共享存储库中，并频繁地进行构建和测试，以确保代码的一致性和稳定性。  

#### 流程

在现在的开发模式中，一般的项目，协同开发是离不开的，这就涉及到多个开发人员编写处理自己负责的功能模块或者某些开发人员共同负责一个模块。于是，通过**版本控制系统**可以将各个开发人员的代码集成在该共享存储库里，在存储库里，每个开发人员根据需求的不同来创建对应的分支，在完成需求后，每个人都需要提交合并将开发分支代码集成在一起，这就需要解决代码冲突，并且如何除了 code review 之外如何确保这些更改对应用没有产生影响，一旦提交请求合并到主分支，**自动化构建工具**就会根据流程自动编译构建安装应用，并执行**单元测试框架**的自动化测试来校验提交的修改。  

#### 关键组件

以下是一些用于构建有效 CI 流程的关键组件：

1.  **版本控制系统（Version Control System，VCS）**：
    *   例如 Git，用于跟踪代码变更，协作开发，并确保团队成员之间的代码同步。
2.  **自动化构建工具**：
    *   如 Jenkins、Travis CI、CircleCI 等，用于在每次代码提交时自动触发构建过程。
3.  **单元测试框架**：
    *   例如 JUnit（Java）、pytest（Python），用于确保代码的基本功能在集成后仍然有效。  
        

#### 作用

*   **减少集成问题**： 在传统的开发模式中，团队成员可能在各自的开发分支上独立工作，而在合并时可能会产生冲突和集成问题。CI 通过持续集成代码，及时发现和解决这些问题，避免了集成地狱。
*   **提高代码质量**： CI 强调自动化测试，包括单元测试、集成测试等。每次代码变更都会触发这些测试，确保新代码不会破坏现有功能，并减少引入 bug 的可能性。这有助于提高整体代码质量。
*   **快速反馈**： CI 通过快速执行自动化构建和测试，提供了即时反馈。开发人员可以在提交代码后迅速得知其是否通过了构建和测试，帮助他们更快速地发现和修复问题。
*   **提高开发效率**： 通过自动化构建、测试和部署，CI 减少了手动操作的需求，提高了开发效率。开发人员可以专注于编写代码而不必花费过多时间在手动构建和测试上。
*   **自动化部署**： 与持续交付（Continuous Delivery）和持续部署（Continuous Deployment）结合，CI 可以实现自动化部署。这意味着经过测试的代码变更可以自动部署到预定环境，实现快速且可靠的交付流程。
*   **团队协作**： CI 鼓励团队成员频繁集成代码，确保大家的工作在一个共享的代码库中协同进行。这促进了团队之间的协作和沟通，减少了因代码集成问题而导致的沟通障碍。
*   **降低风险**： 通过频繁集成和自动测试，CI 减少了发布到生产环境时出现问题的可能性。提前发现和解决问题有助于降低风险，确保稳定的软件交付。  
    

### 持续交付（CD）

#### 概念

持续交付建立在持续集成的基础上，通过自动化的流程确保软件可以随时随地进行部署。  

#### 流程

这时，持续交付后的代码已经在主分支上了，这处于某个版本的待发布的状态，随时可以将开发环境的功能部署到生产环境中（部署到生成环境前还需要在测试环境性能测试、回归测试、自动化测试、人工测试等），运行脚本**构建打包应用**，通过**自动化部署工具**部署到生产环境运行应用，**监控**生产环境指标，如出现问题和错误，可以触发手动或自动**回滚**，如系统正常，则定期回顾，收集反馈，优化，并**持续改进**。  

#### 关键组件

以下是一些用于实现持续交付的关键组件：

1.  **自动化部署工具**：
    *   例如 Docker、Ansible、Kubernetes 等，用于自动化地部署应用程序和其依赖。
2.  **环境配置管理**：
    *   工具如 Terraform，确保不同环境（开发、测试、生产）的一致性。
3.  **持续监控和反馈**：
    *   使用工具如 Prometheus、Grafana，确保在部署后能够监控应用程序的性能和稳定性。  
        

#### 作用

*   **快速交付**： 持续交付强调频繁、快速地将新的代码变更交付到生产环境。这使得团队能够更加迅速地响应用户需求，推出新功能或修复 bug。
*   **稳定交付**： 通过自动化测试、自动化部署和验证流程，持续交付确保每次交付都是经过充分验证的，降低了引入错误的风险，提高了软件的稳定性。
*   **降低发布成本**： 持续交付通过自动化流程降低了发布的人工成本。这意味着开发团队不再需要手动执行繁琐的部署步骤，减少了错误的可能性，提高了整体效率。
*   **支持持续改进**： 持续交付是一个循环过程，通过不断收集用户反馈、监控系统性能和流程改进，团队能够不断优化持续交付流程，提高整体效率和质量。  
    

GitHub Action
-------------

### 概念

采用 CI/CD 可以通过自动化流程和工具自动帮你构建应用、测试应用、部署应用，将你的应用交给流程工具来管理，做到自动触发、验证、部署等功能，从而减省人工成本、提高交付速度，在敏捷开发、DevOps 中扮演着重要的角色。  
GitHub Action 正是这样一个实现持续集成交付的自动化流程工具，是由 GitHub 提供的一个组件。你可以通过 YAML 文件的配置定义工作流程以构建执行 CI/CD 流水线，并可以触发不同事件时（如代码提交 push、Pull Request、schedule）自动执行这些工作流程。  

### 基本要点

#### 工作流（Workflows）

工作流是 GitHub Actions 执行任务的基本单位，你可以为 Git 上不同的事件（如 push、pull、request 等）定义不同的工作流，以响应操作代码的变更。  

#### 任务（Jobs）

工作流程由一个或多个任务组成，每个任务运行在独立的虚拟环境中。任务可以是构建、测试、部署等操作。  

#### 步骤（Steps）

任务由多个步骤组成，每个步骤执行一个操作。一个步骤可以是运行命令、使用某个预定义的操作，或者调用自定义脚本。  

### 实践场景

#### 需求

假如我们对项目中其中一个服务做了修改，添加了某些功能，完成任务后，我们在本地分支通过 Git 提交代码到 Github 项目仓库下的 dev 分支（这里直接本地提交到测试分支，省去测试环境测试的流程），并请求合并到 master 分支，这时，我们希望在合并之前先对该模块进行构建，运行测试来校验代码质量与验证代码是否出错，确保代码的基本功能在集成后仍然有效，测试通过后，提交到打包部署分支 bdeploy 来自动将模块打包成一个容器镜像推送到容器镜像仓库，并将 docker-compose 文件拷贝到远程生产服务器执行部署。

#### 步骤

##### 创建 pull_request 工作流文件

在项目目录下创建. github/workflows 目录  
添加 compile.yml 文件用于构建并测试项目：

```
name: compile

on:
  pull_request:
    paths: #当有 pull request，且文件路径包含 Java 文件或者当前的工作流配置文件时触发。
      - '**.java'
      - .github/workflows/compile.yml
jobs:
  compile: #任务名称
    runs-on: ubuntu-latest
    timeout-minutes: 30

    steps:
      - name: Checkout code
        uses: actions/checkout@v3 #actions/checkout@v3 是 GitHub Actions 中一个常用的操作（Action），用于从存储库中检出代码。@v3 是指定该 Action 的版本号。在这里，v3 表示使用的是版本 3。

      - name: Set up Java
        uses: actions/setup-java@v3 #actions/setup-java@v3 操作被用于设置 Java 运行环境
        with:
          java-version: '11' # 指定所需的 Java 版本
          distribution: 'temurin' #'temurin' 表示使用 Temurin（先前称为 AdoptOpenJDK） 的发行版。Temurin 提供了免费的、社区驱动的 OpenJDK 发行版。

      - name: Build with Maven
        run: mvn clean install

      - name: Run JUnit test
        run: mvn test


```

上面工作流配置文件定义了在 pull_request 的时候会触发任务  
定义了一个任务`compile`的四个步骤：

*   `Checkout code`：从存储库中检出代码
*   `Set up Java`：设置 Java 运行环境
*   `Build with Maven`：构建安装相关依赖
*   `Run JUnit test`：执行单元测试  
    

##### 提交 PR 并观察

将当前分支的代码推送到远程 github 项目仓库的 dev 分支，并提 PR 请求合并到 master 分支。  
![](https://i-blog.csdnimg.cn/blog_migrate/c7d095bd26dc434953be6f8835b5ada6.png)  
提交 PR 后会自动触发执行工作流任务，查看详细：  
依次执行了我们定义的任务，并且设置环境、执行构建和测试通过  
![](https://i-blog.csdnimg.cn/blog_migrate/504c1d857fc714df5179b379f6e58a92.png)  
![](https://i-blog.csdnimg.cn/blog_migrate/a82eacddfb589414225e10fe095352c9.png)  
![](https://i-blog.csdnimg.cn/blog_migrate/cf30404243a9b116187640d26aea7371.png)

之后可以选择合并此 PR 到 master 分支，将修改的代码合并到主分支准备部署。  

##### 创建用于部署的工作流文件

代码合并到 master 分支后，在. github/workflows 目录目录下创建用于部署的工作流文件 bdeploy.yml：

```
name: Build and Deploy for aliyun

on:
  push:
    branches: [bdeploy]

jobs:
  build:
    runs-on: ubuntu-latest
    timeout-minutes: 30
    strategy: #矩阵策略
      matrix:
        java: [ '11' ]

    steps:
      - name: Checkout code
        uses: actions/checkout@v3 #actions/checkout@v3 是 GitHub Actions 中一个常用的操作（Action），用于从存储库中检出代码。@v3 是指定该 Action 的版本号。在这里，v3 表示使用的是版本 3。

      - name: Set up Java
        uses: actions/setup-java@v3 #actions/setup-java@v3 操作被用于设置 Java 运行环境
        with:
          java-version: ${{ matrix.java }} # 指定所需的 Java 版本
          distribution: 'temurin' #'temurin' 表示使用 Temurin（先前称为 AdoptOpenJDK） 的发行版。Temurin 提供了免费的、社区驱动的 OpenJDK 发行版。

      - name: Build base
        run: mvn clean install

      - name: Build container image
        run: mvn clean package -DskipTests jib:build # -Pdocker

      - name: Deploy server
        run: |
          echo -e "[demo] \n${{ secrets.SERVER_DEMO }} ansible_ssh_port=${{ secrets.PORT_DEMO }} ansible_ssh_user=${{ secrets.ACCOUNT_DEMO }} ansible_ssh_pass='${{ secrets.PASSWORD_DEMO }}'" > ./hostfile
          docker run -v $PWD/hostfile:/tmp/hostfile -e ANSIBLE_HOST_KEY_CHECKING=false --rm ghcr.io/yunhorn/ubuntu:ansible ansible -i /tmp/hostfile demo -m shell -a "touch /usr/local/demo/docker-compose-deploy.yml && mv /usr/local/demo/docker-compose-deploy.yml /usr/local/demo/docker-compose-deploy.yml_bak"
          docker run -v $PWD/hostfile:/tmp/hostfile -v $PWD/deploy:/tmp/deploy -e ANSIBLE_HOST_KEY_CHECKING=false --rm ghcr.io/yunhorn/ubuntu:ansible ansible -i /tmp/hostfile demo -m copy -a "src=/tmp/deploy/prod/docker-compose-deploy.yml dest=/usr/local/demo/docker-compose-deploy.yml"
          docker run -v $PWD/hostfile:/tmp/hostfile -e ANSIBLE_HOST_KEY_CHECKING=false --rm ghcr.io/yunhorn/ubuntu:ansible ansible -i /tmp/hostfile demo -m shell -a "docker login --username=${{ secrets.ALINYUN_USERNAME }} --password=${{ secrets.ALINYUN_PASSWORD }} registry.cn-hangzhou.aliyuncs.com"
          docker run -v $PWD/hostfile:/tmp/hostfile -e ANSIBLE_HOST_KEY_CHECKING=false --rm ghcr.io/yunhorn/ubuntu:ansible ansible -i /tmp/hostfile demo -m shell -a "docker-compose -f /usr/local/demo/docker-compose-deploy.yml --compatibility up -d "
          docker run -v $PWD/hostfile:/tmp/hostfile -e ANSIBLE_HOST_KEY_CHECKING=false --rm ghcr.io/yunhorn/ubuntu:ansible ansible -i /tmp/hostfile demo -m shell -a "docker logout registry.cn-hangzhou.aliyuncs.com"
          rm -f hostfile

```

上面工作流文件创建了一个名为 Build and Deploy for aliyun 的工作，定义了一个任务`build`的五个步骤：

*   `Checkout code`：从存储库中检出代码
*   `Set up Java`：设置 Java 运行环境
*   `Build base`：安装依赖构建项目
*   `Build container image`：执行 Google Jib 的 maven 插件将当前项目打包并推送到远程容器镜像仓库。该插件的具体用法可参考我之前写的文章：

[轻松构建 Docker 镜像：无需 Docker 引擎的 Google Jib-CSDN 博客](https://blog.csdn.net/weixin_44268936/article/details/133385011?spm=1001.2014.3001.5501)

*   `Deploy server`：该部分实现了通过拷贝我们即定的 docker-compose 文件到远程服务器上，并在远程服务上拉取该项目的容器镜像，最后启动容器来实现部署。操作远程服务器的行为借助了**自动化运维工具 ansible**。

Ansible 是一种自动化工具，基于 Python 开发，集合了众多运维工具（puppet、chef、func、fabric）的优点，实现了批量系统配置、批量程序部署、批量运行命令等功能。它是一个开源工具，使用简单，无需在被管理的主机上安装客户端，而且支持多云环境和多种操作系统。  
[Ansible is Simple IT Automation](https://www.ansible.com/)

我们来看下`Deploy server`做了什么

`run: |`  
“|”是 YAML 语法中的一个标记，表示执行一个多行字符串块，也称为 “折叠块”（folded block）或“纵向线条”（vertical line），“ | ” 后面的缩进代码块是一个 shell 命令的多行字符串。这样的写法允许你在一个步骤中执行多个命令，而不需要每个命令都单独使用一个步骤。

`echo -e "[demo] \n${{ secrets.SERVER_DEMO }} ansible_ssh_port=${{ secrets.PORT_DEMO }} ansible_ssh_user=${{ secrets.ACCOUNT_DEMO }} ansible_ssh_pass='${{ secrets.PASSWORD_DEMO }}'" > ./hostfile`  
将要登陆的服务信息写入`hostfile`文件，`[ ]`里用于指定一个服务组别。  
写入的格式为 ansible 可识别的主机清单文件格式，格式风格为：

```
[web_servers]
ansible_host ansible_ssh_port=22 ansible_ssh_user=username ansible_ssh_pass=password

```

`${{ }}`为 Github Actions 的`secrets and variables`语法。  
可以在项目的`Setting`的`secrets and variables`的`Actions`下来创建这些`Repository secrets`  
![](https://i-blog.csdnimg.cn/blog_migrate/649dd14076cb3afea4cc676eea08679b.png)

`docker run -v $PWD/hostfile:/tmp/hostfile -e ANSIBLE_HOST_KEY_CHECKING=false --rm ghcr.io/yunhorn/ubuntu:ansible ansible -i /tmp/hostfile demo -m shell -a "touch /usr/local/demo/docker-compose-deploy.yml && mv /usr/local/demo/docker-compose-deploy.yml /usr/local/demo/docker-compose-deploy.yml_bak"`  
将本地主机的 hostfile 文件挂载到容器中的 /tmp/hostfile 目录，以提供 Ansible 主机清单。-i /tmp/hostfile 指定了 Ansible 主机清单文件的路径，demo 是指定的主机组。-m shell：使用 Ansible 的 shell 模块，该模块用于在目标主机上执行 shell 命令。-a “touch /usr/local/demo/docker-compose-deploy.yml && mv /usr/local/demo/docker-compose-deploy.yml /usr/local/demo/docker-compose-deploy.yml_bak”：是 shell 模块的参数，其中包含要执行的 shell 命令，这创建了一个空的 docker-compose-deploy.yml 文件（如果不存在），将现有的 docker-compose-deploy.yml 部署文件备份。

`docker run -v $PWD/hostfile:/tmp/hostfile -v $PWD/deploy:/tmp/deploy -e ANSIBLE_HOST_KEY_CHECKING=false --rm ghcr.io/yunhorn/ubuntu:ansible ansible -i /tmp/hostfile demo -m copy -a "src=/tmp/deploy/prod/docker-compose-deploy.yml dest=/usr/local/demo/docker-compose-deploy.yml"`  
将本地主机的 hostfile 文件挂载到容器中的 /tmp/hostfile 目录，这是为了将本地主机上的 Ansible 主机清单文件提供给容器使用，并且将 deploy 目录下的文件挂载到容器中的 /tmp/deploy 目录，用于传递部署相关的文件。  
-m copy 使用 Ansible 的 copy 模块，该模块用于复制文件。-a “src=/tmp/deploy/prod/docker-compose-deploy.yml dest=/usr/local/demo/docker-compose-deploy.yml” 是 copy 模块的参数，指定了源文件和目标文件的路径，这将刚刚挂载到 ansible 容器内的部署文件复制到远程主机的指定目录文件下，方便后续启动部署的项目容器。

`docker run -v $PWD/hostfile:/tmp/hostfile -e ANSIBLE_HOST_KEY_CHECKING=false --rm ghcr.io/yunhorn/ubuntu:ansible ansible -i /tmp/hostfile demo -m shell -a "docker login --username=${{ secrets.ALINYUN_USERNAME }} --password=${{ secrets.ALINYUN_PASSWORD }} registry.cn-hangzhou.aliyuncs.com"`  
在远程服务上执行 docker login 登录到阿里云的容器镜像仓库。

`docker run -v $PWD/hostfile:/tmp/hostfile -e ANSIBLE_HOST_KEY_CHECKING=false --rm ghcr.io/yunhorn/ubuntu:ansible ansible -i /tmp/hostfile demo -m shell -a "docker-compose -f /usr/local/demo/docker-compose-deploy.yml --compatibility up -d "`  
在远程服务器上执行 docker-compose up 来启动项目容器。

`docker run -v $PWD/hostfile:/tmp/hostfile -e ANSIBLE_HOST_KEY_CHECKING=false --rm ghcr.io/yunhorn/ubuntu:ansible ansible -i /tmp/hostfile demo -m shell -a "docker logout registry.cn-hangzhou.aliyuncs.com"`  
`rm -f hostfile`  
最后退出阿里云的容器镜像仓库并删除本地（Github 项目上）hostfile 文件。  

##### 部署文件

编写用于部署我们提交的项目的 docker-compose 部署相关的文件：

```
version: "3.7"

x-logging:
  &default-logging
  options:
    mode: non-blocking
    max-buffer-size: 1m
    tag: "demo.{{.Name}}"  #配置容器的tag,以demo.为前缀,容器名称为后缀,docker-compose会给容器添加副本后缀

services:
  demo:
    logging: *default-logging
    restart: always
    user: root #该服务内运行的进程将以root用户的身份启动
    image: registry.cn-hangzhou.aliyuncs.com/minggo/demo:0.0.2-SNAPSHOT
    ports:
      - 8081:8081
    environment:
      - server.port=8081

```

##### 提交到部署分支并观察

假设我们的部署分支是`bdeploy`，提交到该分支后会自动触发用于部署的工作流文件的任务。  
在任务里会看到成功构建容器镜像并推送到阿里云容器镜像仓库：  
![](https://i-blog.csdnimg.cn/blog_migrate/9659ccb1c630cbfd35b736b1bc887495.png)  
成功拷贝我们的部署文件到远程服务，并且拉取我们刚刚推送的容器镜像，在服务器上创建了该容器，从而实现了项目的部署。  
![](https://i-blog.csdnimg.cn/blog_migrate/e076279bb1dccf8f4cc5fcff5534b31a.png)  
![](https://i-blog.csdnimg.cn/blog_migrate/bd249c54f7a76f867c97f85e40d9bd60.png)  
我们到服务器上查看

```
docker images

```

![](https://i-blog.csdnimg.cn/blog_migrate/caf7e15ea60c8bc70a6f6494af78fe1d.png)

```
docker ps

```

![](https://i-blog.csdnimg.cn/blog_migrate/94bfaddf4cc7d5e03e5868d2ef499db2.png)

以上就是一个基本的持续集成部署流程的示例，展示了如何使用 GitHub Actions 自动化构建、测试和部署一个 Java 应用项目的过程。当然，GitHub Actions 的功能远不止这些，它提供了丰富的集成和自定义选项，满足各种复杂的自动化需求。  

### 其他概念介绍

#### 概念功能

##### 触发器（Triggers）

GitHub Actions 的工作流程可以通过多种触发器启动。除了常见的 on: push，还有 on: pull_request、on: schedule（定时触发）等。触发器的选择取决于你想要的 CI/CD 触发条件。

```
on:
  push:
    branches:
      - main
  pull_request:
    branches:
      - main
  schedule:
    - cron: '0 0 * * *'

```

##### 环境（Environments）

GitHub Actions 允许你为特定的任务或步骤定义环境。这可以是不同的操作系统（如 Windows、Linux、macOS），也可以是自定义的虚拟环境。这对于需要在不同环境中运行的项目非常有用。

```
jobs:
  build:
    runs-on: ubuntu-latest
  deploy:
    runs-on: windows-latest

```

##### 矩阵构建（Matrix Builds）

矩阵构建允许在不同参数下并行运行同一个工作流。这对于在多个版本、操作系统或配置下测试和构建应用程序非常有用，可以加速整个流程。

```
strategy:
  matrix:
    node-version: [10, 12, 14]

```

##### 缓存（Caching）

GitHub Actions 允许你缓存依赖项，以减少构建和测试的时间。通过缓存，你可以在不重复下载或构建相同依赖项的情况下提高工作流的效率。

```
jobs:
  build:
    runs-on: ubuntu-latest
    steps:
    - name: Checkout code
      uses: actions/checkout@v2
    - name: Set up Node.js
      uses: actions/setup-node@v2
      with:
        node-version: '14'
    - name: Cache dependencies
      uses: actions/cache@v2
      with:
        path: ~/.npm
        key: ${{ runner.os }}-node-${{ hashFiles('**/*.lock') }}
        restore-keys: |
          ${{ runner.os }}-node-
    - name: Install dependencies
      run: npm install

```

##### 自定义操作（Custom Actions）

除了使用 GitHub Actions 提供的内置操作外，你还可以创建自己的自定义操作。这些操作可以在不同的工作流程中重复使用，使得你的配置更加模块化和可维护。  

##### 部署（Deployment）

GitHub Actions 可以与部署目标（如服务器、云服务、容器等）集成，实现自动化部署。使用预定义的 deploy 操作或者自定义脚本，你可以将应用程序快速部署到目标环境。  

##### Secrets

Secrets 允许你安全地存储和使用敏感信息，如 API 密钥、访问令牌等。这些 Secrets 可以在工作流程中被引用，但不会被显示在日志中。

```
jobs:
  deploy:
    runs-on: ubuntu-latest
    steps:
    - name: Deploy to production
      uses: my-custom-deployment-action
      with:
        api-key: ${{ secrets.DEPLOY_API_KEY }}

```

#### 语法

##### 工作流程文件结构

一个 GitHub Actions 的工作流程文件通常包含以下几个部分：

*   **name**： 定义工作流程的名称。

```
name: My CI/CD Workflow

```

*   **on**： 定义触发工作流程的事件，如 push、pull_request 等。

```
on:
  push:
    branches:
      - main

```

*   **jobs**： 定义工作流程中的任务，一个任务可以包含多个步骤。

```
jobs:
  build:
    runs-on: ubuntu-latest
    steps:
      - name: Checkout code
        uses: actions/checkout@v2
      - name: Set up Node.js
        uses: actions/setup-node@v2
        with:
          node-version: '14'

```

##### 步骤（Steps）

步骤定义了工作流程中的具体操作。每个步骤都包含一个或多个命令，可以是运行脚本、使用预定义的操作或自定义的操作。

```
steps:
  - name: Checkout code
    uses: actions/checkout@v2

  - name: Set up Node.js
    uses: actions/setup-node@v2
    with:
      node-version: '14'

```

##### 操作（Actions）

操作是可重用的、独立的任务单元。GitHub Actions 提供了一系列官方操作，也允许用户创建自定义的操作。操作可以通过 uses 字段引入。

```
steps:
  - name: Use a custom action
    uses: ./path/to/my-action

```

##### 矩阵构建（Matrix Builds）

矩阵构建允许在不同参数下并行运行同一个工作流。这在同时测试多个版本或环境时非常有用。

```
strategy:
  matrix:
    node-version: [10, 12, 14]

```

##### 环境变量

可以使用 env 字段定义环境变量，这些变量可以在工作流程的各个步骤中使用。

```
env:
  MY_VARIABLE: 'some value'

```

```
steps:
  - name: Use environment variable
    run: echo $MY_VARIABLE

```

##### 条件（Conditions）

可以使用 if 字段为步骤定义条件，根据条件来决定是否执行该步骤。

```
steps:
  - name: Run only on main branch
    run: echo "Hello, World!"
    if: github.ref == 'refs/heads/main'

```

##### 超时和重试

使用 timeout-minutes 定义步骤的最大执行时间，使用 retry 定义步骤的最大重试次数。

```
steps:
  - name: My step
    run: echo "Hello, World!"
    timeout-minutes: 10
    retries: 3

```

##### 缓存（Caching）

使用 actions/cache 操作可以缓存依赖项，以减少构建和测试的时间。

```
steps:
  - name: Cache dependencies
    uses: actions/cache@v2
    with:
      path: ~/.npm
      key: ${{ runner.os }}-node-${{ hashFiles('**/*.lock') }}
      restore-keys: |
        ${{ runner.os }}-node-

```