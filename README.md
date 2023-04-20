# NanoOJ

尝试一个OJ系统的手写

## 项目概览

目前是纯C/C++ OJ。整个项目使用了 Vue全家桶 + ElementPlus + Golang-Gin + C++ 开发。

>（其实本来想直接用docker当判题机环境的，但是还没做）

- `judge` 判题端，使用 GNU/Linux syscalls 和 C++ 完成，目前还没加输入支持和防止禁用源码系统调用的功能，安全性尚不能用于生产环境。
- `web` 网页端。使用 Vue+ElementUIPlus 搓的单页应用(SPA)。目前完成了大半~~才怪~~，剩下一些~~一堆~~细节。

>想提前试下判题机可以直接把 `judge/core/` 里边那些玩意编译出来（要用Linux编译哦）用用看。
>用法的话，不带参数直接运行程序就会输出。

```tree
NanoOJ/
├── judge/  # 后端目录
│   ├── config/  # 配置文件目录
│   │   ├── config.go  # 读取和处理配置文件的代码
│   │   └── config.yaml  # 配置文件，例如数据库连接信息
│   ├── controller/  # 控制器目录
│   │   ├── auth.go  # 处理用户身份验证相关请求的代码
│   │   ├── problem.go  # 处理问题相关请求的代码
│   │   ├── submission.go  # 处理代码提交相关请求的代码
│   │   ├── testset.go  # 处理测试集相关请求的代码
│   │   └── user.go  # 处理用户相关请求的代码
│   ├── database/  # 数据库目录
│   │   ├── migrate/  # 数据库迁移文件目录
│   │   │   ├── 20220301_init.sql  # 初始化数据库表的 SQL 文件
│   │   │   └── 20220324_add_submission_status_column.sql  # 增加提交状态列的 SQL 文件
│   │   ├── model/  # 数据模型目录
│   │   │   ├── problem.go  # 问题模型
│   │   │   ├── submission.go  # 提交模型
│   │   │   └── user.go  # 用户模型
│   │   ├── database.go  # 数据库连接和初始化代码
│   │   └── repository.go  # 数据库操作的统一接口
│   ├── middleware/  # 中间件目录
│   │   ├── auth.go  # 身份验证中间件
│   │   └── logging.go  # 日志中间件
│   ├── router/  # 路由目录
│   │   ├── middleware.go  # 路由中间件
│   │   └── router.go  # 路由处理代码
│   ├── worker/  # 后台工作目录
│   │   └── judger.go  # 评测机代码
│   ├── main.go  # 后端程序入口文件
│   └── go.mod  # Golang 依赖管理文件
│
└── web/  # 前端目录
    ├── public/  # 静态资源目录
    ├── src/  # 代码目录
    │   ├── api/  # API 相关代码目录
    │   │   ├── auth.js  # 处理用户身份验证相关 API 的代码
    │   │   ├── problem.js  # 处理问题相关 API 的代码
    │   │   ├── submission.js  # 处理提交相关 API 的代码
    │   │   ├── testset.js  # 处理测试集相关 API 的代码
    │   │   └── user.js # 
    │   ├── assets/  # 静态资源目录
    │   ├── components/  # 组件目录
    │   ├── router/  # 路由目录
    │   ├── store/  # Vuex 相关代码目录
    │   │   ├── modules/  # Vuex 模块目录
    │   │   ├── actions.js  # Vuex actions 代码
    │   │   ├── getters.js  # Vuex getters 代码
    │   │   ├── index.js  # Vuex 入口文件
    │   │   ├── mutations.js  # Vuex mutations 代码
    │   │   └── state.js  # Vuex state 代码
    │   ├── views/  # 视图目录
    │   ├── App.vue  # Vue 根组件
    │   ├── main.js  # Vue 程序入口文件
    │   ├── router.js  # Vue 路由配置文件
    │   └── store.js  # Vuex 状态管理文件
    ├── .eslintrc.js  # ESLint 配置文件
    ├── babel.config.js  # Babel 配置文件
    ├── package.json  # NPM 依赖管理文件
    ├── postcss.config.js  # PostCSS 配置文件
    ├── README.md  # 项目说明文件
    ├── vue.config.js  # Vue CLI 配置文件
    └── yarn.lock  # Yarn 依赖管理锁文件
```

### 整体设计

整个系统的开发核心，应该是功能需求。首先根据功能需求来设计前后端的数据库，以及各自的CRUD代码。然后是后端的各功能业务逻辑代码，以及后端的API组件，鉴权组件，日志组件等。前端类似，先是编写好API代码，然后是前端路由，页面数据模型设计，业务代码实现，最后是UI、表单等的实现。

当然，后端还有其他的重要功能，比如数据库的导出导入功能，和命令行的参数交互功能（启动参数，数据库dump参数等），使用简介，检查更新（？）等。

完成后，打包编译。可以将前端编译后的静态文件直接集成到后端中（类似Gitea的实现）。这样可执行文件的部署会非常方便（单文件），也便于部署在Docker镜像中。

