### 一、服务信息
----
services
├── logic-engine-coordinator
└── logic-engine-worker

###### logic-engine-coordinator
数据编排调度服务

###### logic-engine-worker
数据编排worker, 节点执行逻辑都在这

### 二、环境信息和访问方式
-----
#### 开发环境

容器平台地址: 

	logic-engine-worker: https://kubernetes.woa.com/v4/projects/prj8k8cf/workloads/cls-q11hd6gg/ns-prj8k8cf-1184133-test/StatefulSetPlus/logic-engine-worker

	logic-engine-coordinator: https://kubernetes.woa.com/v4/projects/prj8k8cf/workloads/cls-q11hd6gg/ns-prj8k8cf-1184133-test/StatefulSetPlus/logic-engine-coordinator

流水线: 

	https://devops.woa.com/console/pipeline/welink/p-92cd56abe1c0458a93803d4c556ab412/history

发布方式:

1、流水线发布
2、本地发布脚本发布，sh deploy/deploy_tke.sh {服务名}, ex: sh deploy/deploy_tke.sh logic-engine-worker
第二种方式发布服务需要 本地docker登录csig镜像平台

#### 测试环境

容器平台地址: 

	logic-engine-worker: https://rancher.welink.qq.com/p/c-f7g6x:p-5pkzt/workload/statefulset:verify-center:logic-engine-worker	

	logic-engine-coordinator: https://rancher.welink.qq.com/p/c-f7g6x:p-5pkzt/workload/statefulset:verify-center:logic-engine-coordinator

流水线:

	https://devops.woa.com/console/pipeline/welink/p-8ab962a7cf4a45b39a634d3b51e1bbb4/history

发布方式:
开发没有权限， 只能运维同学发布



