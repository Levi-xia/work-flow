
# 工作流引擎

节点公共对象
    key
    name
节点对象
    code
    name
    inputs
    outputs
    execute()
边对象
    code
    name
    source
    target
    expr
    execute()

定义对象
    code
    name
    version
    content
    store
实例对象
    define
    status
    variables
    store
任务对象
    code
    name
    instance
    status
    variables
    store
流程对象
    nodes
    edges
    tasks
    getStart()

执行对象
    processTask
    variables

定义持久化接口
    createProcessDefine()
    deleteProcessDefine()
    getProcessDefine()

实例持久化接口
    createProcessInstance()
    finishProcessInstance()
    getProcessInstance()
    addVariables()
    deleteVariables()

任务持久化接口
    createProcessTask()
    finishProcessTask()
    getProcessTask()



任务创建Hook
任务前后置拦截器
执行权限
Store MySQL实现
Testing 测试