### 后端设计

后端为二合一设计，包含了判题端和网页端的后端。这两种模式根据配置文件来编辑修改，且可以同时开启两种模式，并规定不同的端口。

配置文件不光可以手动修改文件，还可以通过网页端的后台管理界面来修改，且网页端还可以编辑其他判题端的配置文件。

后端支持多种数据库，默认使用内置的SQLite3，也支持MySQL等。后一种模式需要指定正确的数据库账号密码。

后端的判题系统包含一个判题镜像，这镜像中包含了判题核心程序，以及各种编译器。目前的设计是所有语言共用一个判题镜像，后期计划将判题镜像分解成多个，每个语言一个。这样，以后想加上诸如OS-Lab镜像也会很简单。判题的模式大概是，将测试用例和用户代码存储到一个目录中，然后运行判题镜像，并将这个目录挂载在特定位置。判题镜像会自动运行判题核心程序，来判定用户的代码是否正确。

关键就是判题镜像的管理和设计。应该允许自定义判题机，这样可以让用户自己设计判题机，来适应自己的需求。判题机的设计应该是一个可执行文件，它的输入是一个目录，这个目录中包含了测试用例和用户代码。**判题机按照规定格式输出判题结果**到stdout，供后端读取。或者也可以直接输出到一个文件中，供后端读取。

按照上面规定的设计方式，后端的判题系统可以这么设计：包含一个默认的Dockerfile，默认隐藏，在高级选项界面提供自定义选项。给出环境配置的选项（选择安装哪些编译器和工具），并且让用户上传判题机（要求**接收几个固定的参数**来完成判题），根据这些来让管理员创建一个判题机镜像。允许管理员查看判题机镜像列表，以及分布式判题集群状态。从镜像创建实例的时候，挂载待判定的目录，以及测试用例的目录，来让后端直接判定固定目录中的内容，实现镜像自动判题。

这也要求重新设计后端的判题队列分配系统。现在的做法是，开一个检查进程来将所有Pending的提交加入判题队列，再开一个进程，根据队列是否为空来循环启动判题机，直到队列为空。可以更改实现为：前者不变，整个系统使用同一个判题队列。后者在启动判题机的时候，从所有判题服务器中选择一个负载最低的机器（这就要求判题系统增加一个状态查询接口。哦对了，后端进程也加上状态查询接口吧），从队列中拿出一个submission，并以它为参数，启动判题进程。注意，判题进程目前还是运行在后端系统内，目标判题机是作为参数传递给判题进程。判题进程依然位于后端系统上，它负责准备工作，会将待判目录作为判题机的启动参数；在这之后，启动判题机时，才会用到目标判题机这个参数，让目标判题机去启动判题容器实例，并返回判题结果，由后端系统的判题进程负责结果的解析，以及数据的更改和数据库状态的写回。

>关于OS-Lab：这是后期计划添加的一个功能。OS-Lab可以理解成一个虚拟机，其中可以进行OS实验（使用内置的Qemu模拟器）。它的帕判题方式是，让用户在Lab里完成操作系统的代码编写（将所有代码按照规范放置在一个规定的目录），然后后端会启动另一个镜像，并且将用户的代码目录挂载到这个镜像中。另一个镜像中包含了精心设计的判题机，判题机会作为镜像的入口程序自动运行，来判定用户的代码是否正确。这是OS-Lab的大致实现思路。和上面的判题机设计是完全类似的，只不过这次提供的不是一个实体目录，而是为用户准备一个实验环境的Docker容器，在编写完成后，直接将这个容器的目录挂载到判题机镜像实例的目录中，来实现判题。前者是提交代码，后者是提交容器实例。

上面的镜像实例，最好全部保持断网，防止用户在判题机中执行恶意代码。同时，在后端系统中，如果发现判题机实例状态异常（比如持续时间过长，输出格式不对等），直接杀死判题机进程并删除判题机实例，并更新提交状态为TLE（Time Limit Exceed）。这个也可以试试能不能在判题机启动参数里设置。以及资源占用问题，限制判题机实例的资源占用情况，防止出现一个恶意提交卡死整个判题系统的情况出现。

另外，整个后端程序其实也可以运行在Docker容器中。这可以在组织后端系统的架构时，使用docker-compose来组织：后端进程跑在一个机器上，开放一个范围的端口；作为判题机的后端进程跑在其他机器上，开放一个固定的端口；然后将它们的端口映射到一起。数据库如果是SQLite3，可以使用Volume挂载在一个目录下；如果是MySQL，则可以另外再开一个Docker容器，将它的端口映射到后端容器上。然后再将后端的Web端口映射到主机的端口，或者是再用nginx做反向代理/负载均衡等。

后端的编写顺序和前面的总览一样，自底向上开发。

### 前端设计

按照前端单页程序的标准设计思路实现。要求有主页，题库，练习/竞赛，提交记录，用户中心，排行榜，后台管理系统等页面。最后一项可以是另外一个单页应用。交互逻辑就是，用户在主页中选择题目，然后进入题目页面；选择练习或者竞赛，然后进入练习/竞赛页面，再进入题目页面。在题目页面中，可以编写代码，提交代码，查看提交记录等操作。用户中心可以查看用户信息，修改密码，查看提交记录等。后台管理系统可以管理题目，管理用户，管理练习/竞赛，查看提交记录所有信息等。

然后是前端架构设计。一个前端的总数据库，一组和后端通信的API工具，以及一组数据库操作的CRUD函数库（就是Vuex里边那些mutations, actions之类的）。然后各个页面调用这些数据库函数来获取/更新/删除数据。比如，管理员界面要创建一个问题，就首先在页面里用表单等组件构造好一个问题的实体。然后使用API接口发送数据到后端，拿到结果后，如果是200-OK，就再调用API接口获取后端数据库，来更新前端数据库，从而更新界面的数据展示。这一套很不错，但是最大的问题是mutations和actions太多了，很容易就会写一大堆的这些玩意，导致代码臃肿，而且还不咋通用。所以嘛。。如果太复杂了，还是把每个数据结构的定义和更改函数分开到每个文件里去，保证代码结构的清晰，最后再统一导入到index.js中来使用。

此外，前端还有个数据持久化的问题。尤其是API Key，打开页面的时候首先需要校验Key是否过期，如果过期了就删除Key，要求重新登录。否则，直接用Key就行了。同时得向后端申请个新的Key换上，以防用着用着Key突然过期了。

路由设计也是个重点。有些参数（比如id这种）一般通路由参数传递，需要保证参数匹配。

在写完Vuex和API部分的代码，就可以去写前端的逻辑代码和页面代码了。开发应该是以页面为单位，在App.vue中完成初始化，然后用Vue-router来切换不同页面。页面的开发流程基本就是根据功能设计页面的数据结构，然后再编写好页面的CRUD相关的函数，最后编写UI，表单，按钮等组件，最后再进行功能和数据的对接。在处理一些数据的时候，会用到一些辅助函数，可以把它们放到utils下，方便复用。前端的设计也需要考虑到用户的权限问题，根据权限，展示不同的内容。比如管理员就能看到问题提交记录的详细信息，测试用例的内容，系统日志等。

组件化的开发也不错，一个模块只需要定义好数据结构，以及对传入数据的操作，就能完成一个小型的组件。其他模块只需要调用这个组件，就能完成一个模块的功能。这一点，在面对很多表单（比如管理员界面）的时候就很有用，比如对问题/练习/通知/用户/提交记录等都需要进行新增/修改/删除的功能。这就可以通过组件来进行代码复用，来减少代码冗杂度。

前端的开发总体上很简单，用Vuex，Vue-router和UI库elementPlus就能搞定。比较复杂的部分，一个是前端代码编辑器（问题详情提交代码处用的），另一个是前端

下面是各个页面的设计思路。

#### 主页

没啥说的，就是几个组件拿出来拼一下，整个状态展示和公告栏。

#### 题库

就是所有题目列表。顺道，一些题目可以从主题库隐藏，因此可以给题目加上所属组的属性：属于默认组，就会出现在这里。同时，也可以给题目增加标签设计，便于分类查看。另外也可以加上搜素功能（虽然没啥用...），查找题目用。另外这个组件可以做成可复用的形式，便于展示不同过滤选项的问题列表。

#### 排行榜

也是个可复用组件，有总榜（根据用户积分决定），对于每个练习/竞赛有分榜，对于每个题目也有分榜。另外加上封榜功能，为OI赛制提供支持。

#### 题目详情

需要好好设计的组件。对于每个题，有时间/空间限制选项，有提交代码/提交文件的选项，还支持填空题。根据题目可使用语言的限制，提供代码编辑器/文本框（填空题）的选线。同时，展示题目的历史提交记录展示列表，以及当前状态（分数，是否通过等），其他信息（问题描述）的展示。

#### 练习/竞赛详情

基本是几个组件的组合。对于一个练习，有题目列表，练习信息、排行榜等几个组件。另外还有前面提到的封榜功能，再加上做题进度的功能，还有提交记录查看等小功能。都是组件复用，实现起来应该比较快。

## 计划

- 先把数据库设计好吧（摊
- 啊，然后是前后端功能对接，也稍微再完善一下吧
- 然后嘛......还没有写管理员相关的功能，也稍微写一下吧
- 话说 `web/src/views/` 和 `components/` 乱的一团糟，那也稍微收拾下吧
- 还有好多好看的 ElementUIPlus 组件没用呢
- 话说判题端也得加上读取stdin的功能了
- 以及得想办法禁了源码的系统调用，不然服务器一定会炸的
- 多线程支持也加上吧，大概能快一些吧
- ~~总感觉好累啊 那就有时间再摸吧~~
